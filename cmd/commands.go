package cmd

import (
	"fmt"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
)


var gbInstallCmd = &cobra.Command{
	Use:					fmt.Sprintf("install <%s name>", defaults.LanguageContainerName),
	SuggestFor:				[]string{"download" ,"add"},
	Short:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Install a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Install a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch install golang"),
	DisableFlagParsing:		false,
	Run:					gbInstallFunc,
	Args:					cobra.ExactArgs(1),
}
var gbUninstallCmd = &cobra.Command{
	Use:					fmt.Sprintf("uninstall <%s name>", defaults.LanguageContainerName),
	SuggestFor:				[]string{"remove"},
	Short:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Uninstall a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Uninstall a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch uninstall golang"),
	DisableFlagParsing:		false,
	Run:					gbUninstallFunc,
	Args:					cobra.ExactArgs(1),
}
var gbReinstallCmd = &cobra.Command{
	Use:					fmt.Sprintf("update <%s name>", defaults.LanguageContainerName),
	SuggestFor:				[]string{"reinstall"},
	Short:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Update a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Update a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch update golang"),
	DisableFlagParsing:		false,
	Run:					gbReinstallFunc,
	Args:					cobra.ExactArgs(1),
}
var gbCleanCmd = &cobra.Command{
	Use:					fmt.Sprintf("clean <%s name>", defaults.LanguageContainerName),
	SuggestFor:				[]string{},
	Short:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Completely uninstall a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Completely uninstall a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch clean golang"),
	DisableFlagParsing:		false,
	Run:					gbCleanFunc,
	Args:					cobra.ExactArgs(1),
}


var gbListCmd = &cobra.Command{
	Use:					"list",
	Aliases:				[]string{"show"},
	Short:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - List a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - List a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch list all"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbListFunc,
	Args:					cobra.RangeArgs(0, 1),
}
var gbPortsCmd = &cobra.Command{
	Use:					fmt.Sprintf("ports [%s name]", defaults.LanguageContainerName),
	//Aliases:				[]string{"links"},
	Short:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - List ports provided by a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - List ports provided by a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch list ports golang"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbPortsFunc,
	Args:					cobra.RangeArgs(0, 1),
}
var gbLinksCmd = &cobra.Command{
	Use:					fmt.Sprintf("files [%s name]", defaults.LanguageContainerName),
	Aliases:				[]string{"links", "ls"},
	Short:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - List files provided by a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - List files provided by a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch list files golang"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbLinksFunc,
	Args:					cobra.RangeArgs(0, 1),
}
var gbDetailsCmd = &cobra.Command{
	Use:					"all",
	Aliases:				[]string{"details"},
	Short:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - List all details provided by a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - List all details provided by a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch list all golang"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbDetailsFunc,
	Args:					cobra.RangeArgs(0, 1),
}


var gbRunCmd = &cobra.Command{
	Use:					fmt.Sprintf("run <%s name> [%s args]", defaults.LanguageContainerName, defaults.LanguageContainerName),
	Aliases:				[]string{},
	Short:					ux.SprintfMagenta("Execute") + ux.SprintfBlue(" - Run default %s %s command", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Execute") + ux.SprintfBlue(" - Run default %s %s command.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch run golang build"),
	DisableFlagParsing:		true,
	DisableFlagsInUseLine:	true,
	Run:					gbRunFunc,
	Args:					cobra.MinimumNArgs(1),
}
var gbShellCmd = &cobra.Command{
	Use:					fmt.Sprintf("shell <%s name> [command] [args]", defaults.LanguageContainerName),
	Short:					ux.SprintfMagenta("Execute") + ux.SprintfBlue(" - Execute shell in %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Execute") + ux.SprintfBlue(" - Execute shell in %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch shell mysql ps -eaf"),
	DisableFlagParsing:		true,
	DisableFlagsInUseLine:	true,
	Run:					gbShellFunc,
	Args:					cobra.MinimumNArgs(1),
}
var gbUnitTestCmd = &cobra.Command{
	Use:					fmt.Sprintf("test <%s name>", defaults.LanguageContainerName),
	Short:					ux.SprintfMagenta("Execute") + ux.SprintfBlue(" - Execute unit tests in %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Execute") + ux.SprintfBlue(" - Execute unit tests in %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch unit tests terminus"),
	DisableFlagParsing:		true,
	DisableFlagsInUseLine:	true,
	Run:					gbUnitTestFunc,
	Args:					cobra.MinimumNArgs(1),
}


var gbStartCmd = &cobra.Command{
	Use:					fmt.Sprintf("start <%s name>", defaults.LanguageContainerName),
	Short:					ux.SprintfMagenta("Run") + ux.SprintfBlue(" - Start a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Run") + ux.SprintfBlue(" - Start a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch start golang"),
	DisableFlagParsing:		false,
	Run:					gbStartFunc,
	Args:					cobra.ExactArgs(1),
}
var gbStopCmd = &cobra.Command{
	Use:					fmt.Sprintf("stop <%s name>", defaults.LanguageContainerName),
	Short:					ux.SprintfMagenta("Run") + ux.SprintfBlue(" - Stop a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Run") + ux.SprintfBlue(" - Stop a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch stop golang"),
	DisableFlagParsing:		false,
	Run:					gbStopFunc,
	Args:					cobra.ExactArgs(1),
}


var gbBuildCmd = &cobra.Command {
	Use:					fmt.Sprintf("build <%s name>", defaults.LanguageContainerName),
	SuggestFor:				[]string{ "compile", "generate" },
	Short:					ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Build a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Allows building of arbitrary containers as a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch build golang"),
	DisableFlagParsing:		false,
	Run:					gbBuildFunc,
	Args:					cobra.ExactArgs(1),
}
var gbPublishCmd = &cobra.Command{
	Use:					fmt.Sprintf("publish <%s name>", defaults.LanguageContainerName),
	SuggestFor:				[]string{"upload"},
	Short:					ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Publish a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Publish a %s %s to GitHub or DockerHub.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch publish golang"),
	DisableFlagParsing:		false,
	Run:					gbPublishFunc,
	Args:					cobra.ExactArgs(1),
}


var gbSaveCmd = &cobra.Command{
	Use:					fmt.Sprintf("export <%s name>", defaults.LanguageContainerName),
	SuggestFor:				[]string{"save"},
	Short:					ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Save state of a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Save state of a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch save golang"),
	DisableFlagParsing:		false,
	Run:					gbSaveFunc,
	Args:					cobra.ExactArgs(1),
}
var gbLoadCmd = &cobra.Command{
	Use:					fmt.Sprintf("import <%s name>", defaults.LanguageContainerName),
	SuggestFor:				[]string{"load"},
	Short:					ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Load a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:					ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Load a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:				ux.SprintfWhite("launch load golang"),
	DisableFlagParsing:		false,
	Run:					gbLoadFunc,
	Args:					cobra.ExactArgs(1),
}
