package dockerClient

import (
	"context"
	"errors"
	"fmt"
	"launch/defaults"
	"launch/gear/gearJson"
	"launch/only"
	"launch/ux"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"io"
	"os"
	"time"
)

type Container struct {
	ID           string
	Name         string
	Version      string

	Summary      *types.Container
	Details      *types.ContainerJSON
	GearConfig   *gearJson.GearConfig

	_Parent       *DockerGear
}
type Containers []Container


func (me *Container) EnsureNotNil() ux.State {
	var state ux.State

	for range only.Once {
		if me == nil {
			state.SetError("gear is nil")
			break
		}

		if me.ID == "" {
			state.SetError("gear ID is nil")
			break
		}

		if me.Name == "" {
			state.SetError("gear name is nil")
			break
		}

		if me.Version == "" {
			state.SetError("gear version is nil")
			break
		}

		if me._Parent.Client == nil {
			state.SetError("docker client is nil")
			break
		}

		// if me.ctx == nil {
		// 	state = errors.New("ctx is nil")
		// 	break
		// }
		//
		// if me.Summary == nil {
		// 	state = errors.New("container is nil")
		// 	break
		// }
		//
		// if me.Details == "" {
		// 	state = errors.New("container is nil")
		// 	break
		// }
	}

	return state
}


func (me *Container) State() ux.State {
	var state ux.State

	for range only.Once {
		var err error

		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		df := filters.NewArgs()
		df.Add("id", me.ID)

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		var containers []types.Container
		containers, err = me._Parent.Client.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: df})
		if err != nil {
			state.SetError("gear list error: %s", err)
			break
		}
		if len(containers) == 0 {
			state.SetWarning("no gears found")
			break
		}

		me.Summary = &containers[0]

		me.GearConfig, state = gearJson.New(me.Summary.Labels["gearbox.json"])
		if state.IsError() {
			break
		}

		if me.GearConfig.Meta.Organization != defaults.Organization {
			state.SetError("not a Gearbox container")
			break
		}

		d := types.ContainerJSON{}
		d, err = me._Parent.Client.ContainerInspect(ctx, me.ID)
		if err != nil {
			state.SetError("gear inspect error: %s", err)
			break
		}
		me.Details = &d

		state.SetString(me.Details.State.Status)
	}

	if state.IsError() {
		me.Summary = nil
		me.Details = nil
	}

	return state
}


func (me *Container) WaitForState(s string, t time.Duration) ux.State {
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		until := time.Now()
		until.Add(t)

		for now := time.Now(); until.Before(now); now = time.Now() {
			state = me.State()
			if state.IsError() {
				break
			}

			if state.String == s {
				break
			}
		}
	}

	return state
}


