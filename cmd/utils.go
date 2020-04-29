package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"launch/only"
	"launch/ux"
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
	Mount string
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

		//ga.Version, err = cmd.Flags().GetString(argVersion)
		//if err != nil {
		//	ga.Version = "latest"
		//}

		ga.Mount, err = cmd.Flags().GetString(argProject)
		if err != nil {
			ga.Mount = ""
			break
		}

		ga.Mount, err = filepath.Abs(ga.Mount)
		if err != nil {
			ga.Mount = ""
			break
		}

		//gearRef, state = provider.NewGear()
		//if state.IsError() {
		//	break
		//}
	}

	return &ga, state
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
