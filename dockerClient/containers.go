package dockerClient

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/table"
	"launch/defaults"
	"launch/gear/gearJson"
	"launch/only"
	"launch/ux"
	"os"
	"strings"
)

// List and manage containers
// You can use the API to list containers that are running, just like using docker ps:
// func ContainerList(f types.ContainerListOptions) error {
func (gear *DockerGear) ContainerList(f string) (int, ux.State) {
	var state ux.State
	var count int

	for range only.Once {
		var err error

		if gear.Debug {
			fmt.Printf("DEBUG: ContainerList(%s)\n", f)
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		var containers []types.Container
		containers, err = gear.Client.ContainerList(ctx, types.ContainerListOptions{Size: true, All: true})
		if err != nil {
			state.SetError("gear list error: %s", err)
			break
		}

		ux.PrintfCyan("Installed Gearbox gears: ")
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{
			"Name",
			"Class",
			"State",
			"Image",
			"Ports",
			"SSH port",
			"IP Address",
			"Mounts",
			"Size",
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
				//ports += fmt.Sprintf("%s://%s:%d => %d\n", p.Type, p.IP, p.PublicPort, p.PrivatePort)
				if p.IP == "0.0.0.0" {
					ports += fmt.Sprintf("%d => %d\n", p.PublicPort, p.PrivatePort)
				} else {
					ports += fmt.Sprintf("%s://%s:%d => %d\n", p.Type, p.IP, p.PublicPort, p.PrivatePort)
				}
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
			if c.State == ux.StateRunning {
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

		state.ClearError()
		count = t.Length()
		if count == 0 {
			ux.PrintfYellow("None found\n")
			break
		}

		ux.PrintfGreen("%d found\n", count)
		t.Render()
		ux.PrintfWhite("\n")
	}

	return count, state
}

func (gear *DockerGear) FindContainer(gearName string, gearVersion string) (bool, ux.State) {
	var ok bool
	var state ux.State

	for range only.Once {
		var err error

		if gear.Debug {
			fmt.Printf("DEBUG: FindContainer(%s, %s)\n", gearName, gearVersion)
		}

		state = gear.EnsureNotNil()
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
		//noinspection GoDeferInLoop
		defer cancel()

		var containers []types.Container
		containers, err = gear.Client.ContainerList(ctx, types.ContainerListOptions{All: true, Limit: 256})
		if err != nil {
			state.SetError("gear list error: %s", err)
			break
		}

		// Start out with "not found". Will be cleared if found or error occurs.
		state.SetWarning("Gear '%s:%s' doesn't exist.", gearName, gearVersion)

		for _, c := range containers {
			var gc *gearJson.GearConfig
			gc, state = gearJson.New(c.Labels["gearbox.json"])
			if state.IsError() {
				continue
			}

			if gc.Meta.Organization != defaults.Organization {
				continue
			}

			// @TODO - BEGIN - This needs to be refactored!!!
			if gc.Meta.Name != gearName {
				if !defaults.RunAs.AsLink {
					continue
				}

				cs := gc.MatchCommand(gearName)
				if cs == nil {
					continue
				}

				gearName = gc.Meta.Name
			}

			if gearVersion == "latest" {
				gl := gc.Versions.GetLatest()
				if gl == "" {
					continue
				}
				gearVersion = gl
			} else {
				if !gc.Versions.HasVersion(gearVersion) {
					// Finally compare container image name.
					finalCheck := fmt.Sprintf("%s/%s:%s", defaults.Organization, gearName, gearVersion)
					if finalCheck != c.Image {
						continue
					}
				}
			}

			//fmt.Printf("%s => F:%s:F:%s:F\n", gearVersion, c.Labels["gearbox.version"], c.Labels["container.majorversion"])
			if c.Labels["gearbox.version"] == gearVersion {
			} else if c.Labels["container.majorversion"] == gearVersion {
			} else {
				continue
			}
			// @TODO - END - This needs to be refactored!!!

			gear.Container.Name = gearName
			gear.Container.Version = gearVersion
			gear.Container.GearConfig = gc
			gear.Container.Summary = &c
			gear.Container.ID = c.ID
			gear.Container.Name = gc.Meta.Name
			ok = true
			state.ClearAll()

			break
		}

		if state.IsError() {
			state.ClearError()
			break
		}

		if gear.Container.Summary == nil {
			break
		}

		d := types.ContainerJSON{}
		d, err = gear.Client.ContainerInspect(ctx, gear.Container.ID)
		if err != nil {
			state.SetError("gear inspect error: %s", err)
			break
		}
		gear.Container.Details = &d

		state = gear.Container.EnsureNotNil()
		if state.IsError() {
			break
		}
	}

	if gear.Debug {
		state.Print()
	}

	return ok, state
}
