package dockerClient

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"io"
	"launch/defaults"
	"launch/gear/gearJson"
	"launch/only"
	"launch/ux"
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


func (c *Container) EnsureNotNil() ux.State {
	var state ux.State

	for range only.Once {
		if c == nil {
			state.SetError("gear is nil")
			break
		}

		if c.ID == "" {
			state.SetError("gear ID is nil")
			break
		}

		if c.Name == "" {
			state.SetError("gear name is nil")
			break
		}

		if c.Version == "" {
			state.SetError("gear version is nil")
			break
		}

		if c._Parent.Client == nil {
			state.SetError("docker client is nil")
			break
		}

		// if c.ctx == nil {
		// 	state = errors.New("ctx is nil")
		// 	break
		// }
		//
		// if c.Summary == nil {
		// 	state = errors.New("container is nil")
		// 	break
		// }
		//
		// if c.Details == "" {
		// 	state = errors.New("container is nil")
		// 	break
		// }
	}

	return state
}


func (c *Container) State() ux.State {
	var state ux.State

	for range only.Once {
		var err error

		state = c.EnsureNotNil()
		if state.IsError() {
			break
		}

		df := filters.NewArgs()
		df.Add("id", c.ID)

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		var containers []types.Container
		containers, err = c._Parent.Client.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: df})
		if err != nil {
			state.SetError("gear list error: %s", err)
			break
		}
		if len(containers) == 0 {
			state.SetWarning("no gears found")
			break
		}

		c.Summary = &containers[0]

		c.GearConfig, state = gearJson.New(c.Summary.Labels["gearbox.json"])
		if state.IsError() {
			break
		}

		if c.GearConfig.Meta.Organization != defaults.Organization {
			state.SetError("not a Gearbox container")
			break
		}

		d := types.ContainerJSON{}
		d, err = c._Parent.Client.ContainerInspect(ctx, c.ID)
		if err != nil {
			state.SetError("gear inspect error: %s", err)
			break
		}
		c.Details = &d

		state.SetString(c.Details.State.Status)
	}

	if state.IsError() {
		c.Summary = nil
		c.Details = nil
	}

	return state
}


