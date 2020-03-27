package dockerClient

import (
	"gb-launch/only"
	"errors"
	"github.com/docker/docker/client"
	"time"
)

// DOCKER_HOST=tcp://macpro:2376

const (
	Organization = "gearboxworks"
	Timeout = time.Second * 2
	DefaultProject = "/home/gearbox/projects/default"
)

type GearEnv map[string]string
type GearPorts []string
type GearArgs string

type GearVersion struct {
	Base         string `json:"base"`
	Latest       bool   `json:"latest"`
	MajorVersion string `json:"majorversion"`
	Ref          string `json:"ref"`
}
type GearVersions map[string]GearVersion

func (me *GearVersions) GetLatest() string {
	var v string

	var r GearVersion
	for v, r = range *me {
		if r.Latest {
			break
		}
	}

	return v
}

func (me *GearVersions) HasVersion(c string) bool {
	for v, r := range *me {
		if r.Latest && (c == "latest") {
			return true
		}

		if v == c {
			return true
		}
	}
	return false
}


type GearConfig struct {
	Args         GearArgs     `json:"args"`
	Base         bool         `json:"base"`
	Class        string       `json:"class"`
	Env          GearEnv      `json:"env"`
	Maintainer   string       `json:"maintainer"`
	Name         string       `json:"name"`
	Network      string       `json:"network"`
	Organization string       `json:"organization"`
	Ports        GearPorts    `json:"ports"`
	Refurl       string       `json:"refurl"`
	Restart      string       `json:"restart"`
	Run          string       `json:"run"`
	State        string       `json:"state"`
	Versions     GearVersions `json:"versions"`
	Volumes      string       `json:"volumes"`
}
type GearConfigs map[string]GearConfig

type Gear struct {
	Image        Image
	Container    Container
	DockerClient *client.Client
	SshClient    *SSH

	Debug        bool
}

func NewGear(d bool) (*Gear, error) {
	var cli Gear
	var err error

	for range only.Once {
		cli.Debug = d

		//cli.DockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		cli.DockerClient, err = client.NewEnvClient()
		if err != nil {
			break
		}

		cli.Container.client = cli.DockerClient
		cli.Image.client = cli.DockerClient
	}

	return &cli, err
}


func (g *Gear) DecodeError(err error) (bool, error) {
	var ok bool

	for range only.Once {
		switch {
			case err != nil:
				ok = true

			//case client.IsErrContainerNotFound(err):
			case client.IsErrConnectionFailed(err):
			case client.IsErrNotFound(err):
			case client.IsErrPluginPermissionDenied(err):
			case client.IsErrUnauthorized(err):
			default:
		}
	}

	return ok, err
}

func (g *Gear) EnsureNotNil() error {
	var err error

	for range only.Once {
		if g == nil {
			err = errors.New("gear is nil")
			break
		}

		if g.DockerClient == nil {
			err = errors.New("docker client is nil")
			break
		}
	}

	return err
}
