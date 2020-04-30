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

		if found {
			// Create symlinks.
			gearRef.GearConfig.CreateLinks(defaults.RunAs, ga.Name, ga.Version)

			cmdState.SetOk("Gear '%s:%s' already installed.", ga.Name, ga.Version)
			break
		}
		cmdState.ClearAll()

		cmdState = gearRef.Docker.NetworkCreate(defaults.GearboxNetwork)
		if cmdState.IsError() {
			break
		}
		cmdState.ClearAll()

		if !quietFlag {
			ux.Printf("Installing Gear '%s:%s': ", ga.Name, ga.Version)
		}
		cmdState = gearRef.Docker.Container.ContainerCreate(ga.Name, ga.Version, ga.Mount)
		if cmdState.IsError() {
			if !quietFlag {
				ux.PrintfRed("error installing - %s\n", cmdState.Error)
			}
			cmdState.SetError("Gear '%s:%s' install error - %s", ga.Name, ga.Version, cmdState.Error)

		} else if cmdState.IsCreated() {
			if !quietFlag {
				ux.PrintfGreen("OK\n")
			}

			cmdState = gearRef.State()
			if cmdState.IsError() {
				break
			}

			// Create symlinks.
			gearRef.GearConfig.CreateLinks(defaults.RunAs, ga.Name, ga.Version)

			cmdState.SetOk("Gear '%s:%s' installed OK", ga.Name, ga.Version)

		} else {
			if !quietFlag {
				ux.PrintfWarning("cannot be installed\n")
			}
			cmdState.SetWarning("Gear '%s:%s' cannot be installed", ga.Name, ga.Version)
		}
	}
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
		if !found {
			cmdState.SetOk("Gear '%s:%s' already removed.", ga.Name, ga.Version)
			break
		}
		cmdState.ClearAll()

		gbStopFunc(cmd, args)
		if cmdState.IsError() {
			break
		}

		if !quietFlag {
			ux.Printf("Removing gear '%s:%s': ", ga.Name, ga.Version)
		}
		cmdState = gearRef.Docker.Container.Remove()
		if cmdState.IsError() {
			if !quietFlag {
				ux.PrintfRed("error removing - %s\n", cmdState.Error)
			}
			cmdState.SetError("Gear '%s:%s' remove error - %s", ga.Name, ga.Version, cmdState.Error)

		} else if cmdState.IsOk() {
			if !quietFlag {
				ux.PrintfGreen("OK\n")
			}
			cmdState.SetOk("Gear '%s:%s' removed OK", ga.Name, ga.Version)

			// Remove symlinks.
			gearRef.GearConfig.RemoveLinks(defaults.RunAs, ga.Name, ga.Version)

		} else {
			if !quietFlag {
				ux.PrintfWarning("cannot be removed\n")
			}
			cmdState.SetWarning("Gear '%s:%s' cannot be removed", ga.Name, ga.Version)
		}
	}
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
	for range only.Once {
		var ga *gearArgs

		ga, cmdState = getGearArgs(cmd, args)

		gbUninstallFunc(cmd, args)
		if cmdState.IsError() {
			break
		}

		gbInstallFunc(cmd, args)
		if cmdState.IsError() {
			break
		}

		cmdState.SetOk("Gear '%s:%s' reinstalled.", ga.Name, ga.Version)
	}
}
