package gearJson

import (
	"encoding/json"
	"fmt"
	"launch/defaults"
	"launch/only"
	"launch/ux"
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

During runtime(interactive commands):
1. ARG 1 of the command line will be checked against GearRun.Commands and


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

func (me *GearConfig) GetCommand(cmd []string) []string {

	for range only.Once {
		switch {
			case len(cmd) == 0:
				cmd = []string{ defaults.DefaultCommandName }

			case cmd[0] == "":
				cmd = []string{ defaults.DefaultCommandName }
		}

		var c string
		var ok bool
		if c, ok = me.Run.Commands[cmd[0]]; !ok {
			cmd = []string{ "" }
			break
		}

		cmd = append([]string{c}, cmd...)
	}

	return cmd
}

//func (me *GearCommands) Join() string {
//	var c string
//
//	for range only.Once {
//		c = strings.
//	}
//
//	return c
//}


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
