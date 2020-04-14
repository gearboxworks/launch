package githubClient

import (
	"fmt"
	"launch/only"
	"launch/ospaths"
	"launch/ux"
	"github.com/cavaliercoder/grab"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"os"
	"strings"
	"time"
)

const (
	Brandname = "Gearbox"
)

type GitHubRepo struct {
	Map             ReleasesMap
	Latest	        *Release
	Selected        *Release
	BaseDir         *ospaths.Dir
}
type ReleasesMap map[Version]*Release
type Version string

type Release struct {
	Version       Version
	File          *ospaths.File
	Size          int64
	Url           string
	Instance      *github.RepositoryRelease
	DlIndex       int
	IsDownloading bool
}

type ReleaseSelector struct {
	// These are considered to be AND-ed together.
	FromDate        time.Time
	UntilDate       time.Time
	SpecificVersion string
	RegexpVersion   string
	Latest			*bool
}


func New() (*GitHubRepo, ux.State) {
	var ret *GitHubRepo
	var state ux.State

	for range only.Once {
		p := ospaths.New("")

		me := GitHubRepo{}
		me.BaseDir = p.UserConfigDir.AddToPath("iso")
		me.Map = make(ReleasesMap)

		state = me.UpdateReleases()

		ret = &me

		//eblog.Debug(entity.VmBoxEntityName, "created new release structre")
	}

	//eblog.LogIfNil(ret, err)
	//eblog.LogIfError(entity.VmBoxEntityName, err)

	return ret, state
}


func (me *GitHubRepo) ShowReleases() ux.State {
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		fmt.Printf("Latest release: %v\n\n", me.Latest)
		for _, release := range me.Map {
			fmt.Printf("Assets for release:	%v\n", release.Instance.GetName())
			fmt.Printf("UploadURL: 			%v\n", release.Instance.GetUploadURL())
			fmt.Printf("ZipballURL: 			%v\n", release.Instance.GetZipballURL())
			fmt.Printf("TarballURL: 			%v\n", release.Instance.GetTarballURL())
			fmt.Printf("Body: 				%v\n", release.Instance.GetBody())
			fmt.Printf("AssetsURL: 			%v\n", release.Instance.GetAssetsURL())
			fmt.Printf("URL: 				%v\n", release.Instance.GetURL())
			fmt.Printf("HTMLURL:				%v\n", release.Instance.GetHTMLURL())

			for _, asset := range release.Instance.Assets {
				fmt.Printf("	Name:				%v\n", asset.GetName())
				fmt.Printf("	ID:					%v\n", asset.GetID())
				fmt.Printf("	URL:					%v\n", asset.GetURL())
				fmt.Printf("	Size:				%v\n", asset.GetSize())
				fmt.Printf("	CreatedAt:			%v\n", asset.GetCreatedAt())
				fmt.Printf("	UpdatedAt:			%v\n", asset.GetUpdatedAt())
				fmt.Printf("	BrowserDownloadURL:	%v\n", asset.GetBrowserDownloadURL())
				fmt.Printf("	State:				%v\n", asset.GetState())
				fmt.Printf("	ContentType:			%v\n", asset.GetContentType())
				fmt.Printf("	DownloadCount:		%v\n", asset.GetDownloadCount())
				fmt.Printf("	NodeID:				%v\n", asset.GetNodeID())
			}
		}

		//eblog.Debug(entity.VmBoxEntityName, "Showing all ISO releases. Latest == %s", me.Latest)
	}

	//eblog.LogIfNil(me, err)
	//eblog.LogIfError(entity.VmBoxEntityName, err)

	return state
}


func (me *Release) ShowRelease() ux.State {
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		if me.Instance.Name == nil {
			state.SetError("no release version specified")
			break
		}

		fmt.Printf("Assets for release:	%v\n", *me.Instance.Name)
		for _, asset := range me.Instance.Assets {
			fmt.Printf("	Name:				%v\n", asset.GetName())
			fmt.Printf("	ID:					%v\n", asset.GetID())
			fmt.Printf("	URL:					%v\n", asset.GetURL())
			fmt.Printf("	Size:				%v\n", asset.GetSize())
			fmt.Printf("	CreatedAt:			%v\n", asset.GetCreatedAt())
			fmt.Printf("	UpdatedAt:			%v\n", asset.GetUpdatedAt())
			fmt.Printf("	BrowserDownloadURL:	%v\n", asset.GetBrowserDownloadURL())
			fmt.Printf("	State:				%v\n", asset.GetState())
			fmt.Printf("	ContentType:			%v\n", asset.GetContentType())
			fmt.Printf("	DownloadCount:		%v\n", asset.GetDownloadCount())
			fmt.Printf("	NodeID:				%v\n", asset.GetNodeID())
		}

		//eblog.Debug(entity.VmBoxEntityName, "Showing ISO release for v%s", *me.Instance.Name)
	}

	//eblog.LogIfNil(me, err)
	//eblog.LogIfError(entity.VmBoxEntityName, err)

	return state
}


