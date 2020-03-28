package dockerClient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gb-launch/only"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/table"
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
	GearConfig GearConfig

	client     *client.Client
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

		if me.client == nil {
			err = errors.New("client is nil")
			break
		}
	}

	return err
}


// List all images
// List the images on your Engine, similar to docker image ls:
// func ImageList(f types.ImageListOptions) error {
func (me *Gear) ImageList(f string) error {
	var err error

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: ImageList(%s)\n", f)
		}

		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		df := filters.NewArgs()
		//if f != "" {
		//	df.Add("label", f)
		//}

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		var images []types.ImageSummary
		images, err = me.DockerClient.ImageList(ctx, types.ImageListOptions{All: true, Filters: df})
		if err != nil {
			break
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Class", "State", "Image", "Ports", "Size"})

		for _, i := range images {
			var gc *GearConfig
			gc, err = NewGearConfig(i.Labels["gearbox.json"])
			if err != nil {
				continue
			}

			if gc.Organization != Organization {
				continue
			}

			if i.RepoTags[0] == "<none>:<none>" {
				continue
			}

			if f != "" {
				if gc.Name != f {
					continue
				}
			}

			// foo := fmt.Sprintf("%s/%s", gc.Organization, gc.Name)
			t.AppendRow([]interface{}{gc.Class, gc.State, i.RepoTags[0], strings.Join(gc.Ports, " "), humanize.Bytes(uint64(i.Size))})
		}

		// t.AppendFooter(table.Row{"", "", "Total", 10000})
		t.Render()
		err = nil
	}

	return err
}


func (me *Gear) FindImage(gearName string, gearVersion string) (bool, error) {
	var ok bool
	var err error

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: FindImage(%s, %s)\n", gearName, gearVersion)
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

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		var images []types.ImageSummary
		images, err = me.DockerClient.ImageList(ctx, types.ImageListOptions{All: true})
		if err != nil {
			break
		}

		for _, i := range images {
			var gc *GearConfig
			gc, err = NewGearConfig(i.Labels["gearbox.json"])
			if err != nil {
				continue
			}

			if gc.Organization != Organization {
				continue
			}

			if i.RepoTags[0] == "<none>:<none>" {
				continue
			}

			if gc.Name != gearName {
				continue
			}

			if gearVersion == "latest" {
				gl := gc.Versions.GetLatest()
				if gl == "" {
					continue
				}
				// gearVersion = gl
			} else {
				if !gc.Versions.HasVersion(gearVersion) {
					continue
				}
			}

			me.Image.Name = gearName
			me.Image.Version = gearVersion
			me.Image.GearConfig = *gc
			me.Image.Summary = &i
			me.Image.ID = i.ID
			me.Image.client = me.DockerClient
			ok = true

			break
		}

		if err != nil {
			break
		}

		if me.Image.Summary == nil {
			break
		}

		me.Image.Details, _, err = me.DockerClient.ImageInspectWithRaw(ctx, me.Image.ID)
		if err != nil {
			break
		}

		err = me.Image.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return ok, err
}


// Search for an image in remote registry.
func (me *Gear) Search(gearName string, gearVersion string) error {
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		var repo string
		if gearVersion == "" {
			repo = fmt.Sprintf("gearboxworks/%s", gearName)
		} else {
			repo = fmt.Sprintf("gearboxworks/%s:%s", gearName, gearVersion)
		}

		ctx := context.Background()
		//ctx, cancel := context.WithTimeout(context.Background(), Timeout * 1000)
		//defer cancel()

		df := filters.NewArgs()
		//df.Add("name", "terminus")
		repo = gearName

		var images []registry.SearchResult
		images, err = me.DockerClient.ImageSearch(ctx, repo, types.ImageSearchOptions{Filters:df, Limit: 100})
		for _, v := range images {
			if !strings.HasPrefix(v.Name, "gearboxworks/") {
				continue
			}
			fmt.Printf("%s - %s\n", v.Name, v.Description)
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
		out, err = me.client.ImagePull(ctx, repo, types.ImagePullOptions{All: false})
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
func (me *Gear) ImageAuthPull() error {
	var err error

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: ImageAuthPull()\n")
		}

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

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		out, err := me.DockerClient.ImagePull(ctx, "alpine", types.ImagePullOptions{RegistryAuth: authStr})
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

		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()

		options := types.ImageRemoveOptions {
			Force:         false,
			PruneChildren: true,
		}

		_, err = me.client.ImageRemove(ctx, me.ID, options)
		if err != nil {
			break
		}
	}

	return err
}
