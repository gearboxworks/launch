package cmd

import (
	"fmt"
	"github.com/gearboxworks/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
	"strings"
)


// ******************************************************************************** //
var gbRunCmd = &cobra.Command{
	Use:					fmt.Sprintf("run <%s name> [%s args]", defaults.LanguageContainerName, defaults.LanguageContainerName),
	Aliases:				[]string{},
	Short:					ux.SprintfBlue("Run default %s %s command", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("Run default %s %s command.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch run golang build"),
	DisableFlagParsing:		true,
	DisableFlagsInUseLine:	true,
	Run:					gbRunFunc,
	Args:					cobra.MinimumNArgs(1),
}

//goland:noinspection GoUnusedParameter
func gbRunFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		ga.Quiet = true
		// Windows...
		//br, _ := rootCmd.Flags().GetBool(flagQuiet)
		//fmt.Printf("DEBUG: flagQuiet: %v\n", br)
		//fmt.Printf("DEBUG: DisableFlagParsing: %v\n", rootCmd.DisableFlagParsing)

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

		ga.State = ga.Gears.SelectedSsh(false, ga.SshStatus, ga.Mount, ga.Args)
		if !ga.State.IsError() {
			break
		}

		if ga.Temporary {
			ga.State = ga.gbUninstallFunc()
		}
	}

	return ga.State
}


// ******************************************************************************** //
var gbShellCmd = &cobra.Command{
	Use:					fmt.Sprintf("shell <%s name> [command] [args]", defaults.LanguageContainerName),
	Short:					ux.SprintfBlue("Execute shell in %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("Execute shell in %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch shell mysql ps -eaf"),
	DisableFlagParsing:		true,
	DisableFlagsInUseLine:	true,
	Run:					gbShellFunc,
	Args:					cobra.MinimumNArgs(1),
}

//goland:noinspection GoUnusedParameter
func gbShellFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			//if Cmd.State.IsNotOk() {
			//	Cmd.State.PrintResponse()
			//}
			break
		}

		Cmd.State = ga.gbShellFunc()
		if Cmd.State.IsError() {
			//if Cmd.State.IsNotOk() {
			//	Cmd.State.PrintResponse()
			//}
			break
		}
	}
}

func (ga *LaunchArgs) gbShellFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		if len(ga.Args) > 0 {
			ga.Quiet = true
		}

		ga.State = ga.gbStartFunc()
		if !ga.State.IsRunning() {
			ga.State.SetError("Cannot shell out to %s '%s:%s.'", defaults.LanguageContainerName, ga.Name, ga.Version)
			break
		}

		ga.State = ga.Gears.SelectedSsh(true, ga.SshStatus, ga.Mount, ga.Args[1:])

		if ga.Temporary {
			ga.State = ga.gbUninstallFunc()
		}
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}
