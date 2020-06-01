package helperDocker

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

	_Parent      *DockerGear
	Debug        bool
	State        *ux.State
}
type Containers []Container


func NewContainer(debugMode bool) *Container {
	c := &Container{
		ID:         "",
		Name:       "",
		Version:    "",
		Summary:    nil,
		Details:    nil,
		GearConfig: &gearJson.GearConfig{},
		_Parent:    nil,
		State:      ux.NewState(debugMode),
	}
	c.Debug = debugMode

	return c
}

func (i *Container) EnsureNotNil() *Container {
	for range OnlyOnce {
		if i == nil {
			i = NewContainer(false)
		}
		i.State = i.State.EnsureNotNil()
	}
	return i
}

func (c *Container) IsNil() *ux.State {
	if state := ux.IfNilReturnError(c); state.IsError() {
		return state
	}
	c.State = c.State.EnsureNotNil()
	return c.State
}

func (c *Container) IsValid() *ux.State {
	if state := ux.IfNilReturnError(c); state.IsError() {
		return state
	}

	for range OnlyOnce {
		c.State = c.State.EnsureNotNil()

		if c.ID == "" {
			c.State.SetError("gear ID is nil")
			break
		}

		if c.Name == "" {
			c.State.SetError("gear name is nil")
			break
		}

		if c.Version == "" {
			c.State.SetError("gear version is nil")
			break
		}

		if c._Parent.Client == nil {
			c.State.SetError("docker client is nil")
			break
		}
	}

	return c.State
}


func (c *Container) Status() *ux.State {
	if state := c.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
		df := filters.NewArgs()
		df.Add("id", c.ID)

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		var containers []types.Container
		var err error
		containers, err = c._Parent.Client.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: df})
		if err != nil {
			c.State.SetError("gear list error: %s", err)
			break
		}
		if len(containers) == 0 {
			c.State.SetWarning("no gears found")
			break
		}

		c.Summary = &containers[0]

		c.GearConfig = gearJson.New(c.Summary.Labels["gearbox.json"])
		if c.GearConfig.State.IsError() {
			c.State.SetState(c.GearConfig.State)
			break
		}

		if c.GearConfig.Meta.Organization != defaults.Organization {
			c.State.SetError("not a Gearbox container")
			break
		}

		d := types.ContainerJSON{}
		d, err = c._Parent.Client.ContainerInspect(ctx, c.ID)
		if err != nil {
			c.State.SetError("gear inspect error: %s", err)
			break
		}
		c.Details = &d

		c.State.SetRunState(c.Details.State.Status)
	}

	if c.State.IsError() {
		c.Summary = nil
		c.Details = nil
	}

	return c.State
}


func (c *Container) WaitForState(s string, t time.Duration) *ux.State {
	if state := c.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
		until := time.Now()
		until.Add(t)

		for now := time.Now(); until.Before(now); now = time.Now() {
			c.State = c.Status()
			if c.State.IsError() {
				break
			}

			if c.State.RunStateEquals(s) {
				break
			}
		}
	}

	return c.State
}


// Run a container in the background
// You can also run containers in the background, the equivalent of typing docker run -d bfirsh/reticulate-splines:
func (c *Container) Start() *ux.State {
	if state := c.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
		c.State = c.Status()
		if c.State.IsError() {
			break
		}

		if c.State.IsRunning() {
			break
		}

		if !c.State.IsCreated() && !c.State.IsExited() {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		err := c._Parent.Client.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
		if err != nil {
			c.State.SetError("Gear '%s:%s' start error - %s", c.Name, c.Version, err)
			break
		}

		//statusCh, errCh := c._Parent.Client.ContainerWait(ctx, c.ID, "") // container.WaitConditionNotRunning
		//select {
		//	case err := <-errCh:
		//		if err != nil {
		//			c.State.SetError("Docker client error: %s", err)
		//			// fmt.Printf("SC: %s\n", response.Error)
		//			// return false, err
		//		}
		//		break
		//
		//	case status := <-statusCh:
		//		fmt.Printf("status.StatusCode: %#+v\n", status.StatusCode)
		//		break
		//}
		// fmt.Printf("SC: %s\n", status)
		// fmt.Printf("SC: %s\n", err)

		c.State = c.WaitForState(ux.StateRunning, defaults.Timeout)
		if c.State.IsError() {
			break
		}
		if !c.State.IsRunning() {
			c.State.SetError("Gear '%s:%s' failed to start.", c.Name, c.Version)
			break
		}
	}

	return c.State
}


