package cmd

import (
	"fmt"
	"github.com/gearboxworks/scribeHelpers/toolGear"
	"github.com/gearboxworks/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
)


// ******************************************************************************** //
var gbListCmd = &cobra.Command{
	Use:					"list",
	Aliases:				[]string{"show"},
	Short:					ux.SprintfBlue("List a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("List a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch manage list all"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbListFunc,
	Args:					cobra.RangeArgs(0, 1),
}
//gbListFunc takes a pointer to cobra.command and command's argument and
//output appropriate list.
func gbListFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, false)
		if Cmd.State.IsNotOk() {
			break
		}
		//Cmd.SetDebug(ga.Debug)

		switch {
			case len(args) == 0:
				_ = cmd.Help()
		}
	}
}

func (ga *LaunchArgs) gbListFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.Gears.ListImages(ga.Name)
		if ga.State.IsError() {
			break
		}

		ga.State = ga.Gears.ListContainers(ga.Name)
		if ga.State.IsError() {
			break
		}

		ga.State = ga.Gears.NetworkList(toolGear.DefaultNetwork)
		if ga.State.IsError() {
			break
		}
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}


// ******************************************************************************** //
var gbDetailsCmd = &cobra.Command{
	Use:					"all",
	Aliases:				[]string{"details"},
	Short:					ux.SprintfBlue("List all details provided by a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("List all details provided by a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch manage list all golang"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbDetailsFunc,
	Args:					cobra.RangeArgs(0, 1),
}

//gbDetailsFunc takes a pointer to cobra.command and command's arguments and
//output appropriate details message.
//goland:noinspection GoUnusedParameter
func gbDetailsFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}
		//Cmd.SetOptions(ga)
		//Cmd.SetDebug(ga.Debug)

		Cmd.State = ga.gbListFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}


// ******************************************************************************** //
var gbLinksCmd = &cobra.Command{
	Use:					fmt.Sprintf("files [%s name]", defaults.LanguageContainerName),
	Aliases:				[]string{"links", "ls"},
	Short:					ux.SprintfBlue("List files provided by a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("List files provided by a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch manage list files golang"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbLinksFunc,
	Args:					cobra.RangeArgs(0, 1),
}

//gbLinksFunc takes a pointer to cobra.command and command's arguments and
//output appropriate links.
//goland:noinspection GoUnusedParameter
func gbLinksFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}
		//Cmd.SetDebug(ga.Debug)

		Cmd.State = ga.gbLinksFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}

func (ga *LaunchArgs) gbLinksFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		remote := false
		if Cmd.Host != "" {
			// We are remote.
			remote = true
		}

		ga.State = ga.ListLinks(remote)
		if ga.State.IsError() {
			break
		}
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}


// ******************************************************************************** //
var gbPortsCmd = &cobra.Command{
	Use:					fmt.Sprintf("ports [%s name]", defaults.LanguageContainerName),
	//Aliases:				[]string{"links"},
	Short:					ux.SprintfBlue("List ports provided by a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("List ports provided by a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch manage list ports golang"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbPortsFunc,
	Args:					cobra.RangeArgs(0, 1),
}

//gbPortsFunc takes a pointer to cobra.command and command's arguments and
//output appropriate ports.
//goland:noinspection GoUnusedParameter
func gbPortsFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}
		//Cmd.SetDebug(ga.Debug)

		//ga.Gears.Selected.ScanPorts()

		Cmd.State = ga.gbPortsFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}

func (ga *LaunchArgs) gbPortsFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		remote := false
		if Cmd.Host != "" {
			// We are remote.
			remote = true
		}

		ga.State = ga.ListPorts(remote)
		if ga.State.IsError() {
			break
		}
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}
