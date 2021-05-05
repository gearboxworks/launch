package cmd

import (
	"fmt"
	"github.com/gearboxworks/scribeHelpers/toolGear"
	"github.com/gearboxworks/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
	"strings"
)


// ******************************************************************************** //
//var gbManageCmd = &cobra.Command {
//	Use:					"manage",
//	//Aliases:				[]string{"show"},
//	Short:					ux.Sprintf("Manage %s %s", defaults.LanguageAppName, defaults.LanguageContainerPluralName),
//	Long:					ux.SprintfBlue("Manage %s %s.", defaults.LanguageAppName, defaults.LanguageContainerPluralName),
//	Example:				ux.SprintfWhite("launch manage"),
//	DisableFlagParsing:		false,
//	DisableFlagsInUseLine:	false,
//	Run:					gbManageFunc,
//	Args:					cobra.RangeArgs(0, 2),
//}
//
//func gbManageFunc(cmd *cobra.Command, args []string) {
//	for range onlyOnce {
//		var ga LaunchArgs
//
//		Cmd.State = ga.ProcessArgs(rootCmd, args)
//		if Cmd.State.IsNotOk() {
//			break
//		}
//		Cmd.SetDebug(ga.Debug)
//		//switch {
//		//	case len(args) == 0:
//				_ = cmd.Help()
//		//}
//	}
//}


// ******************************************************************************** //
var gbSearchCmd = &cobra.Command {
	Use:					"search",
	Aliases:				[]string{"available", "find", "avail"},
	Short:					ux.SprintfBlue("Search for available %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("Search for available %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch manage search"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbSearchFunc,
	Args:					cobra.RangeArgs(0, 1),
}

//goland:noinspection GoUnusedParameter
func gbSearchFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}
		//Cmd.SetDebug(ga.Debug)

		Cmd.State = ga.gbSearchFunc(args)
		if Cmd.State.IsNotOk() {
			break
		}
	}
}
func (ga *LaunchArgs) gbSearchFunc(args []string) *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		//term := toolGear.DefaultOrganization
		term := ""
		if len(args) > 0 {
			term = args[0]
		}

		ga.State = ga.Gears.SearchPrint(term, "")
		if ga.State.IsNotOk() {
			break
		}
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}


// ******************************************************************************** //
var gbInstallCmd = &cobra.Command {
	Use:					fmt.Sprintf("install <%s name>", defaults.LanguageContainerName),
	SuggestFor:				[]string{"download" ,"add"},
	Short:					ux.SprintfBlue("Install a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("Install a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch install golang"),
	DisableFlagParsing:		false,
	Run:					gbInstallFunc,
	Args:					cobra.ExactArgs(1),
}

//goland:noinspection GoUnusedParameter
func gbInstallFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}
		//Cmd.SetDebug(ga.Debug)

		Cmd.State = ga.gbInstallFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}
