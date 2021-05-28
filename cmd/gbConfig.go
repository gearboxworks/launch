package cmd

import (
	"github.com/gearboxworks/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
)

// ******************************************************************************** //
var gbConfigCmd = &cobra.Command{
	Use:                   "config",
	Short:                 ux.SprintfBlue("Manage config file for %s", defaults.BinaryName),
	Long:                  ux.SprintfBlue("Manage config file for %s", defaults.BinaryName),
	Example:               ux.SprintfWhite("launch config"),
	DisableFlagParsing:    false,
	DisableFlagsInUseLine: false,
	Run:                   gbConfigFunc,
	Args:                  cobra.RangeArgs(0, 2),
}

//gbConfigFunc takes a pointer to cobra.command and
//command arguments to config command state
func gbConfigFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, false)
		if Cmd.State.IsNotOk() {
			break
		}

		//for _, c := range cmd.Commands() {
		//	for k, v := range c.Annotations {
		//		fmt.Printf(`%s - "%s": "%s"`,
		//			c.Use,
		//			k,
		//			v,
		//		)
		//		fmt.Printf("\n")
		//	}
		//}

		switch {
		case len(args) == 0:
			_ = cmd.Help()
		}
	}
}

// ******************************************************************************** //
var gbConfigTimeoutCmd = &cobra.Command{
	Use:                   "timeout",
	Short:                 "Determine the minimum provider timeout",
	Long:                  "Determine the minimum provider timeout",
	Example:               ux.SprintfWhite("launch config timeout"),
	DisableFlagParsing:    false,
	DisableFlagsInUseLine: false,
	Run:                   gbConfigTimeoutFunc,
	//Args:					cobra.RangeArgs(0, 2),
}

//goland:noinspection GoUnusedParameter
//gbConfigTimeoutFunc takes a pointer to cobra.command and
//command arguments to config timeout command state
func gbConfigTimeoutFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var err error
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, false)
		//if Cmd.State.IsNotOk() {
		//	break
		//}

		ux.PrintfOk("Determining the minimum provider timeout")
		Cmd.Timeout = ga.Gears.DetermineTimeout(defaults.DefaultMinTimeout, defaults.DefaultMaxTimeout, 150)
		ga.Timeout = Cmd.Timeout
		Cmd.Runtime.Timeout = Cmd.Timeout

		if Cmd.Timeout == 0 {
			ux.PrintflnError("Cannot find timeout value.")
			break
		}

		rootViper.Set(flagTimeout, Cmd.Timeout.String())

		err = rootViper.WriteConfig()
		if err != nil {
			Cmd.State.SetError(err)
			break
		}

		ux.PrintflnOk("Minimum provider timeout set to %s", Cmd.Timeout.String())

		Cmd.State.SetOk()
	}
}
