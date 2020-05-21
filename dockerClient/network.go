package dockerClient

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/jedib0t/go-pretty/table"
	"launch/defaults"
	"launch/only"
	"launch/ux"
	"os"
)

// List and manage containers
// You can use the API to list containers that are running, just like using docker ps:
// func ContainerList(f types.ContainerListOptions) error {
func (gear *DockerGear) NetworkList(f string) *ux.State {
	if state := gear.IsNil(); state.IsError() {
		return state
	}

	for range only.Once {
		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		df := filters.NewArgs()
		df.Add("name", f)

		nets, err := gear.Client.NetworkList(ctx, types.NetworkListOptions{Filters: df})
		if err != nil {
			gear.State.SetError("error listing networks")
			break
		}

		ux.PrintflnCyan("\nConfigured Gearbox networks:")
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{
			"Name",
			"Driver",
			"Subnet",
		})

		for _, c := range nets {
			n := ""
			if len(c.IPAM.Config) > 0 {
				n = c.IPAM.Config[0].Subnet
			}

			t.AppendRow([]interface{}{
				ux.SprintfWhite(c.Name),
				ux.SprintfWhite(c.Driver),
				ux.SprintfWhite(n),
			})
		}

		t.Render()
	}

	return gear.State
}


func (gear *DockerGear) FindNetwork(netName string) *ux.State {
	if state := gear.IsNil(); state.IsError() {
		return state
	}

	for range only.Once {
		if netName == "" {
			gear.State.SetError("empty gear name")
			break
		}

		df := filters.NewArgs()
		df.Add("name", netName)

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		nets, err := gear.Client.NetworkList(ctx, types.NetworkListOptions{Filters: df})
		if err != nil {
			gear.State.SetError("gear image search error: %s", err)
			break
		}

		for _, c := range nets {
			if c.Name == netName {
				gear.State.SetOk("found")
				break
			}
		}
	}

	return gear.State
}


func (gear *DockerGear) NetworkCreate(netName string) *ux.State {
	if state := gear.IsNil(); state.IsError() {
		return state
	}

	for range only.Once {
		gear.State = gear.FindNetwork(netName)
		if gear.State.IsError() {
			break
		}
		if gear.State.IsOk() {
			break
		}

		netConfig := types.NetworkCreate {
			CheckDuplicate: true,
			Driver:         "bridge",
			Scope:          "local",
			EnableIPv6:     false,
			IPAM:           &network.IPAM {
				Driver:  "default",
				Options: nil,
				Config:  []network.IPAMConfig {
					{
						Subnet: "172.42.0.0/24",
						Gateway: "172.42.0.1",
					},
				},
			},
			Internal:       false,
			Attachable:     false,
			Ingress:        false,
			ConfigOnly:     false,
			ConfigFrom:     nil,
			Options:        nil,
			Labels:         nil,
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		resp, err := gear.Client.NetworkCreate(ctx, netName, netConfig)
		gear.State.SetError(err)
		if gear.State.IsError() {
			break
		}

		if resp.ID == "" {
			gear.State.SetError("cannot create network")
			break
		}
	}

	return gear.State
}
