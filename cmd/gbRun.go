package cmd

import (
	"launch/only"
	"launch/ux"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(gbRunCmd)
	rootCmd.AddCommand(gbShellCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//gbRunCmd.PersistentFlags().String("help", "", "Run and usage")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gbRunCmd.Flags().BoolP("toggle", "t", false, "Run message for toggle")
}


// gbRunCmd represents the gbRun command
var gbRunCmd = &cobra.Command{
	Use:   "run",
	Short: ux.SprintfBlue("Run default Gearbox gear command"),
	Long: ux.SprintfBlue("Run default Gearbox gear command."),
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	Run: gbRunFunc,
}

func gbRunFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		// var err error
		showArgs(cmd, args)
	}
}


// gbShellCmd represents the gbShell command
var gbShellCmd = &cobra.Command{
	Use:   "shell",
	Short: ux.SprintfBlue("Execute shell in Gearbox gear"),
	Long: ux.SprintfBlue("Execute shell in Gearbox gear."),
	DisableFlagParsing: true,
	Run: gbShellFunc,
}

func gbShellFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		// var err error
		showArgs(cmd, args)
	}
}