// Run a container in the background
// You can also run containers in the background, the equivalent of typing docker run -d bfirsh/reticulate-splines:
func (me *Container) Start() ux.State {
	var state ux.State

	for range only.Once {
		var err error

		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		state = me.State()
		if state.IsError() {
			break
		}

		if state.IsRunning() {
			break
		}

		if !state.IsCreated() && !state.IsExited() {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		err = me._Parent.Client.ContainerStart(ctx, me.ID, types.ContainerStartOptions{})
		if err != nil {
			state.SetError("gear start error: %s", err)
			break
		}

		statusCh, errCh := me._Parent.Client.ContainerWait(ctx, me.ID, "") // container.WaitConditionNotRunning
		select {
			case err := <-errCh:
				if err != nil {
					state.SetError("Docker client error: %s", err)
					// fmt.Printf("SC: %s\n", response.Error)
					// return false, err
				}
				break

			case status := <-statusCh:
				fmt.Printf("status.StatusCode: %#+v\n", status.StatusCode)
				break
		}
		// fmt.Printf("SC: %s\n", status)
		// fmt.Printf("SC: %s\n", err)

		state = me.WaitForState(ux.StateRunning, defaults.Timeout)
		if state.IsError() {
			break
		}
		if !state.IsRunning() {
			state.SetError("cannot start gear")
			break
		}

		// "created", "running", "paused", "restarting", "removing", "exited", or "dead"
		// out, err := me.DockerClient.ImagePull(me.Ctx, me.GearConfig.Name, types.ImagePullOptions{})
		// if err != nil {
		// 	break
		// }
		// _, _ = io.Copy(os.Stdout, out)
		//
		// _, err = me.DockerClient.ContainerCreate(me.Ctx, &container.Config{
		// 	Image: me.GearConfig.Name,
		// }, nil, nil, "")
		// if err != nil {
		// 	break
		// }
	}

	return state
}


// Run a container
// This first example shows how to run a container using the Docker API.
// On the command line, you would use the docker run command, but this is just as easy to do from your own apps too.
// This is the equivalent of typing docker run alpine echo hello world at the command prompt:
func (me *Container) ContainerCreate(gearName string, gearVersion string, gearMount string) ux.State {
	var state ux.State

	for range only.Once {
		if me._Parent.Debug {
			fmt.Printf("DEBUG: ContainerCreate(%s, %s, %s)\n", gearName, gearVersion, gearMount)
		}

		if gearName == "" {
			state.SetError("empty gearname")
			break
		}

		if gearVersion == "" {
			gearVersion = "latest"
		}

		var ok bool
		//var err error
		ok, state = me._Parent.FindContainer(gearName, gearVersion)
		if state.IsError() {
			break
		}
		if !ok {
			//response.Error = me.Search(gearName, gearVersion)

			ok, state = me._Parent.FindImage(gearName, gearVersion)
			if (state.IsError()) || (!ok) {
				me._Parent.Image.ID = gearName
				me._Parent.Image.Name = gearName
				me._Parent.Image.Version = gearVersion
				state = me._Parent.Image.Pull()
				if state.IsError() {
					state.SetError("no such gear '%s'", gearName)
					break
				}
			}

			ok, state = me._Parent.FindContainer(gearName, gearVersion)
			if state.IsError() {
				state.SetError("error creating gear: %s", state.Error)
				break
			}
		}

		me.ID = me._Parent.Image.ID
		me.Name = me._Parent.Image.Name
		me.Version = me._Parent.Image.Version

		// me.Image.Details.Container = "gearboxworks/golang:1.14"
		// tag := fmt.Sprintf("", me.Image.Name, me.Image.Version)
		tag := fmt.Sprintf("gearboxworks/%s:%s", me.Name, me.Version)
		gn := fmt.Sprintf("%s-%s", me.Name, me.Version)
		var binds []string
		if gearMount != "" {
			binds = append(binds, fmt.Sprintf("%s:%s", gearMount, defaults.DefaultProject))
		}

		config := container.Config {
			// Hostname:        "",
			// Domainname:      "",
			User:            "root",
			// AttachStdin:     false,
			AttachStdout:    true,
			AttachStderr:    true,
			ExposedPorts:    nil,
			Tty:             false,
			OpenStdin:       false,
			StdinOnce:       false,
			Env:             nil,
			Cmd:             []string{"/init"},
			// Healthcheck:     nil,
			// ArgsEscaped:     false,
			Image:           tag,
			// Volumes:         nil,
			// WorkingDir:      "",
			// Entrypoint:      nil,
			// NetworkDisabled: false,
			// MacAddress:      "",
			// OnBuild:         nil,
			// Labels:          nil,
			// StopSignal:      "",
			// StopTimeout:     nil,
			// Shell:           nil,
		}

		netConfig := network.NetworkingConfig {}

		// DockerMount
		// ms := mount.Mount {
		// 	Type:          "bind",
		// 	Source:        "/Users/mick/Documents/GitHub/containers/docker-golang",
		// 	Target:        "/foo",
		// 	ReadOnly:      false,
		// 	Consistency:   "",
		// 	BindOptions:   nil,
		// 	VolumeOptions: nil,
		// 	TmpfsOptions:  nil,
		// }

		hostConfig := container.HostConfig {
			Binds:           binds,
			ContainerIDFile: "",
			LogConfig:       container.LogConfig{
				Type:   "",
				Config: nil,
			},
			NetworkMode:     defaults.GearboxNetwork,
			PortBindings:    nil,						// @TODO
			RestartPolicy:   container.RestartPolicy {
				Name:              "",
				MaximumRetryCount: 0,
			},
			AutoRemove:      false,
			VolumeDriver:    "",
			VolumesFrom:     nil,
			CapAdd:          nil,
			CapDrop:         nil,
			//Capabilities:    nil,
			//CgroupnsMode:    "",
			DNS:             []string{},
			DNSOptions:      []string{},
			DNSSearch:       []string{},
			ExtraHosts:      nil,
			GroupAdd:        nil,
			IpcMode:         "",
			Cgroup:          "",
			Links:           nil,
			OomScoreAdj:     0,
			PidMode:         "",
			Privileged:      false,
			PublishAllPorts: true,
			ReadonlyRootfs:  false,
			SecurityOpt:     nil,
			StorageOpt:      nil,
			Tmpfs:           nil,
			UTSMode:         "",
			UsernsMode:      "",
			ShmSize:         0,
			Sysctls:         nil,
			Runtime:         "runc",
			ConsoleSize:     [2]uint{},
			Isolation:       "",
			Resources:       container.Resources{},
			Mounts:          []mount.Mount{},
			//MaskedPaths:     nil,
			//ReadonlyPaths:   nil,
			Init:            nil,
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		var resp container.ContainerCreateCreatedBody
		var err error
		resp, err = me._Parent.Client.ContainerCreate(ctx, &config, &hostConfig, &netConfig, gn)
		if err != nil {
			state.SetError("error creating gear: %s", err)
			break
		}

		if resp.ID == "" {
			break
		}

		me.ID = resp.ID
		//me.Container.Name = me.Image.Name
		//me.Container.Version = me.Image.Version

		// var response Response
		state = me.State()
		if state.IsError() {
			break
		}

		if state.IsCreated() {
			break
		}

		//if state.IsRunning() {
		//	break
		//}
		//
		//if state.IsPaused() {
		//	break
		//}
		//
		//if state.IsRestarting() {
		//	break
		//}
	}

	if me._Parent.Debug {
		state.Print()
	}

	return state
}


// Stop all running containers
// Now that you know what containers exist, you can perform operations on them.
// This example stops all running containers.
func (me *Container) Stop() ux.State {
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		err := me._Parent.Client.ContainerStop(ctx, me.ID, nil)
		if err != nil {
			state.SetError("gear stop error: %s", err)
			break
		}

		state = me.State()
	}

	return state
}


// Remove containers
// Now that you know what containers exist, you can perform operations on them.
// This example stops all running containers.
func (me *Container) Remove() ux.State {
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		options := types.ContainerRemoveOptions{
			RemoveVolumes: false,
			RemoveLinks:   false,
			Force:         false,
		}

		err := me._Parent.Client.ContainerRemove(ctx, me.ID, options)
		if err != nil {
			state.SetError("gear remove error: %s", err)
			break
		}

		//state = me.State()
		state.SetOk("OK")
	}

	return state
}


