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

		state = ga.gbStartFunc()
		if state.IsError() {
			if state.IsNotOk() {
				state.PrintResponse()
			}
			break
		}
	}

	_cmdState = state
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

		state = ga.gbStopFunc()
		if state.IsError() {
			if state.IsNotOk() {
				state.PrintResponse()
			}
			break
		}
	}

	_cmdState = state
}


func (ga *GearArgs) gbStartFunc() *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range only.Once {
		var found bool
		found, ga.State = ga.GearRef.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		if !found {
			if !ga.Temporary {
				ga.Quiet = false
			}

			if ga.NoCreate {
				ga.State.SetError("Not creating Gear '%s:%s'.", ga.Name, ga.Version)
				break
			}

			ga.gbInstallFunc()
			if ga.State.IsError() {
				ga.State.SetError("Cannot start Gear '%s:%s'.", ga.Name, ga.Version)
				break
			}

			// Need a better way to handle the "Docker client error: context deadline exceeded" errors.
		}

		if ga.State.IsRunning() {
			ga.State.SetOk("Gear '%s:%s' already started.", ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}


		if !ga.Quiet {
			ux.PrintflnNormal("Starting gear '%s:%s': ", ga.Name, ga.Version)
		}
		ga.State = ga.GearRef.Docker.Container.Start()
		if ga.State.IsError() {
			ga.State.SetError("Gear '%s:%s' start error - %s", ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsRunning() {
			ga.State.SetOk("Gear '%s:%s' started OK", ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetError("Gear '%s:%s' cannot be started", ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}

func (ga *GearArgs) gbStopFunc() *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range only.Once {
		var found bool
		found, ga.State = ga.GearRef.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		if !found {
			break
		}
		if ga.State.IsExited() {
			ga.State.SetOk("Gear image '%s:%s' already stopped.", ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}


		if !ga.Quiet {
			ux.PrintflnNormal("Stopping gear '%s:%s': ", ga.Name, ga.Version)
		}
		ga.State = ga.GearRef.Docker.Container.Stop()
		if ga.State.IsError() {
			ga.State.SetError("Gear '%s:%s' stop error - %s", ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsExited() {
			ga.State.SetOk("Gear '%s:%s' stopped OK", ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		if ga.State.IsCreated() {
			ga.State.SetOk("Gear '%s:%s' stopped OK", ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetWarning("Gear '%s:%s' cannot be stopped", ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}
