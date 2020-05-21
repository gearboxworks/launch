package cmd

import (
	"github.com/spf13/cobra"
	"launch/defaults"
	"launch/only"
	"launch/ux"
	"strings"
)

func init() {
	rootCmd.AddCommand(gbRunCmd)
	rootCmd.AddCommand(gbShellCmd)
	rootCmd.AddCommand(gbUnitTestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//gbRunCmd.PersistentFlags().String("help", "", "Run and usage")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//gbRunCmd.Flags().BoolP("status", "", false, "Show shell status line.")
	//gbRunCmd.Flags().BoolP(argTemporary, "t", false, ux.SprintfBlue("Temporary container - remove after running command."))
}


// gbRunCmd represents the gbRun command
var gbRunCmd = &cobra.Command{
	Use:   "run <gear name> [gear args]",
	Short: ux.SprintfBlue("Run default Gearbox gear command"),
	Long: ux.SprintfBlue("Run default Gearbox gear command."),
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	Example: ux.SprintfWhite("launch run golang build"),
	Run: gbRunFunc,
	Args: cobra.MinimumNArgs(1),
}
func gbRunFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		var ga GearArgs

		state := ga.ProcessArgs(rootCmd, args)
		if state.IsError() {
			break
		}

		state = ga.gbRunFunc()
		if state.IsError() {
			break
		}
	}
}


// gbShellCmd represents the gbShell command
var gbShellCmd = &cobra.Command{
	Use:   "shell <gear name> [command] [args]",
	Short: ux.SprintfBlue("Execute shell in Gearbox gear"),
	Long: ux.SprintfBlue("Execute shell in Gearbox gear."),
	Example: ux.SprintfWhite("launch shell mysql ps -eaf"),
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	Run: gbShellFunc,
	Args: cobra.MinimumNArgs(1),
}
func gbShellFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		var ga GearArgs

		state := ga.ProcessArgs(rootCmd, args)
		if state.IsError() {
			break
		}

		state = ga.gbShellFunc()
		if state.IsError() {
			break
		}
	}
}


// gbShellCmd represents the gbShell command
var gbUnitTestCmd = &cobra.Command{
	Use:   "test <gear name>",
	Short: ux.SprintfBlue("Execute unit tests in Gearbox gear"),
	Long: ux.SprintfBlue("Execute unit tests in Gearbox gear."),
	Example: ux.SprintfWhite("launch unit tests terminus"),
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	Run: gbUnitTestFunc,
	Args: cobra.MinimumNArgs(1),
}
func gbUnitTestFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		var ga GearArgs

		state := ga.ProcessArgs(rootCmd, args)
		if state.IsError() {
			break
		}

		state = ga.gbUnitTestFunc()
		if state.IsError() {
			break
		}
	}
}



func (ga *GearArgs) gbRunFunc() *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range only.Once {
		ga.Quiet = true

		ga.State = ga.gbStartFunc()
		if !ga.State.IsRunning() {
			ga.State.SetError("container not started")
			break
		}

		// Yuck!
		sp := strings.Split(ga.Args[0], ":")
		ga.Args[0] = sp[0]
		ga.Args = ga.GearRef.Docker.Container.GearConfig.GetCommand(ga.Args)
		if len(ga.Args) == 0 {
			ga.State.SetError("ERROR: no default command defined in gearbox.json")
			break
		}

		// Only for the "run" instance - default to mounting CWD.
		if ga.Mount == defaults.DefaultPathNone {
			ga.Mount = DeterminePath(".")
		}

		ga.State = ga.GearRef.Docker.ContainerSsh(false, ga.SshStatus, ga.Mount, ga.Args)
		if !ga.State.IsError() {
			break
		}

		if ga.Temporary {
			ga.State = ga.gbUninstallFunc()
		}
	}

	return ga.State
}

func (ga *GearArgs) gbShellFunc() *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range only.Once {
		ga.State = ga.gbStartFunc()
		if !ga.State.IsRunning() {
			ga.State.SetError("Cannot shell out to Gear '%s:%s.'", ga.Name, ga.Version)
			break
		}

		ga.State = ga.GearRef.Docker.ContainerSsh(true, ga.SshStatus, ga.Mount, ga.Args[1:])

		if ga.Temporary {
			ga.State = ga.gbUninstallFunc()
		}
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}

func (ga *GearArgs) gbUnitTestFunc() *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range only.Once {
		ga.State = ga.gbStartFunc()
		if !ga.State.IsRunning() {
			ga.State.SetError("container not started")
			break
		}

		ga.Args = []string{defaults.DefaultUnitTestCmd}
		ga.State = ga.GearRef.Docker.ContainerSsh(true, ga.SshStatus, ga.Mount, ga.Args)

		if ga.Temporary {
			ga.State = ga.gbUninstallFunc()
		}
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}
