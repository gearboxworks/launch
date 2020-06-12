package cmd

import (
	"github.com/spf13/cobra"
	"github.com/newclarity/scribeHelpers/ux"
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
	Use:   "export <gear name>",
	//Aliases: []string{"save"},
	SuggestFor: []string{"save"},
	Short: ux.SprintfBlue("Save state of a Gearbox gear"),
	Long: ux.SprintfBlue("Save state of a Gearbox gear."),
	Example: ux.SprintfWhite("launch save golang"),
	DisableFlagParsing: false,
	Run: gbSaveFunc,
	Args: cobra.ExactArgs(1),
}
func gbSaveFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		//showArgs(rootCmd, args)
		ux.PrintfWarning("Command not yet implemented.\n")
	}
}


// gbLoadCmd represents the gbLoad command
var gbLoadCmd = &cobra.Command{
	Use:   "import <gear name>",
	//Aliases: []string{"load"},
	SuggestFor: []string{"load"},
	Short: ux.SprintfBlue("Load a Gearbox gear"),
	Long: ux.SprintfBlue("Load a Gearbox gear."),
	Example: ux.SprintfWhite("launch load golang"),
	DisableFlagParsing: false,
	Run: gbLoadFunc,
	Args: cobra.ExactArgs(1),
}
func gbLoadFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		//showArgs(rootCmd, args)
		ux.PrintfWarning("Command not yet implemented.\n")
	}
}
