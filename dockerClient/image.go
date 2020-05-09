package dockerClient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"io"
	"launch/defaults"
	"launch/gear/gearJson"
	"launch/only"
	"launch/ux"
	"os"
	"strings"
)

type Image struct {
	ID         string
	Name       string

	Version    string
	Summary    *types.ImageSummary
	Details    types.ImageInspect
	GearConfig *gearJson.GearConfig

	_Parent    *DockerGear
}


func (image *Image) EnsureNotNil() ux.State {
	var state ux.State

	for range only.Once {
		if image == nil {
			state.SetError("gear is nil")
			break
		}

		if image.ID == "" {
			state.SetError("ID is nil")
			break
		}

		if image.Name == "" {
			state.SetError("name is nil")
			break
		}

		if image.Version == "" {
			state.SetError("version is nil")
			break
		}
	}

	return state
}


func (image *Image) State() ux.State {
	var state ux.State

	for range only.Once {
		var err error

		state = image.EnsureNotNil()
		if state.IsError() {
			break
		}

		if image.Summary == nil {
			ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
			//noinspection GoDeferInLoop
			defer cancel()

			df := filters.NewArgs()
			//df.Add("id", image.ID)
			df.Add("reference", fmt.Sprintf("%s/%s:%s", defaults.Organization, image.Name, image.Version))

			var images []types.ImageSummary
			images, err = image._Parent.Client.ImageList(ctx, types.ImageListOptions{All: true, Filters: df})
			if err != nil {
				state.SetError("gear list error: %s", err)
				break
			}
			if len(images) == 0 {
				state.SetWarning("no gears found")
				break
			}

			image.Summary = &images[0]
			image.Summary.ID = image.ID

			d := types.ImageInspect{}
			d, _, err = image._Parent.Client.ImageInspectWithRaw(ctx, image.ID)
			if err != nil {
				state.SetError("gear inspect error: %s", err)
				break
			}
			image.Details = d
		}

		if image.GearConfig == nil {
			image.GearConfig, state = gearJson.New(image.Summary.Labels["gearbox.json"])
			if state.IsError() {
				break
			}
		}

		if image.GearConfig.Meta.Organization != defaults.Organization {
			state.SetError("not a Gearbox container")
			break
		}

		//state.SetString("")
	}

	if state.IsError() {
		image.Summary = nil
	}

	return state
}


// Pull an image
// Pull an image, like docker pull:
func (image *Image) Pull() ux.State {
	var state ux.State

	for range only.Once {
		state = image.EnsureNotNil()
		if state.IsError() {
			break
		}

		var repo string
		if image.Version == "" {
			repo = fmt.Sprintf("gearboxworks/%s", image.Name)
		} else {
			repo = fmt.Sprintf("gearboxworks/%s:%s", image.Name, image.Version)
		}

		ctx := context.Background()
		//ctx, cancel := context.WithTimeout(context.Background(), Timeout * 1000)
		//defer cancel()

		//df := filters.NewArgs()
		//df.Add("name", "terminus")
		//var results []registry.SearchResult
		//results, err = image.client.ImageSearch(ctx, "", types.ImageSearchOptions{Filters:df})
		//for _, v := range results {
		//	fmt.Printf("%s - %s\n", v.Name, v.Description)
		//}

		var out io.ReadCloser
		var err error
		out, err = image._Parent.Client.ImagePull(ctx, repo, types.ImagePullOptions{All: false})
		if err != nil {
			state.SetError("error pulling gear: %s", err)
			break
		}

		//noinspection GoDeferInLoop
		defer out.Close()

		ux.Printf("pulling Gearbox gear: %s\n", image.Name)
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
				ux.PrintfOk("pulling Gearbox gear %s - OK\n", image.Name)
			} else if strings.HasPrefix(event.Status, "Status: Image is up to date for") {
				// up-to-date
				ux.PrintfOk("pulling Gearbox gear %s - image up to date\n", image.Name)
			} else {
				ux.PrintfWarning("pulling Gearbox gear %s - unknown\n", image.Name)
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
func (image *Image) ImageAuthPull() ux.State {
	var state ux.State

	for range only.Once {
		state = image.EnsureNotNil()
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
		//noinspection GoDeferInLoop
		defer cancel()

		out, err := image._Parent.Client.ImagePull(ctx, "alpine", types.ImagePullOptions{RegistryAuth: authStr})
		if err != nil {
			state.SetError("error pulling gear: %s", err)
			break
		}

		//noinspection GoDeferInLoop
		defer out.Close()

		_, _ = io.Copy(os.Stdout, out)
	}

	return state
}


// Remove containers
// Now that you know what containers exist, you can perform operations on them.
// This example stops all running containers.
func (image *Image) Remove() ux.State {
	var state ux.State

	for range only.Once {
		state = image.EnsureNotNil()
		if state.IsError() {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		options := types.ImageRemoveOptions {
			Force:         true,
			PruneChildren: true,
		}

		_, err := image._Parent.Client.ImageRemove(ctx, image.ID, options)
		if err != nil {
			state.SetError("error removing gear: %s", err)
			break
		}
	}

	return state
}