func (ga *LaunchArgs) gbInstallFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var found bool

		ga.State = ga.Gears.FindImage(ga.Name, ga.Version)
		found = ga.Gears.State.GetResponseAsBool()

		if !found {
			ga.State = ga.Gears.Search(ga.Name, ga.Version)
			if !ga.State.GetResponseAsBool() {
				ga.State.PrintResponse()
				ux.PrintflnBlue("")
				ga.gbSearchFunc([]string{})
				ga.State.SetWarning("%s not found in registry.", ga.Gears.Language.ImageName)
				break
			}
			if ga.State.IsNotOk() {
				break
			}
		}

		found, ga.State = ga.Gears.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}

		if found {
			if !ga.Temporary {
				ga.State = ga.Gears.FindImage(ga.Name, ga.Version)
				found = ga.Gears.State.GetResponseAsBool()
				if found {
					// Create symlinks.
					ga.CreateLinks(ga.Version)
				}
			}

			ga.State.SetOk("%s '%s:%s' already installed.", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State = ga.Gears.NetworkCreate(defaults.GearboxNetwork)
		if ga.State.IsError() {
			break
		}


		if ga.Project != toolGear.DefaultPathNone {
			ga.Gears.SelectedAddVolume(ga.Project, toolGear.DefaultProject)
		}

		if ga.TmpDir != toolGear.DefaultPathNone {
			ga.Gears.SelectedAddVolume(ga.TmpDir, toolGear.DefaultTmpDir)
		}


		if !ga.Quiet {
			ux.PrintflnNormal("Installing %s '%s:%s'.", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.Gears.ContainerCreate(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}

		if ga.State.IsCreated() {
			if ga.Temporary {
				ga.State.Clear()
				break
			}

			// Create symlinks.
			ga.State = ga.CreateLinks(ga.Version)

			ga.State.SetOk("Installed %s '%s:%s' OK.", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetWarning("%s '%s:%s' cannot be installed", defaults.LanguageContainerName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}


// ******************************************************************************** //
var gbUninstallCmd = &cobra.Command {
	Use:					fmt.Sprintf("uninstall <%s name>", defaults.LanguageContainerName),
	SuggestFor:				[]string{"remove"},
	Short:					ux.SprintfBlue("Uninstall a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("Uninstall a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch uninstall golang"),
	DisableFlagParsing:		false,
	Run:					gbUninstallFunc,
	Args:					cobra.ExactArgs(1),
}

//goland:noinspection GoUnusedParameter
func gbUninstallFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}
		//Cmd.SetDebug(ga.Debug)

		Cmd.State = ga.gbUninstallFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}
func (ga *LaunchArgs) gbUninstallFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var found bool
		found, ga.State = ga.Gears.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		if !found {
			if !ga.Temporary {
				ga.State = ga.Gears.FindImage(ga.Name, ga.Version)
				found = ga.Gears.State.GetResponseAsBool()
				if found {
					// Remove symlinks.
					ga.RemoveLinks(ga.Version)
				}
			}

			ga.State.SetOk("%s '%s:%s' already removed.", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}
		ga.State.Clear()

		ga.State = ga.gbStopFunc()
		if ga.State.IsError() {
			break
		}

		if !ga.Quiet {
			ux.PrintflnNormal("Removing %s '%s:%s'.\n", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.Gears.SelectedRemove()
		if ga.State.IsError() {
			ga.State.SetError("%s '%s:%s' remove error - %s", defaults.LanguageContainerName, ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsOk() {
			if ga.Temporary {
				ga.State.Clear()
				break
			}

			// Remove symlinks.
			ga.RemoveLinks(ga.Version)

			ga.State.SetOk("%s '%s:%s' removed OK", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetWarning("%s '%s:%s' cannot be removed", defaults.LanguageContainerName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}


// ******************************************************************************** //
var gbReinstallCmd = &cobra.Command {
	Use:					fmt.Sprintf("update <%s name>", defaults.LanguageContainerName),
	SuggestFor:				[]string{"reinstall"},
	Short:					ux.SprintfBlue("Update a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("Update a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch update golang"),
	DisableFlagParsing:		false,
	Run:					gbReinstallFunc,
	Args:					cobra.ExactArgs(1),
}

//goland:noinspection GoUnusedParameter
func gbReinstallFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}
		//Cmd.SetDebug(ga.Debug)

		Cmd.State = ga.gbReinstallFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}
func (ga *LaunchArgs) gbReinstallFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.gbUninstallFunc()
		if ga.State.IsError() {
			break
		}

		ga.State = ga.gbInstallFunc()
		if ga.State.IsError() {
			break
		}

		ga.State.SetOk("%s '%s:%s' reinstalled.", defaults.LanguageContainerName, ga.Name, ga.Version)
		ga.State.SetOutput("")
	}

	return ga.State
}


// ******************************************************************************** //
var gbCleanCmd = &cobra.Command{
	Use:					fmt.Sprintf("clean <%s name>", defaults.LanguageContainerName),
	SuggestFor:				[]string{},
	Short:					ux.SprintfBlue("Completely uninstall a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("Completely uninstall a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch clean golang"),
	DisableFlagParsing:		false,
	Run:					gbCleanFunc,
	Args:					cobra.ExactArgs(1),
}

//goland:noinspection GoUnusedParameter
func gbCleanFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}
		//Cmd.SetDebug(ga.Debug)

		Cmd.State = ga.gbCleanFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}
func (ga *LaunchArgs) gbCleanFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.gbUninstallFunc()
		if ga.State.IsError() {
			break
		}

		var found bool
		ga.State = ga.Gears.FindImage(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		found = ga.Gears.State.GetResponseAsBool()
		if !found {
			ga.State.SetOk("%s '%s:%s' already removed.", defaults.LanguageImageName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}
		ga.State.Clear()

		if !ga.Quiet {
			ux.PrintflnNormal("Removing %s '%s:%s': ", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.Gears.SelectedImageRemove()
		if ga.State.IsError() {
			ga.State.SetError("%s '%s:%s' remove error - %s", defaults.LanguageImageName, ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsOk() {
			if ga.Temporary {
				ga.State.Clear()
				break
			}

			ga.State.SetOk("%s '%s:%s' removed OK", defaults.LanguageImageName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetWarning("%s '%s:%s' cannot be removed", defaults.LanguageImageName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}


// ******************************************************************************** //
var gbStartCmd = &cobra.Command{
	Use:					fmt.Sprintf("start <%s name>", defaults.LanguageContainerName),
	Short:					ux.SprintfBlue("Start a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("Start a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch start golang"),
	DisableFlagParsing:		false,
	Run:					gbStartFunc,
	Args:					cobra.ExactArgs(1),
}

//goland:noinspection GoUnusedParameter
func gbStartFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}
		//Cmd.SetDebug(ga.Debug)

		Cmd.State = ga.gbStartFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}
func (ga *LaunchArgs) gbStartFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var found bool
		found, ga.State = ga.Gears.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			// Pass-through warnings.
			break
		}
		if !found {
			if !ga.Temporary {
				ga.Quiet = false
			}

			if ga.NoCreate {
				ga.State.SetError("Not creating %s '%s:%s'.", defaults.LanguageContainerName, ga.Name, ga.Version)
				break
			}

			ga.gbInstallFunc()
			if ga.State.IsError() {
				ga.State.SetError("Cannot start %s '%s:%s'.", defaults.LanguageContainerName, ga.Name, ga.Version)
				break
			}
			if ga.State.IsNotOk() {
				break
			}

			// @TODO - Need a better way to handle the "Docker client error: context deadline exceeded" errors.
		}

		if ga.State.IsRunning() {
			ga.State.SetOk("%s '%s:%s' already started.", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}


		if !ga.Quiet {
			ux.PrintflnNormal("Starting %s '%s:%s': ", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.Gears.SelectedStart()
		if ga.State.IsError() {
			if strings.Contains(ga.State.GetError().Error(), "address already in use") {
				//saved := gear.State.GetError()
				ga.Gears.Selected.ListImagePorts()
				ga.State.SetError("Error: There are ports already used.")
				break
			}
			ga.State.SetError("%s '%s:%s' start error - %s", defaults.LanguageContainerName, ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsRunning() {
			ga.State.SetOk("%s '%s:%s' started OK", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetError("%s '%s:%s' cannot be started", defaults.LanguageContainerName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}


// ******************************************************************************** //
var gbStopCmd = &cobra.Command{
	Use:					fmt.Sprintf("stop <%s name>", defaults.LanguageContainerName),
	Short:					ux.SprintfBlue("Stop a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("Stop a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch stop golang"),
	DisableFlagParsing:		false,
	Run:					gbStopFunc,
	Args:					cobra.ExactArgs(1),
}

//goland:noinspection GoUnusedParameter
func gbStopFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}
		//Cmd.SetDebug(ga.Debug)

		Cmd.State = ga.gbStopFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}
func (ga *LaunchArgs) gbStopFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var found bool
		found, ga.State = ga.Gears.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		if !found {
			break
		}
		if ga.State.IsExited() {
			ga.State.SetOk("%s '%s:%s' already stopped.", defaults.LanguageImageName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}


		if !ga.Quiet {
			ux.PrintflnNormal("Stopping %s '%s:%s': ", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.Gears.SelectedStop()
		if ga.State.IsError() {
			ga.State.SetError("%s '%s:%s' stop error - %s", defaults.LanguageContainerName, ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsExited() {
			ga.State.SetOk("%s '%s:%s' stopped OK", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		if ga.State.IsCreated() {
			ga.State.SetOk("%s '%s:%s' stopped OK", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetWarning("%s '%s:%s' cannot be stopped", defaults.LanguageContainerName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}


// ******************************************************************************** //
var gbLogsCmd = &cobra.Command{
	Use:					fmt.Sprintf("log <%s name>", defaults.LanguageContainerName),
	SuggestFor:				[]string{"logs"},
	Short:					ux.SprintfBlue("Show logs of %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfBlue("Show logs of %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch log golang"),
	DisableFlagParsing:		false,
	Run:					gbLogsFunc,
	Args:					cobra.ExactArgs(1),
}

//goland:noinspection GoUnusedParameter
func gbLogsFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}
		//Cmd.SetDebug(ga.Debug)

		if len(args) == 0 {
			Cmd.State.SetError("Need to specify a %s.", defaults.LanguageContainerName)
			break
		}

		Cmd.State = ga.gbLogsFunc()
	}
}
func (ga *LaunchArgs) gbLogsFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var found bool
		found, ga.State = ga.Gears.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		if !found {
			break
		}

		ga.State = ga.Gears.SelectedLogs()
		if ga.State.IsError() {
			ga.State.SetError("%s '%s:%s' stop error - %s", defaults.LanguageContainerName, ga.Name, ga.Version, ga.State.GetError())
			break
		}
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}
