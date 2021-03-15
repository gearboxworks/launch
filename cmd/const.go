package cmd

const onlyOnce = "1"
var onlyTwice = []string{"", ""}


var tmplUsage = `
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
	{{ rpad (SprintfGreen .Name) .NamePadding}}     	- {{ .Short }}{{ end }}
{{- end }}
{{- end }}


{{- if .HasAvailableLocalFlags }}
{{ SprintfBlue "\nType '--flags' to show global flag options" }}
{{- end }}

{{- if .HasAvailableInheritedFlags }}
{{ SprintfBlue "\nType '--flags' to show global flag options" }}
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
{{ GetCmdHelp . }}
{{- end }}
`

var tmplHelp = `{{ GetVersion . }}

{{ SprintfBlue "Commmand:" }} {{ SprintfCyan .Use }}

{{ SprintfBlue "Description:" }} 
	{{ with (or .Long .Short) }}
{{- . | trimTrailingWhitespaces }}
{{- end }}

{{- if or .Runnable .HasSubCommands }}
{{ .UsageString }}
{{- end }}
`

var tmplFlagUsage = `
{{- if .HasAvailableLocalFlags }}
{{ SprintfBlue "\nFlags:" }}
{{ .LocalFlags.FlagUsages | trimTrailingWhitespaces }}
{{- end }}
{{- if .HasAvailableInheritedFlags }}
{{ SprintfBlue "\nGlobal Flags:" }}
{{ .InheritedFlags.FlagUsages | trimTrailingWhitespaces }}
{{- end }}
`

var tmplFlagHelp = `{{ GetVersion . }}

{{ SprintfBlue "Commmand:" }} {{ SprintfCyan .Use }}
{{- if or .Runnable .HasSubCommands }}
{{ .UsageString }}
{{- end }}
`
