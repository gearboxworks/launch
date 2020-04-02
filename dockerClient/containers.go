package dockerClient

import (
	"context"
	"fmt"
	"gb-launch/defaults"
	"gb-launch/gear/gearJson"
	"gb-launch/only"
	"gb-launch/ux"
	"github.com/docker/docker/api/types"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/table"
	"os"
	"strings"
)

// List and manage containers
// You can use the API to list containers that are running, just like using docker ps:
// func ContainerList(f types.ContainerListOptions) error {
func (me *DockerGear) ContainerList(f string) ux.State {
	var state ux.State

	for range only.Once {
		var err error

		if me.Debug {
			fmt.Printf("DEBUG: ContainerList(%s)\n", f)
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		var containers []types.Container
		containers, err = me.Client.ContainerList(ctx, types.ContainerListOptions{Size: true, All: true})
		if err != nil {
			state.SetError("gear list error: %s", err)
			break
		}

		ux.PrintfCyan("\nConfigured Gearbox gears:\n")
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{
			("Name"),
			("Class"),
			("State"),
			("Image"),
			("Ports"),
			("SSH port"),
			("IP Address"),
			("Mounts"),
			("Size"),
		})

		for _, c := range containers {
			var gc *gearJson.GearConfig
			gc, state = gearJson.New(c.Labels["gearbox.json"])
			if state.IsError() {
				continue
			}

			if gc.Meta.Organization != defaults.Organization {
				continue
			}

			if f != "" {
				if gc.Meta.Name != f {
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

			var state string
			if c.State == "running" {
				state = ux.SprintfGreen(c.State)
			} else {
				state = ux.SprintfYellow(c.State)
			}

			t.AppendRow([]interface{}{
				ux.SprintfWhite(name),
				ux.SprintfWhite(gc.Meta.Class),
				state,
				ux.SprintfWhite(c.Image),
				ux.SprintfWhite(ports),
				ux.SprintfWhite(sshPort),
				ux.SprintfWhite(ipAddress),
				ux.SprintfWhite(mounts),
				ux.SprintfWhite(humanize.Bytes(uint64(c.SizeRootFs))),
			})
		}

		t.Render()
		state.ClearError()
	}

	return state
}


func (me *DockerGear) FindContainer(gearName string, gearVersion string) (bool, ux.State) {
	var ok bool
	var state ux.State

	for range only.Once {
		var err error

		if me.Debug {
			fmt.Printf("DEBUG: FindContainer(%s, %s)\n", gearName, gearVersion)
		}

		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		if gearName == "" {
			state.SetError("empty gearname")
			break
		}

		if gearVersion == "" {
			gearVersion = "latest"
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		var containers []types.Container
		containers, err = me.Client.ContainerList(ctx, types.ContainerListOptions{All: true, Limit: 256})
		if err != nil {
			state.SetError("gear list error: %s", err)
			break
		}

		for _, c := range containers {
			var gc *gearJson.GearConfig
			gc, state = gearJson.New(c.Labels["gearbox.json"])
			if state.IsError() {
				continue
			}

			if gc.Meta.Organization != defaults.Organization {
				continue
			}

			if gc.Meta.Name != gearName {
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
			me.Container.Name = gc.Meta.Name
			ok = true

			break
		}

		if state.IsError() {
			state.ClearError()
			break
		}

		if me.Container.Summary == nil {
			break
		}

		d := types.ContainerJSON{}
		d, err = me.Client.ContainerInspect(ctx, me.Container.ID)
		if err != nil {
			state.SetError("gear inspect error: %s", err)
			break
		}
		me.Container.Details = &d

		state = me.Container.EnsureNotNil()
		if state.IsError() {
			break
		}
	}

	if me.Debug {
		state.Print()
	}

	return ok, state
}
