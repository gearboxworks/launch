package cmd

import (
	"fmt"
	"github.com/gearboxworks/scribeHelpers/toolGear"
	"github.com/gearboxworks/scribeHelpers/toolGear/gearConfig"
	"github.com/gearboxworks/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

//goland:noinspection GoUnusedFunction
func showArgs(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		flargs := cmd.Flags().Args()
		if flargs != nil {
			ux.Printf("'%s' called with '%s'\n", cmd.CommandPath(), strings.Join(flargs, " "))
			break
		}

		ux.Printf("'%s' called with '%s'\n", cmd.CommandPath(), strings.Join(args, " "))
		break
	}

	fmt.Println("")
}

type LaunchArgs struct {
	Name      string
	Version   string

	Args      []string
	Project   string
	Mount     string
	TmpDir    string

	Temporary bool
	SshStatus bool
	Quiet     bool
	Debug     bool
	NoCreate  bool
	Timeout   time.Duration

	Provider  *toolGear.Provider
	//GearRef   *toolGear.Gear
	Gears     toolGear.Gears

	Valid     bool
	State     *ux.State
}


func (ga *LaunchArgs) IsNil() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}
	ga.State = ga.State.EnsureNotNil()
	return ga.State
}

func (ga *LaunchArgs) IsValid() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		if !ga.Valid {
			ga.State.SetError("%s args are not valid", defaults.LanguageContainerName)
			break
		}

		ga.State = ga.State.EnsureNotNil()
	}

	return ga.State
}


//goland:noinspection GoUnusedParameter
func (ga *LaunchArgs) ProcessArgs(cmd *cobra.Command, args []string, scan bool) *ux.State {
	for range onlyOnce {
		ga.State = ux.NewState(Cmd.Runtime.CmdName, Cmd.Debug)

		ga.Args = args
		if len(ga.Args) > 0 {
			ga.Name = ga.Args[0]
			if strings.Contains(ga.Name, ":") {
				spl := strings.Split(ga.Name, ":")
				ga.Name = spl[0]
				ga.Version = spl[1]
			}

			if ga.Version == "" {
				ga.Version = gearConfig.LatestName
			}
		}

		if ga.Project != defaults.DefaultPathNone {
			ga.Project = DeterminePath(Cmd.Project)
		}

		if ga.Mount != defaults.DefaultPathNone {
			ga.Mount = DeterminePath(Cmd.Mount)
		}

		if ga.TmpDir != defaults.DefaultPathNone {
			ga.TmpDir = DeterminePath(Cmd.TmpDir)
		}

		ga.Debug = Cmd.Debug
		Cmd.Runtime.Debug = Cmd.Debug
		ga.State.DebugSet(Cmd.Debug)

		ga.Timeout = Cmd.Timeout
		Cmd.Runtime.Timeout = Cmd.Timeout


		//ga.Provider = toolGear.NewProvider(Cmd.Runtime)
		//ga.State = ga.Provider.State
		//if ga.State.IsError() {
		//	break
		//}
		//
		//ga.State = ga.Provider.SetProvider(Cmd.Provider)
		//if ga.State.IsError() {
		//	break
		//}
		//
		//ga.State = ga.Provider.SetHost(Cmd.Host, Cmd.Port)
		//if ga.State.IsError() {
		//	break
		//}

		ga.Gears = toolGear.NewGears(Cmd.Runtime)
		ga.State = ga.Gears.State
		if ga.State.IsError() {
			break
		}

		ga.State = ga.Gears.SetLanguage(defaults.LanguageAppName, defaults.LanguageImageName, defaults.LanguageContainerName)

		ga.State = ga.Gears.SetProvider(Cmd.Provider)
		if ga.State.IsError() {
			break
		}

		ga.State = ga.Gears.SetProviderHost(Cmd.Host, Cmd.Port)
		if ga.State.IsError() {
			break
		}

		if scan {
			ga.State = ga.Gears.Get()
			if ga.State.IsNotOk() {
				break
			}
		}

		//ga.GearRef = toolGear.NewGear(Cmd.Runtime)
		//ga.State = ga.GearRef.State
		//if ga.State.IsError() {
		//	break
		//}

		ga.Valid = true
	}

	return ga.State
}


func (ga *LaunchArgs) ListLinks(create bool) *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.Gears.GetContainers(ga.Name)

		for _, dc := range ga.Gears.Array {
			if dc.Container.ID == "" {
				continue
			}

			if create {
				ga.State = dc.GearConfig.CreateLinks(dc.Container.Version)
			}

			if ga.Name == "" {
				ga.State = dc.ListLinks()	// dc.Container.Version)
				continue
			}

			if ga.Name != dc.Container.Name {
				continue
			}

			if ga.Version == "all" {
				ga.State = dc.ListLinks()	// dc.Container.Version)
				continue
			}

			if ga.Version == gearConfig.LatestName {
				if dc.Container.IsLatest {
					ga.State = dc.ListLinks()	// dc.Container.Version)
				}
				continue
			}

			if ga.Version == "" {
				ga.State = dc.ListLinks()	// dc.Container.Version)
				continue
			}

			if dc.Container.Version == ga.Version {
				ga.State = dc.ListLinks()	// dc.Container.Version)
				continue
			}
		}
	}

	return ga.State
}


//goland:noinspection GoUnusedParameter
func (ga *LaunchArgs) ListPorts(remote bool) *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range onlyOnce {
		for _, dc := range ga.Gears.Array {
			if ga.Name == "" {
				ga.State = dc.ListContainerPorts()
				ga.State = dc.ListImagePorts()
				continue
			}

			if ga.Name != dc.Container.Name {
				continue
			}

			if ga.Version == "all" {
				ga.State = dc.ListContainerPorts()
				continue
			}

			if ga.Version == gearConfig.LatestName {
				if dc.Container.IsLatest {
					ga.State = dc.ListContainerPorts()
				}
				continue
			}

			if ga.Version == "" {
				ga.State = dc.ListContainerPorts()
				continue
			}

			if dc.Container.Version == ga.Version {
				ga.State = dc.ListContainerPorts()
				continue
			}
		}
	}

	return ga.State
}


func (ga *LaunchArgs) CreateLinks(version string) *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.Gears.Selected.CreateLinks(version)
	}

	return ga.State
}


func (ga *LaunchArgs) RemoveLinks(version string) *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.Gears.Selected.RemoveLinks(version)
	}

	return ga.State
}

// DeterminePath takes a string and return a path.
func DeterminePath(mp string) string {
	var ok bool

	for range onlyOnce {
		var err error
		var cwd string

		if mp == defaults.DefaultPathNone {
			break
		}

		switch {
			case mp == defaults.DefaultPathEmpty:
				fallthrough
			case mp == defaults.DefaultPathCwd:
				cwd, err = os.Getwd()
				if err != nil {
					break
				}
				ok = true
				mp = cwd

			case mp == defaults.DefaultPathHome:
				var u *user.User
				u, err = user.Current()
				if err != nil {
					break
				}
				ok = true
				mp = u.HomeDir

			default:
				mp, err = filepath.Abs(mp)
				if err != nil {
					break
				}
				ok = true
		}

		if err != nil {
			break
		}

		if !ok {
			mp = defaults.DefaultPathNone
		}
	}

	return mp
}

//GetLaunchDir return launch directory.
func GetLaunchDir() string {
	var d string

	for range onlyOnce {
		u, _ := user.Current()
		d = filepath.Join(u.HomeDir, ".launch")
	}

	return d
}
