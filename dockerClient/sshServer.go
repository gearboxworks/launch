package dockerClient

import (
	"fmt"
	"io"
	"io/ioutil"
	"launch/only"
	"launch/ux"
	"log"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// GEARBOX_MOUNT_PATH=/Users/mick/.gearbox GEARBOX_MOUNT_HOST=10.0.5.57 GEARBOX_MOUNT_PORT=49410 ./ssh-mount.sh &

//type SshServer struct {
//	Config      *ssh.ServerConfig
//	Listener    net.Listener
//	Connection  net.Conn
//
//	// SSH related.
//	Auth        *SshAuth
//	//Username    string
//	//Password    string
//	//Host        string
//	//Port        string
//	//PublicKey   string
//
//	// FUSE related
//	ReadOnly    bool
//	Debug       bool
//}
//type SshServerArgs SshServer

//func NewSshServer(args ...SshServerArgs) *SshServer {
//	var err error
//	me := &SshServer{}
//
//	for range only.Once {
//		var _args SshServerArgs
//		if len(args) > 0 {
//			_args = args[0]
//		}
//
//		_args.Auth = NewSshAuth(*_args.Auth)
//
//		// An SSH server is represented by a ServerConfig, which holds
//		// certificate details and handles authentication of ServerConns.
//		_args.Config = &ssh.ServerConfig{
//			PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
//				// Should use constant-time compare (or better, salt+hash) in
//				// a production setting.
//				if _args.Debug {
//					_, _ = fmt.Fprintf(os.Stderr, "Login: %s\n", c.User())
//				}
//				if c.User() == DefaultUsername && string(pass) == DefaultPassword {
//					return nil, nil
//				}
//				return nil, fmt.Errorf("password rejected for %q", c.User())
//			},
//		}
//
//		var privateBytes []byte
//		privateBytes, err = ioutil.ReadFile("id_rsa")
//		if err != nil {
//			//log.Fatal("Failed to load private key", err)
//			privateBytes = []byte(SshHostPrivateKey)
//		}
//
//		var private ssh.Signer
//		private, err = ssh.ParsePrivateKey(privateBytes)
//		if err != nil {
//			ux.PrintfRed("SSHFS SERVER: %s\n", err)
//			break
//		}
//
//		_args.ServerConfig.AddHostKey(private)
//
//		*me = SshServer(_args)
//	}
//
//	return me
//}


func SshAuthenticate(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	// Should use constant-time compare (or better, salt+hash) in
	// a production setting.
	//if _args.Debug {
	//	_, _ = fmt.Fprintf(os.Stderr, "Login: %s\n", c.User())
	//}

	if c.User() == DefaultUsername && string(pass) == DefaultPassword {
		return nil, nil
	}

	return nil, nil
	//return nil, fmt.Errorf("password rejected for %q", c.User())
}


// Based on example server code from golang.org/x/crypto/ssh and server_standalone
func (me *Ssh) InitServer() error {
	var err error

	for range only.Once {
		// An SSH server is represented by a ServerConfig, which holds
		// certificate details and handles authentication of ServerConns.
		me.ServerConfig = &ssh.ServerConfig{
			PasswordCallback: SshAuthenticate,
		}

		var privateBytes []byte
		privateBytes, err = ioutil.ReadFile("id_rsa")
		if err != nil {
			//log.Fatal("Failed to load private key", err)
			privateBytes = []byte(SshHostPrivateKey)
		}

		var private ssh.Signer
		private, err = ssh.ParsePrivateKey(privateBytes)
		if err != nil {
			ux.PrintfRed("SSHFS SERVER: %s\n", err)
			break
		}

		me.ServerConfig.AddHostKey(private)

		// Once a ServerConfig has been configured, connections can be
		// accepted.
		me.ServerListener, err = net.Listen("tcp", "0.0.0.0:0")
		if err != nil {
			ux.PrintfRed("SSHFS SERVER: listener ERROR - %s\n", err)
			break
		}
		if me.Debug {
			ux.Printf("Listening on %v\n", me.ServerListener.Addr())
		}

		//err = os.Setenv("SSHFS_HOST", me.ServerListener.Addr().String())
		me.ServerAuth = NewSshAuth()
		//hp := strings.Split(me.ServerListener.Addr().String(), ":")
		switch addr := me.ServerListener.Addr().(type) {
			case *net.UDPAddr:
				me.ServerAuth.Host = "host.docker.internal"	// addr.IP.String()
				me.ServerAuth.Port = fmt.Sprintf("%d", addr.Port)
				//p.DstPort = uint(localAddr.(*net.UDPAddr).Port)
			case *net.TCPAddr:
				me.ServerAuth.Host = "host.docker.internal"	// addr.IP.String()
				me.ServerAuth.Port = fmt.Sprintf("%d", addr.Port)
				//p.DstPort = uint(localAddr.(*net.TCPAddr).Port)
		}
		me.ServerAuth.Username =  me.ClientAuth.Username
		me.ServerAuth.Password =  me.ClientAuth.Password
		me.ServerAuth.PublicKey = me.ClientAuth.PublicKey
	}

	return err
}

func (me *Ssh) StartServer() error {
	var err error

	for range only.Once {
		debugStream := ioutil.Discard
		if me.Debug {
			debugStream = os.Stderr
		}

		me.ServerConnection, err = me.ServerListener.Accept()
		if err != nil {
			ux.PrintfRed("SSHFS SERVER: listener accept ERROR - %s\n", err)
			break
		}

		// Before use, a handshake must be performed on the incoming
		// net.Conn.
		var chans <-chan ssh.NewChannel
		var reqs <-chan *ssh.Request
		_, chans, reqs, err = ssh.NewServerConn(me.ServerConnection, me.ServerConfig)
		if err != nil {
			ux.PrintfRed("SSHFS SERVER: handshake ERROR - %s\n", err)
			break
		}
		fmt.Fprintf(debugStream, "SSH server established\n")

		// The incoming Request channel must be serviced.
		go ssh.DiscardRequests(reqs)

		// Service the incoming Channel channel.
		for newChannel := range chans {
			// Channels have a type, depending on the application level
			// protocol intended. In the case of an SFTP session, this is "subsystem"
			// with a payload string of "<length=4>sftp"
			fmt.Fprintf(debugStream, "Incoming channel: %s\n", newChannel.ChannelType())
			if newChannel.ChannelType() != "session" {
				newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
				fmt.Fprintf(debugStream, "Unknown channel type: %s\n", newChannel.ChannelType())
				continue
			}

			var requests <-chan *ssh.Request
			var channel ssh.Channel
			channel, requests, err = newChannel.Accept()
			if err != nil {
				log.Fatal("could not accept channel.", err)
			}
			fmt.Fprintf(debugStream, "Channel accepted\n")

			// Sessions have out-of-band requests such as "shell",
			// "pty-req" and "env".  Here we handle only the
			// "subsystem" request.
			go func(in <-chan *ssh.Request) {
				for req := range in {
					fmt.Fprintf(debugStream, "Request: %v\n", req.Type)
					ok := false

					switch req.Type {
					case "subsystem":
						fmt.Fprintf(debugStream, "Subsystem: %s\n", req.Payload[4:])
						if string(req.Payload[4:]) == "sftp" {
							ok = true
						}
					}

					fmt.Fprintf(debugStream, " - accepted: %v\n", ok)
					req.Reply(ok, nil)
				}
			}(requests)

			serverOptions := []sftp.ServerOption{
				sftp.WithDebug(debugStream),
			}

			if me.FsReadOnly {
				serverOptions = append(serverOptions, sftp.ReadOnly())
				fmt.Fprintf(debugStream, "Read-only server\n")
			} else {
				fmt.Fprintf(debugStream, "Read write server\n")
			}

			var server *sftp.Server
			server, err = sftp.NewServer(
				channel,
				serverOptions...,
			)
			if err != nil {
				log.Fatal(err)
			}

			err = server.Serve()
			if err == io.EOF {
				server.Close()
				ux.PrintfOk("SSHFS SERVER: exited OK\n")
				break
			} else if err != nil {
				ux.PrintfRed("SSHFS SERVER: exit ERROR - %s\n", err)
				break
			}
		}
	}

	return err
}
