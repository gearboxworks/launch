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
	"github.com/docker/docker/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"launch/defaults"
	"launch/only"
	"launch/ux"
	"net/url"
	"os"
)

const (
	argConfig = "config"
	argHelp = "help"
	argExample = "example"
	argProvider = "provider"
	argHost = "host"
	argPort = "port"
	argProject = "project"
	argCompletion = "completion"
	argVersion = "version"
)

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, argConfig, "/Users/mick/.gearbox/launch.json", ux.SprintfBlue("Config file."))

	//rootCmd.PersistentFlags().BoolP(argHelp, "h", false, ux.SprintfBlue("Short help for command."))
	rootCmd.PersistentFlags().BoolP(argExample, "e", false, ux.SprintfBlue("Help examples for command."))

	rootCmd.PersistentFlags().StringP(argProvider, "", "docker", ux.SprintfBlue("Set virtual provider"))
	rootCmd.PersistentFlags().StringP(argHost, "", "", ux.SprintfBlue("Set virtual provider host."))
	rootCmd.PersistentFlags().StringP(argPort, "", "", ux.SprintfBlue("Set virtual provider port."))
	rootCmd.PersistentFlags().StringP(argProject, "p", "", ux.SprintfBlue("Mount project directory."))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP(argVersion, "v", false, ux.SprintfBlue("Display version of " + defaults.BinaryName))
	rootCmd.Flags().BoolP(argCompletion, "", false, ux.SprintfBlue("Generate BASH completion script."))
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
				fmt.Println(err)
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
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			_ = viper.WriteConfig()
		}
	}
}

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   defaults.BinaryName,
	Short: ux.SprintfBlue("Gearbox gear launcher"),
	Long: ux.SprintfBlue(`Gearbox gear launcher.`),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	TraverseChildren: true,
}


// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	var err error

	for range only.Once {
		//fmt.Printf("%s", rootCmd.UsageTemplate())
		//fmt.Printf("%s", rootCmd.HelpTemplate())

		SetHelp(rootCmd)

		err = rootCmd.Execute()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var host string
		host, err = rootCmd.Flags().GetString("host")

		var port string
		port, err = rootCmd.Flags().GetString("port")

		if (host != "") && (port != "") {
			var urlString *url.URL
			urlString, err = client.ParseHostURL(fmt.Sprintf("tcp://%s:%s", host, port))
			if err != nil {
				break
			}

			err = os.Setenv("DOCKER_HOST", urlString.String())
			if err != nil {
				break
			}
		}

		var ok bool
		ok, err = rootCmd.Flags().GetBool("completion")
		if ok {
			var out bytes.Buffer
			_ = rootCmd.GenBashCompletion(&out)
			fmt.Printf("# Gearbox BASH completion:\n%s\n", out.String())
			os.Exit(0)
			//cmd.CommandPath()
		}

		ok, err = rootCmd.Flags().GetBool("version")
		if err != nil {
			break
		}
		if ok {
			fmt.Printf("# Gearbox version: %s\n", "1.4.2")
			os.Exit(0)
		}
	}

	return err
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

func SetHelp(c *cobra.Command) {
	var tmplHelp string
	var tmplUsage string

	cobra.AddTemplateFunc("SprintfBlue", _SprintfBlue)
	cobra.AddTemplateFunc("SprintfCyan", _SprintfCyan)
	cobra.AddTemplateFunc("SprintfGreen", _SprintfGreen)
	cobra.AddTemplateFunc("SprintfMagenta", _SprintfMagenta)
	cobra.AddTemplateFunc("SprintfRed", _SprintfRed)
	cobra.AddTemplateFunc("SprintfWhite", _SprintfWhite)
	cobra.AddTemplateFunc("SprintfYellow", _SprintfYellow)


	tmplUsage += `
{{ SprintfBlue "Usage: " }}
{{- if .Runnable }}{{ SprintfCyan .UseLine }}{{ end }}
{{- if .HasAvailableSubCommands }}{{ SprintfCyan .UseLine }} {{ SprintfCyan "[command]" }}{{ end }}
{{ if gt (len .Aliases) 0 }}
{{- SprintfBlue "Aliases:" }} {{ .NameAndAliases }}
{{- end }}
{{ if .HasExample }}
{{- SprintfBlue "Examples:" }}
	{{ .Example }}
{{- end }}
{{ if .HasAvailableSubCommands }}
{{- SprintfBlue "Available Commands:" }}
{{- range .Commands }}
{{- if (or .IsAvailableCommand (eq .Name "help")) }}
	{{ rpad (SprintfCyan .Name) .NamePadding}}	- {{ .Short }}{{ end }}
{{- end }}
{{- end }}
{{ if .HasAvailableLocalFlags }}
{{- SprintfBlue "Flags:" }}
{{ .LocalFlags.FlagUsages | trimTrailingWhitespaces }}
{{- end }}
{{ if .HasAvailableInheritedFlags }}
{{- SprintfBlue "Global Flags:" }}
{{ .InheritedFlags.FlagUsages | trimTrailingWhitespaces }}
{{- end }}
{{ if .HasHelpSubCommands }}
{{- SprintfBlue "Additional help topics:" }}
{{- range .Commands }}
{{- if .IsAdditionalHelpTopicCommand }}
	{{ rpad (SprintfCyan .CommandPath) .CommandPathPadding }} {{ .Short }}
{{- end }}
{{- end }}
{{- end }}

{{- if .HasAvailableSubCommands }}
{{ SprintfBlue "Use" }} {{ SprintfCyan .CommandPath }} {{ SprintfCyan "[command] --help" }} {{ SprintfBlue "for more information about a command." }}
{{- end }}
`

	tmplHelp = `
{{- SprintfBlue .Use }}{{- SprintfBlue " - " }}
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
	//fmt.Printf("%s", tmpl)
}

func PrintHelp(c *cobra.Command, args []string) {

}