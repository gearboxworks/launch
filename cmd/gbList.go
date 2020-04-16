package cmd

import (
	"github.com/spf13/cobra"
	"launch/only"
	"launch/ux"
)

func init() {
	rootCmd.AddCommand(gbListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//gbListCmd.PersistentFlags().String("help", "", "List and usage")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gbListCmd.Flags().BoolP("toggle", "t", false, "List message for toggle")
}


// gbListCmd represents the gbList command
var gbListCmd = &cobra.Command{
	Use:   "list",
	Short: ux.SprintfBlue("List a Gearbox gear"),
	Long: ux.SprintfBlue("List a Gearbox gear."),
	DisableFlagParsing: false,
	Run: gbListFunc,
}

func gbListFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		// var err error
		showArgs(cmd, args)

		//fmt.Printf("%s", cmd.UsageTemplate())
	}
}