func (me *GitHubRepo) UpdateReleases() ux.State {

	var rm = make(ReleasesMap)
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		if me.BaseDir == nil {
			p := ospaths.New("")
			me.BaseDir = p.UserConfigDir.AddToPath("iso")
		}

		me.Map = rm

		client := github.NewClient(nil)
		//ctx := context.Background()
		opt := &github.ListOptions{}

		releases, _, err := client.Repositories.ListReleases(context.Background(), "gearboxworks", "docker-os", opt)
		if err != nil {
			state.SetError("can't fetch GitHub releases")
			break
		}

		findFirst := true
		for _, rel := range releases {
			if rel == nil {
				continue
			}

			name := Version(rel.GetName())

			release := Release{
				Version: name,
				Url: "",
				Instance: rel,
			}

			// rm[name].Url/File - Find the first ISO asset.
			for _, asset := range rel.Assets {
				if strings.HasSuffix(asset.GetBrowserDownloadURL(), ".iso") {
					// Return the first ISO found.
					release.Url = asset.GetBrowserDownloadURL()
					release.File = me.BaseDir.AddFileToPath(asset.GetName())
					release.Size = int64(asset.GetSize())
					break
				}
			}

			// rm[name].Version - Copy version name.
			rm[name] = &release

			// rm.Latest - Find first version and select as 'latest'.
			if findFirst {
				me.Latest = &release
				findFirst = false
			}
		}

		//if findFirst == true {
		//	// If we never found a "first", then there was no data.
		//	// So don't update the map.
		//}

		me.Map = rm

		//eblog.Debug(entity.VmBoxEntityName, "Fetching ISO releases. Latest == %s", me.Latest)
	}

	//eblog.LogIfNil(me, err)
	//eblog.LogIfError(entity.VmBoxEntityName, err)

	return state
}


/*
Updates the following:
   me.VmIsoVersion    string
   me.VmIsoFile       string
   me.VmIsoUrl 		string
   me.VmIsoRelease    Release
*/
func (me *GitHubRepo) SelectRelease(selector ReleaseSelector) (*Release, ux.State) {
	var r *Release
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		//err = me.UpdateReleases()
		//if err != nil {
		//	break
		//}

		// For now just select the latest.
		me.Selected = me.Latest
		r = me.Selected

		//eblog.Debug(entity.VmBoxEntityName, "selecting the latest release == %s", me.Latest.Version)
	}

	//eblog.LogIfNil(me, err)
	//eblog.LogIfError(entity.VmBoxEntityName, err)

	return r, state
}


func (me *Release) GetIso() ux.State {
	var state ux.State

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		if me.File.String() == "" {
			state.SetError(fmt.Sprintf("no Gearbox OS iso file defined VmIsoUrl:%s VmIsoFile:%s", me.Url, me.File.String()))
			break
		}

		if me.Url == "" {
			state.SetError(fmt.Sprintf("no Gearbox OS iso url defined VmIsoUrl:%s VmIsoFile:%s", me.Url, me.File.String()))
			break
		}


		var numb int
		numb, state = me.IsIsoFilePresent()
		if numb != IsoFileNeedsToDownload {
			break
		}


		// Start download
		me.DlIndex = 0
		me.IsDownloading = true
		client := grab.NewClient()
		req, _ := grab.NewRequest(me.File.String(), me.Url)
		fmt.Sprintf("downloading ISO from URL %s", req.URL().String())
		resp := client.Do(req)
		// fmt.Printf("  %v\n", resp.HTTPResponse.Status)
		fmt.Printf("%s VM: Downloading ISO from '%s' to '%s'. Size:%d\n",
			Brandname,
			me.Url,
			me.File.String(),
			resp.Size)


		// start UI loop
		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()

		Loop:
			for {
				select {
					case <-t.C:
						me.DlIndex = int(100*resp.Progress())
						//me.publishDownloadState()
						//fmt.Printf("Downloading '%s' transferred %v / %v bytes (%d%%)\n", me.File.String(), resp.BytesComplete(), resp.Size, me.DlIndex)
						fmt.Printf("%s VM: Downloading ISO - %d%% complete.\r",
							Brandname,
							me.DlIndex)

					case <-resp.Done:
						// download is complete
						break Loop
				}
			}

		// check for errors
		if err := resp.Err(); err != nil {
			fmt.Printf("\nDownload failed\n")
			state.SetError(fmt.Sprintf("ISO download failed VmIsoUrl:%s VmIsoFile:%s", me.Url, me.File.String()))
			break
		}
		fmt.Printf("%s VM: Downloaded ISO completed OK.\n",
			Brandname,
		)


		//eblog.Debug(entity.VmBoxEntityName, "ISO fetched from '%s' and saved to '%s'. Size:%d", me.Url, me.File.String(), resp.Size)
		me.DlIndex = 100
		//me.publishDownloadState()
		me.IsDownloading = false
	}

	//eblog.LogIfNil(me, err)
	//eblog.LogIfError(entity.VmBoxEntityName, err)

	return state
}


