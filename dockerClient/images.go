package dockerClient

import (
	"context"
	"fmt"
	"launch/defaults"
	"launch/gear/gearJson"
	"launch/only"
	"launch/ux"
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
func (me *DockerGear) ImageList(f string) ux.State {
	var state ux.State

	for range only.Once {
		var err error

		if me.Debug {
			fmt.Printf("DEBUG: ImageList(%s)\n", f)
		}

		state = me.EnsureNotNil()
		if state.IsError() {
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
			state.SetError("gear image list error: %s", err)
			break
		}

		ux.PrintfCyan("\nConfigured Gearbox images:\n")
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		//t.AppendHeader(table.Row{"Class", "State", "Image", "Ports", "Size"})
		t.AppendHeader(table.Row{"Class", "Image", "Ports", "Size"})

		for _, i := range images {
			var gc *gearJson.GearConfig
			gc, state = gearJson.New(i.Labels["gearbox.json"])
			if state.IsError() {
				continue
			}

			if gc.Meta.Organization != defaults.Organization {
				continue
			}

			if len(i.RepoTags) == 0 {
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
			t.AppendRow([]interface{}{
				ux.SprintfWhite(gc.Meta.Class),
				//ux.SprintfWhite(gc.Meta.State),
				ux.SprintfWhite(i.RepoTags[0]),
				ux.SprintfWhite(gc.Build.Ports.ToString()),
				ux.SprintfWhite(humanize.Bytes(uint64(i.Size))),
			})
		}

		// t.AppendFooter(table.Row{"", "", "Total", 10000})
		t.Render()
		state.ClearError()
	}

	return state
}


func (me *DockerGear) FindImage(gearName string, gearVersion string) (bool, ux.State) {
	var ok bool
	var state ux.State
	//var err error

	for range only.Once {
		if me.Debug {
			fmt.Printf("DEBUG: FindImage(%s, %s)\n", gearName, gearVersion)
		}

		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		if gearName == "" {
			state.SetError("empty gear name")
			break
		}

		if gearVersion == "" {
			gearVersion = "latest"
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		defer cancel()

		var images []types.ImageSummary
		var err error
		images, err = me.Client.ImageList(ctx, types.ImageListOptions{All: true})
		if err != nil {
			state.SetError("gear image list error: %s", err)
			break
		}

		if len(images) == 0 {
			break
		}

		for _, i := range images {
			var gc *gearJson.GearConfig
			gc, state = gearJson.New(i.Labels["gearbox.json"])
			if state.IsError() {
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
			me.Image.GearConfig = gc
			me.Image.Summary = &i
			me.Image.ID = i.ID
			//me.Image.client = me.DockerClient
			ok = true

			break
		}

		if state.IsError() {
			break
		}

		if me.Image.Summary == nil {
			break
		}

		me.Image.Details, _, err = me.Client.ImageInspectWithRaw(ctx, me.Image.ID)
		if err != nil {
			state.SetError("error inspecting gear: %s", err)
			break
		}

		state = me.Image.EnsureNotNil()
		if state.IsError() {
			break
		}

		state.SetOk("found image")
	}

	return ok, state
}


// Search for an image in remote registry.
func (me *DockerGear) Search(gearName string, gearVersion string) ux.State {
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
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
		var err error
		images, err = me.Client.ImageSearch(ctx, repo, types.ImageSearchOptions{Filters:df, Limit: 100})
		if err != nil {
			state.SetError("gear image search error: %s", err)
			break
		}

		for _, v := range images {
			if !strings.HasPrefix(v.Name, "gearboxworks/") {
				continue
			}
			fmt.Printf("%s - %s\n", v.Name, v.Description)
		}
	}

	return state
}
