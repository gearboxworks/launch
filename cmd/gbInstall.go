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

		if found {
			if !gearArgs.Temporary {
				// Create symlinks.
				gearArgs.GearRef.GearConfig.CreateLinks(defaults.RunAs, gearArgs.Name, gearArgs.Version)
			}

			state.SetOk("Gear '%s:%s' already installed.", gearArgs.Name, gearArgs.Version)
			break
		}
		state.ClearAll()


		state = gearArgs.GearRef.Docker.NetworkCreate(defaults.GearboxNetwork)
		if state.IsError() {
			break
		}
		state.ClearAll()

		if !gearArgs.Quiet {
			ux.Printf("Installing Gear '%s:%s': ", gearArgs.Name, gearArgs.Version)
		}
		state = gearArgs.GearRef.Docker.Container.ContainerCreate(gearArgs.Name, gearArgs.Version, gearArgs.Project)
		if state.IsError() {
			if !gearArgs.Quiet {
				ux.PrintfRed("error installing - %s\n", state.Error)
			}

			state.SetError("Gear '%s:%s' install error - %s", gearArgs.Name, gearArgs.Version, state.Error)
			break
		}

		if state.IsCreated() {
			if !gearArgs.Quiet {
				ux.PrintfGreen("OK\n")
			}

			state = gearArgs.GearRef.State()
			if state.IsError() {
				break
			}

			if gearArgs.Temporary {
				state.ClearAll()
				break
			}

			// Create symlinks.
			gearArgs.GearRef.GearConfig.CreateLinks(defaults.RunAs, gearArgs.Name, gearArgs.Version)
			state.SetOk("Gear '%s:%s' installed OK", gearArgs.Name, gearArgs.Version)
			break
		}

		if !gearArgs.Quiet {
			ux.PrintfWarning("cannot be installed\n")
		}
		state.SetWarning("Gear '%s:%s' cannot be installed", gearArgs.Name, gearArgs.Version)
	}

	cmdState = state
}


// gbInstallCmd represents the gbInstall command
var gbUninstallCmd = &cobra.Command{
	Use:   "uninstall <gear name>",
	//Aliases: []string{"remove"},
	SuggestFor: []string{"clean", "remove"},
	Short: ux.SprintfBlue("Uninstall a Gearbox gear"),
	Long: ux.SprintfBlue("Uninstall a Gearbox gear."),
	Example: ux.SprintfWhite("launch uninstall golang"),
	DisableFlagParsing: false,
	Run: gbUninstallFunc,
	Args: cobra.ExactArgs(1),
}

func gbUninstallFunc(cmd *cobra.Command, args []string) {
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
			state.SetOk("Gear '%s:%s' already removed.", gearArgs.Name, gearArgs.Version)
			break
		}
		state.ClearAll()

		gbStopFunc(cmd, args)
		if cmdState.IsError() {
			state = cmdState
			break
		}

		if !gearArgs.Quiet {
			ux.Printf("Removing gear '%s:%s': ", gearArgs.Name, gearArgs.Version)
		}
		state = gearArgs.GearRef.Docker.Container.Remove()
		if state.IsError() {
			if !gearArgs.Quiet {
				ux.PrintfRed("error removing - %s\n", state.Error)
			}

			state.SetError("Gear '%s:%s' remove error - %s", gearArgs.Name, gearArgs.Version, state.Error)
			break
		}

		if state.IsOk() {
			if !gearArgs.Quiet {
				ux.PrintfGreen("OK\n")
			}

			if gearArgs.Temporary {
				state.ClearAll()
				break
			}

			state.SetOk("Gear '%s:%s' removed OK", gearArgs.Name, gearArgs.Version)
			// Remove symlinks.
			gearArgs.GearRef.GearConfig.RemoveLinks(defaults.RunAs, gearArgs.Name, gearArgs.Version)
			break
		}

		if !gearArgs.Quiet {
			ux.PrintfWarning("cannot be removed\n")
		}
		state.SetWarning("Gear '%s:%s' cannot be removed", gearArgs.Name, gearArgs.Version)
	}

	cmdState = state
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
	var state ux.State

	for range only.Once {
		state = gearArgs.ProcessArgs(cmd, args)
		if state.IsError() {
			break
		}

		gbUninstallFunc(cmd, args)
		if cmdState.IsError() {
			state = cmdState
			break
		}

		gbInstallFunc(cmd, args)
		if cmdState.IsError() {
			state = cmdState
			break
		}

		state.SetOk("Gear '%s:%s' reinstalled.", gearArgs.Name, gearArgs.Version)
	}

	cmdState = state
}
