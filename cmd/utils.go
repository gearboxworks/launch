package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"launch/defaults"
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
		if !debugFlag {
			break
		}

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

type gearArgs struct {
	Name string
	Version string
	Project string
	Mount string
	Temporary bool
	SshStatus bool
}

func getGearArgs(cmd *cobra.Command, args []string) (*gearArgs, ux.State) {
	var ga gearArgs
	var state ux.State

	for range only.Once {
		var err error

		showArgs(cmd, args)

		if len(args) == 0 {
			state.SetError("no args")
			break
		}

		ga.Name = args[0]
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


		ga.Project, err = cmd.Flags().GetString(argProject)
		if err != nil {
			ga.Project = defaults.DefaultPathNone
			break
		}
		ga.Project = DeterminePath(ga.Project)


		ga.Mount, err = cmd.Flags().GetString(argMount)
		if err != nil {
			ga.Mount = defaults.DefaultPathNone
			break
		}
		ga.Mount = DeterminePath(ga.Mount)


		ga.SshStatus, err = cmd.Flags().GetBool(argMount)
		if err != nil {
			ga.SshStatus = false
			break
		}


		ga.Temporary, err = cmd.Flags().GetBool(argMount)
		if err != nil {
			ga.Temporary = false
			break
		}
	}

	return &ga, state
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