//func (me *Release) publishDownloadState() {
//
//	client := messages.MessageAddress(entity.VmUpdateEntityName)
//	state := states.New(&client, &client, entity.VmBoxEntityName)
//	state.SetWant("100%")
//	state.SetCurrent(states.State(fmt.Sprintf("%d%%", me.DlIndex)))
//
//	f := messages.MessageAddress(states.ActionUpdate)
//	msg := f.ConstructMessage(entity.BroadcastEntityName, states.ActionStatus, state.ToMessageText())
//	_ = me.channels.Publish(msg)
//}


const IsoFileNeedsToDownload	= 0
const IsoFileIsDownloading		= 1
const IsoFileDownloaded			= 2
func (me *Release) IsIsoFilePresent() (int, ux.State) {
	var state ux.State
	var ret int
	var stat os.FileInfo

	for range only.Once {
		state = me.EnsureNotNil()
		if state.IsError() {
			break
		}

		if me.File.String() == "" {
			state.SetError( fmt.Sprintf("no Gearbox OS iso file defined VmIsoUrl:%s VmIsoFile:%s", me.Url, me.File.String()))
			break
		}

		var err error
		stat, err = os.Stat(me.File.String())
		if os.IsNotExist(err) {
			state.SetError("ISO file needs to download from GitHub VmIsoUrl:%s VmIsoFile:%s", me.Url, me.File.String())
			ret = IsoFileNeedsToDownload
			break
		}

		if me.IsDownloading {
			state.SetError("ISO file still downloading VmIsoUrl:%s VmIsoFile:%s Percent:%d", me.Url, me.File.String(), me.DlIndex)
			ret = IsoFileIsDownloading
			break
		}

		if stat.Size() != me.Size {
			state.SetError("ISO file needs to re-download from GitHub VmIsoUrl:%s VmIsoFile:%s", me.Url, me.File.String())
			ret = IsoFileNeedsToDownload
			break
		}

		//if me.DlIndex < 100 {
		//	err = errors.New("ISO file needs to re-download from GitHub VmIsoUrl:%s VmIsoFile:%s", me.Url, me.File.String())
		//	ret = IsoFileNeedsToDownload
		//	break
		//}

		ret = IsoFileDownloaded
		me.DlIndex = 100
		//eblog.Debug(entity.VmBoxEntityName, "ISO already fetched from '%s' and saved to '%s'", me.Url, me.File.String())
	}

	//eblog.LogIfNil(me, err)
	//eblog.LogIfError(entity.VmBoxEntityName, err)

	return ret, state
}


func (me *GitHubRepo) EnsureNotNil() ux.State {
	var state ux.State

	for range only.Once {
		if me == nil {
			state.SetError("releases is nil")
			break
		}
	}

	return state
}

func EnsureReleasesNotNil(me *GitHubRepo) ux.State {
	return me.EnsureNotNil()
}


func (me *ReleasesMap) EnsureNotNil() ux.State {
	var state ux.State

	for range only.Once {
		if me == nil {
			state.SetError("Release is nil")
			break
		}
	}

	return state
}

func EnsureReleasesMapNotNil(me *ReleasesMap) ux.State {
	return me.EnsureNotNil()
}


func (me *Release) EnsureNotNil() ux.State {
	var state ux.State

	for range only.Once {
		if me == nil {
			state.SetError("Release is nil")
			break
		}
	}

	return state
}

func EnsureReleaseNotNil(me *Release) ux.State {
	return me.EnsureNotNil()
}




//func EnsureReleaseNotNil(rm *Release) (sts status.Status) {
//	if rm == nil {
//		sts = status.Fail(&status.Args{
//			Message: "unexpected error",
//			Help:    help.ContactSupportHelp(), // @TODO need better support here
//			Data:    VmStateUnknown,
//		})
//	}
//
//	return sts
//}

