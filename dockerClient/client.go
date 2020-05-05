package dockerClient

import (
	"context"
	"launch/defaults"
	"launch/only"
	"launch/ux"
	"github.com/docker/docker/client"
)


type DockerGear struct {
	Image     Image
	Container Container
	//RunAs     *defaults.ExecCommand

	Client    *client.Client
	Ssh       *Ssh

	Debug     bool
}


func New() (*DockerGear, ux.State) {
	var cli DockerGear
	var state ux.State

	for range only.Once {
		//cli.Debug = d

		var err error
		cli.Client, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		//cli.DockerClient, err = client.NewEnvClient()
		if err != nil {
			state.SetError("Docker client error: %s", err)
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		//var result types.Ping
		_, err = cli.Client.Ping(ctx)
		if err != nil {
			state.SetError("Docker client error: %s", err)
			break
		}
		//fmt.Printf("PING: %v", result)

		cli.Image._Parent = &cli
		cli.Container._Parent = &cli
	}

	return &cli, state
}


func (gear *DockerGear) EnsureNotNil() ux.State {
	var state ux.State

	for range only.Once {
		if gear == nil {
			state.SetError("gear is nil")
			break
		}

		if gear.Client == nil {
			state.SetError("docker client is nil")
			break
		}
	}

	return state
}

func (gear *DockerGear) SetSshStatusLine(s bool) {
	gear.Ssh.StatusLine.Enable = s
}

func (gear *DockerGear) SetSshShell(s bool) {
	gear.Ssh.Shell = s
}
