package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"launch/only"
	"launch/ux"
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
			state.SetError("no ")
			break
		}
		ga.Name = args[0]
		ga.Version, err = cmd.Flags().GetString(argVersion)
		if err != nil {
			ga.Version = ""
		}
		ga.Mount, err = cmd.Flags().GetString(argProject)
		if err != nil {
			ga.Mount = ""
		}

		//gearRef, state = provider.NewGear()
		//if state.IsError() {
		//	break
		//}
	}

	return &ga, state
}