// Print the logs of a specific container
// You can also perform actions on individual containers.
// This example prints the logs of a container given its ID.
// You need to modify the code before running it to change the hard-coded ID of the container to print the logs for.
func (me *Container) Logs() ux.State {
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		options := types.ContainerLogsOptions{ShowStdout: true}
		// Replace this ID with a container that really exists
		out, err := me._Parent.Client.ContainerLogs(ctx, me.ID, options)
		if err != nil {
			state.SetError("gear logs error: %s", err)
			break
		}

		_, _ = io.Copy(os.Stdout, out)

		state = me.State()
	}

	return state
}


// Commit a container
// Commit a container to create an image from its contents:
func (me *Container) Commit() ux.State {
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		createResp, err := me._Parent.Client.ContainerCreate(ctx, &container.Config{
			Image: "alpine",
			Cmd:   []string{"touch", "/helloworld"},
		}, nil, nil, "")
		if err != nil {
			break
		}

		if err := me._Parent.Client.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{}); err != nil {
			break
		}

		statusCh, errCh := me._Parent.Client.ContainerWait(ctx, createResp.ID, container.WaitConditionNotRunning)
		select {
			case err := <-errCh:
				if err != nil {
					//response.State.SetError("gear stop error: %s", err)
					break
				}
			case <-statusCh:
		}

		commitResp, err := me._Parent.Client.ContainerCommit(ctx, createResp.ID, types.ContainerCommitOptions{Reference: "helloworld"})
		if err != nil {
			break
		}

		fmt.Println(commitResp.ID)
	}

	return state
}


func findSshPort(c types.Container) error {
	var err error
	err = errors.New("no container ssh")

	for range only.Once {
		for _, p := range c.Ports {
			if p.PrivatePort == 22 {
				err = nil
				break
			}
		}
	}

	return err
}

func (me *Container) GetContainerSsh() (string, ux.State) {
	var port string
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		var found bool
		for _, p := range me.Summary.Ports {
			if p.PrivatePort == 22 {
				port = fmt.Sprintf("%d", p.PublicPort)
				found = true
				break
			}
		}

		if !found {
			state.SetError("no SSH port")
		}
	}

	return port, state
}

// func (me *Gear) GetContainerSsh(f string) (string, error) {
// 	var port string
// 	var err error
//
// 	for range only.Once {
// 		err = me.EnsureNotNil()
// 		if err != nil {
// 			break
// 		}
//
// 		var containers []types.Container
// 		containers, err = me.DockerClient.ContainerList(me.Ctx, types.ContainerListOptions{All: true})
// 		if err != nil {
// 			break
// 		}
//
// 		for _, c := range containers {
// 			var gc *GearConfig
// 			gc, err = NewGearConfig(c.Labels["gearbox.json"])
// 			if err != nil {
// 				continue
// 			}
//
// 			if gc.Organization != Organization {
// 				continue
// 			}
//
// 			if gc.Name != f {
// 				continue
// 			}
//
// 			err = findSshPort(c)
// 			if err != nil {
// 				continue
// 			}
// 			// for _, p := range c.Ports {
// 			// 	if p.PrivatePort == 22 {
// 			// 		port = fmt.Sprintf("%d", p.PublicPort)
// 			// 		break
// 			// 	}
// 			// }
// 			err = nil
// 			break
// 		}
// 	}
//
// 	return port, err
// }
