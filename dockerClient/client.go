package dockerClient

import (
	"errors"
	"gb-launch/only"
	"github.com/docker/docker/client"
)


type DockerGear struct {
	Image     Image
	Container Container

	Client    *client.Client
	Ssh       *SSH

	Debug     bool
}


func New(d bool) (*DockerGear, error) {
	var cli DockerGear
	var err error

	for range only.Once {
		cli.Debug = d

		cli.Client, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		//cli.DockerClient, err = client.NewEnvClient()
		if err != nil {
			break
		}

		cli.Image._Parent = &cli
		cli.Container._Parent = &cli
	}

	return &cli, err
}


func (gear *DockerGear) EnsureNotNil() error {
	var err error

	for range only.Once {
		if gear == nil {
			err = errors.New("gear is nil")
			break
		}

		if gear.Client == nil {
			err = errors.New("docker client is nil")
			break
		}
	}

	return err
}
