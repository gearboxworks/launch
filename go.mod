module launch

go 1.14

replace github.com/newclarity/scribeHelpers/loadTools => ../scribeHelpers/loadTools

replace github.com/newclarity/scribeHelpers/ux => ../scribeHelpers/ux

replace github.com/newclarity/scribeHelpers/toolCopy => ../scribeHelpers/toolCopy

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190717161051-705d9623b7c1

replace github.com/newclarity/scribeHelpers/toolDocker => ../scribeHelpers/toolDocker

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

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/getlantern/context v0.0.0-20190109183933-c447772a6520 // indirect
	github.com/getlantern/hex v0.0.0-20190417191902-c6586a6fe0b7 // indirect
	github.com/getlantern/hidden v0.0.0-20190325191715-f02dbb02be55 // indirect
	github.com/getlantern/ops v0.0.0-20190325191751-d70cb0d6f85f // indirect
	github.com/go-openapi/strfmt v0.19.5 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/jedib0t/go-pretty v4.3.0+incompatible // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/newclarity/scribeHelpers/toolGear v0.0.0-00010101000000-000000000000
	github.com/newclarity/scribeHelpers/toolRuntime v0.0.0-20200604000029-dbb313f0fedc
	github.com/newclarity/scribeHelpers/ux v0.0.0-20200604000029-dbb313f0fedc
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/sftp v1.11.0 // indirect
	github.com/rhysd/go-github-selfupdate v1.2.2
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.4.0
)
