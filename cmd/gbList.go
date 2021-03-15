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

		switch {
			case len(args) == 0:
				_ = cmd.Help()
		}
	}
}


func gbDetailsFunc(cmd *cobra.Command, args []string) {
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

func gbPortsFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsError() {
			if Cmd.State.IsNotOk() {
				Cmd.State.PrintResponse()
			}
			break
		}

		Cmd.State = ga.gbPortsFunc()
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
		ga.State = ga.Gears.ListImages(ga.Name)
		if ga.State.IsError() {
			break
		}

		ga.State = ga.Gears.ListContainers(ga.Name)
		if ga.State.IsError() {
			break
		}

		ga.State = ga.Gears.NetworkList("gearboxnet")
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

func (ga *LaunchArgs) gbPortsFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		remote := false
		if Cmd.Host != "" {
			// We are remote.
			remote = true
		}

		ga.State = ga.ListPorts(remote)
		if ga.State.IsError() {
			break
		}
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}
