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
	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"launch/defaults"
	"launch/ux"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	argConfig = "config"
	//argHelp = "help"
	argDebug = "debug"
	argNoCreate = "no-create"
	argExample = "example"

	argProvider = "provider"
	argProviderDefault = "docker"

	argHost = "host"
	argPort = "port"
	argProject = "project"
	argMount = "mount"
	argCompletion = "completion"
	argVersion = "version"
	argQuiet = "quiet"
	argTemporary = "temporary"
	argStatus = "status"
)

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, argConfig, GetLaunchConfig(), ux.SprintfBlue("Config file."))

	rootCmd.PersistentFlags().BoolP(argExample, "e", false, ux.SprintfBlue("Help examples for command."))
	rootCmd.PersistentFlags().BoolP(argNoCreate, "n", false, ux.SprintfBlue("Don't create container."))

	rootCmd.PersistentFlags().StringP(argProvider, "", defaults.DefaultProvider, ux.SprintfBlue("Set virtual provider"))
	rootCmd.PersistentFlags().StringP(argHost, "", "", ux.SprintfBlue("Set virtual provider host."))
	rootCmd.PersistentFlags().StringP(argPort, "", "", ux.SprintfBlue("Set virtual provider port."))
	rootCmd.PersistentFlags().StringP(argProject, "p", defaults.DefaultPathNone, ux.SprintfBlue("Mount project directory."))
	rootCmd.PersistentFlags().StringP(argMount, "m", defaults.DefaultPathNone, ux.SprintfBlue("Mount arbitrary directory via SSHFS."))

	rootCmd.Flags().BoolP(argTemporary, "t", false, ux.SprintfBlue("Temporary container - remove after running command."))
	rootCmd.Flags().BoolP(argStatus, "s", false, ux.SprintfBlue("Show shell status line."))
	rootCmd.Flags().BoolP(argDebug, "d", false, ux.SprintfBlue("Debug mode."))
	rootCmd.Flags().BoolP(argQuiet, "q", false, ux.SprintfBlue("Silence all Gearbox messsages."))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP(argVersion, "v", false, ux.SprintfBlue("Display version of " + defaults.BinaryName))
	rootCmd.Flags().BoolP(argCompletion, "b", false, ux.SprintfBlue("Generate BASH completion script."))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error

	for range OnlyOnce {
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
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

var _cmdState *ux.State
var cfgFile string


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command {
	Use:   defaults.BinaryName,
	Short: ux.SprintfBlue("Gearbox gear launcher"),
	Long: ux.SprintfBlue(`Gearbox gear launcher.`),
	Run: gbRootFunc,
	TraverseChildren: true,
	//ValidArgs: []string{"run", "shell", "test"},
}


func gbRootFunc(cmd *cobra.Command, args []string) {
	_cmdState = ux.NewState(false)

	for range OnlyOnce {
		var err error
		fl := cmd.Flags()
		//ux.Printf("F5: %v\n", fl.Args())
		//ux.Printf("F6: %v\n", args)

		//quietFlag, _ = fl.GetBool(argQuiet)

		var debugFlag bool
		debugFlag, _ = fl.GetBool(argDebug)
		_cmdState.DebugSet(debugFlag)
		if debugFlag {
			showArgs(cmd, args)
			flargs := fl.Args()
			ux.Printf("flargs: %s\n", strings.Join(flargs, " "))
			ux.Printf("args: %s\n", strings.Join(args, " "))
		}

		//tempFlag, _ = cmd.Flags().GetBool(argTemporary)

		// Produce BASH completion script.
		var ok bool
		ok, err = fl.GetBool("completion")
		if ok {
			var out bytes.Buffer
			_ = cmd.GenBashCompletion(&out)
			fmt.Printf("# Gearbox BASH completion:\n%s\n", out.String())
			_cmdState.Clear()
			break
			//os.Exit(0)
		}


		// Show version.
		ok, err = fl.GetBool("version")
		if err != nil {
			_cmdState.SetError("%s", err)
			break
		}
		if ok {
			ux.Printf("%s: v%s\n", defaults.BinaryName, defaults.BinaryVersion)
			_cmdState.Clear()
			break
			//os.Exit(0)
		}


		// Create new provider connection.
		//provider.Debug = debugFlag
		//provider.Name, _ = fl.GetString(argProvider)
		//provider.Host, _ = fl.GetString(argHost)
		//provider.Port, _ = fl.GetString(argPort)
		//provider.Project, _ = fl.GetString(argProject)
		//state = provider.NewProvider()


		// Show help if no commands specified.
		if len(args) == 0 {
			_ = cmd.Help()
			_cmdState.Clear()
			break
			//os.Exit(0)
		}
	}

	if _cmdState.IsNotOk() {
		_cmdState.PrintResponse()
	}
}


// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() *ux.State {
	_cmdState = _cmdState.EnsureNotNil()

	for range OnlyOnce {
		var err error

		SetHelp(rootCmd)

		// WARNING: Critical code area.
		// Support for running launch via symlink.
		//defaults.RunAs.FullPath, err = filepath.Abs(os.Args[0])
		defaults.RunAs.FullPath, err = osext.Executable()
		if err != nil {
			_cmdState.SetError("%s", err)
			break
		}
		//defaults.RunAs.FullPath = "/Users/mick/Documents/GitHub/gb-launch/bin/psql-9.4.26"
		//defaults.RunAs.FullPath = "/Users/mick/Documents/GitHub/gb-launch/bin/postgresql-9.4.26"

		defaults.RunAs.Dir, defaults.RunAs.File = filepath.Split(defaults.RunAs.FullPath)

		ok, _ := regexp.MatchString("^" + defaults.BinaryName, defaults.RunAs.File)
		if !ok {
			defaults.RunAs.AsLink = true
			defaults.RunAs.File = strings.ReplaceAll(defaults.RunAs.File, "-", ":")
			newArgs := []string{"run", defaults.RunAs.File}
			newArgs = append(newArgs, os.Args[1:]...)
			rootCmd.SetArgs(newArgs)

			_ = rootCmd.Flags().Set(argQuiet, "true")
			rootCmd.DisableFlagParsing = true
		}
		// WARNING: Critical code area.

		err = rootCmd.Execute()
		if err != nil {
			_cmdState.SetError("%s", err)
			break
		}
	}

	return _cmdState
}


