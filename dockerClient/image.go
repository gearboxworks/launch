package dockerClient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"launch/defaults"
	"launch/gear/gearJson"
	"launch/only"
	"launch/ux"
	"github.com/docker/docker/api/types"
	"io"
	"os"
	"strings"
)

type Image struct {
	ID         string
	Name       string

	Version    string
	Summary    *types.ImageSummary
	Details    types.ImageInspect
	GearConfig gearJson.GearConfig

	_Parent    *DockerGear
}


func (me *Image) EnsureNotNil() ux.State {
	var state ux.State

	for range only.Once {
		if me == nil {
			state.SetError("gear is nil")
			break
		}

		if me.ID == "" {
			state.SetError("ID is nil")
			break
		}

		if me.Name == "" {
			state.SetError("name is nil")
			break
		}

		if me.Version == "" {
			state.SetError("version is nil")
			break
		}
	}

	return state
}


// Pull an image
// Pull an image, like docker pull:
func (me *Image) Pull() ux.State {
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		var repo string
		if me.Version == "" {
			repo = fmt.Sprintf("gearboxworks/%s", me.Name)
		} else {
			repo = fmt.Sprintf("gearboxworks/%s:%s", me.Name, me.Version)
		}

		ctx := context.Background()
		//ctx, cancel := context.WithTimeout(context.Background(), Timeout * 1000)
		//defer cancel()

		//df := filters.NewArgs()
		//df.Add("name", "terminus")
		//var results []registry.SearchResult
		//results, err = me.client.ImageSearch(ctx, "", types.ImageSearchOptions{Filters:df})
		//for _, v := range results {
		//	fmt.Printf("%s - %s\n", v.Name, v.Description)
		//}

		var out io.ReadCloser
		var err error
		out, err = me._Parent.Client.ImagePull(ctx, repo, types.ImagePullOptions{All: false})
		if err != nil {
			state.SetError("error pulling gear: %s", err)
			break
		}

		defer out.Close()

		ux.Printf("pulling Gearbox gear: %s\n", me.Name)
		d := json.NewDecoder(out)
		var event *PullEvent
		for {
			err := d.Decode(&event)
			if err != nil {
				if err == io.EOF {
					break
				}

				state.SetError("error pulling gear: %s", err)
				break
			}

			// fmt.Printf("EVENT: %+v\n", event)
			ux.Printf("%+v\r", event.Progress)
		}
		ux.Printf("\n")

		if state.IsError() {
			break
		}

		// Latest event for new image
		// EVENT: {Status:Status: Downloaded newer image for busybox:latest Error: Progress:[==================================================>]  699.2kB/699.2kB ProgressDetail:{Current:699243 Total:699243}}
		// Latest event for up-to-date image
		// EVENT: {Status:Status: Image is up to date for busybox:latest Error: Progress: ProgressDetail:{Current:0 Total:0}}
		if event != nil {
			if strings.HasPrefix(event.Status, "Status: Downloaded newer image for") {
				// new
				ux.PrintfOk("pulling Gearbox gear %s - OK\n", me.Name)
			} else if strings.HasPrefix(event.Status, "Status: Image is up to date for") {
				// up-to-date
				ux.PrintfOk("pulling Gearbox gear %s - image up to date\n", me.Name)
			} else {
				ux.PrintfWarning("pulling Gearbox gear %s - unknown\n", me.Name)
			}
		}
		//ux.Printf("\nGear image pull OK: %+v\n", event)
		ux.Printf("%s\n", event.Status)

		//buf := new(bytes.Buffer)
		//_, err = buf.ReadFrom(out)
		//fmt.Printf("%s", buf.String())
		//_, _ = io.Copy(os.Stdout, out)
	}

	return state
}

type PullEvent struct {
	Status         string `json:"status"`
	Error          string `json:"error"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}


// Pull an image with authentication
// Pull an image, like docker pull, with authentication:
func (me *Image) ImageAuthPull() ux.State {
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		authConfig := types.AuthConfig{
			Username: "username",
			Password: "password",
		}

		var err error
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			state.SetError("error pulling gear: %s", err)
			break
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		out, err := me._Parent.Client.ImagePull(ctx, "alpine", types.ImagePullOptions{RegistryAuth: authStr})
		if err != nil {
			state.SetError("error pulling gear: %s", err)
			break
		}

		defer out.Close()

		_, _ = io.Copy(os.Stdout, out)
	}

	return state
}


// Remove containers
// Now that you know what containers exist, you can perform operations on them.
// This example stops all running containers.
func (me *Image) Remove() ux.State {
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		options := types.ImageRemoveOptions {
			Force:         true,
			PruneChildren: true,
		}

		_, err := me._Parent.Client.ImageRemove(ctx, me.ID, options)
		if err != nil {
			state.SetError("error removing gear: %s", err)
			break
		}
	}

	return state
}
