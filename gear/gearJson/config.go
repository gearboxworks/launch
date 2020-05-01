package gearJson

import (
	"encoding/json"
	"fmt"
	"launch/defaults"
	"launch/only"
	"launch/ux"
	"os"
	"path/filepath"
)

type GearConfig struct {
	Meta       GearMeta       `json:"meta"`
	Build      GearBuild      `json:"build"`
	Run        GearRun        `json:"run"`
	Project    GearProject    `json:"project"`
	Extensions GearExtensions `json:"extensions"`
	Versions   GearVersions   `json:"versions"`

	Schema     string         `json:"schema"`
}
type GearConfigs map[string]GearConfig

func (me *GearConfig) GetName() string {
	return me.Meta.Name
}

func (me *GearConfig) GetCommand(cmd []string) []string {
	var retCmd []string

	for range only.Once {
		var cmdExec string
		switch {
			case len(cmd) == 0:
				cmdExec = defaults.DefaultCommandName

			case cmd[0] == "":
				cmdExec = defaults.DefaultCommandName

			case cmd[0] == me.Meta.Name:
				cmdExec = defaults.DefaultCommandName

			case cmd[0] != "":
				cmdExec = cmd[0]

			default:
				//cmdExec = cmd[0]
				cmdExec = defaults.DefaultCommandName
		}

		c := me.MatchCommand(cmdExec)
		if c == nil {
			retCmd = []string{}
			break
		}

		retCmd = append([]string{*c}, cmd[1:]...)
	}

	return retCmd
}

func (me *GearConfig) MatchCommand(cmd string) *string {
	var c *string

	for range only.Once {
		if c2, ok := me.Run.Commands[cmd]; ok {
			c = &c2
			break
		}
	}

	return c
}

func (me *GearConfig) CreateLinks(c defaults.ExecCommand, name string, version string) ux.State {
	var state ux.State

	for range only.Once {
		state = me.ValidateGearConfig()
		if state.IsError() {
			break
		}

		var created bool
		for k, v := range me.Run.Commands {
			var err error
			var dstFile string
			var linkStat os.FileInfo

			if k == "default" {
				k = filepath.Base(v)
			}

			if version == "latest" {
				dstFile, err = filepath.Abs(fmt.Sprintf("%s%c%s", c.Dir, filepath.Separator, k))
			} else {
				dstFile, err = filepath.Abs(fmt.Sprintf("%s%c%s-%s", c.Dir, filepath.Separator, k, version))
			}
			if err != nil {
				continue
			}

			linkStat, err = os.Lstat(dstFile)
			if linkStat == nil {
				created = true

				// Symlink doesn't exist - create.
				err = os.Symlink(c.File, dstFile)
				if err != nil {
					continue
				}

				//continue
				linkStat, err = os.Lstat(dstFile)
				if linkStat == nil {
					continue
				}
			}

			// Symlink exists - validate.
			l, _ := os.Readlink(dstFile)
			//if !filepath.IsAbs(l) {
			//	l, _ = filepath.Abs(fmt.Sprintf("%s%c%s", c.Dir, filepath.Separator, l))
			//}
			if l == "" {

			}

			//fmt.Printf("'%s' (%s) => '%s'\n", k, dstFile, v)
			//fmt.Printf("\tReadlink() => %s\n", l)
			//fmt.Printf("\tLstat() => %s	%s	%s	%s	%d\n",
			//	linkStat.Name(),
			//	linkStat.IsDir(),
			//	linkStat.Mode().String(),
			//	linkStat.ModTime().String(),
			//	linkStat.Size(),
			//)
			//fmt.Printf("\n")
		}

		if created {
			ux.PrintfOk("Created application links.\n")
		}
	}

	return state
}

func (me *GearConfig) RemoveLinks(c defaults.ExecCommand, name string, version string) ux.State {
	var state ux.State

	for range only.Once {
		state = me.ValidateGearConfig()
		if state.IsError() {
			break
		}

		var removed bool
		for k, _ := range me.Run.Commands {
			var err error
			var dstFile string
			var linkStat os.FileInfo

			if k == "default" {
				continue
			}

			dstFile, err = filepath.Abs(fmt.Sprintf("%s%c%s-%s", c.Dir, filepath.Separator, k, version))
			if err != nil {
				continue
			}

			linkStat, err = os.Lstat(dstFile)
			if err != nil {
				continue
			}

			if linkStat == nil {
				// Symlink doesn't exist.
				continue
			}

			removed = true

			l, _ := os.Readlink(dstFile)
			if l == defaults.BinaryName {
				// Symlink exists - remove.
				//if !filepath.IsAbs(l) {
				//	l, _ = filepath.Abs(fmt.Sprintf("%s%c%s", c.Dir, filepath.Separator, l))
				//}
				err = os.Remove(dstFile)
			}
		}

		if removed {
			ux.PrintfOk("Removed application links.\n")
		}
	}

	return state
}


