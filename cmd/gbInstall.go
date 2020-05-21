package cmd

import (
	"github.com/spf13/cobra"
	"launch/defaults"
	"launch/only"
	"launch/ux"
)

func init() {
	rootCmd.AddCommand(gbInstallCmd)
	rootCmd.AddCommand(gbUninstallCmd)
	rootCmd.AddCommand(gbReinstallCmd)
	rootCmd.AddCommand(gbCleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//gbInstallCmd.PersistentFlags().String("help", "", "Remove and usage")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gbInstallCmd.Flags().BoolP("toggle", "t", false, "Remove message for toggle")
}


// gbInstallCmd represents the gbInstall command
var gbInstallCmd = &cobra.Command{
	Use:   "install <gear name>",
	//Aliases: []string{"add"},
	SuggestFor: []string{"download" ,"add"},
	Short: ux.SprintfBlue("Install a Gearbox gear"),
	Long: ux.SprintfBlue("Install a Gearbox gear."),
	Example: ux.SprintfWhite("launch install golang"),
	DisableFlagParsing: false,
	Run: gbInstallFunc,
	Args: cobra.ExactArgs(1),
}
func gbInstallFunc(cmd *cobra.Command, args []string) {
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

		state = ga.gbInstallFunc()
		if state.IsError() {
			if state.IsNotOk() {
				state.PrintResponse()
			}
			break
		}
	}

	_cmdState = state
}


// gbInstallCmd represents the gbInstall command
var gbUninstallCmd = &cobra.Command{
	Use:   "uninstall <gear name>",
	//Aliases: []string{"remove"},
	SuggestFor: []string{"remove"},
	Short: ux.SprintfBlue("Uninstall a Gearbox gear"),
	Long: ux.SprintfBlue("Uninstall a Gearbox gear."),
	Example: ux.SprintfWhite("launch uninstall golang"),
	DisableFlagParsing: false,
	Run: gbUninstallFunc,
	Args: cobra.ExactArgs(1),
}
func gbUninstallFunc(cmd *cobra.Command, args []string) {
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

		state = ga.gbUninstallFunc()
		if state.IsError() {
			if state.IsNotOk() {
				state.PrintResponse()
			}
			break
		}
	}

	_cmdState = state
}


// gbReinstallCmd represents the gbInstall command
var gbReinstallCmd = &cobra.Command{
	Use:   "reinstall <gear name>",
	//Aliases: []string{"update"},
	SuggestFor: []string{"update"},
	Short: ux.SprintfBlue("Update a Gearbox gear"),
	Long: ux.SprintfBlue("Update a Gearbox gear."),
	Example: ux.SprintfWhite("launch reinstall golang"),
	DisableFlagParsing: false,
	Run: gbReinstallFunc,
	Args: cobra.ExactArgs(1),
}
func gbReinstallFunc(cmd *cobra.Command, args []string) {
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

		state = ga.gbReinstallFunc()
		if state.IsError() {
			if state.IsNotOk() {
				state.PrintResponse()
			}
			break
		}
	}

	_cmdState = state
}


// gbInstallCmd represents the gbInstall command
var gbCleanCmd = &cobra.Command{
	Use:   "clean <gear name>",
	//Aliases: []string{"remove"},
	SuggestFor: []string{},
	Short: ux.SprintfBlue("Completely uninstall a Gearbox gear"),
	Long: ux.SprintfBlue("Completely uninstall a Gearbox gear."),
	Example: ux.SprintfWhite("launch uninstall golang"),
	DisableFlagParsing: false,
	Run: gbCleanFunc,
	Args: cobra.ExactArgs(1),
}
func gbCleanFunc(cmd *cobra.Command, args []string) {
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

		state = ga.gbCleanFunc()
		if state.IsError() {
			if state.IsNotOk() {
				state.PrintResponse()
			}
			break
		}
	}

	_cmdState = state
}


func (ga *GearArgs) gbInstallFunc() *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range only.Once {
		var found bool
		found, ga.State = ga.GearRef.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}

		if found {
			if !ga.Temporary {
				// Create symlinks.
				//ga.GearRef.GearConfig.CreateLinks(defaults.RunAs, ga.Name, ga.Version)
				ga.State = ga.GearRef.GearConfig.CreateLinks(defaults.RunAs, ga.Version)
			}

			ga.State.SetOk("Gear '%s:%s' already installed.", ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}
		//ga.State.Clear()


		ga.State = ga.GearRef.Docker.NetworkCreate(defaults.GearboxNetwork)
		if ga.State.IsError() {
			break
		}
		//ga.State.Clear()


		if !ga.Quiet {
			ux.PrintflnNormal("Installing Gear '%s:%s'.", ga.Name, ga.Version)
		}
		ga.State = ga.GearRef.Docker.Container.ContainerCreate(ga.Name, ga.Version, ga.Project)
		if ga.State.IsError() {
			break
		}

		if ga.State.IsCreated() {
			ga.State = ga.GearRef.Status()
			if ga.State.IsError() {
				break
			}

			if ga.Temporary {
				ga.State.Clear()
				break
			}

			// Create symlinks.
			ga.State = ga.GearRef.GearConfig.CreateLinks(defaults.RunAs, ga.Version)

			ga.State.SetOk("Installed Gear '%s:%s' OK.", ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetWarning("Gear '%s:%s' cannot be installed", ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}


func (ga *GearArgs) gbUninstallFunc() *ux.State {
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
			ga.State.SetOk("Gear '%s:%s' already removed.", ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}
		ga.State.Clear()

		ga.State = ga.gbStopFunc()
		if ga.State.IsError() {
			break
		}

		if !ga.Quiet {
			ux.PrintflnNormal("Removing gear '%s:%s'.\n", ga.Name, ga.Version)
		}
		ga.State = ga.GearRef.Docker.Container.Remove()
		if ga.State.IsError() {
			ga.State.SetError("Gear '%s:%s' remove error - %s", ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsOk() {
			if ga.Temporary {
				ga.State.Clear()
				break
			}

			// Remove symlinks.
			ga.GearRef.GearConfig.RemoveLinks(defaults.RunAs, ga.Version)

			ga.State.SetOk("Gear '%s:%s' removed OK", ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetWarning("Gear '%s:%s' cannot be removed", ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}


func (ga *GearArgs) gbReinstallFunc() *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range only.Once {
		ga.State = ga.gbUninstallFunc()
		if ga.State.IsError() {
			break
		}

		ga.State = ga.gbInstallFunc()
		if ga.State.IsError() {
			break
		}

		ga.State.SetOk("Gear '%s:%s' reinstalled.", ga.Name, ga.Version)
		ga.State.SetOutput("")
	}

	return ga.State
}


func (ga *GearArgs) gbCleanFunc() *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range only.Once {
		ga.State = ga.gbUninstallFunc()
		if ga.State.IsError() {
			break
		}

		var found bool
		found, ga.State = ga.GearRef.FindImage(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		if !found {
			ga.State.SetOk("Gear image '%s:%s' already removed.", ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}
		ga.State.Clear()

		if !ga.Quiet {
			ux.PrintflnNormal("Removing gear '%s:%s': ", ga.Name, ga.Version)
		}
		ga.State = ga.GearRef.Docker.Image.Remove()
		if ga.State.IsError() {
			ga.State.SetError("Gear image '%s:%s' remove error - %s", ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsOk() {
			if ga.Temporary {
				ga.State.Clear()
				break
			}

			ga.State.SetOk("Gear image '%s:%s' removed OK", ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetWarning("Gear image '%s:%s' cannot be removed", ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}
