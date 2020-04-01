package dockerClient

import (
	"context"
	"errors"
	"fmt"
	"gb-launch/defaults"
	"gb-launch/gear/gearJson"
	"gb-launch/only"
	"github.com/docker/docker/api/types"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/table"
	"os"
	"strings"
)

// List and manage containers
// You can use the API to list containers that are running, just like using docker ps:
// func ContainerList(f types.ContainerListOptions) error {
func (me *DockerGear) ContainerList(f string) error {
	var err error

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: ContainerList(%s)\n", f)
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		var containers []types.Container
		containers, err = me.Client.ContainerList(ctx, types.ContainerListOptions{Size: true, All: true})
		if err != nil {
			break
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Name", "Class", "State", "Image", "Ports", "SSH port", "IP Address", "Mounts", "Size"})

		for _, c := range containers {
			var gc *gearJson.GearConfig
			gc, err = gearJson.New(c.Labels["gearbox.json"])
			if err != nil {
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

			t.AppendRow([]interface{}{name, gc.Meta.Class, c.State, c.Image, ports, sshPort, ipAddress, mounts, humanize.Bytes(uint64(c.SizeRootFs))})
			err = nil
		}

		t.Render()
		err = nil
	}

	return err
}


func (me *DockerGear) FindContainer(gearName string, gearVersion string) (bool, error) {
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

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		var containers []types.Container
		containers, err = me.Client.ContainerList(ctx, types.ContainerListOptions{All: true, Limit: 256})
		if err != nil {
			break
		}

		for _, c := range containers {
			var gc *gearJson.GearConfig
			gc, err = gearJson.New(c.Labels["gearbox.json"])
			if err != nil {
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

		if err != nil {
			err = nil
			break
		}

		if me.Container.Summary == nil {
			break
		}

		d := types.ContainerJSON{}
		d, err = me.Client.ContainerInspect(ctx, me.Container.ID)
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
		fmt.Printf("DEBUG: FindContainer() error: %s\n", err)
	}

	return ok, err
}
