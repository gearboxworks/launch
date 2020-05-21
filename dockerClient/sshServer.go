package dockerClient

import (
	"fmt"
	"io"
	"io/ioutil"
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
//	for range OnlyOnce {
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
//			ux.PrintflnRed("SSHFS SERVER: %s\n", err)
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
func (s *Ssh) InitServer() error {
	var err error

	for range OnlyOnce {
		// An SSH server is represented by a ServerConfig, which holds
		// certificate details and handles authentication of ServerConns.
		s.ServerConfig = &ssh.ServerConfig{
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
			ux.PrintflnRed("SSHFS SERVER: %s", err)
			break
		}

		s.ServerConfig.AddHostKey(private)

		// Once a ServerConfig has been configured, connections can be
		// accepted.
		s.ServerListener, err = net.Listen("tcp", "0.0.0.0:0")
		if err != nil {
			ux.PrintflnRed("SSHFS SERVER: listener ERROR - %s", err)
			break
		}
		if s.Debug {
			ux.Printf("Listening on %v\n", s.ServerListener.Addr())
		}

		//err = os.Setenv("SSHFS_HOST", s.ServerListener.Addr().String())
		s.ServerAuth = NewSshAuth()
		//hp := strings.Split(s.ServerListener.Addr().String(), ":")
		switch addr := s.ServerListener.Addr().(type) {
			case *net.UDPAddr:
				s.ServerAuth.Host = "host.docker.internal" // addr.IP.String()
				s.ServerAuth.Port = fmt.Sprintf("%d", addr.Port)
				//p.DstPort = uint(localAddr.(*net.UDPAddr).Port)
			case *net.TCPAddr:
				s.ServerAuth.Host = "host.docker.internal" // addr.IP.String()
				s.ServerAuth.Port = fmt.Sprintf("%d", addr.Port)
				//p.DstPort = uint(localAddr.(*net.TCPAddr).Port)
		}
		s.ServerAuth.Username =  s.ClientAuth.Username
		s.ServerAuth.Password =  s.ClientAuth.Password
		s.ServerAuth.PublicKey = s.ClientAuth.PublicKey
	}

	return err
}

func (s *Ssh) StartServer() error {
	var err error

	for range OnlyOnce {
		debugStream := ioutil.Discard
		if s.Debug {
			debugStream = os.Stderr
		}

		s.ServerConnection, err = s.ServerListener.Accept()
		if err != nil {
			ux.PrintflnRed("SSHFS SERVER: listener accept ERROR - %s", err)
			break
		}

		// Before use, a handshake must be performed on the incoming
		// net.Conn.
		var chans <-chan ssh.NewChannel
		var reqs <-chan *ssh.Request
		_, chans, reqs, err = ssh.NewServerConn(s.ServerConnection, s.ServerConfig)
		if err != nil {
			ux.PrintflnRed("SSHFS SERVER: handshake ERROR - %s", err)
			break
		}
		_, _ = fmt.Fprintf(debugStream, "SSH server established\n")

		// The incoming Request channel must be serviced.
		go ssh.DiscardRequests(reqs)

		// Service the incoming Channel channel.
		for newChannel := range chans {
			// Channels have a type, depending on the application level
			// protocol intended. In the case of an SFTP session, this is "subsystem"
			// with a payload string of "<length=4>sftp"
			_, _ = fmt.Fprintf(debugStream, "Incoming channel: %s\n", newChannel.ChannelType())
			if newChannel.ChannelType() != "session" {
				_ = newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
				_, _ = fmt.Fprintf(debugStream, "Unknown channel type: %s\n", newChannel.ChannelType())
				continue
			}

			var requests <-chan *ssh.Request
			var channel ssh.Channel
			channel, requests, err = newChannel.Accept()
			if err != nil {
				log.Fatal("could not accept channel.", err)
			}
			_, _ = fmt.Fprintf(debugStream, "Channel accepted\n")

			// Sessions have out-of-band requests such as "shell",
			// "pty-req" and "env".  Here we handle only the
			// "subsystem" request.
			go func(in <-chan *ssh.Request) {
				for req := range in {
					_, _ = fmt.Fprintf(debugStream, "Request: %v\n", req.Type)
					ok := false

					switch req.Type {
						case "subsystem":
							_, _ = fmt.Fprintf(debugStream, "Subsystem: %s\n", req.Payload[4:])
							if string(req.Payload[4:]) == "sftp" {
								ok = true
							}
					}

					_, _ = fmt.Fprintf(debugStream, " - accepted: %v\n", ok)
					_ = req.Reply(ok, nil)
				}
			}(requests)

			serverOptions := []sftp.ServerOption{
				sftp.WithDebug(debugStream),
			}

			if s.FsReadOnly {
				serverOptions = append(serverOptions, sftp.ReadOnly())
				_, _ = fmt.Fprintf(debugStream, "Read-only server\n")
			} else {
				_, _ = fmt.Fprintf(debugStream, "Read write server\n")
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
				_ = server.Close()
				ux.PrintfOk("SSHFS SERVER: exited OK\n")
				break
			} else if err != nil {
				ux.PrintflnRed("SSHFS SERVER: exit ERROR - %s", err)
				break
			}
		}
	}

	return err
}
