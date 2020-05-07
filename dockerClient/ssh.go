package dockerClient

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"launch/defaults"
	"launch/only"
	"launch/ux"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
)

//type Ssh struct {
//	Server *SshServer
//	Client *SshClient
//}

type SshAuth struct {
	// SSH related.
	Username    string
	Password    string
	Host        string
	Port        string
	PublicKey   string
}
type SshAuthArgs SshAuth

func NewSshAuth(args ...SshAuth) *SshAuth {

	var _args SshAuth
	if len(args) > 0 {
		_args = args[0]
	}

	if _args.Username == "" {
		_args.Username = DefaultUsername
	}

	if _args.Password == "" {
		_args.Password = DefaultPassword
	}

	if _args.PublicKey == "" {
		_args.PublicKey = DefaultKeyFile
	}

	if _args.Host == "" {
		_args.Host = DefaultSshHost
	}

	if _args.Port == "" {
		_args.Port = DefaultSshPort
	}

	//sshAuth := &SshAuth{}
	//*sshAuth = SshAuth(_args)

	// Query VB to see if it exists.
	// If not return nil.

	return &_args
}


func (me *DockerGear) ContainerSsh(interactive bool, statusLine bool, mountPath string, cmdArgs []string) ux.State {
	var state ux.State

	for range only.Once {
		// Get Docker container SSH port.
		var clientPort string
		clientPort, state = me.Container.GetContainerSsh()
		if state.IsError() {
			break
		}
		if clientPort == "" {
			state.SetError("no SSH port in gear")
			break
		}

		u := url.URL{}
		var err error
		err = u.UnmarshalBinary([]byte(me.Client.DaemonHost()))
		if err != nil {
			state.SetError("error finding SSH port: %s", err)
			break
		}


		// Create SSH client config.
		// fmt.Printf("Connect to %s:%s\n", u.Hostname(), port)
		me.Ssh = NewSshClient(SshClientArgs {
			ClientAuth: &SshAuth {
				Host:      u.Hostname(),
				Port:      clientPort,
				Username:  DefaultUsername,
				Password:  DefaultPassword,
			},
			StatusLine: StatusLine {
				Enable: statusLine,
			},
			Shell: interactive,
			GearName: me.Container.Name,
			GearVersion: me.Container.Version,
			CmdArgs: cmdArgs,
		})


		// Run server for SSHFS if required.
		if me.SetMountPath(mountPath) {
			err = me.Ssh.InitServer()
			if err == nil {
				go me.Ssh.StartServer()

				// GEARBOX_MOUNT_HOST=10.0.5.57
				// GEARBOX_MOUNT_PATH=/Users/mick/.gearbox
				// GEARBOX_MOUNT_PORT=49410
				//time.Sleep(time.Second * 5)
				//for ; me.Ssh.ServerAuth == nil; {
				//	time.Sleep(time.Second)
				//}

				err = os.Setenv("GEARBOX_MOUNT_HOST", me.Ssh.ServerAuth.Host)
				err = os.Setenv("GEARBOX_MOUNT_PORT", me.Ssh.ServerAuth.Port)
				err = os.Setenv("GEARBOX_MOUNT_USER", me.Ssh.ServerAuth.Username)
				err = os.Setenv("GEARBOX_MOUNT_PASSWORD", me.Ssh.ServerAuth.Password)
				err = os.Setenv("GEARBOX_MOUNT_PATH", me.Ssh.FsMount)
			}
		}


		// Process env
		err = me.Ssh.getEnv()
		if err != nil {
			break
		}


		// Connect to container.
		err = me.Ssh.Connect()
		if err != nil {
			switch v := err.(type) {
				case *ssh.ExitError:
					state.SetExitCode(v.Waitmsg.ExitStatus())
					if len(cmdArgs) == 0 {
						state.SetError("Command exited with error code %d", v.Waitmsg.ExitStatus())
					} else {
						state.SetError("Command '%s' exited with error code %d", cmdArgs[0], v.Waitmsg.ExitStatus())
					}
			}
			break
		}
	}

	return state
}


func (me *DockerGear) SetMountPath(mp string) bool {
	var ok bool

	for range only.Once {
		var err error
		var cwd string

		if mp == defaults.DefaultPathNone {
			break
		}

		switch {
			case mp == defaults.DefaultPathEmpty:
				fallthrough
			case mp == defaults.DefaultPathCwd:
				cwd, err = os.Getwd()
				if err != nil {
					break
				}
				ok = true
				mp = cwd

			case mp == defaults.DefaultPathHome:
				var u *user.User
				u, err = user.Current()
				if err != nil {
					break
				}
				ok = true
				mp = u.HomeDir

			default:
				mp, err = filepath.Abs(mp)
				if err != nil {
					break
				}
				ok = true
		}

		if err != nil {
			break
		}

		if ok == true {
			me.Ssh.FsMount = mp
		}
	}

	return ok
}

func readPublicKeyFile(file string) (ssh.AuthMethod, error) {

	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		// fmt.Printf("# Error reading file '%s': %s\n", file, err)
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		// fmt.Printf("# Error parsing key '%s': %s\n", signer, err)
		return nil, err
	}

	return ssh.PublicKeys(signer), err
}
