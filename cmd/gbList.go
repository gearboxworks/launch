package cmd

import (
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
)


func gbListFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		Cmd.State = ga.gbListFunc()
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}
	}
}


func gbLinksFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		Cmd.State = ga.gbLinksFunc()
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}
	}
}


func (ga *LaunchArgs) gbListFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.GearRef.Docker.List(ga.Name)
		if ga.State.IsError() {
			break
		}
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}


func (ga *LaunchArgs) gbLinksFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		//ga.State = ga.GearRef.GearConfig.CheckFile("latest", "mb1", "/usr/local/bin/mb1")
		//ga.State.PrintResponse()
		//ga.State = ga.GearRef.GearConfig.CheckFile("latest", "mb2", "/usr/local/bin/mb2")
		//ga.State.PrintResponse()
		//ga.State = ga.GearRef.GearConfig.CheckFile("latest", "mb3", "/usr/local/bin/mb3")
		//ga.State.PrintResponse()
		//ga.State = ga.GearRef.GearConfig.CheckFile("latest", "mb4", "/usr/local/bin/mb4")
		//ga.State.PrintResponse()
		//ga.State = ga.GearRef.GearConfig.CheckFile("latest", "mb5", "/usr/local/bin/mb5")
		//ga.State.PrintResponse()
		//ga.State = ga.GearRef.GearConfig.CheckFile("latest", "mb6", "/usr/local/bin/mb6")
		//ga.State.PrintResponse()
		//ga.State = ga.GearRef.GearConfig.CheckFile("latest", "mb7", "/usr/local/bin/mb7")
		//ga.State.PrintResponse()
		//ga.State = ga.GearRef.GearConfig.CheckFile("latest", "mb8", "/usr/local/bin/mb8")
		//ga.State.PrintResponse()
		//ga.State = ga.GearRef.GearConfig.CheckFile("latest", "mb9", "/usr/local/bin/mb9")
		//ga.State.PrintResponse()
		//
		//ga.Name = "mountebank"
		//ga.Version = "2.4.0"
		//ga.GearRef.FindContainer(ga.Name, ga.Version)
		////ga.State = ga.CreateLinks(ga.Version)
		//ga.State = ga.RemoveLinks(ga.Version)

		remote := false
		if Cmd.Host != "" {
			// We are remote.
			remote = true
		}

		ga.State = ga.ListLinks(remote)
		if ga.State.IsError() {
			break
		}
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}
