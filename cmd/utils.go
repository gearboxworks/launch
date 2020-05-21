package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"launch/defaults"
	"launch/gear"
	"launch/ux"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func showArgs(cmd *cobra.Command, args []string) {
	// var err error

	for range OnlyOnce {
		//if !debugFlag {
		//	break
		//}

		flargs := cmd.Flags().Args()
		if flargs != nil {
			ux.Printf("'%s' called with '%s'\n", cmd.CommandPath(), strings.Join(flargs, " "))
			break
		}

		ux.Printf("'%s' called with '%s'\n", cmd.CommandPath(), strings.Join(args, " "))
		break
	}

	fmt.Println("")
}

type GearArgs struct {
	Name      string
	Version   string

	Args      []string
	Project   string
	Mount     string
	Temporary bool
	SshStatus bool
	Quiet     bool
	Debug     bool
	NoCreate  bool

	Provider  gear.Provider
	GearRef   gear.Gear

	Valid     bool
	State     *ux.State
}


func (ga *GearArgs) IsNil() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}
	ga.State = ga.State.EnsureNotNil()
	return ga.State
}

func (ga *GearArgs) IsValid() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range OnlyOnce {
		if !ga.Valid {
			ga.State.SetError("gear args is not valid")
			break
		}

		ga.State = ga.State.EnsureNotNil()
	}

	return ga.State
}


func (ga *GearArgs) ProcessArgs(cmd *cobra.Command, args []string) *ux.State {
	for range OnlyOnce {
		ga.State = ux.NewState(false)

		if ga.Valid {
			break
		}

		var err error
		fl := cmd.Flags()

		ga.Debug, err = fl.GetBool(argDebug)
		if err != nil {
			ga.Debug = false
		}
		ga.State.DebugSet(ga.Debug)


		ga.Args = args
		if len(ga.Args) > 0 {
			ga.Name = ga.Args[0]
			if strings.Contains(ga.Name, ":") {
				spl := strings.Split(ga.Name, ":")
				ga.Name = spl[0]
				ga.Version = spl[1]
				//} else if strings.Contains(ga.Name, "-") {
				//	spl := strings.Split(ga.Name, "-")
				//	ga.Name = spl[0]
				//	ga.Version = spl[1]
			}

			if ga.Version == "" {
				ga.Version = "latest"
			}
		}


		ga.Project, err = fl.GetString(argProject)
		if err != nil {
			ga.Project = defaults.DefaultPathNone
		} else {
			ga.Project = DeterminePath(ga.Project)
		}


		ga.Mount, err = fl.GetString(argMount)
		if err != nil {
			ga.Mount = defaults.DefaultPathNone
		} else {
			ga.Mount = DeterminePath(ga.Mount)
		}


		ga.Quiet, err = fl.GetBool(argQuiet)
		if err != nil {
			ga.Quiet = false
		}


		ga.SshStatus, err = fl.GetBool(argStatus)
		if err != nil {
			ga.SshStatus = false
		}


		ga.Temporary, err = fl.GetBool(argTemporary)
		if err != nil {
			ga.Temporary = false
		}

		ga.NoCreate, err = fl.GetBool(argNoCreate)
		if err != nil {
			ga.NoCreate = false
		}


		ga.Provider.Debug = ga.Debug
		ga.Provider.Name, err = fl.GetString(argProvider)
		if err != nil {
			ga.Provider.Name = defaults.DefaultProvider
		}

		ga.Provider.Host, err = fl.GetString(argHost)
		if err != nil {
			ga.Provider.Host = ""
		}

		ga.Provider.Port, err = fl.GetString(argPort)
		if err != nil {
			ga.Provider.Port = ""
		}

		ga.Provider.Project, err = fl.GetString(argProject)
		if err != nil {
			ga.Provider.Project = ""
		}


		ga.State = ga.Provider.NewProvider(ga.Debug)
		if ga.State.IsError() {
			break
		}

		ga.State = ga.GearRef.NewGear(ga.Debug)
		//ga.GearRef, state = ga.Provider.NewGear()
		if ga.State.IsError() {
			break
		}


		//if len(args) == 0 {
		//	state.SetWarning("no args")
		//	break
		//}

		ga.Valid = true
	}

	return ga.State
}


func (ga *GearArgs) CreateLinks(c defaults.ExecCommand, version string) *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
		ga.State = ga.GearRef.GearConfig.CreateLinks(c, version)
	}

	return ga.State
}


func DeterminePath(mp string) string {
	var ok bool

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

		if !ok {
			mp = defaults.DefaultPathNone
		}
	}

	return mp
}

func IsNoCreate(cmd *cobra.Command) bool {
	var ok bool

	for range OnlyOnce {
		var err error

		ok, err = cmd.Flags().GetBool(argNoCreate)
		if err != nil {
			ok = false
			break
		}
	}

	return ok
}

func GetGearboxDir() string {
	var d string

	for range OnlyOnce {
		u, _ := user.Current()
		d = filepath.Join(u.HomeDir, ".gearbox")
	}

	return d
}

func GetLaunchConfig() string {
	var d string

	for range OnlyOnce {
		u, _ := user.Current()
		d = filepath.Join(u.HomeDir, ".gearbox", "launch.json")
	}

	return d
}
