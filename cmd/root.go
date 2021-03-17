/*
terminus == launch --name terminus -- --flag1 arg1 --flag2 arg2 --flag3 --flag4

launch --gb-name terminus --gb-version 2.4.0 site:list --format=json --gb-project foo

launch --gb-name terminus --gb-version 2.4.0 --gb-project foo site:list --format=json

launch --gb-name terminus --gb-version 2.4.0 --gb-project site:list --format=json

launch
	[verb modifiers] uninstall	<gear name> <gear args>
	[verb modifiers] install	<gear name> <gear args>
	[verb modifiers] list		<gear name> <gear args>
	[verb modifiers] import		<gear name> <gear args>
	[verb modifiers] export		<gear name> <gear args>

    [verb modifiers] start		<gear name> <gear args>
	[verb modifiers] stop		<gear name> <gear args>

	[verb modifiers] run		<gear name> <gear args>
	[verb modifiers] shell		<gear name> <gear args>

	[verb modifiers] build		<gear name>
	[verb modifiers] publish	<gear name>
	[verb modifiers] help

[verb modifiers]
	--help
	--examples
	--version
	--provider	- docker, aws, vm, etc, (default:docker local socket)
	--host		- default:localhost
	--port		- default:2356
	--project
*/
package cmd

import (
	"fmt"
	"github.com/newclarity/scribeHelpers/loadTools"
	"github.com/newclarity/scribeHelpers/toolCobraHelp"
	"github.com/newclarity/scribeHelpers/toolGear"
	"github.com/newclarity/scribeHelpers/toolSelfUpdate"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"launch/defaults"
	"os"
	"strings"
)


var Cmd *TypeLaunchArgs

var CobraHelp *toolCobraHelp.TypeCommands

//noinspection ALL
var CmdSelfUpdate *toolSelfUpdate.TypeSelfUpdate
//noinspection ALL
var CmdScribe *loadTools.TypeScribeArgs

var ConfigFile string
const flagConfigFile  	= "config"

const DefaultJsonFile = "gearbox.json"
const DefaultJsonString = "{}"
const DefaultTemplateFile = "gearbox.tmpl"
//const DefaultTemplateString = `
//{{- $gear := NewGear }}
//{{- $gear.ParseGearConfig .Json }}
//{{- $gear.PrintGearConfig }}
//`
const DefaultTemplateString = ""


type redirectHelp struct {
	Text    string
	MovedTo string
}
var courtesyHelp = map[string]redirectHelp {
	// Moved to sub-command "manage"
	"install": { Text: "manage install", MovedTo: "manage"},     // gbManageCmd },
	"uninstall": { Text: "manage uninstall", MovedTo: "manage"}, // gbManageCmd },
	"reinstall": { Text: "manage reinstall", MovedTo: "manage"}, // gbManageCmd },
	"clean": { Text: "manage clean", MovedTo: "manage"},         // gbManageCmd },
	"list": { Text: "manage list", MovedTo: "manage"},           // gbListCmd },
	"start": { Text: "manage start", MovedTo: "manage"},         // gbManageCmd },
	"stop": { Text: "manage stop", MovedTo: "manage"},           // gbManageCmd },

	// Moved to sub-command "create"
	"build": { Text: "create build", MovedTo: "create"},     // gbCreateCmd },
	"publish": { Text: "create publish", MovedTo: "create"}, // gbCreateCmd },
	"save": { Text: "create save", MovedTo: "create"},       // gbCreateCmd },
	"load": { Text: "create load", MovedTo: "create"},       // gbCreateCmd },
}


