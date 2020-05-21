package gear

import (
	"github.com/docker/docker/client"
	//"github.com/docker/docker/integration-cli/cli"
	"launch/dockerClient"
	"launch/gear/gearJson"
	"launch/githubClient"
	"launch/ux"
	"strings"
)
// DOCKER_HOST=tcp://macpro:2376


type Gear struct {
	Repo         *githubClient.GitHubRepo
	Docker       *dockerClient.DockerGear
	GearConfig   *gearJson.GearConfig

	Debug        bool
	State        *ux.State
}


func (gear *Gear) NewGear(debugMode bool) *ux.State {
	for range OnlyOnce {
		gear.State = ux.NewState(debugMode)
		gear.State.DebugSet(debugMode)
		gear.Debug = debugMode

		if gear.Docker == nil {
			gear.Docker, gear.State = dockerClient.New(debugMode)
			if gear.State.IsError() {
				gear.State.SetError("can not connect to Docker service provider")
				break
			}
		}

		if gear.Repo == nil {
			gear.Repo = githubClient.New(debugMode)
			if gear.Repo.State.IsError() {
				break
			}

			gear.State.ClearError()
			//if state.IsError() {
			//	break
			//}
		}
	}

	return gear.State
}


func (gear *Gear) Status() *ux.State {
	if state := gear.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
		gear.State = gear.Docker.Container.Status()
		if gear.State.IsError() {
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

		state2 := gear.Docker.Image.Status()
		if state2.IsError() {
			break
		}

		//state = runState

		//state = gear.Docker.Image.State()
		//if state.IsError() {
		//	break
		//}
	}

	return gear.State
}


func (gear *Gear) FindContainer(gearName string, gearVersion string) (bool, *ux.State) {
	var found bool
	if state := gear.IsNil(); state.IsError() {
		return false, state
	}

	for range OnlyOnce {
		found, gear.State = gear.Docker.FindContainer(gearName, gearVersion)
		if !found {
			break
		}
		if gear.State.IsError() {
			break
		}

		gear.State = gear.Status()
		if gear.State.IsError() {
			break
		}
	}

	return found, gear.State
}


func (gear *Gear) FindImage(gearName string, gearVersion string) (bool, *ux.State) {
	var found bool
	if state := gear.IsNil(); state.IsError() {
		return false, state
	}

	for range OnlyOnce {
		found, gear.State = gear.Docker.FindImage(gearName, gearVersion)
		if !found {
			//state.ClearError()
			break
		}
		if gear.State.IsError() {
			break
		}

		//@TODO - TO CHECK
		//state = gear.Status()
		//if state.IsError() {
		//	break
		//}
	}

	return found, gear.State
}


func (gear *Gear) DecodeError(err error) (bool, *ux.State) {
	var ok bool
	if state := gear.IsNil(); state.IsError() {
		return false, state
	}

	for range OnlyOnce {
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

	return ok, gear.State
}

func (gear *Gear) IsNil() *ux.State {
	if state := ux.IfNilReturnError(gear); state.IsError() {
		return state
	}

	for range OnlyOnce {
		gear.State = gear.State.EnsureNotNil()

		gear.State = gear.Docker.IsNil()
		if gear.State.IsNotOk() {
			break
		}

		gear.State = gear.Repo.IsNil()
		if gear.State.IsNotOk() {
			break
		}
	}

	return gear.State
}

func (gear *Gear) IsValid() *ux.State {
	if state := ux.IfNilReturnError(gear); state.IsError() {
		return state
	}

	for range OnlyOnce {
		gear.State = gear.State.EnsureNotNil()
	}

	return gear.State
}
