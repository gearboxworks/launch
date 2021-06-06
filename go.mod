module launch

go 1.14

replace github.com/gearboxworks/scribeHelpers/loadTools => ../scribeHelpers/loadTools

replace github.com/gearboxworks/scribeHelpers/ux => ../scribeHelpers/ux

replace github.com/gearboxworks/scribeHelpers/toolCopy => ../scribeHelpers/toolCopy

// replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190717161051-705d9623b7c1
//replace github.com/docker/docker => github.com/docker/engine v1.13.1
//replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.8.1

replace github.com/gearboxworks/scribeHelpers/toolExec => ../scribeHelpers/toolExec

replace github.com/gearboxworks/scribeHelpers/toolGear => ../scribeHelpers/toolGear

replace github.com/gearboxworks/scribeHelpers/toolGit => ../scribeHelpers/toolGit

replace github.com/gearboxworks/scribeHelpers/toolGitHub => ../scribeHelpers/toolGitHub

replace github.com/gearboxworks/scribeHelpers/toolPath => ../scribeHelpers/toolPath

replace github.com/gearboxworks/scribeHelpers/toolPrompt => ../scribeHelpers/toolPrompt

replace github.com/gearboxworks/scribeHelpers/toolService => ../scribeHelpers/toolService

replace github.com/gearboxworks/scribeHelpers/toolSystem => ../scribeHelpers/toolSystem

replace github.com/gearboxworks/scribeHelpers/toolTypes => ../scribeHelpers/toolTypes

replace github.com/gearboxworks/scribeHelpers/toolUx => ../scribeHelpers/toolUx

replace github.com/gearboxworks/scribeHelpers/toolRuntime => ../scribeHelpers/toolRuntime

replace github.com/gearboxworks/scribeHelpers/toolSelfUpdate => ../scribeHelpers/toolSelfUpdate

replace github.com/gearboxworks/scribeHelpers/toolGhr => ../scribeHelpers/toolGhr

replace github.com/gearboxworks/scribeHelpers/toolCobraHelp => ../scribeHelpers/toolCobraHelp

replace github.com/gearboxworks/scribeHelpers/toolNetwork => ../scribeHelpers/toolNetwork

require (
	github.com/containerd/containerd v1.5.2 // indirect
	github.com/docker/docker v20.10.6+incompatible // indirect
	github.com/gearboxworks/scribeHelpers/loadTools v0.0.0-20200621234507-ba6f08c6b68d
	github.com/gearboxworks/scribeHelpers/toolCobraHelp v0.0.0-00010101000000-000000000000
	github.com/gearboxworks/scribeHelpers/toolCopy v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/gearboxworks/scribeHelpers/toolGear v0.0.0-20200621234507-ba6f08c6b68d
	github.com/gearboxworks/scribeHelpers/toolGhr v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/gearboxworks/scribeHelpers/toolGit v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/gearboxworks/scribeHelpers/toolGitHub v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/gearboxworks/scribeHelpers/toolNetwork v0.0.0-00010101000000-000000000000 // indirect
	github.com/gearboxworks/scribeHelpers/toolRuntime v0.0.0
	github.com/gearboxworks/scribeHelpers/toolSelfUpdate v0.0.0-20200621234507-ba6f08c6b68d
	github.com/gearboxworks/scribeHelpers/toolService v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/gearboxworks/scribeHelpers/toolSystem v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/gearboxworks/scribeHelpers/toolUx v0.0.0-20200621234507-ba6f08c6b68d // indirect
	github.com/gearboxworks/scribeHelpers/ux v0.0.0
	github.com/go-openapi/errors v0.19.6 // indirect
	github.com/mitchellh/mapstructure v1.3.2 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect; indirectx`
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/pkg/profile v1.5.0
	github.com/spf13/afero v1.3.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	gopkg.in/ini.v1 v1.57.0 // indirect
)
