package cmd

import (
	"bytes"
	"fmt"
	"github.com/gearboxworks/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
	"strings"
)


var tmplUsageCmd = `
{{ SprintfBlue "Usage: " }}    	{{ GetUsage . }}

{{- if gt (len .Aliases) 0 }}
{{ SprintfBlue "\nAliases:" }}	{{ .NameAndAliases }}
{{- end }}

{{- if .HasExample }}
{{ SprintfBlue "\nExamples:" }}
	{{ .Example }}
{{- end }}

{{- if .HasAvailableSubCommands }}
{{ SprintfBlue "\nWhere " }}{{ SprintfGreen "[command]" }}{{ SprintfBlue " is one of:" }}
{{- range $key, $value := .Commands }}
{{- if (and $value.IsAvailableCommand (ne $value.Name "help")) }}
{{- with $value.Annotations }}
	{{ rpad (SprintfGreen $value.Name) $value.NamePadding}}     	{{ with $value.Annotations.area }}- {{ SprintfMagenta . }}  	{{ end }}- {{ SprintfWhite $value.Short }}{{ end }}
{{- end }}
{{- end }}
{{- end }}

{{ SprintfBlue "For assistance use " }}
{{- SprintfCyan "assist flags" }}{{ SprintfBlue ", " }}
{{- SprintfCyan "assist examples" }}{{ SprintfBlue ", " }}
{{- SprintfCyan "assist basic" }}{{ SprintfBlue ", " }}
{{- SprintfCyan "assist advanced" }}{{ SprintfBlue ", " }}
{{- SprintfCyan "assist all" }}{{ SprintfBlue "." }}
{{- if .HasHelpSubCommands }}
{{ SprintfBlue "Additional help topics:" }}
{{- range .Commands }}
{{- if .IsAdditionalHelpTopicCommand }}
	{{ rpad (SprintfGreen .CommandPath) .CommandPathPadding }} {{ .Short }}
{{- end }}
{{- end }}
{{- end }}
{{- if .HasAvailableSubCommands }}
{{- GetCmdHelp . }}
{{- end }}
`

var tmplHelpCmd = `{{ GetVersion . }}

{{ SprintfBlue "Command:" }}	{{ SprintfCyan .CommandPath }}

{{ SprintfBlue "Description:" }}	{{ with (or .Long .Short) }}
{{- . | trimTrailingWhitespaces }}
{{- end }}

{{- if or .Runnable .HasSubCommands }}
{{ .UsageString }}
{{- end }}
`


var tmplUsageSmall = `{{ SprintfGreen .CommandPath }}
	- {{ GetUsage . }}
	- {{ with (or .Long .Short) }}{{- . | trimTrailingWhitespaces }}{{- end }}

{{- if .HasAvailableSubCommands }}
{{ SprintfBlue "\nWhere " }}{{ SprintfGreen "[command]" }}{{ SprintfBlue " is one of:" }}
{{- range $key, $value := .Commands }}
{{- if (and $value.IsAvailableCommand (ne $value.Name "help")) }}
{{- with $value.Annotations }}
	{{ rpad (SprintfGreen $value.Name) $value.NamePadding}}     	{{ with $value.Annotations.area }}- {{ SprintfMagenta . }}  	{{ end }}- {{ SprintfWhite $value.Short }}{{ end }}
{{- end }}
{{- end }}
{{- end }}

{{- if .HasHelpSubCommands }}
{{ SprintfBlue "\nAdditional help topics:" }}
{{- range .Commands }}
{{- if .IsAdditionalHelpTopicCommand }}
	{{ rpad (SprintfGreen .CommandPath) .CommandPathPadding }} {{ .Short }}
{{- end }}
{{- end }}
{{- end }}
`

var tmplHelpSmall = `{{- if or .Runnable .HasSubCommands }}
{{ .UsageString }}
{{- end }}
`


var tmplFlagUsage = `
{{ SprintfBlue "Usage: " }}
	{{ GetUsage . }}

{{- if .HasAvailableLocalFlags }}
{{ SprintfBlue "\nFlags:" }}
{{ .LocalFlags.FlagUsages | trimTrailingWhitespaces }}
{{- end }}
{{- if .HasAvailableInheritedFlags }}
{{ SprintfBlue "\nGlobal Flags:" }}
{{ .InheritedFlags.FlagUsages | trimTrailingWhitespaces }}
{{- end }}

{{ SprintfBlue "For assistance use " }}
{{- SprintfCyan "assist flags" }}{{ SprintfBlue ", " }}
{{- SprintfCyan "assist examples" }}{{ SprintfBlue ", " }}
{{- SprintfCyan "assist basic" }}{{ SprintfBlue ", " }}
{{- SprintfCyan "assist advanced" }}{{ SprintfBlue ", " }}
{{- SprintfCyan "assist all" }}{{ SprintfBlue "." }}
{{- if .HasAvailableSubCommands }}
{{- GetCmdHelp . }}
{{- end }}
`