//type ReleaseAsset struct {
//	ID                 *int64     `json:"id,omitempty"`
//	URL                *string    `json:"url,omitempty"`
//	Name               *string    `json:"name,omitempty"`
//	Label              *string    `json:"label,omitempty"`
//	State              *string    `json:"state,omitempty"`
//	ContentType        *string    `json:"content_type,omitempty"`
//	Size               *int       `json:"size,omitempty"`
//	DownloadCount      *int       `json:"download_count,omitempty"`
//	CreatedAt          *Timestamp `json:"created_at,omitempty"`
//	UpdatedAt          *Timestamp `json:"updated_at,omitempty"`
//	BrowserDownloadURL *string    `json:"browser_download_url,omitempty"`
//	Uploader           *User      `json:"uploader,omitempty"`
//	NodeID             *string    `json:"node_id,omitempty"`
//}
//
//type RepositoryRelease struct {
//	ID              *int64         `json:"id,omitempty"`
//	TagName         *string        `json:"tag_name,omitempty"`
//	TargetCommitish *string        `json:"target_commitish,omitempty"`
//	Name            *string        `json:"name,omitempty"`
//	Body            *string        `json:"body,omitempty"`
//	Draft           *bool          `json:"draft,omitempty"`
//	Prerelease      *bool          `json:"prerelease,omitempty"`
//	CreatedAt       *Timestamp     `json:"created_at,omitempty"`
//	PublishedAt     *Timestamp     `json:"published_at,omitempty"`
//	URL             *string        `json:"url,omitempty"`
//	HTMLURL         *string        `json:"html_url,omitempty"`
//	AssetsURL       *string        `json:"assets_url,omitempty"`
//	Assets          []ReleaseAsset `json:"assets,omitempty"`
//	UploadURL       *string        `json:"upload_url,omitempty"`
//	ZipballURL      *string        `json:"zipball_url,omitempty"`
//	TarballURL      *string        `json:"tarball_url,omitempty"`
//	Author          *User          `json:"author,omitempty"`
//	NodeID          *string        `json:"node_id,omitempty"`
//}
//
//
//Data returned:
//
//release.ID=0xc000289538
//release.TagName=0xc0002964c0
//release.TargetCommitish=0xc0002964d0
//release.Name=0xc0002964e0
//release.Body=0xc000296770
//release.Draft=0xc00028955b
//release.Prerelease=0xc00028957d
//release.CreatedAt=2019-05-23 02:34:10 +0000 UTC
//release.PublishedAt=2019-05-23 02:43:04 +0000 UTC
//release.URL=0xc000296470
//release.HTMLURL=0xc0002964a0
//release.AssetsURL=0xc000296480
//release.Assets=[github.ReleaseAsset{
//	ID:12825393,
//	URL:"https://api.github.com/repos/gearboxworks/gearbox-os/releases/assets/12825393",
//	Name:"gearbox-0.5.0.iso",
//	State:"uploaded",
//	ContentType:"application/octet-stream",
//	Size:67108864,
//	DownloadCount:0,
//	CreatedAt:github.Timestamp{2019-05-23 02:37:48 +0000 UTC},
//	UpdatedAt:github.Timestamp{2019-05-23 02:42:56 +0000 UTC},
//	BrowserDownloadURL:"https://github.com/gearboxworks/gearbox-os/releases/download/0.5.0/gearbox-0.5.0.iso",
//	Uploader:github.User{
//		Login:"MickMake",
//		ID:17118367,
//		NodeID:"MDQ6VXNlcjE3MTE4MzY3",
//		AvatarURL:"https://avatars0.githubusercontent.com/u/17118367?v=4",
//		HTMLURL:"https://github.com/MickMake",
//		GravatarID:"",
//		Type:"User",
//		SiteAdmin:false,
//		URL:"https://api.github.com/users/MickMake",
//		EventsURL:"https://api.github.com/users/MickMake/events{/privacy}",
//		FollowingURL:"https://api.github.com/users/MickMake/following{/other_user}",
//		FollowersURL:"https://api.github.com/users/MickMake/followers",
//		GistsURL:"https://api.github.com/users/MickMake/gists{/gist_id}",
//		OrganizationsURL:"https://api.github.com/users/MickMake/orgs",
//		ReceivedEventsURL:"https://api.github.com/users/MickMake/received_events",
//		ReposURL:"https://api.github.com/users/MickMake/repos",
//		StarredURL:"https://api.github.com/users/MickMake/starred{/owner}{/repo}",
//		SubscriptionsURL:"https://api.github.com/users/MickMake/subscriptions"
//		},
//	NodeID:"MDEyOlJlbGVhc2VBc3NldDEyODI1Mzkz"
//	}]
//release.UploadURL=0xc000296490
//release.ZipballURL=0xc000296760
//release.TarballURL=0xc000296750
//release.Author=github.User{Login:"MickMake", ID:17118367, NodeID:"MDQ6VXNlcjE3MTE4MzY3", AvatarURL:"https://avatars0.githubusercontent.com/u/17118367?v=4", HTMLURL:"https://github.com/MickMake", GravatarID:"", Type:"User", SiteAdmin:false, URL:"https://api.github.com/users/MickMake", EventsURL:"https://api.github.com/users/MickMake/events{/privacy}", FollowingURL:"https://api.github.com/users/MickMake/following{/other_user}", FollowersURL:"https://api.github.com/users/MickMake/followers", GistsURL:"https://api.github.com/users/MickMake/gists{/gist_id}", OrganizationsURL:"https://api.github.com/users/MickMake/orgs", ReceivedEventsURL:"https://api.github.com/users/MickMake/received_events", ReposURL:"https://api.github.com/users/MickMake/repos", StarredURL:"https://api.github.com/users/MickMake/starred{/owner}{/repo}", SubscriptionsURL:"https://api.github.com/users/MickMake/subscriptions"}
//release.NodeID=0xc0002964b0
//
//
//type Release struct {
//	Name string
//	UploadURL string
//	ZipballURL string
//	TarballURL string
//	Body string
//	AssetsURL string
//	URL string
//	HTMLURL string
//	Name string
//    Assets
//}
//type Releases []Release
//
//type Asset struct {
//      Name
//      ID
//      URL
//      Size
//      CreatedAt
//      UpdatedAt
//      BrowserDownloadURL
//      State
//      ContentType
//      DownloadCount
//      NodeID
//}
//type Assets []Asset
//
//
//
//
//   Assets for release:	0.5.0
//   UploadURL: 			https://uploads.github.com/repos/gearboxworks/gearbox-os/releases/17531887/assets{?name,label}
//   ZipballURL: 			https://api.github.com/repos/gearboxworks/gearbox-os/zipball/0.5.0
//   TarballURL: 			https://api.github.com/repos/gearboxworks/gearbox-os/tarball/0.5.0
//   Body: 				0.5.0 pre-release
//   AssetsURL: 			https://api.github.com/repos/gearboxworks/gearbox-os/releases/17531887/assets
//   URL: 				https://api.github.com/repos/gearboxworks/gearbox-os/releases/17531887
//   HTMLURL:				https://github.com/gearboxworks/gearbox-os/releases/tag/0.5.0
//   foo: 				0.5.0
//   Name:				gearbox-0.5.0.iso
//   ID:					12825393
//   URL:					https://api.github.com/repos/gearboxworks/gearbox-os/releases/assets/12825393
//   Size:				67108864
//   CreatedAt:			2019-05-23 02:37:48 +0000 UTC
//   UpdatedAt:			2019-05-23 02:42:56 +0000 UTC
//   BrowserDownloadURL:	https://github.com/gearboxworks/gearbox-os/releases/download/0.5.0/gearbox-0.5.0.iso
//   State:				uploaded
//   ContentType:			application/octet-stream
//   DownloadCount:		0
//   NodeID:				MDEyOlJlbGVhc2VBc3NldDEyODI1Mzkz
//
//	fmt.Printf(`
//		release.ID=%v
//		release.TagName=%v
//		release.TargetCommitish=%v
//		release.Name=%v
//		release.Body=%v
//		release.Draft=%v
//		release.Prerelease=%v
//		release.CreatedAt=%v
//		release.PublishedAt=%v
//		release.URL=%v
//		release.HTMLURL=%v
//		release.AssetsURL=%v
//		release.Assets=%v
//		release.UploadURL=%v
//		release.ZipballURL=%v
//		release.TarballURL=%v
//		release.Author=%v
//		release.NodeID=%v\n`,
//		release.ID,
//		release.TagName,
//		release.TargetCommitish,
//		release.Name,
//		release.Body,
//		release.Draft,
//		release.Prerelease,
//		release.CreatedAt,
//		release.PublishedAt,
//		release.URL,
//		release.HTMLURL,
//		release.AssetsURL,
//		release.Assets,
//		release.UploadURL,
//		release.ZipballURL,
//		release.TarballURL,
//		release.Author,
//		release.NodeID,
//		)

