package helperDocker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/table"
	"launch/defaults"
	"launch/gear/gearJson"
	"launch/ux"
	"os"
	"strings"
)

// List and manage containers
// You can use the API to list containers that are running, just like using docker ps:
// func ContainerList(f types.ContainerListOptions) error {
func (gear *DockerGear) ContainerList(f string) (int, *ux.State) {
	var count int
	if state := gear.IsNil(); state.IsError() {
		return 0, state
	}

	for range OnlyOnce {
		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		var containers []types.Container
		var err error
		containers, err = gear.Client.ContainerList(ctx, types.ContainerListOptions{Size: true, All: true})
		if err != nil {
			gear.State.SetError("gear list error: %s", err)
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
			gc = gearJson.New(c.Labels["gearbox.json"])
			if gc.State.IsError() {
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

		gear.State.ClearError()
		count = t.Length()
		if count == 0 {
			ux.PrintfYellow("None found\n")
			break
		}

		ux.PrintflnGreen("%d found", count)
		t.Render()
		ux.PrintflnBlue("")
	}

	return count, gear.State
}

func (gear *DockerGear) FindContainer(gearName string, gearVersion string) (bool, *ux.State) {
	var ok bool
	if state := gear.IsNil(); state.IsError() {
		return false, state
	}

	for range OnlyOnce {
		if gearName == "" {
			gear.State.SetError("empty gearname")
			break
		}

		if gearVersion == "" {
			gearVersion = "latest"
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		var containers []types.Container
		var err error
		containers, err = gear.Client.ContainerList(ctx, types.ContainerListOptions{All: true, Limit: 256})
		if err != nil {
			gear.State.SetError("gear list error: %s", err)
			break
		}

		// Start out with "not found". Will be cleared if found or error occurs.
		gear.State.SetWarning("Gear '%s:%s' doesn't exist.", gearName, gearVersion)

		for _, c := range containers {
			var gc *gearJson.GearConfig
			ok, gc = MatchContainer(&c, defaults.Organization, gearName, gearVersion)
			if !ok {
				continue
			}

			gear.Container.Name = gearName
			gear.Container.Version = gearVersion
			gear.Container.GearConfig = gc
			gear.Container.Summary = &c
			gear.Container.ID = c.ID
			gear.Container.Name = gc.Meta.Name
			gear.Container.State = gear.Container.State.EnsureNotNil()
			ok = true
			gear.State.SetOk("Found Gear '%s:%s'.", gearName, gearVersion)

			break
		}

		if gear.State.IsNotOk() {
			if !ok {
				gear.State.ClearError()
			}
			break
		}

		if gear.Container.Summary == nil {
			break
		}

		ctx2, cancel2 := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel2()
		d := types.ContainerJSON{}
		d, err = gear.Client.ContainerInspect(ctx2, gear.Container.ID)
		if err != nil {
			gear.State.SetError("gear inspect error: %s", err)
			break
		}
		gear.Container.Details = &d
	}

	return ok, gear.State
}


func MatchContainer(m *types.Container, gearOrg string, gearName string, gearVersion string) (bool, *gearJson.GearConfig) {
	var ok bool
	var gc *gearJson.GearConfig

	for range OnlyOnce {
		if MatchTag("<none>:<none>", m.Names) {
			ok = false
			break
		}

		gc = gearJson.New(m.Labels["gearbox.json"])
		if gc.State.IsError() {
			ok = false
			break
		}

		if gc.Meta.Organization != defaults.Organization {
			ok = false
			break
		}

		tagCheck := fmt.Sprintf("%s/%s:%s", gearOrg, gearName, gearVersion)
		if m.Image == tagCheck {
			ok = true
			break
		}

		if gc.Meta.Name != gearName {
			if !defaults.RunAs.AsLink {
				ok = false
				break
			}

			cs := gc.MatchCommand(gearName)
			if cs == nil {
				ok = false
				break
			}

			gearName = gc.Meta.Name
		}

		if !gc.Versions.HasVersion(gearVersion) {
			ok = false
			break
		}

		if gearVersion == "latest" {
			gl := gc.Versions.GetLatest()
			if gearVersion != "" {
				gearVersion = gl
			}
		}
		for range OnlyOnce {
			if m.Labels["gearbox.version"] == gearVersion {
				ok = true
				break
			}

			if m.Labels["container.majorversion"] == gearVersion {
				ok = true
				break
			}

			ok = false
		}
		break
	}

	return ok, gc
}