func (c *Container) WaitForState(s string, t time.Duration) ux.State {
	var state ux.State

	for range only.Once {
		state = c.EnsureNotNil()
		if state.IsError() {
			break
		}

		until := time.Now()
		until.Add(t)

		for now := time.Now(); until.Before(now); now = time.Now() {
			state = c.State()
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
func (c *Container) Start() ux.State {
	var state ux.State

	for range only.Once {
		var err error

		state = c.EnsureNotNil()
		if state.IsError() {
			break
		}

		state = c.State()
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
		//noinspection GoDeferInLoop
		defer cancel()

		err = c._Parent.Client.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
		if err != nil {
			state.SetError("gear start error: %s", err)
			break
		}

		statusCh, errCh := c._Parent.Client.ContainerWait(ctx, c.ID, "") // container.WaitConditionNotRunning
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

		state = c.WaitForState(ux.StateRunning, defaults.Timeout)
		if state.IsError() {
			break
		}
		if !state.IsRunning() {
			state.SetError("cannot start gear")
			break
		}

		// "created", "running", "paused", "restarting", "removing", "exited", or "dead"
		// out, err := c.DockerClient.ImagePull(c.Ctx, c.GearConfig.Name, types.ImagePullOptions{})
		// if err != nil {
		// 	break
		// }
		// _, _ = io.Copy(os.Stdout, out)
		//
		// _, err = c.DockerClient.ContainerCreate(c.Ctx, &container.Config{
		// 	Image: c.GearConfig.Name,
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
func (c *Container) ContainerCreate(gearName string, gearVersion string, gearMount string) ux.State {
	var state ux.State

	for range only.Once {
		if c._Parent.Debug {
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
		ok, state = c._Parent.FindContainer(gearName, gearVersion)
		if state.IsError() {
			break
		}
		if !ok {
			//response.Error = c.Search(gearName, gearVersion)

			ok, state = c._Parent.FindImage(gearName, gearVersion)
			if (state.IsError()) || (!ok) {
				c._Parent.Image.ID = gearName
				c._Parent.Image.Name = gearName
				c._Parent.Image.Version = gearVersion
				state = c._Parent.Image.Pull()
				if state.IsError() {
					state.SetError("no such gear '%s'", gearName)
					break
				}
			}
			state.ClearAll()

			ok, state = c._Parent.FindContainer(gearName, gearVersion)
			if state.IsError() {
				state.SetError("error creating gear: %s", state.Error)
				break
			}
		}

		c.ID = c._Parent.Image.ID
		c.Name = c._Parent.Image.Name
		c.Version = c._Parent.Image.Version

		// c.Image.Details.Container = "gearboxworks/golang:1.14"
		// tag := fmt.Sprintf("", c.Image.Name, c.Image.Version)
		tag := fmt.Sprintf("gearboxworks/%s:%s", c.Name, c.Version)
		gn := fmt.Sprintf("%s-%s", c.Name, c.Version)
		var binds []string
		if gearMount != defaults.DefaultPathNone {
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
			Privileged:      true,
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
		//noinspection GoDeferInLoop
		defer cancel()

		var resp container.ContainerCreateCreatedBody
		var err error
		resp, err = c._Parent.Client.ContainerCreate(ctx, &config, &hostConfig, &netConfig, gn)
		if err != nil {
			state.SetError("error creating gear: %s", err)
			break
		}

		if resp.ID == "" {
			break
		}

		c.ID = resp.ID
		//c.Container.Name = c.Image.Name
		//c.Container.Version = c.Image.Version

		// var response Response
		state = c.State()
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

	if c._Parent.Debug {
		state.Print()
	}

	return state
}


// Stop all running containers
// Now that you know what containers exist, you can perform operations on them.
// This example stops all running containers.
func (c *Container) Stop() ux.State {
	var state ux.State

	for range only.Once {
		state = c.EnsureNotNil()
		if state.IsError() {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		err := c._Parent.Client.ContainerStop(ctx, c.ID, nil)
		if err != nil {
			state.SetError("gear stop error: %s", err)
			break
		}

		state = c.State()
	}

	return state
}


// Remove containers
// Now that you know what containers exist, you can perform operations on them.
// This example stops all running containers.
func (c *Container) Remove() ux.State {
	var state ux.State

	for range only.Once {
		state = c.EnsureNotNil()
		if state.IsError() {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		options := types.ContainerRemoveOptions{
			RemoveVolumes: false,
			RemoveLinks:   false,
			Force:         false,
		}

		err := c._Parent.Client.ContainerRemove(ctx, c.ID, options)
		if err != nil {
			state.SetError("gear remove error: %s", err)
			break
		}

		//state = c.State()
		state.SetOk("OK")
	}

	return state
}


// Print the logs of a specific container
// You can also perform actions on individual containers.
// This example prints the logs of a container given its ID.
// You need to modify the code before running it to change the hard-coded ID of the container to print the logs for.
func (c *Container) Logs() ux.State {
	var state ux.State

	for range only.Once {
		state = c.EnsureNotNil()
		if state.IsError() {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		options := types.ContainerLogsOptions{ShowStdout: true}
		// Replace this ID with a container that really exists
		out, err := c._Parent.Client.ContainerLogs(ctx, c.ID, options)
		if err != nil {
			state.SetError("gear logs error: %s", err)
			break
		}

		_, _ = io.Copy(os.Stdout, out)

		state = c.State()
	}

	return state
}


// Commit a container
// Commit a container to create an image from its contents:
func (c *Container) Commit() ux.State {
	var state ux.State

	for range only.Once {
		state = c.EnsureNotNil()
		if state.IsError() {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		createResp, err := c._Parent.Client.ContainerCreate(ctx, &container.Config{
			Image: "alpine",
			Cmd:   []string{"touch", "/helloworld"},
		}, nil, nil, "")
		if err != nil {
			break
		}

		if err := c._Parent.Client.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{}); err != nil {
			break
		}

		statusCh, errCh := c._Parent.Client.ContainerWait(ctx, createResp.ID, container.WaitConditionNotRunning)
		select {
			case err := <-errCh:
				if err != nil {
					//response.State.SetError("gear stop error: %s", err)
					break
				}
			case <-statusCh:
		}

		commitResp, err := c._Parent.Client.ContainerCommit(ctx, createResp.ID, types.ContainerCommitOptions{Reference: "helloworld"})
		if err != nil {
			break
		}

		fmt.Println(commitResp.ID)
	}

	return state
}


//func findSshPort(c types.Container) error {
//	var err error
//	err = errors.New("no container ssh")
//
//	for range only.Once {
//		for _, p := range c.Ports {
//			if p.PrivatePort == 22 {
//				err = nil
//				break
//			}
//		}
//	}
//
//	return err
//}


func (c *Container) GetContainerSsh() (string, ux.State) {
	var port string
	var state ux.State

	for range only.Once {
		state = c.EnsureNotNil()
		if state.IsError() {
			break
		}

		var found bool
		for _, p := range c.Summary.Ports {
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
