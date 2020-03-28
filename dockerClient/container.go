package dockerClient

import (
	"gb-launch/only"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/table"
	"io"
	"os"
	"strings"
)

type Container struct {
	ID           string
	Name         string
	Version      string

	Summary      *types.Container
	Details      *types.ContainerJSON
	GearConfig   *GearConfig

	// ctx          *context.Context
	client       *client.Client
}
type Containers []Container


func (me *Container) EnsureNotNil() error {
	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("gear is nil")
			break
		}

		if me.ID == "" {
			err = errors.New("ID is nil")
			break
		}

		if me.Name == "" {
			err = errors.New("name is nil")
			break
		}

		if me.Version == "" {
			err = errors.New("version is nil")
			break
		}

		if me.client == nil {
			err = errors.New("client is nil")
			break
		}

		// if me.ctx == nil {
		// 	err = errors.New("ctx is nil")
		// 	break
		// }
		//
		// if me.Summary == nil {
		// 	err = errors.New("container is nil")
		// 	break
		// }
		//
		// if me.Details == "" {
		// 	err = errors.New("container is nil")
		// 	break
		// }
	}

	return err
}


// List and manage containers
// You can use the API to list containers that are running, just like using docker ps:
// func ContainerList(f types.ContainerListOptions) error {
func (me *Gear) ContainerList(f string) error {
	var err error

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: ContainerList(%s)\n", f)
		}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		var containers []types.Container
		containers, err = me.DockerClient.ContainerList(ctx, types.ContainerListOptions{Size: true, All: true})
		if err != nil {
			break
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Name", "Class", "State", "Image", "Ports", "SSH port", "IP Address", "Mounts", "Size"})

		for _, c := range containers {
			var gc *GearConfig
			gc, err = NewGearConfig(c.Labels["gearbox.json"])
			if err != nil {
				continue
			}

			if gc.Organization != Organization {
				continue
			}

			if f != "" {
				if gc.Name != f {
					continue
				}
			}

			name := strings.TrimPrefix(c.Names[0], "/")

			sshPort := ""
			var ports string
			for _, p := range c.Ports {
				if p.PrivatePort == 22 {
					sshPort = fmt.Sprintf("%d", p.PublicPort)
					continue
				}
				ports += fmt.Sprintf("%s://%s:%d => %d\n", p.Type, p.IP, p.PublicPort, p.PrivatePort)
			}
			if sshPort == "0" {
				sshPort = "none"
			}

			var mounts string
			for _, m := range c.Mounts {
				// ms += fmt.Sprintf("%s(%s) host:%s => container:%s (RW:%v)\n", m.Name, m.Type, m.Source, m.Destination, m.RW)
				mounts += fmt.Sprintf("host:%s\n\t=> container:%s (RW:%v)\n", m.Source, m.Destination, m.RW)
			}

			var ipAddress string
			for k, n := range c.NetworkSettings.Networks {
				ipAddress += fmt.Sprintf("(%s) %s\n", k, n.IPAddress)
			}

			t.AppendRow([]interface{}{name, gc.Class, c.State, c.Image, ports, sshPort, ipAddress, mounts, humanize.Bytes(uint64(c.SizeRootFs))})
			err = nil
		}

		t.Render()
	}

	return err
}


func (me *Gear) FindContainer(gearName string, gearVersion string) (bool, error) {
	var ok bool
	var err error

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: FindContainer(%s, %s)\n", gearName, gearVersion)
		}

		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if gearName == "" {
			err = errors.New("empty gearname")
			break
		}

		if gearVersion == "" {
			gearVersion = "latest"
		}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		var containers []types.Container
		containers, err = me.DockerClient.ContainerList(ctx, types.ContainerListOptions{All: true})
		if err != nil {
			break
		}

		for _, c := range containers {
			var gc *GearConfig
			gc, err = NewGearConfig(c.Labels["gearbox.json"])
			if err != nil {
				continue
			}

			if gc.Organization != Organization {
				continue
			}

			if gc.Name != gearName {
				continue
			}

			if gearVersion == "latest" {
				gl := gc.Versions.GetLatest()
				if gl == "" {
					continue
				}
				gearVersion = gl
			} else {
				if !gc.Versions.HasVersion(gearVersion) {
					continue
				}
			}

			me.Container.Name = gearName
			me.Container.Version = gearVersion
			me.Container.GearConfig = gc
			me.Container.Summary = &c
			me.Container.ID = c.ID
			me.Container.Name = gc.Name
			ok = true

			break
		}

		if err != nil {
			break
		}

		if me.Container.Summary == nil {
			break
		}

		d := types.ContainerJSON{}
		d, err = me.DockerClient.ContainerInspect(ctx, me.Container.ID)
		if err != nil {
			break
		}
		me.Container.Details = &d

		err = me.Container.EnsureNotNil()
		if err != nil {
			break
		}
	}

	if me.Debug {
		fmt.Printf("DEBUG: Error - %s\n", err)
	}

	return ok, err
}