func init() {
	SetCmd()

	CobraHelp.ChangeHelp(rootCmd, tmplUsage, tmplHelp)

	// Level 1 commands.
	CobraHelp.AddCommands("Manage", rootCmd, gbSearchCmd, gbListCmd, gbInstallCmd, gbUninstallCmd, gbReinstallCmd, gbCleanCmd, gbLogsCmd, gbStartCmd, gbStopCmd)
	CobraHelp.AddCommands("Execute", rootCmd, gbRunCmd, gbShellCmd)
	CobraHelp.AddCommands("Build", rootCmd, gbBuildCmd)
	CobraHelp.AddCommands("Guru", rootCmd, gbCompletionCmd)
	CobraHelp.AddCommands("Help", rootCmd, gbHelpCmd)

	// Level 2 commands.
	CobraHelp.AddCommands("Build", gbBuildCmd, gbBuildCreateCmd, gbUnitTestCmd, gbPublishCmd, gbBuildCleanCmd, gbSaveCmd, gbLoadCmd)
	CobraHelp.AddCommands("List", gbListCmd, gbDetailsCmd, gbLinksCmd, gbPortsCmd)
	CobraHelp.AddCommands("Help", gbHelpCmd, gbHelpFlagsCmd, gbHelpExamplesCmd, gbHelpBasicCmd, gbHelpAdvancedCmd, gbHelpAllCmd)

	CobraHelp.SetLevelDefault(gbRunCmd, gbShellCmd, gbListCmd, gbSearchCmd, CmdSelfUpdate.GetCmd())
	CobraHelp.SetLevelAdvanced(gbInstallCmd, gbUninstallCmd, gbReinstallCmd, gbCleanCmd, gbLogsCmd, gbStartCmd, gbStopCmd)
	CobraHelp.SetLevelAdvanced(gbBuildCmd, gbCompletionCmd, CmdScribe.GetCmd())
	CobraHelp.SetType("Guru", CmdScribe.GetCmd(), CmdSelfUpdate.GetCmd())

	cobra.OnInitialize(initConfig)
	cobra.EnableCommandSorting = false

	rootCmd.PersistentFlags().StringVar(&ConfigFile, flagConfigFile, fmt.Sprintf("%s-config.json", defaults.BinaryName), ux.SprintfBlue("%s: config file.", defaults.BinaryName))
	_ = rootCmd.PersistentFlags().MarkHidden(flagConfigFile)

	rootCmd.Flags().BoolVarP(&Cmd.HelpExamples, flagExample, "", false, ux.SprintfBlue("Help examples for command."))
	rootCmd.Flags().BoolVarP(&Cmd.HelpFlags, flagHelp, "", false, ux.SprintfBlue("Show help on flags."))

	rootCmd.Flags().BoolVarP(&Cmd.NoCreate, flagNoCreate, "n", false, ux.SprintfBlue("Don't create container."))

	rootCmd.Flags().StringVarP(&Cmd.Provider, flagProvider, "", defaults.DefaultProvider, ux.SprintfBlue("Set virtual provider"))
	rootCmd.Flags().StringVarP(&Cmd.Host, flagHost, "", "", ux.SprintfBlue("Set virtual provider host."))
	rootCmd.Flags().StringVarP(&Cmd.Port, flagPort, "", "", ux.SprintfBlue("Set virtual provider port."))
	rootCmd.Flags().StringVarP(&Cmd.Project, flagProject, "p", defaults.DefaultPathNone, ux.SprintfBlue("Mount project directory."))
	rootCmd.Flags().StringVarP(&Cmd.Mount, flagMount, "m", defaults.DefaultPathNone, ux.SprintfBlue("Mount arbitrary directory via SSHFS."))
	rootCmd.Flags().StringVarP(&Cmd.TmpDir, flagTmpDir, "", defaults.DefaultPathNone, ux.SprintfBlue("Alternate TMP dir mount point."))

	rootCmd.Flags().BoolVarP(&Cmd.Temporary, flagTemporary, "t", false, ux.SprintfBlue("Temporary container - remove after running command."))
	rootCmd.Flags().BoolVarP(&Cmd.Status, flagStatus, "s", false, ux.SprintfBlue("Show shell status line."))
	rootCmd.Flags().BoolVarP(&Cmd.Debug, flagDebug, "d", false, ux.SprintfBlue("Debug mode."))
	rootCmd.Flags().BoolVarP(&Cmd.Quiet, flagQuiet, "q", false, ux.SprintfBlue("Silence all launch messages."))

	//rootCmd.Flags().BoolVarP(&Cmd.Completion, flagCompletion, "b", false, ux.SprintfBlue("Generate BASH completion script."))
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error

	for range onlyOnce {
		if Cmd.Config != "" {
			// Use config file from the flag.
			viper.SetConfigFile(Cmd.Config)
		} else {
			// Search config in home directory with name "launch" (without extension).
			viper.AddConfigPath(GetGearboxDir())
			viper.SetConfigName("launch")
		}

		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		err = viper.ReadInConfig()
		if err == nil {
			//ux.Printf("using config file '%s'\n", viper.ConfigFileUsed())
		} else {
			_ = viper.WriteConfig()
		}
	}
}


func SetCmd() {
	for range onlyOnce {
		if Cmd == nil {
			Cmd = New()
		}

		if CobraHelp == nil {
			CobraHelp = toolCobraHelp.New(Cmd.Runtime)
			CobraHelp.SetHelp(rootCmd)
		}

		if CmdScribe == nil {
			CmdScribe = loadTools.New(defaults.BinaryName, defaults.BinaryVersion, false)
			CmdScribe.Runtime.SetRepos(defaults.SourceRepo, defaults.BinaryRepo)
			if CmdScribe.State.IsNotOk() {
				break
			}

			CmdScribe.Json.SetDefaults(DefaultJsonFile, DefaultJsonString)
			CmdScribe.Template.SetDefaults(DefaultTemplateFile, DefaultTemplateString)

			// Import additional tools.
			CmdScribe.ImportTools(&toolGear.GetTools)
			if CmdScribe.State.IsNotOk() {
				break
			}

			CmdScribe.LoadCommands(rootCmd, true)
			if CmdScribe.State.IsNotOk() {
				break
			}

			CmdScribe.AddConfigOption(false, false)
			if CmdScribe.State.IsNotOk() {
				break
			}
		}

		if CmdSelfUpdate == nil {
			CmdSelfUpdate = toolSelfUpdate.New(CmdScribe.Runtime)
			CmdSelfUpdate.LoadCommands(rootCmd, false)
			if CmdSelfUpdate.State.IsNotOk() {
				break
			}
		}
	}
}


