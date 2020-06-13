package cmd

import (
	"github.com/spf13/cobra"
	"github.com/newclarity/scribeHelpers/ux"
)


func gbSaveFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		//showArgs(rootCmd, args)
		ux.PrintfWarning("Command not yet implemented.\n")
	}
}


func gbLoadFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		//showArgs(rootCmd, args)
		ux.PrintfWarning("Command not yet implemented.\n")
	}
}
