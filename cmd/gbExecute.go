package cmd

import (
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
	"strings"
)


func gbRunFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			break
		}

		Cmd.State = ga.gbRunFunc()
		if Cmd.State.IsError() {
			break
		}
	}
}

func gbShellFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		Cmd.State = ga.gbShellFunc()
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}
	}
}


func (ga *LaunchArgs) gbRunFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.Quiet = true

		ga.State = ga.gbStartFunc()
		if !ga.State.IsRunning() {
			ga.State.SetError("%s not started", defaults.LanguageContainerName)
			break
		}

		// Yuck!
		sp := strings.Split(ga.Args[0], ":")
		ga.Args[0] = sp[0]
		ga.Args = ga.Gears.Selected.GetCommand(ga.Args)
		if len(ga.Args) == 0 {
			ga.State.SetError("ERROR: no default command defined in gearbox.json")
			break
		}

		// Only for the "run" instance - default to mounting CWD.
		if ga.Mount == defaults.DefaultPathNone {
			ga.Mount = DeterminePath(".")
		}

		ga.State = ga.Gears.Selected.ContainerSsh(false, ga.SshStatus, ga.Mount, ga.Args)
		if !ga.State.IsError() {
			break
		}

		if ga.Temporary {
			ga.State = ga.gbUninstallFunc()
		}
	}

	return ga.State
}

func (ga *LaunchArgs) gbShellFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.gbStartFunc()
		if !ga.State.IsRunning() {
			ga.State.SetError("Cannot shell out to %s '%s:%s.'", defaults.LanguageContainerName, ga.Name, ga.Version)
			break
		}

		ga.State = ga.Gears.Selected.ContainerSsh(true, ga.SshStatus, ga.Mount, ga.Args[1:])

		if ga.Temporary {
			ga.State = ga.gbUninstallFunc()
		}
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}
