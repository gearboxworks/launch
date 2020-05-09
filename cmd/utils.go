package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"launch/defaults"
	"launch/gear"
	"launch/only"
	"launch/ux"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func showArgs(cmd *cobra.Command, args []string) {
	// var err error

	for range only.Once {
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
	Valid bool

	Name string
	Version string
	Project string
	Mount string
	Temporary bool
	SshStatus bool
	Quiet bool
	Debug bool
	Provider gear.Provider
	GearRef gear.Gear
}

func (gear *GearArgs) ProcessArgs(cmd *cobra.Command, args []string) ux.State {
	var state ux.State

	for range only.Once {
		if gear.Valid {
			break
		}

		var err error
		fl := cmd.Flags()

		//showArgs(cmd, args)

		if len(args) > 0 {
			gear.Name = args[0]
			if strings.Contains(gear.Name, ":") {
				spl := strings.Split(gear.Name, ":")
				gear.Name = spl[0]
				gear.Version = spl[1]
				//} else if strings.Contains(gear.Name, "-") {
				//	spl := strings.Split(gear.Name, "-")
				//	gear.Name = spl[0]
				//	gear.Version = spl[1]
			}

			if gear.Version == "" {
				gear.Version = "latest"
			}
		}


		gear.Project, err = fl.GetString(argProject)
		if err != nil {
			gear.Project = defaults.DefaultPathNone
		} else {
			gear.Project = DeterminePath(gear.Project)
		}


		gear.Mount, err = fl.GetString(argMount)
		if err != nil {
			gear.Mount = defaults.DefaultPathNone
		} else {
			gear.Mount = DeterminePath(gear.Mount)
		}


		gear.Debug, err = fl.GetBool(argDebug)
		if err != nil {
			gear.Debug = false
		}


		gear.Quiet, err = fl.GetBool(argQuiet)
		if err != nil {
			gear.Quiet = false
		}


		gear.SshStatus, err = fl.GetBool(argStatus)
		if err != nil {
			gear.SshStatus = false
		}


		gear.Temporary, err = fl.GetBool(argTemporary)
		if err != nil {
			gear.Temporary = false
		}

		gear.Provider.Debug = gear.Debug
		gear.Provider.Name, _ = fl.GetString(argProvider)
		if err != nil {
			gear.Provider.Name = ""
		}

		gear.Provider.Host, _ = fl.GetString(argHost)
		if err != nil {
			gear.Provider.Host = ""
		}

		gear.Provider.Port, _ = fl.GetString(argPort)
		if err != nil {
			gear.Provider.Port = ""
		}

		gear.Provider.Project, _ = fl.GetString(argProject)
		if err != nil {
			gear.Provider.Project = ""
		}

		state = gear.Provider.NewProvider()
		if state.IsError() {
			break
		}

		state = gear.GearRef.NewGear()
		//gear.GearRef, state = gear.Provider.NewGear()
		if state.IsError() {
			break
		}


		//if len(args) == 0 {
		//	state.SetWarning("no args")
		//	break
		//}

		gear.Valid = true
	}

	return state
}

func DeterminePath(mp string) string {
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

		if !ok {
			mp = defaults.DefaultPathNone
		}
	}

	return mp
}

func IsNoCreate(cmd *cobra.Command) bool {
	var ok bool

	for range only.Once {
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

	for range only.Once {
		u, _ := user.Current()
		d = filepath.Join(u.HomeDir, ".gearbox")
	}

	return d
}

func GetLaunchConfig() string {
	var d string

	for range only.Once {
		u, _ := user.Current()
		d = filepath.Join(u.HomeDir, ".gearbox", "launch.json")
	}

	return d
}
