package cmd

import (
	"github.com/newclarity/scribeHelpers/toolGear"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
)


func gbInstallFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		Cmd.State = ga.gbInstallFunc()
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}
	}
}


func gbUninstallFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		Cmd.State = ga.gbUninstallFunc()
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}
	}
}


func gbReinstallFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		Cmd.State = ga.gbReinstallFunc()
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}
	}
}


func gbCleanFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		Cmd.State = ga.gbCleanFunc()
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}
	}
}


func (ga *LaunchArgs) gbInstallFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var found bool
		found, ga.State = ga.GearRef.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}

		if found {
			if !ga.Temporary {
				found, ga.State = ga.GearRef.FindImage(ga.Name, ga.Version)
				if found {
					// Create symlinks.
					ga.CreateLinks(ga.Version)
				}
			}

			ga.State.SetOk("%s '%s:%s' already installed.", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State = ga.GearRef.Docker.NetworkCreate(defaults.GearboxNetwork)
		if ga.State.IsError() {
			break
		}


		if ga.Project != toolGear.DefaultPathNone {
			ga.GearRef.AddVolume(ga.Project, toolGear.DefaultProject)
		}

		if ga.TmpDir != toolGear.DefaultPathNone {
			ga.GearRef.AddVolume(ga.TmpDir, toolGear.DefaultTmpDir)
		}


		if !ga.Quiet {
			ux.PrintflnNormal("Installing %s '%s:%s'.", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.GearRef.ContainerCreate(ga.Name, ga.Version)
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
		found, ga.State = ga.GearRef.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		if !found {
			if !ga.Temporary {
				found, ga.State = ga.GearRef.FindImage(ga.Name, ga.Version)
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
		ga.State = ga.GearRef.Docker.Container.Remove()
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
		found, ga.State = ga.GearRef.FindImage(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		if !found {
			ga.State.SetOk("%s '%s:%s' already removed.", defaults.LanguageImageName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}
		ga.State.Clear()

		if !ga.Quiet {
			ux.PrintflnNormal("Removing %s '%s:%s': ", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.GearRef.Docker.Image.Remove()
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
