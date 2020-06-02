package cmd

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"launch/defaults"
	"github.com/newclarity/scribeHelpers/ux"
	"strings"
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
				VersionShow()

			case args[0] == "update":
				VersionUpdate(cmd, args)

			default:
				VersionShow()
		}
	}
}


func VersionShow() {
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
		fmt.Printf("Checking for more recent version of %s %s at '%s'\n",
			ux.SprintfBlue(defaults.BinaryName),
			ux.SprintfCyan("v%s", defaults.BinaryVersion),
			ux.SprintfGreen(defaults.BinaryRepo),
		)

		repo := StripUrlPrefix(defaults.BinaryRepo)
		previous := semver.MustParse(defaults.BinaryVersion)

		latest, err := selfupdate.UpdateSelf(previous, repo)
		if err != nil {
			CmdState.SetError(err)
		}

		if previous.Equals(latest.Version) {
			ux.PrintflnOk("%s is up to date: v%s", defaults.BinaryName, defaults.BinaryVersion)
		} else {
			ux.PrintflnOk("%s updated to v%s", defaults.BinaryName, latest.Version)
			if latest.ReleaseNotes != "" {
				ux.PrintflnOk("%s %s Release Notes:\n%s", defaults.BinaryName, latest.Version, latest.ReleaseNotes)
			}
		}
	}
}

func StripUrlPrefix(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "github.com/")
	url = strings.TrimSuffix(url, "/")
	url = strings.TrimSpace(url)

	return url
}
