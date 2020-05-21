package dockerClient

import (
	"context"
	"github.com/docker/docker/client"
	"launch/defaults"
	"launch/ux"
)


type DockerGear struct {
	Image     Image
	Container Container
	//RunAs     *defaults.ExecCommand

	Client    *client.Client
	Ssh       *Ssh

	Debug     bool
	State     *ux.State
}

func New(debugMode bool) (*DockerGear, *ux.State) {
	var gear DockerGear

	for range OnlyOnce {
		gear.State = gear.State.EnsureNotNil()
		gear.State.DebugSet(debugMode)
		gear.Debug = debugMode

		gear.Image = *gear.Image.EnsureNotNil()
		gear.State.DebugSet(debugMode)

		gear.Container = *gear.Container.EnsureNotNil()
		gear.State.DebugSet(debugMode)

		var err error
		gear.Client, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		//cli.DockerClient, err = client.NewEnvClient()
		if err != nil {
			gear.State.SetError("Docker client error: %s", err)
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		//var result types.Ping
		_, err = gear.Client.Ping(ctx)
		if err != nil {
			gear.State.SetError("Docker client error: %s", err)
			break
		}
		//fmt.Printf("PING: %v", result)

		gear.Image._Parent = &gear
		gear.Container._Parent = &gear
	}

	return &gear, gear.State
}


func (gear *DockerGear) IsNil() *ux.State {
	if state := ux.IfNilReturnError(gear); state.IsError() {
		return state
	}
	gear.State = gear.State.EnsureNotNil()
	return gear.State
}

func (gear *DockerGear) IsValid() *ux.State {
	if state := ux.IfNilReturnError(gear); state.IsError() {
		return state
	}

	for range OnlyOnce {
		gear.State = gear.State.EnsureNotNil()

		if gear.Client == nil {
			gear.State.SetError("docker client is nil")
			break
		}
	}

	return gear.State
}


func (gear *DockerGear) SetSshStatusLine(s bool) {
	gear.Ssh.StatusLine.Enable = s
}

func (gear *DockerGear) SetSshShell(s bool) {
	gear.Ssh.Shell = s
}