var rootCmd = &cobra.Command {
	Use:   defaults.BinaryName,
	Short: ux.SprintfBlue("%s %s launcher", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long: ux.SprintfBlue("%s %s launcher", defaults.LanguageAppName, defaults.LanguageContainerName),
	Run: gbRootFunc,
	TraverseChildren: true,
}


func gbRootFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		if CmdSelfUpdate.FlagCheckVersion(nil) {
			CmdScribe.State.SetOk()
			break
		}

		//fl := cmd.Flags()
		//
		//// Produce BASH completion script.
		//ok, _ := fl.GetBool(flagCompletion)
		//if ok {
		//	var out bytes.Buffer
		//	_ = cmd.GenBashCompletion(&out)
		//	fmt.Printf("# %s BASH completion:\n%s\n", defaults.LanguageAppName, out.String())
		//	Cmd.State.SetOk()
		//	break
		//}
		//
		//// Show flag help.
		//ok, _ := fl.GetBool(flagHelp)
		//if ok {
		//	CobraHelp.ChangeHelp(cmd, tmplFlagUsage, tmplFlagHelp)
		//	_ = cmd.Help()
		//	Cmd.State.SetOk()
		//	break
		//}

		Cmd.State.SetOk()

		if len(args) == 0 {
			_ = cmd.Help()
			break
		}


		cla := ""
		if len(args) > 1 {
			cla = " " + strings.Join(args[1:], " ")
		}

		if _, ok := courtesyHelp[args[0]]; ok {
			ux.PrintflnWarning("Did you mean '%s%s'?", courtesyHelp[args[0]].Text, cla)

			for _, v := range cmd.Commands() {
				if v.Use != courtesyHelp[args[0]].MovedTo {
					continue
				}

				_ = v.Help()
				break
			}

			Cmd.State.SetExitCode(2)
			break
		}

		ux.PrintflnError("Unknown command: '%s'", strings.Join(args, " "))
		Cmd.State.SetError()
		Cmd.State.SetExitCode(1)

		//switch {
		//	case len(args) == 0:
		//		_ = cmd.Help()
		//
		//	case len(args) == 1:
		//		if _, ok := courtesyHelp[args[0]]; ok {
		//			ux.PrintflnWarning("Did you mean '%s'?", courtesyHelp[args[0]].Text)
		//		} else {
		//			ux.PrintflnWarning("Unknown command: '%s'", strings.Join(args, " "))
		//		}
		//
		//	case len(args) > 1:
		//		if _, ok := courtesyHelp[args[0]]; ok {
		//			ux.PrintflnWarning("Did you mean '%s'?", courtesyHelp[args[0]].Text, strings.Join(args[1:], " "))
		//		} else {
		//			ux.PrintflnWarning("Unknown command: '%s'", strings.Join(args, " "))
		//		}
		//}
		//
		//break
	}

	//if Cmd.State.IsNotOk() {
	//	Cmd.State.PrintResponse()
	//}
}


// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() *ux.State {
	Cmd.State = Cmd.State.EnsureNotNil()

	for range onlyOnce {
		// WARNING: Critical code area.
		// Support for running launch via symlink.
		if !Cmd.Runtime.IsRunningAs(defaults.BinaryName) {
			//defaults.RunAs.AsLink = true
			Cmd.Runtime.CmdFile = strings.ReplaceAll(Cmd.Runtime.CmdFile, "-", ":")
			newArgs := []string{"run", Cmd.Runtime.CmdFile}
			newArgs = append(newArgs, os.Args[1:]...)
			rootCmd.SetArgs(newArgs)

			_ = rootCmd.Flags().Set(flagQuiet, "true")
			rootCmd.DisableFlagParsing = true
		}
		// WARNING: Critical code area.

		err := rootCmd.Execute()
		if err != nil {
			Cmd.State.SetError("%s", err)
			break
		}

		Cmd.State = CheckReturns()
	}

	return Cmd.State
}


func CheckReturns() *ux.State {
	state := Cmd.State
	for range onlyOnce {
		if Cmd.State.IsNotOk() {
			state = Cmd.State
			break
		}

		if CmdScribe.State.IsNotOk() {
			state = CmdScribe.State
			break
		}

		if CmdSelfUpdate.State.IsNotOk() {
			state = CmdSelfUpdate.State
			break
		}

		if CobraHelp.State.IsNotOk() {
			state = CobraHelp.State
			break
		}
	}
	return state
}
