package defaults


//const onlyOnce = "1"
////var onlyTwice = []string{"", ""}
//var rootCmd *cobra.Command
//var CmdScribe *loadTools.TypeScribeArgs
//
//
//func New(cmd *cobra.Command, scribe *loadTools.TypeScribeArgs) {
//	rootCmd = cmd
//	CmdScribe = scribe
//
//	rootCmd.AddCommand(versionCmd)
//	rootCmd.AddCommand(selfUpdateCmd)
//
//	versionCmd.AddCommand(versionCheckCmd)
//	versionCmd.AddCommand(versionListCmd)
//	versionCmd.AddCommand(versionInfoCmd)
//	versionCmd.AddCommand(versionLatestCmd)
//	versionCmd.AddCommand(versionUpdateCmd)
//}
//
//
//var versionCmd = &cobra.Command{
//	Use:   loadTools.CmdVersion,
//	Short: ux.SprintfMagenta(BinaryName) + ux.SprintfBlue(" - Self-manage executable."),
//	Long:  ux.SprintfMagenta(BinaryName) + ux.SprintfBlue(" - Self-manage executable."),
//	DisableFlagParsing: true,
//	DisableFlagsInUseLine: true,
//	Run: func(cmd *cobra.Command, args []string) {
//		CmdScribe.State = Version(cmd, args...)
//	},
//}
//var selfUpdateCmd = &cobra.Command{
//	Use:   loadTools.CmdSelfUpdate,
//	Short: ux.SprintfMagenta(BinaryName) + ux.SprintfBlue(" - Update version of executable."),
//	Long: ux.SprintfMagenta(BinaryName) + ux.SprintfBlue(" - Check and update the latest version."),
//	DisableFlagParsing: true,
//	DisableFlagsInUseLine: true,
//	Run: func(cmd *cobra.Command, args []string) {
//		CmdScribe.State = VersionUpdate()
//	},
//}
//
//
//var versionCheckCmd = &cobra.Command{
//	Use:   loadTools.CmdVersionCheck,
//	Short: ux.SprintfMagenta(loadTools.CmdVersion) + ux.SprintfBlue(" - Show any version updates."),
//	Long:  ux.SprintfMagenta(loadTools.CmdVersion) + ux.SprintfBlue(" - Show any version updates."),
//	DisableFlagParsing: true,
//	DisableFlagsInUseLine: true,
//	Run: func(cmd *cobra.Command, args []string) {
//		CmdScribe.State = VersionCheck()
//	},
//}
//var versionListCmd = &cobra.Command{
//	Use:   loadTools.CmdVersionList,
//	Short: ux.SprintfMagenta(loadTools.CmdVersion) + ux.SprintfBlue(" - Self-manage executable."),
//	Long:  ux.SprintfMagenta(loadTools.CmdVersion) + ux.SprintfBlue(" - Self-manage executable."),
//	DisableFlagParsing: true,
//	DisableFlagsInUseLine: true,
//	Run: func(cmd *cobra.Command, args []string) {
//		CmdScribe.State = VersionList(args...)
//	},
//}
//var versionInfoCmd = &cobra.Command{
//	Use:   loadTools.CmdVersionInfo,
//	Short: ux.SprintfMagenta(loadTools.CmdVersion) + ux.SprintfBlue(" - Info on current version."),
//	Long:  ux.SprintfMagenta(loadTools.CmdVersion) + ux.SprintfBlue(" - Info on current version."),
//	DisableFlagParsing: true,
//	DisableFlagsInUseLine: true,
//	Run: func(cmd *cobra.Command, args []string) {
//		//CmdScribe.State = VersionInfo(CmdScribe.Runtime.GetSemVer())
//		CmdScribe.State = VersionInfo(args...)
//	},
//}
//var versionLatestCmd = &cobra.Command{
//	Use:   loadTools.CmdVersionLatest,
//	Short: ux.SprintfMagenta(loadTools.CmdVersion) + ux.SprintfBlue(" - Info on latest version."),
//	Long:  ux.SprintfMagenta(loadTools.CmdVersion) + ux.SprintfBlue(" - Info on latest version."),
//	DisableFlagParsing: true,
//	DisableFlagsInUseLine: true,
//	Run: func(cmd *cobra.Command, args []string) {
//		CmdScribe.State = VersionInfo(loadTools.CmdVersionLatest)
//	},
//}
//var versionUpdateCmd = &cobra.Command{
//	Use:   loadTools.CmdVersionUpdate,
//	Short: ux.SprintfMagenta(loadTools.CmdVersion) + ux.SprintfBlue(" - Update version of executable."),
//	Long: ux.SprintfMagenta(loadTools.CmdVersion) + ux.SprintfBlue(" - Check and update the latest version."),
//	DisableFlagParsing: true,
//	DisableFlagsInUseLine: true,
//	Run: func(cmd *cobra.Command, args []string) {
//		CmdScribe.State = VersionUpdate()
//	},
//}
//
//
//func Version(cmd *cobra.Command, args ...string) *ux.State {
//	VersionShow()
//	SetHelp(cmd)
//	_ = cmd.Help()
//	CmdScribe.State.Clear()
//	CmdScribe.State.Clear()
//	return CmdScribe.State
//}
//
//
//func VersionShow() *ux.State {
//	CmdScribe.Runtime.PrintNameVersion()
//	CmdScribe.State.Clear()
//	return CmdScribe.State
//}
//
//
////func VersionInfo(v *toolRuntime.VersionValue) *ux.State {
//func VersionInfo(args ...string) *ux.State {
//	state := CmdScribe.State
//
//	for range onlyOnce {
//		if len(args) == 0 {
//			args = []string{loadTools.CmdVersionLatest}
//		}
//
//		update := toolSelfUpdate.New(CmdScribe.Runtime)
//		if update.State.IsError() {
//			state = update.State
//			break
//		}
//
//		for _, v := range args {
//			state = update.PrintVersion(toolSelfUpdate.GetSemVer(v))
//			if state.IsNotOk() {
//				state = update.State
//				break
//			}
//		}
//	}
//
//	return state
//}
//
//
//func VersionList(args ...string) *ux.State {
//	state := CmdScribe.State
//
//	for range onlyOnce {
//		if len(args) == 0 {
//			// @TODO = Obtain full list of versions.
//			args = []string{loadTools.CmdVersionLatest}
//		}
//
//		update := toolSelfUpdate.New(CmdScribe.Runtime)
//		if update.State.IsError() {
//			state = update.State
//			break
//		}
//
//		for _, v := range args {
//			state = update.PrintVersion(toolSelfUpdate.GetSemVer(v))
//			if state.IsNotOk() {
//				state = update.State
//				break
//			}
//		}
//	}
//
//	return state
//}
//
//
//func VersionCheck() *ux.State {
//	state := CmdScribe.State
//
//	for range onlyOnce {
//		update := toolSelfUpdate.New(CmdScribe.Runtime)
//		if update.State.IsError() {
//			state = update.State
//			break
//		}
//
//		state = update.IsUpdated(true)
//		if update.State.IsError() {
//			break
//		}
//	}
//
//	return state
//}
//
//
//func VersionUpdate() *ux.State {
//	state := CmdScribe.State
//
//	for range onlyOnce {
//		update := toolSelfUpdate.New(CmdScribe.Runtime)
//		if update.State.IsError() {
//			state = update.State
//			break
//		}
//
//		state = update.IsUpdated(true)
//		if update.State.IsError() {
//			break
//		}
//
//		state = update.Update()
//		if state.IsNotOk() {
//			break
//		}
//	}
//
//	return state
//}
//
//
//
//func _GetUsage(c *cobra.Command) string {
//	var str string
//
//	if c.Parent() == nil {
//		str += ux.SprintfCyan("%s ", c.Name())
//	} else {
//		str += ux.SprintfCyan("%s ", c.Parent().Name())
//		str += ux.SprintfGreen("%s ", c.Use)
//	}
//
//	if c.HasAvailableSubCommands() {
//		str += ux.SprintfGreen("[command] ")
//		str += ux.SprintfCyan("<args> ")
//	}
//
//	return str
//}
//
//
//func _GetVersion(c *cobra.Command) string {
//	var str string
//
//	if c.Parent() == nil {
//		str = ux.SprintfBlue("%s ", CmdScribe.Runtime.CmdName)
//		str += ux.SprintfCyan("v%s", CmdScribe.Runtime.CmdVersion)
//	}
//
//	return str
//}
//
//
//func SetHelp(c *cobra.Command) {
//	var tmplHelp string
//	var tmplUsage string
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
//{{ SprintfBlue "Commmand:" }} {{ SprintfCyan .Use }}
//
//{{ SprintfBlue "Description:" }}
//	{{ with (or .Long .Short) }}
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
//
//type Example struct {
//	Command string
//	Args []string
//	Info string
//}
//type Examples []Example
//
//
//func HelpExamples() {
//	CmdScribe.State.Clear()
//	return
//
//	for range onlyOnce {
//		var examples Examples
//
//
//		examples = append(examples, Example {
//			Command: "load",
//			Args:    []string{"-json", "config.json", "-template", "'{{ .Json.dir }}'"},
//			Info:    "Print to STDOUT the .dir key from config.json.",
//		})
//		examples = append(examples, Example {
//			Command: "load",
//			Args:    []string{"-template", "'{{ .Json.dir }}'", "config.json"},
//			Info:    "The same thing, but with less arguments.",
//		})
//
//		examples = append(examples, Example {
//			Command: "load",
//			Args:    []string{"-template", "'{{ .Json.hello }}'", "-json", `'{ "hello": "world" }'`},
//			Info:    "Template and JSON arguments can be either string or file reference.",
//		})
//		examples = append(examples, Example {
//			Command: "load",
//			Args:    []string{"-template", "hello_world.tmpl", "-json", `'{ "hello": "world" }'`},
//			Info:    "The same again...",
//		})
//		examples = append(examples, Example {
//			Command: "load",
//			Args:    []string{"-template", "'{{ .Json.hello }}'", "-json", "hello.json"},
//			Info:    "The same again...",
//		})
//
//
//		examples = append(examples, Example {
//			Command: "load",
//			Args:    []string{"-json", "config.json", "-template", "DockerFile.tmpl", "-out", "Dockerfile"},
//			Info:    "Process Dockerfile.tmpl file and output to Dockerfile.",
//		})
//		examples = append(examples, Example {
//			Command: "load",
//			Args:    []string{"-out", "Dockerfile", "config.json", "DockerFile.tmpl"},
//			Info:    "And again with less arguments..",
//		})
//		examples = append(examples, Example {
//			Command: "convert",
//			Args:    []string{"config.json", "DockerFile.tmpl"},
//			Info:    "'convert' does the same , but removes the template file afterwards...",
//		})
//
//
//		examples = append(examples, Example {
//			Command: "load",
//			Args:    []string{"-out", "MyScript.sh", "MyScript.sh.tmpl", "config.json"},
//			Info:    "Process the MyScript.sh.tmpl template file and write the result to MyScript.sh.",
//		})
//		examples = append(examples, Example {
//			Command: "convert",
//			Args:    []string{"MyScript.sh.tmpl", "config.json"},
//			Info:    "Same again using 'convert'. Template and json files can be in any order.",
//		})
//		examples = append(examples, Example {
//			Command: "run",
//			Args:    []string{"MyScript.sh.tmpl", "config.json"},
//			Info:    "Same again using 'run'. This will execute the MyScript.sh output file afterwards.",
//		})
//
//
//		ux.PrintflnBlue("Examples:")
//		for _, v := range examples {
//			fmt.Printf("# %s\n\t%s %s\n\n",
//				ux.SprintfBlue(v.Info),
//				ux.SprintfCyan("%s %s", BinaryName, v.Command),
//				ux.SprintfWhite(strings.Join(v.Args, " ")),
//			)
//		}
//	}
//
//	CmdScribe.State.Clear()
//}
