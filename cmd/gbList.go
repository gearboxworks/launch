package cmd

import (
	"github.com/spf13/cobra"
	"launch/defaults"
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
	Use:   "list [gear name]",
	Aliases: []string{"ls", "show"},
	//SuggestFor: []string{"ls", "show"},
	Short: ux.SprintfBlue("List a Gearbox gear"),
	Long: ux.SprintfBlue("List a Gearbox gear."),
	Example: ux.SprintfWhite("launch list golang"),
	DisableFlagParsing: false,
	Run: gbListFunc,
	Args: cobra.RangeArgs(0, 1),
}

func gbListFunc(cmd *cobra.Command, args []string) {
	var state ux.State

	for range only.Once {
		state = gearArgs.ProcessArgs(cmd, args)
		if state.IsError() {
			break
		}

		_, state = gearArgs.GearRef.Docker.ImageList(gearArgs.Name)
		if state.IsError() {
			break
		}

		_, state = gearArgs.GearRef.Docker.ContainerList(gearArgs.Name)
		if state.IsError() {
			break
		}

		state = gearArgs.GearRef.Docker.NetworkList(defaults.GearboxNetwork)
	}

	cmdState = state
}
