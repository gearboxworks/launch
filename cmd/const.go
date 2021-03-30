package cmd

const onlyOnce = "1"
var onlyTwice = []string{"", ""}

var tmplUsage = `
{{ SprintfBlue "Usage: " }}    	{{ GetUsage . }}

{{- if gt (len .Aliases) 0 }}	{{ SprintfBlue "\nAliases:" }} {{ .NameAndAliases }}
{{- end }}

{{- if .HasExample }}
{{ SprintfBlue "\nExamples:" }}
	{{ .Example }}
{{- end }}

{{- if .HasAvailableSubCommands }}
{{ SprintfBlue "\nWhere " }}{{ SprintfGreen "[command]" }}{{ SprintfBlue " is one of:" }}
{{- range $key, $value := .Commands }}
{{- if (or $value.IsAvailableCommand (eq $value.Name "help")) }}
{{- with $value.Annotations }}
{{- if ne .level "advanced" }}
	{{ rpad (SprintfGreen $value.Name) $value.NamePadding}}     	{{ with $value.Annotations.area }}- {{ SprintfMagenta . }}  	{{ end }}- {{ SprintfWhite $value.Short }}{{ end }}
{{- end }}
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

var tmplHelp = `{{ GetVersion . }}

{{ SprintfBlue "Command:" }}    	{{ SprintfCyan .CommandPath }}

{{ SprintfBlue "Description:" }}	{{ with (or .Long .Short) }}
{{- . | trimTrailingWhitespaces }}
{{- end }}

{{- if or .Runnable .HasSubCommands }}
{{ .UsageString }}
{{- end }}
`
