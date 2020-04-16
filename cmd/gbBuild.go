package cmd

import (
	"launch/only"
	"launch/ux"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(gbBuildCmd)
	rootCmd.AddCommand(gbPublishCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//gbBuildCmd.PersistentFlags().String("help", "", "Build and usage")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gbBuildCmd.Flags().BoolP("toggle", "t", false, "Build message for toggle")
}


// gbBuildCmd represents the gbBuild command
var gbBuildCmd = &cobra.Command {
	Use:   "build",
	SuggestFor: []string{ "compile", "generate" },
	Short: ux.SprintfBlue("Build a Gearbox gear"),
	Long: ux.SprintfBlue("Allows building of arbitrary containers as a Gearbox container, (called gears)."),
	//Example: "",
	DisableFlagParsing: false,
	Run: gbBuildFunc,
}

func gbBuildFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		//var err error
		showArgs(cmd, args)

		//var ok bool
		//ok, err = cmd.Flags().GetBool("help")
		//if err != nil {
		//	break
		//}
		//if ok {
		//	_ = cmd.Help()
		//}
	}
}


// gbPublishCmd represents the gbPublish command
var gbPublishCmd = &cobra.Command{
	Use:   "publish",
	Short: ux.SprintfBlue("Publish a Gearbox gear"),
	Long: ux.SprintfBlue("Publish a Gearbox gear to GitHub or DockerHub."),
	DisableFlagParsing: false,
	Run: gbPublishFunc,
}

func gbPublishFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		// var err error
		showArgs(cmd, args)
	}
}
