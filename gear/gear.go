package gear

import (
	"errors"
	"gb-launch/dockerClient"
	"gb-launch/gear/gearJson"
	"gb-launch/githubClient"
	"gb-launch/only"
	"github.com/docker/docker/client"
)
// DOCKER_HOST=tcp://macpro:2376


type Gear struct {
	Repo         *githubClient.GitHubRepo
	Docker       *dockerClient.DockerGear
	GearConfig   *gearJson.GearConfig
}


func NewGear(d bool) (*Gear, error) {
	var cli Gear
	var err error

	for range only.Once {
		cli.Docker, err = dockerClient.New(d)
		if err != nil {
			break
		}

		cli.Repo, err = githubClient.New()
		if err != nil {
			break
		}
	}

	return &cli, err
}


func (gear *Gear) DecodeError(err error) (bool, error) {
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

func (gear *Gear) EnsureNotNil() error {
	var err error

	for range only.Once {
		if gear == nil {
			err = errors.New("gear is nil")
			break
		}

		err = gear.Docker.EnsureNotNil()
		if err != nil {
			break
		}

		err = gear.Repo.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return err
}
