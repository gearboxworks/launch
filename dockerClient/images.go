package dockerClient

import (
	"context"
	"errors"
	"fmt"
	"gb-launch/defaults"
	"gb-launch/gear/gearJson"
	"gb-launch/only"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/registry"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/table"
	"os"
	"strings"
)


// List all images
// List the images on your Engine, similar to docker image ls:
// func ImageList(f types.ImageListOptions) error {
func (me *DockerGear) ImageList(f string) error {
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

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		var images []types.ImageSummary
		images, err = me.Client.ImageList(ctx, types.ImageListOptions{All: true, Filters: df})
		if err != nil {
			break
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Class", "State", "Image", "Ports", "Size"})

		for _, i := range images {
			var gc *gearJson.GearConfig
			gc, err = gearJson.New(i.Labels["gearbox.json"])
			if err != nil {
				continue
			}

			if gc.Meta.Organization != defaults.Organization {
				continue
			}

			if i.RepoTags[0] == "<none>:<none>" {
				continue
			}

			if f != "" {
				if gc.Meta.Name != f {
					continue
				}
			}

			// foo := fmt.Sprintf("%s/%s", gc.Organization, gc.Name)
			t.AppendRow([]interface{}{gc.Meta.Class, gc.Meta.State, i.RepoTags[0], strings.Join(gc.Build.Ports, " "), humanize.Bytes(uint64(i.Size))})
		}

		// t.AppendFooter(table.Row{"", "", "Total", 10000})
		t.Render()
		err = nil
	}

	return err
}


func (me *DockerGear) FindImage(gearName string, gearVersion string) (bool, error) {
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

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		var images []types.ImageSummary
		images, err = me.Client.ImageList(ctx, types.ImageListOptions{All: true})
		if err != nil {
			break
		}

		for _, i := range images {
			var gc *gearJson.GearConfig
			gc, err = gearJson.New(i.Labels["gearbox.json"])
			if err != nil {
				continue
			}

			if gc.Meta.Organization != defaults.Organization {
				continue
			}

			if i.RepoTags[0] == "<none>:<none>" {
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
			//me.Image.client = me.DockerClient
			ok = true

			break
		}

		if err != nil {
			break
		}

		if me.Image.Summary == nil {
			break
		}

		me.Image.Details, _, err = me.Client.ImageInspectWithRaw(ctx, me.Image.ID)
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
func (me *DockerGear) Search(gearName string, gearVersion string) error {
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
		images, err = me.Client.ImageSearch(ctx, repo, types.ImageSearchOptions{Filters:df, Limit: 100})
		for _, v := range images {
			if !strings.HasPrefix(v.Name, "gearboxworks/") {
				continue
			}
			fmt.Printf("%s - %s\n", v.Name, v.Description)
		}

	}

	return err
}
