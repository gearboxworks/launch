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
	"bytes"
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
var CmdSelfUpdate *toolSelfUpdate.TypeSelfUpdate
var CmdScribe *loadTools.TypeScribeArgs

var ConfigFile string
const flagConfigFile  	= "config"

const DefaultJsonFile = "gearbox.json"
const DefaultJsonString = "{}"
const DefaultTemplateFile = "gearbox.tmpl"
//const DefaultTemplateString = "{{ $gear := NewGear }}{{ $gear.List }}"
const DefaultTemplateString = `
{{- $gear := NewGear }}
{{- $gear.ParseGearConfig .Json }}
{{- $gear.PrintGearConfig }}`
// {{ $gear := NewGear }}{{ $gear.ParseGearConfig .Json }}{{ $gear.PrintGearConfig }}

func init() {
	SetCmd()

	CobraHelp.AddCommands("Manage", rootCmd, gbInstallCmd, gbUninstallCmd, gbReinstallCmd, gbCleanCmd, gbListCmd)
	CobraHelp.AddCommands("Execute", rootCmd, gbRunCmd, gbShellCmd, gbUnitTestCmd)
	CobraHelp.AddCommands("Run", rootCmd, gbStartCmd, gbStopCmd)
	CobraHelp.AddCommands("Create", rootCmd, gbBuildCmd, gbPublishCmd, gbSaveCmd, gbLoadCmd)

	cobra.OnInitialize(initConfig)
	cobra.EnableCommandSorting = false

	//rootCmd.PersistentFlags().StringVarP(&Cmd.Config, flagConfig, "c", GetLaunchConfig(), ux.SprintfBlue("Config file."))
	rootCmd.PersistentFlags().StringVar(&ConfigFile, flagConfigFile, fmt.Sprintf("%s-config.json", defaults.BinaryName), ux.SprintfBlue("%s: config file.", defaults.BinaryName))
	_ = rootCmd.PersistentFlags().MarkHidden(flagConfigFile)

	rootCmd.Flags().BoolVarP(&Cmd.HelpExamples, flagExample, "e", false, ux.SprintfBlue("Help examples for command."))
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

	rootCmd.Flags().BoolVarP(&Cmd.Completion, flagCompletion, "b", false, ux.SprintfBlue("Generate BASH completion script."))
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

			////CmdScribe.FlagHide(loadTools.FlagScribeFile)
			//CmdScribe.FlagSetDefault(loadTools.FlagJsonFile, "gearbox.json")
			//CmdScribe.FlagSetDefault(loadTools.FlagTemplateFile, DefaultTemplateString)

			CmdScribe.LoadCommands(rootCmd, true)
			if CmdScribe.State.IsNotOk() {
				break
			}

			CmdScribe.AddConfigOption(false, false)
			if CmdScribe.State.IsNotOk() {
				break
			}

			////CmdScribe.FlagHide(loadTools.FlagScribeFile)
			//CmdScribe.FlagSetDefault(loadTools.FlagJsonFile, "gearbox.json")
			//CmdScribe.FlagSetDefault(loadTools.FlagTemplateFile, DefaultTemplateString)

			//CmdScribe.Template.Ignore()
			//if CmdScribe.State.IsNotOk() {
			//	break
			//}
			//
			//CmdScribe.Json.Ignore()
			//if CmdScribe.State.IsNotOk() {
			//	break
			//}
			//
			//CmdScribe.Scribe.Ignore()
			//if CmdScribe.State.IsNotOk() {
			//	break
			//}
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
	Short: ux.SprintfBlue("Gearbox gear launcher"),
	Long: ux.SprintfBlue(`Gearbox gear launcher.`),
	Run: gbRootFunc,
	TraverseChildren: true,
}


func gbRootFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		fl := cmd.Flags()

		if CmdSelfUpdate.FlagCheckVersion(nil) {
			CmdScribe.State.SetOk()
			break
		}

		// Produce BASH completion script.
		ok, _ := fl.GetBool("completion")
		if ok {
			var out bytes.Buffer
			_ = cmd.GenBashCompletion(&out)
			fmt.Printf("# Gearbox BASH completion:\n%s\n", out.String())
			Cmd.State.SetOk()
			break
		}

		// Show help if no commands specified.
		if len(args) == 0 {
			_ = cmd.Help()
			Cmd.State.SetOk()
			break
		}
	}

	if Cmd.State.IsNotOk() {
		Cmd.State.PrintResponse()
	}
}


// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() *ux.State {
	Cmd.State = Cmd.State.EnsureNotNil()

	for range onlyOnce {
		//SetHelp(rootCmd)
		//SetCmd()

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


//func _GetUsage(c *cobra.Command) string {
//	var str string
//
//	if c.Parent() == nil {
//		str += ux.SprintfCyan("%s [flags] ", c.Name())
//	} else {
//		str += ux.SprintfCyan("%s [flags] ", c.Parent().Name())
//		str += ux.SprintfGreen("%s ", c.Use)
//	}
//
//	if c.HasAvailableSubCommands() {
//		str += ux.SprintfGreen("[command] ")
//		str += ux.SprintfCyan("<gear name> ")
//	}
//
//	//foo := c.Use
//	//if c.Args == 1 {
//	//	str += ux.SprintfWhite("")
//	//} else {
//	//	str += ux.SprintfWhite("%v", foo)
//	//}
//
//	// .
//
//	return str
//}
//
//func _GetVersion(c *cobra.Command) string {
//	var str string
//
//	if c.Parent() == nil {
//		str += ux.SprintfWhite("%s: v%s", defaults.BinaryName, defaults.BinaryVersion)
//	}
//
//	return str
//}
//
//func SetHelp(c *cobra.Command) {
//	var tmplHelp string
//	var tmplUsage string
//
//	//fmt.Printf("%s", rootCmd.UsageTemplate())
//	//fmt.Printf("%s", rootCmd.HelpTemplate())
//
//	cobra.AddTemplateFunc("GetUsage", _GetUsage)
//	cobra.AddTemplateFunc("GetVersion", _GetVersion)
//
//	cobra.AddTemplateFunc("SprintfBlue", ux.SprintfBlue)
//	cobra.AddTemplateFunc("SprintfCyan", ux.SprintfCyan)
//	cobra.AddTemplateFunc("SprintfGreen", ux.SprintfGreen)
//	cobra.AddTemplateFunc("SprintfMagenta", ux.SprintfMagenta)
//	cobra.AddTemplateFunc("SprintfRed", ux.SprintfRed)
//	cobra.AddTemplateFunc("SprintfWhite", ux.SprintfWhite)
//	cobra.AddTemplateFunc("SprintfYellow", ux.SprintfYellow)
//
//	// 	{{ with .Parent }}{{ SprintfCyan .Name }}{{ end }} {{ SprintfGreen .Name }} {{ if .HasAvailableSubCommands }}{{ SprintfGreen "[command]" }}{{ end }}
//
//	tmplUsage += `
//{{ SprintfBlue "Usage: " }}
//	{{ GetUsage . }}
//
//{{- if gt (len .Aliases) 0 }}
//{{ SprintfBlue "\nAliases:" }} {{ .NameAndAliases }}
//{{- end }}
//
//{{- if .HasExample }}
//{{ SprintfBlue "\nExamples:" }}
//	{{ .Example }}
//{{- end }}
//
//{{- if .HasAvailableSubCommands }}
//{{ SprintfBlue "\nWhere " }}{{ SprintfGreen "[command]" }}{{ SprintfBlue " is one of:" }}
//{{- range .Commands }}
//{{- if (or .IsAvailableCommand (eq .Name "help")) }}
//	{{ rpad (SprintfGreen .Name) .NamePadding}}	- {{ .Short }}{{ end }}
//{{- end }}
//{{- end }}
//
//{{- if .HasAvailableLocalFlags }}
//{{ SprintfBlue "\nFlags:" }}
//{{ .LocalFlags.FlagUsages | trimTrailingWhitespaces }}
//{{- end }}
//
//{{- if .HasAvailableInheritedFlags }}
//{{ SprintfBlue "\nGlobal Flags:" }}
//{{ .InheritedFlags.FlagUsages | trimTrailingWhitespaces }}
//{{- end }}
//
//{{- if .HasHelpSubCommands }}
//{{- SprintfBlue "\nAdditional help topics:" }}
//{{- range .Commands }}
//{{- if .IsAdditionalHelpTopicCommand }}
//	{{ rpad (SprintfGreen .CommandPath) .CommandPathPadding }} {{ .Short }}
//{{- end }}
//{{- end }}
//{{- end }}
//
//{{- if .HasAvailableSubCommands }}
//{{ SprintfBlue "\nUse" }} {{ SprintfCyan .CommandPath }} {{ SprintfCyan "help" }} {{ SprintfGreen "[command]" }} {{ SprintfBlue "for more information about a command." }}
//{{- end }}
//`
//
//	tmplHelp = `{{ GetVersion . }}
//
//{{ SprintfBlue "Description:" }}
//	{{ SprintfBlue .Use }}{{- SprintfBlue " - " }}
//{{- with (or .Long .Short) }}
//{{- . | trimTrailingWhitespaces }}
//{{- end }}
//
//{{- if or .Runnable .HasSubCommands }}
//{{ .UsageString }}
//{{- end }}
//`
//
//	//c.SetHelpCommand(c)
//	//c.SetHelpFunc(PrintHelp)
//	c.SetHelpTemplate(tmplHelp)
//	c.SetUsageTemplate(tmplUsage)
//}
//
//func PrintHelp(c *cobra.Command, args []string) {
//
//}
//
//func GetState() *ux.State {
//	return Cmd.State
//}