type GearMeta struct {
	State        string `json:"state"`
	Organization string `json:"organization"`
	Name         string `json:"name"`
	Maintainer   string `json:"maintainer"`
	Class        string `json:"class"`
	Refurl       string `json:"refurl"`
}


/*
Command execution is complex and there's several steps to the logic.
Essentially, the ENTRYPOINT image definition is converted to an S6 service.

During build:
1. The ENTRYPOINT definition within a Docker image needs to be pulled in to a Gearbox image.
2. This is used to start any service that was defined with ENTRYPOINT within the original image.
3. GearBuild.Run will contain index 0 of the ENTRYPOINT array.
4. GearBuild.Args will contain slice [1:] of the ENTRYPOINT array.

During runtime(boot):
1. GEARBOX_ENTRYPOINT, (aka GearBuild.Run), will be checked and executed as part of an S6 service.
2. GEARBOX_ENTRYPOINT_ARGS, (GearBuild.Args), will be appended and the whole service started.

During runtime(interactive command via symlink):
1. ARG 1 of the command line will be checked against GearRun.Commands for every container on the system.
2. When found, will execute.

During runtime(other interactive commands):
1.


GearBuild.Run
	This will default to GEARBOX_ENTRYPOINT env from within the image build process.
It is generated from the command: `docker inspect --format '{{ with .ContainerConfig.Entrypoint}} {{ index . 0 }}{{ end }}'`

GearBuild.Args
	This will default to GEARBOX_ENTRYPOINT_ARGS env from within the image build process.
It is generated from the command: `docker inspect --format '{{ join .ContainerConfig.Entrypoint " " }}'`
Any additional arguments provided by the user will be appended to this at runtime.
*/

type GearBuild struct {
	Ports        GearPorts    `json:"ports"`
	Run          string       `json:"run"`		//
	Args         GearArgs     `json:"args"`		//
	Env          GearEnv      `json:"env"`
	Network      string       `json:"network"`
}


type GearRun struct {
	Ports        GearPorts    `json:"ports"`
	Env          GearEnv      `json:"env"`
	Volumes      string       `json:"volumes"`
	Network      string       `json:"network"`
	Commands     GearCommands `json:"commands"`
}
//type GearCommand string
//type GearCommands map[string]GearCommand
type GearCommands map[string]string

type GearProject struct {
}

type GearExtensions struct {
}

type GearEnv map[string]string
type GearPorts map[string]string
type GearArgs string

func (ports *GearPorts) ToString() string {
	var p string

	for k, v := range *ports {
		p = fmt.Sprintf("%s %s:%s\n", p, k, v)
	}

	return p
}

type GearVersion struct {
	MajorVersion string `json:"majorversion"`
	Latest       bool   `json:"latest"`
	Ref          string `json:"ref"`
	Base         string `json:"base"`
}
type GearVersions map[string]GearVersion

func (vers *GearVersions) GetLatest() string {
	var v string

	var r GearVersion
	for v, r = range *vers {
		if r.Latest {
			break
		}
	}

	return v
}

func (vers *GearVersions) HasVersion(c string) bool {
	for v, r := range *vers {
		if r.Latest && (c == "latest") {
			return true
		}

		if v == c {
			return true
		}
	}
	return false
}


func New(cs string) (*GearConfig, ux.State) {
	var gc GearConfig
	var state ux.State

	for range only.Once {
		if cs == "" {
			state.SetError("gear config is empty")
			break
		}

		js := []byte(cs)
		if js == nil {
			state.SetError("gear config json is nil")
			break
		}

		err := json.Unmarshal(js, &gc)
		if err != nil {
			state.SetError("gearbox.json schema unknown: %s", err)
			break
		}

		state = gc.ValidateGearConfig()
		if state.IsError() {
			break
		}
	}

	return &gc, state
}

func (me *GearConfig) ValidateGearConfig() ux.State {
	var state ux.State

	for range only.Once {
		if me == nil {
			state.SetError("gear config is nil")
			break
		}
	}

	return state
}
