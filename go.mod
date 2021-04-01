module launch

go 1.14

replace github.com/newclarity/scribeHelpers/loadTools => ../scribeHelpers/loadTools

replace github.com/newclarity/scribeHelpers/ux => ../scribeHelpers/ux

replace github.com/newclarity/scribeHelpers/toolCopy => ../scribeHelpers/toolCopy

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190717161051-705d9623b7c1

// replace github.com/newclarity/scribeHelpers/toolDocker => ../scribeHelpers/toolDocker

replace github.com/newclarity/scribeHelpers/toolExec => ../scribeHelpers/toolExec

replace github.com/newclarity/scribeHelpers/toolGear => ../scribeHelpers/toolGear

replace github.com/newclarity/scribeHelpers/toolGit => ../scribeHelpers/toolGit

replace github.com/newclarity/scribeHelpers/toolGitHub => ../scribeHelpers/toolGitHub

replace github.com/newclarity/scribeHelpers/toolPath => ../scribeHelpers/toolPath

replace github.com/newclarity/scribeHelpers/toolPrompt => ../scribeHelpers/toolPrompt

replace github.com/newclarity/scribeHelpers/toolService => ../scribeHelpers/toolService

replace github.com/newclarity/scribeHelpers/toolSystem => ../scribeHelpers/toolSystem

replace github.com/newclarity/scribeHelpers/toolTypes => ../scribeHelpers/toolTypes

replace github.com/newclarity/scribeHelpers/toolUx => ../scribeHelpers/toolUx

replace github.com/newclarity/scribeHelpers/toolRuntime => ../scribeHelpers/toolRuntime

replace github.com/newclarity/scribeHelpers/toolSelfUpdate => ../scribeHelpers/toolSelfUpdate

replace github.com/newclarity/scribeHelpers/toolGhr => ../scribeHelpers/toolGhr

replace github.com/newclarity/scribeHelpers/toolCobraHelp => ../scribeHelpers/toolCobraHelp

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-openapi/errors v0.19.6 // indirect
	github.com/mitchellh/mapstructure v1.3.2 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/newclarity/scribeHelpers/loadTools v0.0.0-20200621234507-ba6f08c6b68d
	github.com/newclarity/scribeHelpers/toolCobraHelp v0.0.0-00010101000000-000000000000
	github.com/newclarity/scribeHelpers/toolCopy v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/newclarity/scribeHelpers/toolExec v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/newclarity/scribeHelpers/toolGear v0.0.0-20200621234507-ba6f08c6b68d
	github.com/newclarity/scribeHelpers/toolGhr v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/newclarity/scribeHelpers/toolGit v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/newclarity/scribeHelpers/toolGitHub v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/newclarity/scribeHelpers/toolRuntime v0.0.0-20200623081955-45abb1cbefe9
	github.com/newclarity/scribeHelpers/toolSelfUpdate v0.0.0-20200621234507-ba6f08c6b68d
	github.com/newclarity/scribeHelpers/toolService v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/newclarity/scribeHelpers/toolSystem v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/newclarity/scribeHelpers/toolUx v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/newclarity/scribeHelpers/ux v0.0.0-20200623081955-45abb1cbefe9
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pkg/profile v1.5.0
	github.com/spf13/afero v1.3.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.0
	google.golang.org/genproto v0.0.0-20200620020550-bd6e04640131 // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)
