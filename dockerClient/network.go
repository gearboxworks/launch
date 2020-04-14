package dockerClient

import (
	"context"
	"errors"
	"fmt"
	"launch/defaults"
	"launch/only"
	"launch/ux"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/jedib0t/go-pretty/table"
	"os"
)

// List and manage containers
// You can use the API to list containers that are running, just like using docker ps:
// func ContainerList(f types.ContainerListOptions) error {
func (me *DockerGear) NetworkList(f string) ux.State {
	var state ux.State

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: NetworkList(%s)\n", f)
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		df := filters.NewArgs()
		df.Add("name", f)

		var nets []types.NetworkResource
		var err error
		nets, err = me.Client.NetworkList(ctx, types.NetworkListOptions{Filters: df})
		if err != nil {
			state.SetError("error listing networks")
			break
		}

		ux.PrintfCyan("\nConfigured Gearbox networks:\n")
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

	return state
}


func (me *DockerGear) FindNetwork(netName string) ux.State {
	var state ux.State

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: FindNetwork(%s)\n", netName)
		}

		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		if netName == "" {
			state.SetError("empty gear name")
			break
		}

		df := filters.NewArgs()
		df.Add("name", netName)

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		var nets []types.NetworkResource
		var err error
		nets, err = me.Client.NetworkList(ctx, types.NetworkListOptions{Filters: df})
		if err != nil {
			state.SetError("gear image search error: %s", err)
			break
		}

		for _, c := range nets {
			if c.Name == netName {
				state.SetOk("found")
				break
			}
		}
	}

	if me.Debug {
		state.Print()
	}

	return state
}


func (me *DockerGear) NetworkCreate(netName string) ux.State {
	var state ux.State

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: ContainerCreate(%s)\n", netName)
		}

		if netName == "" {
			state.SetError("empty netName")
			break
		}

		state = me.FindNetwork(netName)
		if state.IsError() {
			break
		}
		if state.IsOk() {
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
		defer cancel()

		var resp types.NetworkCreateResponse
		resp, state.Error = me.Client.NetworkCreate(ctx, netName, netConfig)
		if state.Error != nil {
			break
		}

		if resp.ID == "" {
			state.Error = errors.New("cannot create network")
			break
		}
	}

	if me.Debug {
		fmt.Printf("DEBUG: Error - %s\n", state.Error)
	}

	return state
}
