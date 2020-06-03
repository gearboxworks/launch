package cmd

import (
	"github.com/newclarity/scribeHelpers/helperRuntime"
	"github.com/newclarity/scribeHelpers/ux"
	"launch/defaults"
)


const (
	flagConfig   = "config"
	flagDebug    = "debug"
	flagNoCreate = "no-create"

	flagExample  = "help-examples"

	flagProvider        = "provider"
	flagProviderDefault = "docker"

	flagHost       = "host"
	flagPort       = "port"
	flagProject    = "project"
	flagMount      = "mount"
	flagTmpDir     = "tmp"
	flagCompletion = "completion"
	flagVersion    = "version"
	flagQuiet      = "quiet"
	flagTemporary  = "temporary"
	flagStatus     = "status"
)


type TypeLaunchArgs struct {
	Completion bool   // Flag: --completion
	Version    bool   // Flag: --version

	Config     string // Flag: --config

	Provider   string // Flag: --provider
	Host       string // Flag: --host
	Port       string // Flag: --port

	Project    string // Flag: --project
	Mount      string // Flag: --mount
	TmpDir     string // Flag: --tmp

	NoCreate   bool   // Flag: --no-create
	Debug      bool   // Flag: --debug
	Quiet      bool   // Flag: --quiet
	Temporary  bool   // Flag: --temporary
	Status     bool   // Flag: --status

	HelpAll        bool
	HelpExamples   bool

	Runtime        *helperRuntime.TypeRuntime
	State          *ux.State
	valid          bool
}


func New() *TypeLaunchArgs {

	la := TypeLaunchArgs{
		Completion:   false,
		Version:      false,

		Config:       "",

		Provider:     "",
		Host:         "",
		Port:         "",

		Project:      "",
		Mount:        "",
		TmpDir:       "",

		NoCreate:     false,
		Debug:        false,
		Quiet:        false,
		Temporary:    false,
		Status:       false,

		HelpAll:      false,
		HelpExamples: false,

		Runtime:        helperRuntime.New(defaults.BinaryName, defaults.BinaryVersion, false),
		State:          ux.NewState(defaults.BinaryName, false),
		valid:          false,
	}

	return &la
}