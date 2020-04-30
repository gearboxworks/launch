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
	for range only.Once {
		var ga *gearArgs

		ga, cmdState = getGearArgs(cmd, args)

		gearRef, cmdState = provider.NewGear()
		if cmdState.IsError() {
			break
		}

		var found bool
		found, cmdState = gearRef.FindContainer(ga.Name, ga.Version)
		if cmdState.IsError() {
			break
		}
		//cmdState.ClearAll()

		if !found {
			quietFlag = false
			if IsNoCreate(cmd) {
				cmdState.SetError("Gear '%s:%s' doesn't exist.", ga.Name, ga.Version)
				break
			}

			gbInstallFunc(cmd, args)
			if cmdState.IsError() {
				break
			}
			//cmdState.ClearAll()

			found, cmdState = gearRef.FindContainer(ga.Name, ga.Version)
			if cmdState.IsError() {
				break
			}
		}

		cmdState = gearRef.State()

		if cmdState.IsRunning() {
			break
		}


		if !quietFlag {
			ux.Printf("Starting gear '%s:%s': ", ga.Name, ga.Version)
		}
		cmdState = gearRef.Docker.Container.Start()
		if cmdState.IsError() {
			cmdState.SetError("Gear '%s:%s' start error - %s", ga.Name, ga.Version, cmdState.Error)
			ux.PrintfRed("%s\n", cmdState.Error)
		} else if cmdState.IsRunning() {
			if !quietFlag {
				ux.PrintfGreen("OK\n")
			}
		} else {
			cmdState.SetWarning("Gear '%s:%s' cannot be started", ga.Name, ga.Version)
			ux.PrintfWarning("%s\n", cmdState.Warning)
		}
	}
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
	for range only.Once {
		var ga *gearArgs

		ga, cmdState = getGearArgs(cmd, args)

		gearRef, cmdState = provider.NewGear()
		if cmdState.IsError() {
			break
		}

		var found bool
		found, cmdState = gearRef.FindContainer(ga.Name, ga.Version)
		if cmdState.IsError() {
			break
		}
		cmdState.ClearAll()

		if !found {
			cmdState.SetWarning("Gear '%s:%s' doesn't exist.", ga.Name, ga.Version)
			break
		}


		if !quietFlag {
			ux.Printf("Stopping gear '%s:%s': ", ga.Name, ga.Version)
		}
		cmdState = gearRef.Docker.Container.Stop()
		if cmdState.IsError() {
			if !quietFlag {
				ux.PrintfRed("error stopping - %s\n", cmdState.Error)
			}
			cmdState.SetError("Gear '%s:%s' stop error - %s", ga.Name, ga.Version, cmdState.Error)
		} else if cmdState.IsExited() {
			if !quietFlag {
				ux.PrintfGreen("OK\n")
			}
		} else {
			if !quietFlag {
				ux.PrintfWarning("cannot be stopped\n")
			}
			cmdState.SetWarning("Gear '%s:%s' cannot be stopped", ga.Name, ga.Version)
		}

		break
	}
}
