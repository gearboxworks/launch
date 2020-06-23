package cmd

import (
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
)


var gbInstallCmd = &cobra.Command{
	Use:   "install <gear name>",
	//Aliases: []string{"add"},
	SuggestFor: []string{"download" ,"add"},
	Short: ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Install a Gearbox gear"),
	Long: ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Install a Gearbox gear."),
	Example: ux.SprintfWhite("launch install golang"),
	DisableFlagParsing: false,
	Run: gbInstallFunc,
	Args: cobra.ExactArgs(1),
}
var gbUninstallCmd = &cobra.Command{
	Use:   "uninstall <gear name>",
	//Aliases: []string{"remove"},
	SuggestFor: []string{"remove"},
	Short: ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Uninstall a Gearbox gear"),
	Long: ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Uninstall a Gearbox gear."),
	Example: ux.SprintfWhite("launch uninstall golang"),
	DisableFlagParsing: false,
	Run: gbUninstallFunc,
	Args: cobra.ExactArgs(1),
}
var gbReinstallCmd = &cobra.Command{
	Use:   "reinstall <gear name>",
	//Aliases: []string{"update"},
	SuggestFor: []string{"update"},
	Short: ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Update a Gearbox gear"),
	Long: ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Update a Gearbox gear."),
	Example: ux.SprintfWhite("launch reinstall golang"),
	DisableFlagParsing: false,
	Run: gbReinstallFunc,
	Args: cobra.ExactArgs(1),
}
var gbCleanCmd = &cobra.Command{
	Use:   "clean <gear name>",
	//Aliases: []string{"remove"},
	SuggestFor: []string{},
	Short: ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Completely uninstall a Gearbox gear"),
	Long: ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - Completely uninstall a Gearbox gear."),
	Example: ux.SprintfWhite("launch uninstall golang"),
	DisableFlagParsing: false,
	Run: gbCleanFunc,
	Args: cobra.ExactArgs(1),
}


var gbListCmd = &cobra.Command{
	Use:   "list [gear name]",
	Aliases: []string{"ls", "show"},
	Short: ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - List a Gearbox gear"),
	Long: ux.SprintfMagenta("Manage") + ux.SprintfBlue(" - List a Gearbox gear."),
	Example: ux.SprintfWhite("launch list golang"),
	DisableFlagParsing: false,
	DisableFlagsInUseLine: false,
	Run: gbListFunc,
	Args: cobra.RangeArgs(0, 1),
}


var gbRunCmd = &cobra.Command{
	Use:   "run <gear name> [gear args]",
	Aliases: []string{},
	Short: ux.SprintfMagenta("Execute") + ux.SprintfBlue(" - Run default Gearbox gear command"),
	Long: ux.SprintfMagenta("Execute") + ux.SprintfBlue(" - Run default Gearbox gear command."),
	Example: ux.SprintfWhite("launch run golang build"),
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	Run: gbRunFunc,
	Args: cobra.MinimumNArgs(1),
}
var gbShellCmd = &cobra.Command{
	Use:   "shell <gear name> [command] [args]",
	Short: ux.SprintfMagenta("Execute") + ux.SprintfBlue(" - Execute shell in Gearbox gear"),
	Long: ux.SprintfMagenta("Execute") + ux.SprintfBlue(" - Execute shell in Gearbox gear."),
	Example: ux.SprintfWhite("launch shell mysql ps -eaf"),
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	Run: gbShellFunc,
	Args: cobra.MinimumNArgs(1),
}
var gbUnitTestCmd = &cobra.Command{
	Use:   "test <gear name>",
	Short: ux.SprintfMagenta("Execute") + ux.SprintfBlue(" - Execute unit tests in Gearbox gear"),
	Long: ux.SprintfMagenta("Execute") + ux.SprintfBlue(" - Execute unit tests in Gearbox gear."),
	Example: ux.SprintfWhite("launch unit tests terminus"),
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	Run: gbUnitTestFunc,
	Args: cobra.MinimumNArgs(1),
}


var gbStartCmd = &cobra.Command{
	Use:   "start <gear name>",
	Short: ux.SprintfMagenta("Run") + ux.SprintfBlue(" - Start a Gearbox gear"),
	Long: ux.SprintfMagenta("Run") + ux.SprintfBlue(" - Start a Gearbox gear."),
	Example: ux.SprintfWhite("launch start golang"),
	DisableFlagParsing: false,
	Run: gbStartFunc,
	Args: cobra.ExactArgs(1),
}
var gbStopCmd = &cobra.Command{
	Use:   "stop <gear name>",
	Short: ux.SprintfMagenta("Run") + ux.SprintfBlue(" - Stop a Gearbox gear"),
	Long: ux.SprintfMagenta("Run") + ux.SprintfBlue(" - Stop a Gearbox gear."),
	Example: ux.SprintfWhite("launch stop golang"),
	DisableFlagParsing: false,
	Run: gbStopFunc,
	Args: cobra.ExactArgs(1),
}


var gbBuildCmd = &cobra.Command {
	Use:   "build <gear name>",
	SuggestFor: []string{ "compile", "generate" },
	Short: ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Build a Gearbox gear"),
	Long: ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Allows building of arbitrary containers as a Gearbox container, (called gears)."),
	Example: ux.SprintfWhite("launch build golang"),
	DisableFlagParsing: false,
	Run: gbBuildFunc,
	Args: cobra.ExactArgs(1),
}
var gbPublishCmd = &cobra.Command{
	Use:   "publish <gear name>",
	SuggestFor: []string{"upload"},
	Short: ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Publish a Gearbox gear"),
	Long: ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Publish a Gearbox gear to GitHub or DockerHub."),
	Example: ux.SprintfWhite("launch publish golang"),
	DisableFlagParsing: false,
	Run: gbPublishFunc,
	Args: cobra.ExactArgs(1),
}


var gbSaveCmd = &cobra.Command{
	Use:   "export <gear name>",
	//Aliases: []string{"save"},
	SuggestFor: []string{"save"},
	Short: ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Save state of a Gearbox gear"),
	Long: ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Save state of a Gearbox gear."),
	Example: ux.SprintfWhite("launch save golang"),
	DisableFlagParsing: false,
	Run: gbSaveFunc,
	Args: cobra.ExactArgs(1),
}
var gbLoadCmd = &cobra.Command{
	Use:   "import <gear name>",
	//Aliases: []string{"load"},
	SuggestFor: []string{"load"},
	Short: ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Load a Gearbox gear"),
	Long: ux.SprintfMagenta("Create") + ux.SprintfBlue(" - Load a Gearbox gear."),
	Example: ux.SprintfWhite("launch load golang"),
	DisableFlagParsing: false,
	Run: gbLoadFunc,
	Args: cobra.ExactArgs(1),
}