var tmplFlagHelp = `{{ GetVersion . }}

{{ SprintfBlue "Command:" }}	{{ SprintfCyan .CommandPath }}

{{ SprintfBlue "Description:" }}	{{ with (or .Long .Short) }}
{{- . | trimTrailingWhitespaces }}
{{- end }}

{{- if or .Runnable .HasSubCommands }}
{{ .UsageString }}
{{- end }}
`

var tmplFlagUsageSmall = `
{{ SprintfBlue "Usage: " }}    	{{ GetUsage . }}

{{- if .HasAvailableLocalFlags }}
{{ SprintfBlue "\nFlags:" }}
{{ .LocalFlags.FlagUsages | trimTrailingWhitespaces }}
{{- end }}
{{- if .HasAvailableInheritedFlags }}
{{ SprintfBlue "\nGlobal Flags:" }}
{{ .InheritedFlags.FlagUsages | trimTrailingWhitespaces }}
{{- end }}

{{- if .HasAvailableSubCommands }}
{{ GetCmdHelp . }}
{{- end }}
`

var tmplFlagHelpSmall = `{{ SprintfWhite "####################################################################" }}
{{- if or .Runnable .HasSubCommands }}
{{ .UsageString }}
{{- end }}
`


// ******************************************************************************** //
var gbHelpCmd = &cobra.Command {
	Use:					"assist",
	//Aliases:				[]string{"flags"},
	Short:					"Show additional help",
	Long:					"Show additional help",
	Example:				ux.SprintfWhite("launch assist"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbHelpFunc,
	Args:					cobra.RangeArgs(0, 1),
}

func gbHelpFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}
		Cmd.SetDebug(ga.Debug)

		if len(args) == 0 {
			CobraHelp.ChangeHelp(rootCmd, tmplUsageCmd, tmplHelpCmd)
			_ = rootCmd.Help()
			ux.PrintfWhite("####################################################################\n")
			_ = cmd.Help()
			break
		}

		switch args[0] {
			case "scribe":
				c := CmdScribe.GetCmd()
				//CobraHelp.ChangeHelp(c, tmplUsage, tmplHelp)
				_ = c.Help()
				CobraHelp.ChangeHelp(c, tmplUsageSmall, tmplHelpSmall)
				DangerousHelpLoop(c, false)

			case "version":
				c := CmdSelfUpdate.GetCmd()
				CobraHelp.ChangeHelp(c, tmplUsage, tmplHelp)
				_ = c.Help()
				CobraHelp.ChangeHelp(c, tmplUsageSmall, tmplHelpSmall)
				DangerousHelpLoop(c, false)

			default:
				CobraHelp.ChangeHelp(rootCmd, tmplUsageCmd, tmplHelpCmd)
				_ = rootCmd.Help()
				ux.PrintfWhite("####################################################################\n")
				_ = cmd.Help()
		}
	}
	Cmd.State.SetOk()
}


// ******************************************************************************** //
var gbHelpAllCmd = &cobra.Command {
	Use:					"all",
	//Aliases:				[]string{"flags"},
	Short:					"Show all help",
	Long:					"Show all help",
	Example:				ux.SprintfWhite("launch assist all"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbHelpAllFunc,
	//Args:					cobra.RangeArgs(0, 2),
}

func gbHelpAllFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}
		Cmd.SetDebug(ga.Debug)

		parent := cmd.Root()

		ux.PrintfWhite("####################################################################\n")
		CobraHelp.ChangeHelp(rootCmd, tmplUsageCmd, tmplHelpCmd)
		_ = parent.Help()
		CobraHelp.ChangeHelp(rootCmd, tmplFlagUsageSmall, tmplFlagHelpSmall)
		_ = parent.Help()

		CobraHelp.ChangeHelp(rootCmd, tmplUsageSmall, tmplHelpSmall)
		for _, c := range parent.Commands() {
			if strings.HasPrefix(c.CommandPath(), Cmd.Runtime.CmdName + " help") {
				continue
			}
			if strings.HasPrefix(c.CommandPath(), Cmd.Runtime.CmdName + " scribe") {
				continue
			}
			DangerousHelpLoop(c, true)
			ux.PrintfWhite("####################################################################\n")
		}

		//gbHelpFlagsFunc(rootCmd, args)
		//gbHelpExamplesFunc(rootCmd, args)
	}
}

