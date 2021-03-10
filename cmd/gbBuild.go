package cmd

import (
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
)


func gbBuildFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		Cmd.State = ga.gbBuildFunc()
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}
	}
}


func gbPublishFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		Cmd.State = ga.gbPublishFunc()
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}
	}
}


func (ga *LaunchArgs) gbBuildFunc() *ux.State {
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
		ga.State = ga.Gears.Selected.Container.Start()
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


func (ga *LaunchArgs) gbPublishFunc() *ux.State {
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
		ga.State = ga.Gears.Selected.Container.Stop()
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
