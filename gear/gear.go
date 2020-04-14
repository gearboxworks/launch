package gear

import (
	"launch/dockerClient"
	"launch/gear/gearJson"
	"launch/githubClient"
	"launch/only"
	"launch/ux"
	"github.com/docker/docker/client"
)
// DOCKER_HOST=tcp://macpro:2376


type Gear struct {
	Repo         *githubClient.GitHubRepo
	Docker       *dockerClient.DockerGear
	GearConfig   *gearJson.GearConfig
}


func NewGear(d bool) (*Gear, ux.State) {
	var cli Gear
	var state ux.State

	for range only.Once {
		cli.Docker, state = dockerClient.New(d)
		if state.IsError() {
			state.SetError("can not connect to Docker service")
			break
		}

		cli.Repo, state = githubClient.New()
		state.ClearError()
		//if state.IsError() {
		//	break
		//}
	}

	return &cli, state
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

func (gear *Gear) EnsureNotNil() ux.State {
	var state ux.State

	for range only.Once {
		if gear == nil {
			state.SetError("gear is nil")
			break
		}

		state = gear.Docker.EnsureNotNil()
		if state.IsError() {
			break
		}

		state = gear.Repo.EnsureNotNil()
		if state.IsError() {
			break
		}
	}

	return state
}
