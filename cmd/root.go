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
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"launch/defaults"
	"launch/gear"
	"launch/only"
	"launch/ospaths"
	"launch/ux"
	"os"
	"regexp"
	"strings"
)

const (
	argConfig = "config"
	argHelp = "help"
	argDebug = "debug"
	argExample = "example"
	argProvider = "provider"
	argHost = "host"
	argPort = "port"
	argProject = "project"
	argCompletion = "completion"
	argVersion = "version"
	argQuiet = "quiet"
)

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, argConfig, "/Users/mick/.gearbox/launch.json", ux.SprintfBlue("Config file."))

	//rootCmd.PersistentFlags().BoolP(argHelp, "h", false, ux.SprintfBlue("Short help for command."))
	rootCmd.PersistentFlags().BoolP(argExample, "e", false, ux.SprintfBlue("Help examples for command."))
	rootCmd.PersistentFlags().BoolP(argDebug, "d", false, ux.SprintfBlue("Debug mode."))

	rootCmd.PersistentFlags().StringP(argProvider, "", "docker", ux.SprintfBlue("Set virtual provider"))
	rootCmd.PersistentFlags().StringP(argHost, "", "", ux.SprintfBlue("Set virtual provider host."))
	rootCmd.PersistentFlags().StringP(argPort, "", "", ux.SprintfBlue("Set virtual provider port."))
	rootCmd.PersistentFlags().StringP(argProject, "p", "", ux.SprintfBlue("Mount project directory."))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP(argVersion, "v", false, ux.SprintfBlue("Display version of " + defaults.BinaryName))
	rootCmd.Flags().BoolP(argCompletion, "", false, ux.SprintfBlue("Generate BASH completion script."))
	rootCmd.Flags().BoolP(argQuiet, "", false, ux.SprintfBlue("Make everything quiet."))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error

	for range only.Once {
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			var home string
			home, err = homedir.Dir()
			if err != nil {
				ux.PrintError(err)
				os.Exit(1)
			}

			// Search config in home directory with name "launch" (without extension).
			viper.AddConfigPath(home + "/.gearbox")
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

var cmdExec string
var cmdState ux.State
var debugFlag bool
var quietFlag bool
var provider gear.Provider
var gearRef *gear.Gear
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command {
	Use:   defaults.BinaryName,
	Short: ux.SprintfBlue("Gearbox gear launcher"),
	Long: ux.SprintfBlue(`Gearbox gear launcher.`),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	Run: gbRootFunc,
	TraverseChildren: true,
}

func gbRootFunc(cmd *cobra.Command, args []string) {
	for range only.Once {
		var err error
		fl := cmd.Flags()
		//ux.Printf("F5: %v\n", fl.Args())
		//ux.Printf("F6: %v\n", args)

		//quietFlag, _ = fl.GetBool(argQuiet)

		debugFlag, _ = fl.GetBool(argDebug)
		if debugFlag {
			showArgs(cmd, args)
			flargs := fl.Args()
			ux.Printf("flargs: %s\n", strings.Join(flargs, " "))
			ux.Printf("args: %s\n", strings.Join(args, " "))
		}


		// Produce BASH completion script.
		var ok bool
		ok, err = fl.GetBool("completion")
		if ok {
			var out bytes.Buffer
			_ = cmd.GenBashCompletion(&out)
			fmt.Printf("# Gearbox BASH completion:\n%s\n", out.String())
			cmdState.ClearAll()
			break
			//os.Exit(0)
		}


		// Show version.
		ok, err = fl.GetBool("version")
		if err != nil {
			cmdState.SetError("%s", err)
			break
		}
		if ok {
			ux.Printf("version: %s\n", "1.4.2")
			cmdState.ClearAll()
			break
			//os.Exit(0)
		}


		// Create new provider connection.
		provider.Debug = debugFlag
		provider.Name, _ = fl.GetString(argProvider)
		provider.Host, _ = fl.GetString(argHost)
		provider.Port, _ = fl.GetString(argPort)
		provider.Project, _ = fl.GetString(argProject)
		cmdState = provider.NewProvider()


		// Show help if no commands specified.
		if len(args) == 0 {
			_ = cmd.Help()
			cmdState.ClearAll()
			break
			//os.Exit(0)
		}
	}
}


// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() ux.State {
	var state ux.State

	for range only.Once {
		SetHelp(rootCmd)

		foo := ospaths.Split(os.Args[0])
		cmdExec = foo.File.String()
		ok, _ := regexp.MatchString("^" + defaults.BinaryName, cmdExec)
		if ok {
			cmdExec = ""
		} else {
			newArgs := []string{"run", cmdExec}
			newArgs = append(newArgs, os.Args[1:]...)
			rootCmd.SetArgs(newArgs)

			_ = rootCmd.Flags().Set(argQuiet, "true")
			quietFlag = true
			rootCmd.DisableFlagParsing = true
		}

		err := rootCmd.Execute()
		if err != nil {
			fmt.Printf("F4: %v\n", err)
			cmdState.SetError("%s", err)
			break
		}

		if cmdExec == "" {
			break
		}
	}

	return state
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

func SetHelp(c *cobra.Command) {
	var tmplHelp string
	var tmplUsage string

	//fmt.Printf("%s", rootCmd.UsageTemplate())
	//fmt.Printf("%s", rootCmd.HelpTemplate())

	cobra.AddTemplateFunc("GetUsage", _GetUsage)

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
{{ SprintfBlue "\nUse" }} {{ SprintfCyan .CommandPath }} {{ SprintfGreen "[command]" }} {{ SprintfCyan "--help" }} {{ SprintfBlue "for more information about a command." }}
{{- end }}
`

	tmplHelp = `
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

func PrintHelp(c *cobra.Command, args []string) {

}

func GetState() ux.State {
	return cmdState
}