func (me *Container) State() State {
	var state State

	for range only.Once {
		state.Error = me.EnsureNotNil()
		if state.Error != nil {
			break
		}

		df := filters.NewArgs()
		df.Add("id", me.ID)

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		var containers []types.Container
		containers, state.Error = me.client.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: df})
		if state.Error != nil {
			break
		}
		me.Summary = &containers[0]

		me.GearConfig, state.Error = NewGearConfig(me.Summary.Labels["gearbox.json"])
		if state.Error != nil {
			break
		}

		if me.GearConfig.Organization != Organization {
			state.Error = errors.New("not a Gearbox container")
			break
		}

		d := types.ContainerJSON{}
		d, state.Error = me.client.ContainerInspect(ctx, me.ID)
		if state.Error != nil {
			break
		}
		me.Details = &d

		state.SetState(me.Details.State.Status)
	}

	if state.Error != nil {
		me.Summary = nil
		me.Details = nil
	}

	return state
}


// Run a container in the background
// You can also run containers in the background, the equivalent of typing docker run -d bfirsh/reticulate-splines:
func (me *Container) Start() State {
	var state State

	for range only.Once {
		state.Error = me.EnsureNotNil()
		if state.Error != nil {
			break
		}

		state = me.State()
		if state.Error != nil {
			break
		}

		if state.IsRunning() {
			break
		}

		if !state.IsCreated() && !state.IsExited() {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		state.Error = me.client.ContainerStart(ctx, me.ID, types.ContainerStartOptions{})
		if state.Error != nil {
			break
		}

		statusCh, errCh := me.client.ContainerWait(ctx, me.ID, "") // container.WaitConditionNotRunning
		select {
			case err := <-errCh:
				if err != nil {
					// fmt.Printf("SC: %s\n", state.Error)
					// return false, err
				}
				break

			case status := <-statusCh:
				fmt.Printf("status.StatusCode: %#+v\n", status.StatusCode)
				break
		}
		// fmt.Printf("SC: %s\n", status)
		// fmt.Printf("SC: %s\n", err)

		state = me.State()
		if state.Error != nil {
			break
		}
		if !state.IsRunning() {
			state.Error = errors.New("cannot start container")
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
func (me *Gear) ContainerCreate(gearName string, gearVersion string, gearMount string) State {
	var state State

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: ContainerCreate(%s, %s, %s)\n", gearName, gearVersion, gearMount)
		}

		state.Error = me.EnsureNotNil()
		if state.Error != nil {
			break
		}

		if gearName == "" {
			state.Error = errors.New("empty gearname")
			break
		}

		if gearVersion == "" {
			gearVersion = "latest"
		}

		var ok bool
		ok, state.Error = me.FindContainer(gearName, gearVersion)
		if state.Error != nil {
			break
		}
		if !ok {
			//state.Error = me.Search(gearName, gearVersion)

			ok, state.Error = me.FindImage(gearName, gearVersion)
			if state.Error != nil {
				me.Image.ID = gearName
				me.Image.Name = gearName
				me.Image.Version = gearVersion
				state.Error = me.Image.Pull()
				if state.Error != nil {
					state.Error = errors.New(fmt.Sprintf("no such image '%s'", gearName))
					break
				}
			}

			ok, state.Error = me.FindContainer(gearName, gearVersion)
			if state.Error != nil {
				break
			}
		}

		// me.Image.Details.Container = "gearboxworks/golang:1.14"
		// tag := fmt.Sprintf("", me.Image.Name, me.Image.Version)
		tag := fmt.Sprintf("gearboxworks/%s:%s", me.Image.Name, me.Image.Version)
		gn := fmt.Sprintf("%s-%s", me.Image.Name, me.Image.Version)
		var binds []string
		if gearMount != "" {
			binds = append(binds, fmt.Sprintf("%s:%s", gearMount, DefaultProject))
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
			NetworkMode:     "gearboxnet",
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

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		var resp container.ContainerCreateCreatedBody
		resp, state.Error = me.DockerClient.ContainerCreate(ctx, &config, &hostConfig, &netConfig, gn)
		if state.Error != nil {
			break
		}

		if resp.ID == "" {
			break
		}

		me.Container.ID = resp.ID
		me.Container.Name = me.Image.Name
		me.Container.Version = me.Image.Version

		// var state State
		state = me.Container.State()
		if state.Error != nil {
			break
		}

		if state.IsRunning() {
			break
		}

		if state.IsPaused() {
			break
		}

		if state.IsRestarting() {
			break
		}

		if state.IsCreated() {
			break
		}

		// url := fmt.Sprintf("docker.io/library/%s", image)
		// reader, err := me.DockerClient.ImagePull(me.Ctx, url, types.ImagePullOptions{})
		// if err != nil {
		// 	break
		// }
		// _, _ = io.Copy(os.Stdout, reader)
		//
		// if err := me.client.ContainerStart(*me.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		// 	break
		// }
		//
		// statusCh, errCh := me.client.ContainerWait(*me.ctx, resp.ID, container.WaitConditionNotRunning)
		// select {
		// 	case err := <-errCh:
		// 		if err != nil {
		// 			break
		// 		}
		// 	case <-statusCh:
		// }
		//
		// out, err := me.client.ContainerLogs(*me.ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		// if err != nil {
		// 	break
		// }
		//
		// stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}

	if me.Debug {
		fmt.Printf("DEBUG: Error - %s\n", state.Error)
	}

	return state
}


// Stop all running containers
// Now that you know what containers exist, you can perform operations on them.
// This example stops all running containers.
func (me *Container) Stop() error {
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		err = me.client.ContainerStop(ctx, me.ID, nil)
		if err != nil {
			break
		}
	}

	return err
}


// Remove containers
// Now that you know what containers exist, you can perform operations on them.
// This example stops all running containers.
func (me *Container) Remove() error {
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		options := types.ContainerRemoveOptions{
			RemoveVolumes: false,
			RemoveLinks:   false,
			Force:         false,
		}

		err = me.client.ContainerRemove(ctx, me.ID, options)
		if err != nil {
			break
		}
	}

	return err
}


// Print the logs of a specific container
// You can also perform actions on individual containers.
// This example prints the logs of a container given its ID.
// You need to modify the code before running it to change the hard-coded ID of the container to print the logs for.
func (me *Container) Logs() error {
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		options := types.ContainerLogsOptions{ShowStdout: true}
		// Replace this ID with a container that really exists
		out, err := me.client.ContainerLogs(ctx, "f1064a8a4c82", options)
		if err != nil {
			break
		}

		_, _ = io.Copy(os.Stdout, out)
	}

	return err
}


// Commit a container
// Commit a container to create an image from its contents:
func (me *Container) Commit() error {
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		createResp, err := me.client.ContainerCreate(ctx, &container.Config{
			Image: "alpine",
			Cmd:   []string{"touch", "/helloworld"},
		}, nil, nil, "")
		if err != nil {
			break
		}

		if err := me.client.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{}); err != nil {
			break
		}

		statusCh, errCh := me.client.ContainerWait(ctx, createResp.ID, container.WaitConditionNotRunning)
		select {
			case err := <-errCh:
				if err != nil {
					break
				}
			case <-statusCh:
		}

		commitResp, err := me.client.ContainerCommit(ctx, createResp.ID, types.ContainerCommitOptions{Reference: "helloworld"})
		if err != nil {
			break
		}

		fmt.Println(commitResp.ID)
	}

	return err
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

func (me *Gear) GetContainerSsh() (string, error) {
	var port string
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = errors.New("no SSH port")

		for _, p := range me.Container.Summary.Ports {
			if p.PrivatePort == 22 {
				port = fmt.Sprintf("%d", p.PublicPort)
				err = nil
				break
			}
		}
	}

	return port, err
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
