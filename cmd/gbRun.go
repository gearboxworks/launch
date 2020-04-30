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
	gbRunCmd.Flags().BoolP("status", "", false, "Show shell status line.")
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
		var sshStatus bool
		var err error

		//quietFlag, err = cmd.Flags().GetBool(argQuiet)
		//if err != nil {
		//	quietFlag = false
		//}
		quietFlag = true

		gbStartFunc(cmd, args)
		if !cmdState.IsRunning() {
			cmdState.SetError("container not started")
			break
		}

		sshStatus, err = cmd.Flags().GetBool("status")
		if err != nil {
			sshStatus = false
		}

		// Yuck!
		sp := strings.Split(args[0], ":")
		args[0] = sp[0]
		args = gearRef.Docker.Container.GearConfig.GetCommand(args)
		if len(args) == 0 {
			ux.PrintfError("ERROR: no default command defined in gearbox.json")
			break
		}

		cmdState = gearRef.Docker.ContainerSsh(false, sshStatus, args...)
		break
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
		var sshStatus bool
		var err error

		gbStartFunc(cmd, args)
		if !cmdState.IsRunning() {
			cmdState.SetError("container not started")
			break
		}

		sshStatus, err = cmd.Flags().GetBool("status")
		if err != nil {
			sshStatus = false
		}

		cmdState = gearRef.Docker.ContainerSsh(true, sshStatus, args[1:]...)
		break
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
		var sshStatus bool
		var err error

		gbStartFunc(cmd, args)
		if !cmdState.IsRunning() {
			cmdState.SetError("container not started")
			break
		}

		sshStatus, err = cmd.Flags().GetBool("status")
		if err != nil {
			sshStatus = false
		}

		cmdState = gearRef.Docker.ContainerSsh(true, sshStatus, defaults.DefaultUnitTestCmd)
		break
	}
}
