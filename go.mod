module launch

go 1.14

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/cavaliercoder/grab v2.0.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/gearboxworks/go-osbridge v0.0.0-20190605062119-0e1c68c1c70f
	github.com/getlantern/context v0.0.0-20190109183933-c447772a6520 // indirect
	github.com/getlantern/errors v0.0.0-20190325191628-abdb3e3e36f7
	github.com/getlantern/hex v0.0.0-20190417191902-c6586a6fe0b7 // indirect
	github.com/getlantern/hidden v0.0.0-20190325191715-f02dbb02be55 // indirect
	github.com/getlantern/ops v0.0.0-20190325191751-d70cb0d6f85f // indirect
	github.com/go-openapi/strfmt v0.19.5 // indirect
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/jedib0t/go-pretty v4.3.0+incompatible // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/newclarity/scribeHelpers/helperDocker v0.0.0-00010101000000-000000000000
	github.com/newclarity/scribeHelpers/ux v0.0.0-00010101000000-000000000000
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/sftp v1.11.0 // indirect
	github.com/rhysd/go-github-selfupdate v1.2.2
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.4.0
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
)

replace github.com/newclarity/scribeHelpers/ux => ../scribeHelpers/ux

replace github.com/newclarity/scribeHelpers/scribeLoader => ../scribeHelpers/scribeLoader

replace github.com/newclarity/scribeHelpers/helperCopy => ../scribeHelpers/helperCopy

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190717161051-705d9623b7c1
replace github.com/newclarity/scribeHelpers/helperDocker => ../scribeHelpers/helperDocker

replace github.com/newclarity/scribeHelpers/helperExec => ../scribeHelpers/helperExec

replace github.com/newclarity/scribeHelpers/helperGit => ../scribeHelpers/helperGit

replace github.com/newclarity/scribeHelpers/helperGitHub => ../scribeHelpers/helperGitHub

replace github.com/newclarity/scribeHelpers/helperPath => ../scribeHelpers/helperPath

replace github.com/newclarity/scribeHelpers/helperPrompt => ../scribeHelpers/helperPrompt

replace github.com/newclarity/scribeHelpers/helperService => ../scribeHelpers/helperService

replace github.com/newclarity/scribeHelpers/helperSystem => ../scribeHelpers/helperSystem

replace github.com/newclarity/scribeHelpers/helperTypes => ../scribeHelpers/helperTypes

replace github.com/newclarity/scribeHelpers/helperUx => ../scribeHelpers/helperUx

replace github.com/newclarity/scribeHelpers/helperRuntime => ../scribeHelpers/helperRuntime
