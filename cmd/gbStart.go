package cmd

import (
	"launch/only"
	"launch/ux"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(gbStartCmd)
	rootCmd.AddCommand(gbStopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//gbStartCmd.PersistentFlags().String("help", "", "Start and usage")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gbStartCmd.Flags().BoolP("toggle", "t", false, "Start message for toggle")
}


// gbStartCmd represents the gbStart command
var gbStartCmd = &cobra.Command{
	Use:   "start",
	Short: ux.SprintfBlue("Start a Gearbox gear"),
	Long: ux.SprintfBlue("Start a Gearbox gear."),
	DisableFlagParsing: false,
	Run: gbStartFunc,
}

func gbStartFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		// var err error
		showArgs(cmd, args)
	}
}


// gbStatusCmd represents the gbStatus command
var gbStopCmd = &cobra.Command{
	Use:   "stop",
	Short: ux.SprintfBlue("Stop a Gearbox gear"),
	Long: ux.SprintfBlue("Stop a Gearbox gear."),
	DisableFlagParsing: false,
	Run: gbStopFunc,
}

func gbStopFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		// var err error
		showArgs(cmd, args)
	}
}
