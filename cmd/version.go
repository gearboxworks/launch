package cmd

import (
	"github.com/newclarity/scribeHelpers/loadTools"
	"github.com/newclarity/scribeHelpers/toolRuntime"
	"github.com/newclarity/scribeHelpers/toolSelfUpdate"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
)


func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(selfUpdateCmd)
}
var versionCmd = &cobra.Command{
	Use:   loadTools.CmdVersion,
	Short: ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Show version of executable."),
	Long:  ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Show version of executable."),
	Run: func(cmd *cobra.Command, args []string) {
		Cmd.State = ProcessArgs(CmdScribe, cmd, args)
		if Cmd.State.IsNotOk() {
			return
		}

		Cmd.State = Version(args...)
	},
}
var selfUpdateCmd = &cobra.Command{
	Use:   loadTools.CmdSelfUpdate,
	Short: ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Update version of executable."),
	Long: ux.SprintfMagenta(defaults.BinaryName) + ux.SprintfBlue(" - Check and update the latest version."),
	Run: func(cmd *cobra.Command, args []string) {
		Cmd.State = ProcessArgs(CmdScribe, cmd, args)
		if Cmd.State.IsNotOk() {
			return
		}

		Cmd.State = VersionUpdate()
	},
}


func Version(args ...string) *ux.State {
	state := Cmd.State

	for range onlyOnce {
		state = VersionShow()

		switch {
			case len(args) == 0:
				state = VersionHelp()

			case args[0] == "info":
				state = VersionInfo(Cmd.Runtime.GetSemVer())

			case args[0] == "latest":
				state = VersionInfo(nil)

			case args[0] == "check":
				state = VersionCheck()

			case args[0] == "update":
				state = VersionUpdate()

			default:
				state = VersionHelp()
		}
	}

	return state
}


func VersionHelp() *ux.State {
	ux.PrintflnYellow("Need to supply one of:")
	ux.PrintflnYellow("\t'info' - Show info on current version.")
	ux.PrintflnYellow("\t'latest' - Show info on latest version.")
	ux.PrintflnYellow("\t'check' - Check for version updates.")
	ux.PrintflnYellow("\t'update' - Update current executable.")
	Cmd.State.Clear()
	return Cmd.State
}


func VersionShow() *ux.State {
	Cmd.Runtime.PrintNameVersion()
	Cmd.State.Clear()
	return Cmd.State
}


func VersionInfo(v *toolRuntime.VersionValue) *ux.State {
	state := Cmd.State

	for range onlyOnce {
		update := toolSelfUpdate.New(Cmd.Runtime)
		if update.State.IsError() {
			state = update.State
			break
		}

		state = update.PrintVersion((*toolSelfUpdate.VersionValue)(v))
		if state.IsNotOk() {
			state = update.State
			break
		}
	}

	return state
}


func VersionCheck() *ux.State {
	state := Cmd.State

	for range onlyOnce {
		update := toolSelfUpdate.New(Cmd.Runtime)
		if update.State.IsError() {
			state = update.State
			break
		}

		state = update.IsUpdated(true)
		if update.State.IsError() {
			break
		}
	}

	return state
}


func VersionUpdate() *ux.State {
	state := Cmd.State

	for range onlyOnce {
		update := toolSelfUpdate.New(Cmd.Runtime)
		if update.State.IsError() {
			state = update.State
			break
		}

		state = update.IsUpdated(true)
		if update.State.IsError() {
			break
		}

		state = update.Update()
		if state.IsNotOk() {
			break
		}
	}

	return state
}