func DangerousHelpLoop(cmd *cobra.Command, skip bool) *cobra.Command {
	for range onlyOnce {
		if cmd == nil {
			break
		}

		if skip {
			if strings.HasPrefix(cmd.CommandPath(), Cmd.Runtime.CmdName + " help") {
				break
			}
			if strings.HasPrefix(cmd.CommandPath(), Cmd.Runtime.CmdName + " version") {
				break
			}
			if strings.HasPrefix(cmd.CommandPath(), Cmd.Runtime.CmdName + " scribe") {
				break
			}
		}

		_ = cmd.Help()
		if cmd.HasSubCommands() {
			for _, c := range  cmd.Commands() {
				DangerousHelpLoop(c, skip)
			}
		}
	}
	return cmd
}


// ******************************************************************************** //
var gbHelpBasicCmd = &cobra.Command {
	Use:					"basic",
	//Aliases:				[]string{"flags"},
	Short:					"Show basic help",
	Long:					"Show basic help",
	Example:				ux.SprintfWhite("launch assist basic"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbHelpBasicFunc,
	//Args:					cobra.RangeArgs(0, 2),
}

func gbHelpBasicFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}
		Cmd.SetDebug(ga.Debug)

		parent := cmd.Root()
		for _, v := range parent.Commands() {
			if CobraHelp.IsBasic(v) {
				ux.PrintfWhite("####################################################################\n")
				_ = v.Help()
			}
		}
	}
}


// ******************************************************************************** //
var gbHelpAdvancedCmd = &cobra.Command {
	Use:					"advanced",
	Aliases:				[]string{"guru"},
	Short:					"Show advanced help",
	Long:					"Show advanced help",
	Example:				ux.SprintfWhite("launch assist advanced"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbHelpAdvancedFunc,
	//Args:					cobra.RangeArgs(0, 2),
}

func gbHelpAdvancedFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}
		Cmd.SetDebug(ga.Debug)

		parent := cmd.Root()
		for _, v := range parent.Commands() {
			if CobraHelp.IsAdvanced(v) {
				ux.PrintfWhite("####################################################################\n")
				_ = v.Help()
			}
		}
	}
}


// ******************************************************************************** //
var gbHelpFlagsCmd = &cobra.Command {
	Use:					"flags",
	//Aliases:				[]string{"flags"},
	Short:					"Show additional flags",
	Long:					"Show additional flags",
	Example:				ux.SprintfWhite("launch assist flags"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbHelpFlagsFunc,
	//Args:					cobra.RangeArgs(0, 2),
}

//goland:noinspection GoUnusedParameter
func gbHelpFlagsFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}
		Cmd.SetDebug(ga.Debug)

		CobraHelp.ChangeHelp(rootCmd, tmplFlagUsage, tmplFlagHelp)
		_ = rootCmd.Help()
		Cmd.State.SetOk()
	}
}


// ******************************************************************************** //
var gbHelpExamplesCmd = &cobra.Command {
	Use:					"examples",
	//Aliases:				[]string{"flags"},
	Short:					"Show examples",
	Long:					"Show examples",
	Example:				ux.SprintfWhite("launch assist examples"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbHelpExamplesFunc,
	//Args:					cobra.RangeArgs(0, 2),
}

//goland:noinspection GoUnusedParameter
func gbHelpExamplesFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}
		Cmd.SetDebug(ga.Debug)

		//CobraHelp.ChangeHelp(rootCmd, tmplFlagUsage, tmplFlagHelp)
		//_ = rootCmd.Help()
		ux.PrintfWarning("Command not yet implemented.\n")
		Cmd.State.SetOk()
	}
}


// ******************************************************************************** //
var gbCompletionCmd = &cobra.Command {
	Use:					"completion",
	//Aliases:				[]string{"bash"},
	Short:					"Generate BASH completion file",
	Long:					"Generate BASH completion file",
	Example:				ux.SprintfWhite("launch completion"),
	DisableFlagParsing:		false,
	DisableFlagsInUseLine:	false,
	Run:					gbCompletionFunc,
	//Args:					cobra.RangeArgs(0, 2),
}

func gbCompletionFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}
		Cmd.SetDebug(ga.Debug)

		var out bytes.Buffer
		_ = cmd.GenBashCompletion(&out)
		fmt.Printf("# %s BASH completion:\n%s\n", defaults.LanguageAppName, out.String())
		Cmd.State.SetOk()
	}
}
