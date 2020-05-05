package gear

import (
	"github.com/docker/docker/client"
	//"github.com/docker/docker/integration-cli/cli"
	"launch/dockerClient"
	"launch/gear/gearJson"
	"launch/githubClient"
	"launch/only"
	"launch/ux"
	"strings"
)
// DOCKER_HOST=tcp://macpro:2376


type Gear struct {
	Repo         *githubClient.GitHubRepo
	Docker       *dockerClient.DockerGear
	GearConfig   *gearJson.GearConfig
	Debug        bool
}


func (me *Gear) NewGear() ux.State {
	//var g Gear
	var state ux.State

	for range only.Once {
		if me.Docker == nil {
			me.Docker, state = dockerClient.New()
			if state.IsError() {
				state.SetError("can not connect to Docker service provider")
				break
			}
			me.Docker.Debug = me.Debug
		}

		if me.Repo == nil {
			me.Repo, state = githubClient.New()
			state.ClearError()
			//if state.IsError() {
			//	break
			//}
		}
	}

	return state
}


func (gear *Gear) State() ux.State {
	var state ux.State

	for range only.Once {
		state = gear.EnsureNotNil()
		if state.IsError() {
			break
		}

		runState := gear.Docker.Container.State()
		if state.IsError() {
			break
		}

		if gear.GearConfig == nil {
			if gear.Docker.Container.GearConfig != nil {
				gear.GearConfig = gear.Docker.Container.GearConfig
			} else if gear.Docker.Image.GearConfig != nil {
				gear.GearConfig = gear.Docker.Image.GearConfig
			}

		}

		if gear.Docker.Image.ID == "" {
			gear.Docker.Image.ID = strings.TrimPrefix(gear.Docker.Container.Details.Image, "sha256:")
			gear.Docker.Image.Name = gear.Docker.Container.Name
			gear.Docker.Image.Version = gear.Docker.Container.Version
		}

		state = gear.Docker.Image.State()
		if state.IsError() {
			break
		}

		state = runState

		//state = gear.Docker.Image.State()
		//if state.IsError() {
		//	break
		//}
	}

	return state
}


func (gear *Gear) FindContainer(gearName string, gearVersion string) (bool, ux.State) {
	var found bool
	var state ux.State

	for range only.Once {
		found, state = gear.Docker.FindContainer(gearName, gearVersion)
		if !found {
			break
		}
		if state.IsError() {
			break
		}

		state = gear.State()
		if state.IsError() {
			break
		}
	}

	return found, state
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
