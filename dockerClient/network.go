package dockerClient

import (
	"context"
	"errors"
	"fmt"
	"gb-launch/defaults"
	"gb-launch/only"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/jedib0t/go-pretty/table"
	"os"
)

// List and manage containers
// You can use the API to list containers that are running, just like using docker ps:
// func ContainerList(f types.ContainerListOptions) error {
func (me *DockerGear) NetworkList(f string) error {
	var err error

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: NetworkList(%s)\n", f)
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		df := filters.NewArgs()
		df.Add("name", f)

		var nets []types.NetworkResource
		nets, err = me.Client.NetworkList(ctx, types.NetworkListOptions{Filters: df})
		if err != nil {
			break
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Name", "Driver"})

		for _, c := range nets {
			t.AppendRow([]interface{}{c.Name, c.Driver})
		}

		t.Render()
	}

	return err
}


func (me *DockerGear) FindNetwork(netName string) (bool, error) {
	var ok bool
	var err error

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: FindNetwork(%s)\n", netName)
		}

		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if netName == "" {
			err = errors.New("empty gearname")
			break
		}

		df := filters.NewArgs()
		df.Add("name", netName)

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		var nets []types.NetworkResource
		nets, err = me.Client.NetworkList(ctx, types.NetworkListOptions{Filters: df})
		if err != nil {
			break
		}

		for _, c := range nets {
			if c.Name == netName {
				ok = true
				break
			}
		}
	}

	if me.Debug {
		fmt.Printf("DEBUG: FindNetwork() error: %s\n", err)
	}

	return ok, err
}


func (me *DockerGear) NetworkCreate(netName string) State {
	var state State

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: ContainerCreate(%s)\n", netName)
		}

		if netName == "" {
			state.Error = errors.New("empty netName")
			break
		}

		var ok bool
		ok, state.Error = me.FindNetwork(netName)
		if state.Error != nil {
			break
		}
		if ok {
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
