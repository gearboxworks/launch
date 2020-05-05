package cmd

import (
	"github.com/spf13/cobra"
	"launch/only"
	"launch/ux"
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
	Use:   "start <gear name>",
	Short: ux.SprintfBlue("Start a Gearbox gear"),
	Long: ux.SprintfBlue("Start a Gearbox gear."),
	Example: ux.SprintfWhite("launch start golang"),
	DisableFlagParsing: false,
	Run: gbStartFunc,
	Args: cobra.ExactArgs(1),
}

func gbStartFunc(cmd *cobra.Command, args []string) {
	var state ux.State

	for range only.Once {
		state = gearArgs.ProcessArgs(cmd, args)
		if state.IsError() {
			break
		}

		var found bool
		found, state = gearArgs.GearRef.FindContainer(gearArgs.Name, gearArgs.Version)
		if state.IsError() {
			break
		}
		if !found {
			if !gearArgs.Temporary {
				gearArgs.Quiet = false
			}

			if IsNoCreate(cmd) {
				state.SetError("Gear '%s:%s' doesn't exist.", gearArgs.Name, gearArgs.Version)
				break
			}

			gbInstallFunc(cmd, args)
			if cmdState.IsError() {
				state = cmdState
				break
			}
			state.ClearAll()

			found, state = gearArgs.GearRef.FindContainer(gearArgs.Name, gearArgs.Version)
			if state.IsError() {
				break
			}
		}

		//state = gearArgs.GearRef.State()
		if state.IsRunning() {
			if !gearArgs.Quiet {
				ux.PrintfGreen("Gear '%s:%s' already started\n", gearArgs.Name, gearArgs.Version)
			}
			break
		}


		if !gearArgs.Quiet {
			ux.Printf("Starting gear '%s:%s': ", gearArgs.Name, gearArgs.Version)
		}
		state = gearArgs.GearRef.Docker.Container.Start()
		if state.IsError() {
			state.SetError("Gear '%s:%s' start error - %s", gearArgs.Name, gearArgs.Version, state.Error)
			ux.PrintfRed("%s\n", state.Error)
			break
		}

		if state.IsRunning() {
			if !gearArgs.Quiet {
				ux.PrintfGreen("OK\n")
			}
			break
		}

		state.SetWarning("Gear '%s:%s' cannot be started", gearArgs.Name, gearArgs.Version)
		ux.PrintfWarning("%s\n", state.Warning)
	}

	cmdState = state
}


// gbStatusCmd represents the gbStatus command
var gbStopCmd = &cobra.Command{
	Use:   "stop <gear name>",
	Short: ux.SprintfBlue("Stop a Gearbox gear"),
	Long: ux.SprintfBlue("Stop a Gearbox gear."),
	Example: ux.SprintfWhite("launch stop golang"),
	DisableFlagParsing: false,
	Run: gbStopFunc,
	Args: cobra.ExactArgs(1),
}

func gbStopFunc(cmd *cobra.Command, args []string) {
	var state ux.State

	for range only.Once {
		state = gearArgs.ProcessArgs(cmd, args)
		if state.IsError() {
			break
		}

		var found bool
		found, state = gearArgs.GearRef.FindContainer(gearArgs.Name, gearArgs.Version)
		if state.IsError() {
			break
		}
		if !found {
			break
		}
		if state.IsExited() {
			if !gearArgs.Quiet {
				ux.PrintfGreen("Gear '%s:%s' already stopped\n", gearArgs.Name, gearArgs.Version)
			}
			break
		}


		if !gearArgs.Quiet {
			ux.Printf("Stopping gear '%s:%s': ", gearArgs.Name, gearArgs.Version)
		}
		state = gearArgs.GearRef.Docker.Container.Stop()
		if state.IsError() {
			if !gearArgs.Quiet {
				ux.PrintfRed("error stopping - %s\n", state.Error)
			}
			state.SetError("Gear '%s:%s' stop error - %s", gearArgs.Name, gearArgs.Version, state.Error)
			break
		}

		if state.IsExited() {
			if !gearArgs.Quiet {
				ux.PrintfGreen("OK\n")
			}
			break
		}

		if state.IsCreated() {
			if !gearArgs.Quiet {
				ux.PrintfGreen("OK\n")
			}
			break
		}

		if !gearArgs.Quiet {
			ux.PrintfWarning("cannot be stopped\n")
		}
		state.SetWarning("Gear '%s:%s' cannot be stopped", gearArgs.Name, gearArgs.Version)
	}

	cmdState = state
}
