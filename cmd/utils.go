package cmd

import (
	"fmt"
	"github.com/newclarity/scribeHelpers/helperGear"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
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

type LaunchArgs struct {
	Name      string
	Version   string

	Args      []string
	Project   string
	Mount     string
	TmpDir    string

	Temporary bool
	SshStatus bool
	Quiet     bool
	Debug     bool
	NoCreate  bool

	Provider  *helperGear.Provider
	GearRef   *helperGear.Gear

	Valid     bool
	State     *ux.State
}


func (ga *LaunchArgs) IsNil() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}
	ga.State = ga.State.EnsureNotNil()
	return ga.State
}

func (ga *LaunchArgs) IsValid() *ux.State {
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


func (ga *LaunchArgs) ProcessArgs(cmd *cobra.Command, args []string) *ux.State {
	for range OnlyOnce {
		ga.State = ux.NewState(Cmd.Runtime.CmdName, Cmd.Debug)

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

		if ga.Project != defaults.DefaultPathNone {
			ga.Project = DeterminePath(Cmd.Project)
		}

		if ga.Mount != defaults.DefaultPathNone {
			ga.Mount = DeterminePath(Cmd.Mount)
		}

		if ga.TmpDir != defaults.DefaultPathNone {
			ga.TmpDir = DeterminePath(Cmd.TmpDir)
		}

		ga.Provider = helperGear.NewProvider(Cmd.Runtime)
		ga.State = ga.Provider.State
		if ga.State.IsError() {
			break
		}
		ga.State = ga.Provider.SetProvider(Cmd.Provider)

		ga.GearRef = helperGear.NewGear(Cmd.Runtime)
		ga.State = ga.GearRef.State
		if ga.State.IsError() {
			break
		}

		ga.Valid = true
	}

	return ga.State
}


func (ga *LaunchArgs) CreateLinks(version string) *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
		ga.State = ga.GearRef.CreateLinks(version)
	}

	return ga.State
}


func (ga *LaunchArgs) RemoveLinks(version string) *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
		ga.State = ga.GearRef.RemoveLinks(version)
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

//func IsNoCreate(cmd *cobra.Command) bool {
//	var ok bool
//
//	for range OnlyOnce {
//		var err error
//
//		ok, err = cmd.Flags().GetBool(flagNoCreate)
//		if err != nil {
//			ok = false
//			break
//		}
//	}
//
//	return ok
//}

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
