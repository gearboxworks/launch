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

func (me *GearArgs) ProcessArgs(cmd *cobra.Command, args []string) ux.State {
	var state ux.State

	for range only.Once {
		if me.Valid {
			break
		}

		var err error
		fl := cmd.Flags()

		//showArgs(cmd, args)

		if len(args) > 1 {
			me.Name = args[0]
			if strings.Contains(me.Name, ":") {
				spl := strings.Split(me.Name, ":")
				me.Name = spl[0]
				me.Version = spl[1]
				//} else if strings.Contains(me.Name, "-") {
				//	spl := strings.Split(me.Name, "-")
				//	me.Name = spl[0]
				//	me.Version = spl[1]
			}

			if me.Version == "" {
				me.Version = "latest"
			}
		}


		me.Project, err = fl.GetString(argProject)
		if err != nil {
			me.Project = defaults.DefaultPathNone
		} else {
			me.Project = DeterminePath(me.Project)
		}


		me.Mount, err = fl.GetString(argMount)
		if err != nil {
			me.Mount = defaults.DefaultPathNone
		} else {
			me.Mount = DeterminePath(me.Mount)
		}


		me.Debug, err = fl.GetBool(argDebug)
		if err != nil {
			me.Debug = false
		}


		me.Quiet, err = fl.GetBool(argQuiet)
		if err != nil {
			me.Quiet = false
		}


		me.SshStatus, err = fl.GetBool(argStatus)
		if err != nil {
			me.SshStatus = false
		}


		me.Temporary, err = fl.GetBool(argTemporary)
		if err != nil {
			me.Temporary = false
		}

		me.Provider.Debug = me.Debug
		me.Provider.Name, _ = fl.GetString(argProvider)
		if err != nil {
			me.Provider.Name = ""
		}

		me.Provider.Host, _ = fl.GetString(argHost)
		if err != nil {
			me.Provider.Host = ""
		}

		me.Provider.Port, _ = fl.GetString(argPort)
		if err != nil {
			me.Provider.Port = ""
		}

		me.Provider.Project, _ = fl.GetString(argProject)
		if err != nil {
			me.Provider.Project = ""
		}

		state = me.Provider.NewProvider()
		if state.IsError() {
			break
		}

		state = me.GearRef.NewGear()
		//me.GearRef, state = me.Provider.NewGear()
		if state.IsError() {
			break
		}


		//if len(args) == 0 {
		//	state.SetWarning("no args")
		//	break
		//}

		me.Valid = true
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
