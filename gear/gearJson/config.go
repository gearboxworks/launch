package gearJson

import (
	"encoding/json"
	"errors"
	"gb-launch/only"
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
type GearPorts []string
type GearArgs string

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


func New(cs string) (*GearConfig, error) {
	var gc GearConfig
	var err error

	for range only.Once {
		if cs == "" {
			err = errors.New("gear config is nil")
			break
		}

		js := []byte(cs)
		if js == nil {
			err = errors.New("gear config is nil")
			break
		}

		err = json.Unmarshal(js, &gc)
		if err != nil {
			err = errors.New("gearbox.json schema unknown")
			break
		}

		err = gc.ValidateGearConfig()
		if err != nil {
			break
		}
	}

	return &gc, err
}

func (me *GearConfig) ValidateGearConfig() error {
	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("gear config is nil")
			break
		}
	}

	return err
}