// Run a container
// This first example shows how to run a container using the Docker API.
// On the command line, you would use the docker run command, but this is just as easy to do from your own apps too.
// This is the equivalent of typing docker run alpine echo hello world at the command prompt:
func (c *Container) ContainerCreate(gearName string, gearVersion string, gearMount string) *ux.State {
	if state := c.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
		if c._Parent.Debug {
			fmt.Printf("DEBUG: ContainerCreate(%s, %s, %s)\n", gearName, gearVersion, gearMount)
		}

		if gearName == "" {
			c.State.SetError("empty gearname")
			break
		}

		if gearVersion == "" {
			gearVersion = "latest"
		}

		var ok bool
		ok, c.State = c._Parent.FindContainer(gearName, gearVersion)
		if c.State.IsError() {
			break
		}
		if !ok {
			// Find Gear image since we don't have a container.
			for range OnlyOnce {
				ok, c.State = c._Parent.FindImage(gearName, gearVersion)
				if c.State.IsError() {
					ok = false
					break
				}
				if ok {
					break
				}

				ux.PrintflnNormal("Downloading Gear '%s:%s'.", gearName, gearVersion)

				// Pull Gear image.
				c._Parent.Image.ID = gearName
				c._Parent.Image.Name = gearName
				c._Parent.Image.Version = gearVersion
				c.State = c._Parent.Image.Pull()
				if c.State.IsError() {
					c.State.SetError("no such gear '%s'", gearName)
					break
				}

				// Confirm it's there.
				ok, c.State = c._Parent.FindImage(gearName, gearVersion)
				if c.State.IsError() {
					ok = false
					break
				}
			}
			if !ok {
				c.State.SetError("Cannot install Gear image '%s:%s' - %s.", gearName, gearVersion, c.State.GetError())
				break
			}
			//c.State.Clear()
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
			c.State.SetError("error creating gear: %s", err)
			break
		}

		if resp.ID == "" {
			break
		}

		c.ID = resp.ID
		//c.Container.Name = c.Image.Name
		//c.Container.Version = c.Image.Version

		// var response Response
		c.State = c.Status()
		if c.State.IsError() {
			break
		}

		if c.State.IsCreated() {
			break
		}

		//if c.State.IsRunning() {
		//	break
		//}
		//
		//if c.State.IsPaused() {
		//	break
		//}
		//
		//if c.State.IsRestarting() {
		//	break
		//}
	}

	return c.State
}


// Stop all running containers
// Now that you know what containers exist, you can perform operations on them.
// This example stops all running containers.
func (c *Container) Stop() *ux.State {
	if state := c.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		err := c._Parent.Client.ContainerStop(ctx, c.ID, nil)
		if err != nil {
			c.State.SetError("gear stop error: %s", err)
			break
		}

		c.State = c.Status()
	}

	return c.State
}


// Remove containers
// Now that you know what containers exist, you can perform operations on them.
// This example stops all running containers.
func (c *Container) Remove() *ux.State {
	if state := c.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
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
			c.State.SetError("gear remove error: %s", err)
			break
		}

		//c.State = c.State()
		c.State.SetOk("OK")
	}

	return c.State
}


// Print the logs of a specific container
// You can also perform actions on individual containers.
// This example prints the logs of a container given its ID.
// You need to modify the code before running it to change the hard-coded ID of the container to print the logs for.
func (c *Container) Logs() *ux.State {
	if state := c.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		options := types.ContainerLogsOptions{ShowStdout: true}
		// Replace this ID with a container that really exists
		out, err := c._Parent.Client.ContainerLogs(ctx, c.ID, options)
		if err != nil {
			c.State.SetError("gear logs error: %s", err)
			break
		}

		_, _ = io.Copy(os.Stdout, out)

		c.State = c.Status()
	}

	return c.State
}


// Commit a container
// Commit a container to create an image from its contents:
func (c *Container) Commit() *ux.State {
	if state := c.IsNil(); state.IsError() {
		return state
	}

	for range OnlyOnce {
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

	return c.State
}


func (c *Container) GetContainerSsh() (string, *ux.State) {
	var port string
	if state := c.IsNil(); state.IsError() {
		return "", state
	}

	for range OnlyOnce {
		var found bool
		for _, p := range c.Summary.Ports {
			if p.PrivatePort == 22 {
				port = fmt.Sprintf("%d", p.PublicPort)
				found = true
				break
			}
		}

		if !found {
			c.State.SetError("no SSH port")
		}
	}

	return port, c.State
}
