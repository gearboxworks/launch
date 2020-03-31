package dockerClient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gb-launch/defaults"
	"gb-launch/gear/gearJson"
	"gb-launch/only"
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


func (me *Image) EnsureNotNil() error {
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
	}

	return err
}


// Pull an image
// Pull an image, like docker pull:
func (me *Image) Pull() error {
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
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
		out, err = me._Parent.Client.ImagePull(ctx, repo, types.ImagePullOptions{All: false})
		if err != nil {
			break
		}

		defer out.Close()

		fmt.Printf("Pulling Gearbox container: %s\n", me.Name)
		d := json.NewDecoder(out)
		var event *PullEvent
		for {
			if err := d.Decode(&event); err != nil {
				if err == io.EOF {
					break
				}

				panic(err)
			}

			// fmt.Printf("EVENT: %+v\n", event)
			fmt.Printf("%+v\r", event.Progress)
		}
		fmt.Printf("\n%s\n", event.Status)

		// Latest event for new image
		// EVENT: {Status:Status: Downloaded newer image for busybox:latest Error: Progress:[==================================================>]  699.2kB/699.2kB ProgressDetail:{Current:699243 Total:699243}}
		// Latest event for up-to-date image
		// EVENT: {Status:Status: Image is up to date for busybox:latest Error: Progress: ProgressDetail:{Current:0 Total:0}}
		if event != nil {
			if strings.Contains(event.Status, fmt.Sprintf("Downloaded newer image for %s", me.Name)) {
				// new
				fmt.Println("\nnew")
			}

			if strings.Contains(event.Status, fmt.Sprintf("Image is up to date for %s", me.Name)) {
				// up-to-date
				fmt.Println("\nup-to-date")
			}
		}

		//buf := new(bytes.Buffer)
		//_, err = buf.ReadFrom(out)
		//fmt.Printf("%s", buf.String())
		//_, _ = io.Copy(os.Stdout, out)
	}

	return err
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
func (me *Image) ImageAuthPull() error {
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		authConfig := types.AuthConfig{
			Username: "username",
			Password: "password",
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			break
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		out, err := me._Parent.Client.ImagePull(ctx, "alpine", types.ImagePullOptions{RegistryAuth: authStr})
		if err != nil {
			break
		}

		defer out.Close()

		_, _ = io.Copy(os.Stdout, out)
	}

	return err
}


// Remove containers
// Now that you know what containers exist, you can perform operations on them.
// This example stops all running containers.
func (me *Image) Remove() error {
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		options := types.ImageRemoveOptions {
			Force:         true,
			PruneChildren: true,
		}

		_, err = me._Parent.Client.ImageRemove(ctx, me.ID, options)
		if err != nil {
			break
		}
	}

	return err
}
