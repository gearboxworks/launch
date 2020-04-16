package cmd

import (
	"launch/only"
	"launch/ux"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(gbSaveCmd)
	rootCmd.AddCommand(gbLoadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//gbSaveCmd.PersistentFlags().String("help", "", "Save and usage")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gbSaveCmd.Flags().BoolP("toggle", "t", false, "Save message for toggle")
}


// gbSaveCmd represents the gbSave command
var gbSaveCmd = &cobra.Command{
	Use:   "export",
	Aliases: []string{"save"},
	Short: ux.SprintfBlue("Save state of a Gearbox gear"),
	Long: ux.SprintfBlue("Save state of a Gearbox gear."),
	DisableFlagParsing: false,
	Run: gbSaveFunc,
}

func gbSaveFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		// var err error
		showArgs(cmd, args)
	}
}


// gbLoadCmd represents the gbLoad command
var gbLoadCmd = &cobra.Command{
	Use:   "import",
	Aliases: []string{"load"},
	Short: ux.SprintfBlue("Load a Gearbox gear"),
	Long: ux.SprintfBlue("Load a Gearbox gear."),
	DisableFlagParsing: false,
	Run: gbLoadFunc,
}

func gbLoadFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		// var err error
		showArgs(cmd, args)
	}
}
