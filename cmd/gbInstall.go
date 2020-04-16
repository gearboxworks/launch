package cmd

import (
	"launch/only"
	"launch/ux"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(gbInstallCmd)
	rootCmd.AddCommand(gbUninstallCmd)
	rootCmd.AddCommand(gbReinstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//gbInstallCmd.PersistentFlags().String("help", "", "Remove and usage")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gbInstallCmd.Flags().BoolP("toggle", "t", false, "Remove message for toggle")
}


// gbInstallCmd represents the gbInstall command
var gbInstallCmd = &cobra.Command{
	Use:   "install",
	Aliases: []string{"add"},
	SuggestFor: []string{"download"},
	Short: ux.SprintfBlue("Install a Gearbox gear"),
	Long: ux.SprintfBlue("Install a Gearbox gear."),
	DisableFlagParsing: false,
	Run: gbInstallFunc,
}

func gbInstallFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		// var err error
		showArgs(cmd, args)
	}
}


// gbInstallCmd represents the gbInstall command
var gbUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Aliases: []string{"remove"},
	SuggestFor: []string{"clean"},
	Short: ux.SprintfBlue("Uninstall a Gearbox gear"),
	Long: ux.SprintfBlue("Uninstall a Gearbox gear."),
	DisableFlagParsing: false,
	Run: gbUninstallFunc,
}

func gbUninstallFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		// var err error
		showArgs(cmd, args)
	}
}


// gbReinstallCmd represents the gbInstall command
var gbReinstallCmd = &cobra.Command{
	Use:   "reinstall",
	Aliases: []string{"update"},
	SuggestFor: []string{""},
	Short: ux.SprintfBlue("Update a Gearbox gear"),
	Long: ux.SprintfBlue("Update a Gearbox gear."),
	DisableFlagParsing: false,
	Run: gbReinstallFunc,
}

func gbReinstallFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		// var err error
		showArgs(cmd, args)
	}
}
