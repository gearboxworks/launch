package helperDocker

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"launch/defaults"
	"launch/ux"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"
)


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

	return &_args
}


func (gear *DockerGear) ContainerSsh(interactive bool, statusLine bool, mountPath string, cmdArgs []string) *ux.State {
	if state := gear.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
		// Get Docker container SSH port.
		var clientPort string
		clientPort, gear.State = gear.Container.GetContainerSsh()
		if gear.State.IsError() {
			break
		}
		if clientPort == "" {
			gear.State.SetError("no SSH port in gear")
			break
		}

		u := url.URL{}
		err := u.UnmarshalBinary([]byte(gear.Client.DaemonHost()))
		if err != nil {
			gear.State.SetError("error finding SSH port: %s", err)
			break
		}


		// Create SSH client config.
		// fmt.Printf("Connect to %s:%s\n", u.Hostname(), port)
		gear.Ssh = NewSshClient(SshClientArgs {
			ClientAuth: &SshAuth {
				Host:      u.Hostname(),
				Port:      clientPort,
				Username:  DefaultUsername,
				Password:  DefaultPassword,
			},
			StatusLine: StatusLine {
				Enable: statusLine,
			},
			Shell:       interactive,
			GearName:    gear.Container.Name,
			GearVersion: gear.Container.Version,
			CmdArgs:     cmdArgs,
			State:       ux.NewState(gear.Debug),
		})


		// @TODO - Add remote host capability here!
		// Run server for SSHFS if required.
		gear.State = gear.SetMountPath(mountPath)
		if gear.State.IsOk() {
			err = gear.Ssh.InitServer()
			if err == nil {
				//noinspection GoUnhandledErrorResult
				go gear.Ssh.StartServer()

				// GEARBOX_MOUNT_HOST=10.0.5.57
				// GEARBOX_MOUNT_PATH=/Users/mick/.gearbox
				// GEARBOX_MOUNT_PORT=49410
				//time.Sleep(time.Second * 5)
				//for ; gear.Ssh.ServerAuth == nil; {
				//	time.Sleep(time.Second)
				//}

				err = os.Setenv("GEARBOX_MOUNT_HOST", gear.Ssh.ServerAuth.Host)
				err = os.Setenv("GEARBOX_MOUNT_PORT", gear.Ssh.ServerAuth.Port)
				err = os.Setenv("GEARBOX_MOUNT_USER", gear.Ssh.ServerAuth.Username)
				err = os.Setenv("GEARBOX_MOUNT_PASSWORD", gear.Ssh.ServerAuth.Password)
				err = os.Setenv("GEARBOX_MOUNT_PATH", gear.Ssh.FsMount)
			}
		}


		// Process env
		gear.State = gear.Ssh.getEnv()
		if err != nil {
			break
		}


		// Connect to container SSH - retry 5 times.
		for i := 0; i < 5; i++ {
			gear.State.ClearError()
			err = gear.Ssh.Connect()
			if err == nil {
				break
			}

			switch v := err.(type) {
				case *ssh.ExitError:
					gear.State.SetExitCode(v.Waitmsg.ExitStatus())
					if len(cmdArgs) == 0 {
						gear.State.SetError("Command exited with error code %d", v.Waitmsg.ExitStatus())
					} else {
						gear.State.SetError("Command '%s' exited with error code %d", cmdArgs[0], v.Waitmsg.ExitStatus())
					}
					i = 5
					continue

				default:
					gear.State.SetError("SSH to Gear %s:%s failed.", gear.Container.Name, gear.Container.Version)
			}
			time.Sleep(time.Second)
		}
	}

	return gear.State
}


func (gear *DockerGear) SetMountPath(mp string) *ux.State {
	if state := gear.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
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
					gear.State.SetError(err)
					break
				}
				gear.State.SetOk()
				gear.Ssh.FsMount = cwd

			case mp == defaults.DefaultPathHome:
				var u *user.User
				u, err = user.Current()
				if err != nil {
					gear.State.SetError(err)
					break
				}
				gear.State.SetOk()
				gear.Ssh.FsMount = u.HomeDir

			default:
				mp, err = filepath.Abs(mp)
				if err != nil {
					gear.State.SetError(err)
					break
				}
				gear.State.SetOk()
				gear.Ssh.FsMount = mp
		}
	}

	return gear.State
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
