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
	var state *ux.State

	for range only.Once {
		var ga GearArgs

		state = ga.ProcessArgs(rootCmd, args)
		if state.IsError() {
			if state.IsNotOk() {
				state.PrintResponse()
			}
			break
		}

		state = ga.gbListFunc()
		if state.IsError() {
			if state.IsNotOk() {
				state.PrintResponse()
			}
			break
		}
	}

	_cmdState = state
}


func (ga *GearArgs) gbListFunc() *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range only.Once {
		_, ga.State = ga.GearRef.Docker.ImageList(ga.Name)
		if ga.State.IsError() {
			break
		}

		_, ga.State = ga.GearRef.Docker.ContainerList(ga.Name)
		if ga.State.IsError() {
			break
		}

		ga.State = ga.GearRef.Docker.NetworkList(defaults.GearboxNetwork)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}
