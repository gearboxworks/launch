package cmd

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"launch/defaults"
	"launch/ux"
)


func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(selfUpdateCmd)
}


var versionCmd = &cobra.Command{
	Use:   "version",
	Short: ux.SprintfBlue("Show version of %s.", defaults.BinaryName),
	Long:  ux.SprintfBlue("Show version of %s.", defaults.BinaryName),
	Run:   Version,
}
func Version(cmd *cobra.Command, args []string) {
	for range OnlyOnce {
		// Add sub-commands.
		//cmd.AddCommand(versionUpdateCmd)
		//cmd.AddCommand(versionPublishCmd)
		//cmd.AddCommand(versionShowCmd)

		switch {
			case len(args) == 0:
				VersionShow(cmd, args)

			case args[0] == "update":
				VersionUpdate(cmd, args)

			case args[0] == "release":
				VersionPublish(cmd, args)

			default:
				VersionShow(cmd, args)
		}
	}
}


var versionShowCmd = &cobra.Command{
	Use:   "show",
	Short: ux.SprintfBlue("Show version of %s.", defaults.BinaryName),
	Long:  ux.SprintfBlue("Show version of %s.", defaults.BinaryName),
	Run:   VersionShow,
}
func VersionShow(cmd *cobra.Command, args []string) {
	fmt.Printf("%s %s\n",
		ux.SprintfBlue(defaults.BinaryName),
		ux.SprintfCyan("v%s", defaults.BinaryVersion),
	)
}


var selfUpdateCmd = &cobra.Command{
	Use:   "selfupdate",
	Short: ux.SprintfBlue("Update version of %s.", defaults.BinaryName),
	Long:  ux.SprintfBlue("Check and update the latest version of %s.", defaults.BinaryName),
	Run:   SelfUpdate,
}
func SelfUpdate(cmd *cobra.Command, args []string) {
	VersionUpdate(cmd, args)
}


var versionUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: ux.SprintfBlue("Update version of %s.", defaults.BinaryName),
	Long:  ux.SprintfBlue("Check and update the latest version of %s.", defaults.BinaryName),
	Run:   VersionUpdate,
}
func VersionUpdate(cmd *cobra.Command, args []string) {
	for range OnlyOnce {
		fmt.Printf("%s %s\n",
			ux.SprintfBlue(defaults.BinaryName),
			ux.SprintfCyan("v%s", defaults.BinaryVersion),
		)

		err := selfUpdate(defaults.BinaryRepo)
		if err != nil {
			_cmdState.SetError(err)
		}
	}
}


func selfUpdate(slug string) error {
	//selfupdate.EnableLog()

	ux.PrintflnOk("Checking for more recent version: v%s", defaults.BinaryVersion)
	previous := semver.MustParse(defaults.BinaryVersion)
	latest, err := selfupdate.UpdateSelf(previous, slug)
	if err != nil {
		return err
	}

	if previous.Equals(latest.Version) {
		ux.PrintflnOk("%s is up to date: v%s", defaults.BinaryName, defaults.BinaryVersion)
	} else {
		ux.PrintflnOk("%s updated to v%s", defaults.BinaryName, latest.Version)
		if latest.ReleaseNotes != "" {
			ux.PrintflnOk("%s %s Release Notes:\n%s", defaults.BinaryName, latest.Version, latest.ReleaseNotes)
		}
	}

	return nil
}


var versionPublishCmd = &cobra.Command{
	Use:   "release",
	Short: ux.SprintfBlue("Publish release version of %s.", defaults.BinaryName),
	Long:  ux.SprintfBlue("Publish release version of %s.", defaults.BinaryName),
	Run:   VersionPublish,
}
func VersionPublish(cmd *cobra.Command, args []string) {
	for range OnlyOnce {
		fmt.Printf("%s %s\n",
			ux.SprintfBlue(defaults.BinaryName),
			ux.SprintfCyan("v%s", defaults.BinaryVersion),
		)


		// Use defaults.SourceRepo
		//args = []string{
		//	"-u", "gearboxworks",
		//	"-r", "jtc",
		//	defaults.BinaryVersion,
		//}
		//cli := &gitRelease.CLI{OutStream: os.Stdout, ErrStream: os.Stderr}
		//cli.Run(args)


	}
}
