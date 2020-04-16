package cmd

import (
	"github.com/spf13/cobra"
	"launch/only"
	"launch/ux"
)

func init() {
	rootCmd.AddCommand(gbRunCmd)
	rootCmd.AddCommand(gbShellCmd)

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

		gbStartFunc(cmd, args)
		if !cmdState.IsRunning() {
			cmdState.SetError("container not started")
			break
		}

		sshStatus, err = cmd.Flags().GetBool("status")
		if err != nil {
			sshStatus = false
		}

		cmdState = gearRef.Docker.ContainerSsh(false, sshStatus, args[1:]...)
		break
	}
}


// gbShellCmd represents the gbShell command
var gbShellCmd = &cobra.Command{
	Use:   "shell <gear name> [gear args]",
	Short: ux.SprintfBlue("Execute shell in Gearbox gear"),
	Long: ux.SprintfBlue("Execute shell in Gearbox gear."),
	Example: ux.SprintfWhite("launch shell golang ps -eaf"),
	DisableFlagParsing: true,
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