func _SprintfBlue(c string) string {
	return ux.SprintfBlue(c)
}

func _SprintfCyan(c string) string {
	return ux.SprintfCyan(c)
}

func _SprintfWhite(c string) string {
	return ux.SprintfWhite(c)
}

func _SprintfGreen(c string) string {
	return ux.SprintfGreen(c)
}

func _SprintfMagenta(c string) string {
	return ux.SprintfMagenta(c)
}

func _SprintfRed(c string) string {
	return ux.SprintfRed(c)
}

func _SprintfYellow(c string) string {
	return ux.SprintfYellow(c)
}

func _GetUsage(c *cobra.Command) string {
	var str string

	if c.Parent() == nil {
		str += ux.SprintfCyan("%s [flags] ", c.Name())
	} else {
		str += ux.SprintfCyan("%s [flags] ", c.Parent().Name())
		str += ux.SprintfGreen("%s ", c.Use)
	}

	if c.HasAvailableSubCommands() {
		str += ux.SprintfGreen("[command] ")
		str += ux.SprintfCyan("<gear name> ")
	}

	//foo := c.Use
	//if c.Args == 1 {
	//	str += ux.SprintfWhite("")
	//} else {
	//	str += ux.SprintfWhite("%v", foo)
	//}

	// .

	return str
}

func _GetVersion(c *cobra.Command) string {
	var str string

	if c.Parent() == nil {
		str += ux.SprintfWhite("%s: v%s", defaults.BinaryName, defaults.BinaryVersion)
	}

	return str
}

