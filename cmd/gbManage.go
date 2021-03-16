package cmd

import (
	"github.com/newclarity/scribeHelpers/toolGear"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
)


func gbManageFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		//switch {
		//	case len(args) == 0:
				_ = cmd.Help()
		//}
	}
}


func gbInstallFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbInstallFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}

func gbUninstallFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbUninstallFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}

func gbReinstallFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbReinstallFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}

func gbCleanFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbCleanFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}

func gbStartFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbStartFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}

func gbStopFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbStopFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}

func gbLogsFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		if len(args) == 0 {
			Cmd.State.SetError("Need to specify a %s.", defaults.LanguageContainerName)
			break
		}

		Cmd.State = ga.gbLogsFunc()
	}
}


func (ga *LaunchArgs) gbInstallFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var found bool
		found, ga.State = ga.Gears.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}

		if found {
			if !ga.Temporary {
				ga.State = ga.Gears.FindImage(ga.Name, ga.Version)
				found = ga.Gears.State.GetResponseAsBool()
				if found {
					// Create symlinks.
					ga.CreateLinks(ga.Version)
				}
			}

			ga.State.SetOk("%s '%s:%s' already installed.", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State = ga.Gears.NetworkCreate(defaults.GearboxNetwork)
		if ga.State.IsError() {
			break
		}


		if ga.Project != toolGear.DefaultPathNone {
			ga.Gears.Selected.AddVolume(ga.Project, toolGear.DefaultProject)
		}

		if ga.TmpDir != toolGear.DefaultPathNone {
			ga.Gears.Selected.AddVolume(ga.TmpDir, toolGear.DefaultTmpDir)
		}


		if !ga.Quiet {
			ux.PrintflnNormal("Installing %s '%s:%s'.", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.Gears.ContainerCreate(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}

		if ga.State.IsCreated() {
			if ga.Temporary {
				ga.State.Clear()
				break
			}

			// Create symlinks.
			ga.State = ga.CreateLinks(ga.Version)

			ga.State.SetOk("Installed %s '%s:%s' OK.", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetWarning("%s '%s:%s' cannot be installed", defaults.LanguageContainerName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}

func (ga *LaunchArgs) gbUninstallFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var found bool
		found, ga.State = ga.Gears.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		if !found {
			if !ga.Temporary {
				ga.State = ga.Gears.FindImage(ga.Name, ga.Version)
				found = ga.Gears.State.GetResponseAsBool()
				if found {
					// Remove symlinks.
					ga.RemoveLinks(ga.Version)
				}
			}

			ga.State.SetOk("%s '%s:%s' already removed.", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}
		ga.State.Clear()

		ga.State = ga.gbStopFunc()
		if ga.State.IsError() {
			break
		}

		if !ga.Quiet {
			ux.PrintflnNormal("Removing %s '%s:%s'.\n", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.Gears.Selected.Container.Remove()
		if ga.State.IsError() {
			ga.State.SetError("%s '%s:%s' remove error - %s", defaults.LanguageContainerName, ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsOk() {
			if ga.Temporary {
				ga.State.Clear()
				break
			}

			// Remove symlinks.
			ga.RemoveLinks(ga.Version)

			ga.State.SetOk("%s '%s:%s' removed OK", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetWarning("%s '%s:%s' cannot be removed", defaults.LanguageContainerName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}

func (ga *LaunchArgs) gbReinstallFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.gbUninstallFunc()
		if ga.State.IsError() {
			break
		}

		ga.State = ga.gbInstallFunc()
		if ga.State.IsError() {
			break
		}

		ga.State.SetOk("%s '%s:%s' reinstalled.", defaults.LanguageContainerName, ga.Name, ga.Version)
		ga.State.SetOutput("")
	}

	return ga.State
}

func (ga *LaunchArgs) gbCleanFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.gbUninstallFunc()
		if ga.State.IsError() {
			break
		}

		var found bool
		ga.State = ga.Gears.FindImage(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		found = ga.Gears.State.GetResponseAsBool()
		if !found {
			ga.State.SetOk("%s '%s:%s' already removed.", defaults.LanguageImageName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}
		ga.State.Clear()

		if !ga.Quiet {
			ux.PrintflnNormal("Removing %s '%s:%s': ", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.Gears.Selected.Image.Remove()
		if ga.State.IsError() {
			ga.State.SetError("%s '%s:%s' remove error - %s", defaults.LanguageImageName, ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsOk() {
			if ga.Temporary {
				ga.State.Clear()
				break
			}

			ga.State.SetOk("%s '%s:%s' removed OK", defaults.LanguageImageName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetWarning("%s '%s:%s' cannot be removed", defaults.LanguageImageName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}

func (ga *LaunchArgs) gbStartFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var found bool
		found, ga.State = ga.Gears.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		if !found {
			if !ga.Temporary {
				ga.Quiet = false
			}

			if ga.NoCreate {
				ga.State.SetError("Not creating %s '%s:%s'.", defaults.LanguageContainerName, ga.Name, ga.Version)
				break
			}

			ga.gbInstallFunc()
			if ga.State.IsError() {
				ga.State.SetError("Cannot start %s '%s:%s'.", defaults.LanguageContainerName, ga.Name, ga.Version)
				break
			}

			// Need a better way to handle the "Docker client error: context deadline exceeded" errors.
		}

		if ga.State.IsRunning() {
			ga.State.SetOk("%s '%s:%s' already started.", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}


		if !ga.Quiet {
			ux.PrintflnNormal("Starting %s '%s:%s': ", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.Gears.Selected.Start()
		if ga.State.IsError() {
			ga.State.SetError("%s '%s:%s' start error - %s", defaults.LanguageContainerName, ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsRunning() {
			ga.State.SetOk("%s '%s:%s' started OK", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetError("%s '%s:%s' cannot be started", defaults.LanguageContainerName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}

func (ga *LaunchArgs) gbStopFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var found bool
		found, ga.State = ga.Gears.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		if !found {
			break
		}
		if ga.State.IsExited() {
			ga.State.SetOk("%s '%s:%s' already stopped.", defaults.LanguageImageName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}


		if !ga.Quiet {
			ux.PrintflnNormal("Stopping %s '%s:%s': ", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.Gears.Selected.Stop()
		if ga.State.IsError() {
			ga.State.SetError("%s '%s:%s' stop error - %s", defaults.LanguageContainerName, ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsExited() {
			ga.State.SetOk("%s '%s:%s' stopped OK", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		if ga.State.IsCreated() {
			ga.State.SetOk("%s '%s:%s' stopped OK", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetWarning("%s '%s:%s' cannot be stopped", defaults.LanguageContainerName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}

func (ga *LaunchArgs) gbLogsFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var found bool
		found, ga.State = ga.Gears.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		if !found {
			break
		}
		//if ga.State.IsExited() {
		//	ga.State.SetOk("%s '%s:%s' already stopped.", defaults.LanguageImageName, ga.Name, ga.Version)
		//	ga.State.SetOutput("")
		//	break
		//}


		//if !ga.Quiet {
		//	ux.PrintflnNormal("Stopping %s '%s:%s': ", defaults.LanguageContainerName, ga.Name, ga.Version)
		//}
		ga.State = ga.Gears.Selected.Logs()
		if ga.State.IsError() {
			ga.State.SetError("%s '%s:%s' stop error - %s", defaults.LanguageContainerName, ga.Name, ga.Version, ga.State.GetError())
			break
		}

		//if ga.State.IsExited() {
		//	ga.State.SetOk("%s '%s:%s' stopped OK", defaults.LanguageContainerName, ga.Name, ga.Version)
		//	ga.State.SetOutput("")
		//	break
		//}
		//
		//if ga.State.IsCreated() {
		//	ga.State.SetOk("%s '%s:%s' stopped OK", defaults.LanguageContainerName, ga.Name, ga.Version)
		//	ga.State.SetOutput("")
		//	break
		//}
		//
		//ga.State.SetWarning("%s '%s:%s' cannot be stopped", defaults.LanguageContainerName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}
