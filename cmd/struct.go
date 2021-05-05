package cmd

import (
	"github.com/gearboxworks/scribeHelpers/toolGear"
	"github.com/gearboxworks/scribeHelpers/toolRuntime"
	"github.com/gearboxworks/scribeHelpers/ux"
	"launch/defaults"
	"time"
)


const (
	//flagConfig   = "config"
	flagDebug    = "debug"
	flagNoCreate = "no-create"

	flagExample  = "help-examples"
	flagHelp  = "flags"

	flagProvider        = "provider"
	flagProviderDefault = "docker"

	flagHost       = "host"
	flagPort       = "port"
	flagProject    = "project"
	flagMount      = "mount"
	flagTmpDir     = "tmp"
	//flagCompletion = "completion"
	flagQuiet      = "quiet"
	flagTemporary  = "temporary"
	flagStatus     = "status"
	flagTimeout    = "timeout"
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

	NoCreate   bool   			// Flag: --no-create
	Debug      bool   			// Flag: --debug
	Quiet      bool   			// Flag: --quiet
	Temporary  bool   			// Flag: --temporary
	Status     bool   			// Flag: --status
	Timeout    time.Duration  	// Flag: --timeout

	HelpAll        bool
	HelpExamples   bool
	HelpFlags      bool

	Runtime        *toolRuntime.TypeRuntime
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
		Timeout:      toolGear.DefaultTimeout,

		HelpAll:      false,
		HelpExamples: false,
		HelpFlags:    false,

		Runtime:        toolRuntime.New(defaults.BinaryName, defaults.BinaryVersion, false),
		State:          ux.NewState(defaults.BinaryName, false),
		valid:          false,
	}

	return &la
}

func (la *TypeLaunchArgs) _SetDebug(d bool) {
	la.Debug = d
	if la.Runtime != nil {
		la.Runtime.Debug = d
	}
	if la.State != nil {
		la.State.DebugSet(d)
	}
}

func (la *TypeLaunchArgs) _SetOptions(ga LaunchArgs) {
	la._SetDebug(ga.Debug)
	la.Timeout = ga.Timeout
	if la.Runtime.Timeout == 0 {
		la.Runtime.Timeout = ga.Timeout
	}
}