func SetHelp(c *cobra.Command) {
	var tmplHelp string
	var tmplUsage string

	//fmt.Printf("%s", rootCmd.UsageTemplate())
	//fmt.Printf("%s", rootCmd.HelpTemplate())

	cobra.AddTemplateFunc("GetUsage", _GetUsage)
	cobra.AddTemplateFunc("GetVersion", _GetVersion)

	cobra.AddTemplateFunc("SprintfBlue", _SprintfBlue)
	cobra.AddTemplateFunc("SprintfCyan", _SprintfCyan)
	cobra.AddTemplateFunc("SprintfGreen", _SprintfGreen)
	cobra.AddTemplateFunc("SprintfMagenta", _SprintfMagenta)
	cobra.AddTemplateFunc("SprintfRed", _SprintfRed)
	cobra.AddTemplateFunc("SprintfWhite", _SprintfWhite)
	cobra.AddTemplateFunc("SprintfYellow", _SprintfYellow)

	// 	{{ with .Parent }}{{ SprintfCyan .Name }}{{ end }} {{ SprintfGreen .Name }} {{ if .HasAvailableSubCommands }}{{ SprintfGreen "[command]" }}{{ end }}

	tmplUsage += `
{{ SprintfBlue "Usage: " }}
	{{ GetUsage . }}

{{- if gt (len .Aliases) 0 }}
{{ SprintfBlue "\nAliases:" }} {{ .NameAndAliases }}
{{- end }}

{{- if .HasExample }}
{{ SprintfBlue "\nExamples:" }}
	{{ .Example }}
{{- end }}

{{- if .HasAvailableSubCommands }}
{{ SprintfBlue "\nWhere " }}{{ SprintfGreen "[command]" }}{{ SprintfBlue " is one of:" }}
{{- range .Commands }}
{{- if (or .IsAvailableCommand (eq .Name "help")) }}
	{{ rpad (SprintfGreen .Name) .NamePadding}}	- {{ .Short }}{{ end }}
{{- end }}
{{- end }}

{{- if .HasAvailableLocalFlags }}
{{ SprintfBlue "\nFlags:" }}
{{ .LocalFlags.FlagUsages | trimTrailingWhitespaces }}
{{- end }}

{{- if .HasAvailableInheritedFlags }}
{{ SprintfBlue "\nGlobal Flags:" }}
{{ .InheritedFlags.FlagUsages | trimTrailingWhitespaces }}
{{- end }}

{{- if .HasHelpSubCommands }}
{{- SprintfBlue "\nAdditional help topics:" }}
{{- range .Commands }}
{{- if .IsAdditionalHelpTopicCommand }}
	{{ rpad (SprintfGreen .CommandPath) .CommandPathPadding }} {{ .Short }}
{{- end }}
{{- end }}
{{- end }}

{{- if .HasAvailableSubCommands }}
{{ SprintfBlue "\nUse" }} {{ SprintfCyan .CommandPath }} {{ SprintfCyan "help" }} {{ SprintfGreen "[command]" }} {{ SprintfBlue "for more information about a command." }}
{{- end }}
`

	tmplHelp = `{{ GetVersion . }}

{{ SprintfBlue "Description:" }} 
	{{ SprintfBlue .Use }}{{- SprintfBlue " - " }}
{{- with (or .Long .Short) }}
{{- . | trimTrailingWhitespaces }}
{{- end }}

{{- if or .Runnable .HasSubCommands }}
{{ .UsageString }}
{{- end }}
`

	//c.SetHelpCommand(c)
	//c.SetHelpFunc(PrintHelp)
	c.SetHelpTemplate(tmplHelp)
	c.SetUsageTemplate(tmplUsage)
}

//func PrintHelp(c *cobra.Command, args []string) {
//
//}

func GetState() *ux.State {
	return _cmdState
}
