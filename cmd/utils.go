package cmd

import (
	"fmt"
	"github.com/newclarity/scribeHelpers/toolGear"
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

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

	Provider  *toolGear.Provider
	GearRef   *toolGear.Gear

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


func (ga *LaunchArgs) ProcessArgs(cmd *cobra.Command, args []string) *ux.State {
	for range onlyOnce {
		ga.State = ux.NewState(Cmd.Runtime.CmdName, Cmd.Debug)

		//fmt.Printf("cmd.Args:%v\nargs:%v\nCmd.Runtime.Args:%v\nCmd.Runtime.FullArgs:%v\n",
		//	cmd.Args,
		//	args,
		//	Cmd.Runtime.Args,
		//	Cmd.Runtime.FullArgs,
		//	)
		ga.Args = args
		if len(ga.Args) > 0 {
			ga.Name = ga.Args[0]
			if strings.Contains(ga.Name, ":") {
				spl := strings.Split(ga.Name, ":")
				ga.Name = spl[0]
				ga.Version = spl[1]
			}

			if ga.Version == "" {
				ga.Version = "latest"
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

		ga.Provider = toolGear.NewProvider(Cmd.Runtime)
		ga.State = ga.Provider.State
		if ga.State.IsError() {
			break
		}

		ga.State = ga.Provider.SetProvider(Cmd.Provider)
		if ga.State.IsError() {
			break
		}

		ga.State = ga.Provider.SetHost(Cmd.Host, Cmd.Port)
		if ga.State.IsError() {
			break
		}

		ga.GearRef = toolGear.NewGear(Cmd.Runtime)
		ga.State = ga.GearRef.State
		if ga.State.IsError() {
			break
		}

		ga.Valid = true
	}

	return ga.State
}


func (ga *LaunchArgs) ListLinks(version string) *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range onlyOnce {
		var dcs toolGear.TypeDockerGears
		dcs, ga.State = ga.GearRef.Docker.GetContainers(ga.Name)

		for _, dc := range dcs {
			ga.State = dc.Container.GearConfig.ListLinks(dc.Container.Version)
		}
	}

	return ga.State
}


func (ga *LaunchArgs) CreateLinks(version string) *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.GearRef.CreateLinks(version)
	}

	return ga.State
}


func (ga *LaunchArgs) RemoveLinks(version string) *ux.State {
	if state := ga.IsNil(); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.GearRef.RemoveLinks(version)
	}

	return ga.State
}


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


func GetGearboxDir() string {
	var d string

	for range onlyOnce {
		u, _ := user.Current()
		d = filepath.Join(u.HomeDir, ".gearbox")
	}

	return d
}
