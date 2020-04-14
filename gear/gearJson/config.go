package gearJson

import (
	"encoding/json"
	"fmt"
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

	//Args         GearArgs     `json:"args"`
	//Base         bool         `json:"base"`
	//Class        string       `json:"class"`
	//Env          GearEnv      `json:"env"`
	//Maintainer   string       `json:"maintainer"`
	//Name         string       `json:"name"`
	//Network      string       `json:"network"`
	//Organization string       `json:"organization"`
	//Ports        GearPorts    `json:"ports"`
	//Refurl       string       `json:"refurl"`
	//Restart      string       `json:"restart"`
	//Run          string       `json:"run"`
	//State        string       `json:"state"`
	//Versions     GearVersions `json:"versions"`
	//Volumes      string       `json:"volumes"`
}
type GearConfigs map[string]GearConfig

type GearMeta struct {
	State        string `json:"state"`
	Organization string `json:"organization"`
	Name         string `json:"name"`
	Maintainer   string `json:"maintainer"`
	Class        string `json:"class"`
	Refurl       string `json:"refurl"`
}

type GearBuild struct {
	Ports        GearPorts    `json:"ports"`
	Run          string       `json:"run"`
	Args         GearArgs     `json:"args"`
	Env          GearEnv      `json:"env"`
	Network      string       `json:"network"`
}

type GearRun struct {
	Ports        GearPorts    `json:"ports"`
	Env          GearEnv      `json:"env"`
	Volumes      string       `json:"volumes"`
	Network      string       `json:"network"`
}

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
		p = fmt.Sprintf("%s %s:%s", p, k, v)
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
