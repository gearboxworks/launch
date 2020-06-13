package cmd

import (
	"github.com/spf13/cobra"
	"launch/defaults"
	"github.com/newclarity/scribeHelpers/ux"
)


func gbListFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		Cmd.State = ga.gbListFunc()
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}
	}
}


func (ga *LaunchArgs) gbListFunc() *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range onlyOnce {
		_, ga.State = ga.GearRef.Docker.ImageList(ga.Name)
		if ga.State.IsError() {
			break
		}

		_, ga.State = ga.GearRef.Docker.ContainerList(ga.Name)
		if ga.State.IsError() {
			break
		}

		ga.State = ga.GearRef.Docker.NetworkList(defaults.GearboxNetwork)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}
