package dockerClient

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/table"
	"launch/defaults"
	"launch/gear/gearJson"
	"launch/only"
	"launch/ux"
	"os"
	"strings"
)


// List all images
// List the images on your Engine, similar to docker image ls:
// func ImageList(f types.ImageListOptions) error {
func (gear *DockerGear) ImageList(f string) (int, *ux.State) {
	var count int
	if state := gear.IsNil(); state.IsError() {
		return 0, state
	}

	for range only.Once {
		df := filters.NewArgs()
		//if f != "" {
		//	df.Add("label", f)
		//}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		images, err := gear.Client.ImageList(ctx, types.ImageListOptions{All: true, Filters: df})
		if err != nil {
			gear.State.SetError("gear image list error: %s", err)
			break
		}

		ux.PrintfCyan("Downloaded Gearbox images: ")
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Class", "Image", "Ports", "Size"})

		for _, i := range images {
			var gc *gearJson.GearConfig
			gc = gearJson.New(i.Labels["gearbox.json"])
			if gc.State.IsError() {
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

		gear.State.ClearError()
		count = t.Length()
		if count == 0 {
			ux.PrintfYellow("None found\n")
			break
		}

		ux.PrintflnGreen("%d found", count)
		t.Render()
		ux.PrintflnBlue("")
	}

	return count, gear.State
}


func (gear *DockerGear) FindImage(gearName string, gearVersion string) (bool, *ux.State) {
	var ok bool
	if state := gear.IsNil(); state.IsError() {
		return false, state
	}

	for range only.Once {
		if gearName == "" {
			gear.State.SetError("empty gear name")
			break
		}

		if gearVersion == "" {
			gearVersion = "latest"
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel()

		images, err := gear.Client.ImageList(ctx, types.ImageListOptions{All: true})
		if err != nil {
			gear.State.SetError("gear image list error: %s", err)
			break
		}

		if len(images) == 0 {
			break
		}

		// Start out with "not found". Will be cleared if found or error occurs.
		gear.State.SetWarning("Gear image '%s:%s' doesn't exist.", gearName, gearVersion)

		for _, i := range images {
			var gc *gearJson.GearConfig
			ok, gc = MatchImage(&i, defaults.Organization, gearName, gearVersion)
			if !ok {
				continue
			}

			gear.Image.Name = gearName
			gear.Image.Version = gearVersion
			gear.Image.GearConfig = gc
			gear.Image.Summary = &i
			gear.Image.ID = i.ID
			gear.Image.State = gear.Image.State.EnsureNotNil()
			//gear.Image.client = gear.DockerClient
			ok = true

			break
		}

		if gear.State.IsNotOk() {
			if !ok {
				gear.State.ClearError()
			}
			break
		}

		if gear.Image.Summary == nil {
			break
		}

		ctx2, cancel2 := context.WithTimeout(context.Background(), defaults.Timeout)
		//noinspection GoDeferInLoop
		defer cancel2()
		gear.Image.Details, _, err = gear.Client.ImageInspectWithRaw(ctx2, gear.Image.ID)
		if err != nil {
			gear.State.SetError("error inspecting gear: %s", err)
			break
		}

		gear.State.SetOk("found image")
	}

	return ok, gear.State
}


// Search for an image in remote registry.
func (gear *DockerGear) Search(gearName string, gearVersion string) *ux.State {
	if state := gear.IsNil(); state.IsError() {
		return state
	}

	for range only.Once {
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

		images, err := gear.Client.ImageSearch(ctx, repo, types.ImageSearchOptions{Filters: df, Limit: 100})
		if err != nil {
			gear.State.SetError("gear image search error: %s", err)
			break
		}

		for _, v := range images {
			if !strings.HasPrefix(v.Name, "gearboxworks/") {
				continue
			}
			fmt.Printf("%s - %s\n", v.Name, v.Description)
		}
	}

	return gear.State
}


func MatchImage(m *types.ImageSummary, gearOrg string, gearName string, gearVersion string) (bool, *gearJson.GearConfig) {
	var ok bool
	var gc *gearJson.GearConfig

	for range only.Once {
		if MatchTag("<none>:<none>", m.RepoTags) {
			ok = false
			break
		}

		gc = gearJson.New(m.Labels["gearbox.json"])
		if gc.State.IsError() {
			ok = false
			break
		}

		if gc.Meta.Organization != defaults.Organization {
			ok = false
			break
		}

		tagCheck := fmt.Sprintf("%s/%s:%s", gearOrg, gearName, gearVersion)
		if !MatchTag(tagCheck, m.RepoTags) {
			ok = false
			break
		}

		if gc.Meta.Name != gearName {
			if !defaults.RunAs.AsLink {
				ok = false
				break
			}

			cs := gc.MatchCommand(gearName)
			if cs == nil {
				ok = false
				break
			}

			gearName = gc.Meta.Name
		}

		if !gc.Versions.HasVersion(gearVersion) {
			ok = false
			break
		}

		if gearVersion == "latest" {
			gl := gc.Versions.GetLatest()
			if gearVersion != "" {
				gearVersion = gl
			}
		}
		for range only.Once {
			if m.Labels["gearbox.version"] == gearVersion {
				ok = true
				break
			}

			if m.Labels["container.majorversion"] == gearVersion {
				ok = true
				break
			}

			ok = false
		}
		break
	}

	return ok, gc
}


func MatchTag(match string, tags []string) bool {
	var ok bool

	for _, s := range tags {
		if s == match {
			ok = true
			break
		}
	}

	return ok
}
