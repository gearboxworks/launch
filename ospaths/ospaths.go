// A simple wrapper around osbridge.OsBridger.
// This makes it much easier to separate the EventBroker code into it's own package later on.
package ospaths

import (
	"fmt"
	"github.com/gearboxworks/go-osbridge"
	"github.com/gearboxworks/go-status/only"
	"github.com/getlantern/errors"
	"launch/defaults"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	DefaultBaseDir = "app/dist/eventbroker"

	defaultLogBaseDir = "logs"
	defaultEtcBaseDir = "etc"
)

type Name string

type Path struct {
	Dir Dir
	File File
}
type Paths []Path

type Dir  string
type Dirs []Dir

type File string
type Files []File

type BasePaths struct {
	UserHomeDir           Dir
	ProjectBaseDir      Dir
	UserConfigDir         Dir
	AdminRootDir          Dir
	CacheDir              Dir
	EventBrokerDir        Dir
	EventBrokerWorkingDir Dir
	EventBrokerLogDir     Dir
	EventBrokerEtcDir     Dir
	LocalDir              Dir

	osBridger             osbridge.OsBridger
	mutex                 sync.RWMutex
}
//type OsBridge     osbridge.OsBridger



func New(subdir string) *BasePaths {

	var ret BasePaths

	if subdir == "" {
		subdir = DefaultBaseDir
	}

	//foo := ret.osBridger.GetOsBridge(global.Brandname, global.UserDataPath)
	//
	//fmt.Printf("TEST: %s\n", foo)
	//foo.GetProjectDir()

	//ret.osBridger = GetOsBridge(defaults.BrandName, Dir(defaults.DefaultProject))
	ret.osBridger = GetOsBridge(defaults.BrandName, defaults.DefaultProject)

	ret.UserHomeDir = Dir(ret.osBridger.GetUserHomeDir())
	ret.ProjectBaseDir = Dir(ret.osBridger.GetProjectDir())
	ret.UserConfigDir = Dir(ret.osBridger.GetUserConfigDir())
	ret.AdminRootDir = Dir(ret.osBridger.GetAdminRootDir())
	ret.CacheDir = Dir(ret.osBridger.GetCacheDir())

	ret.LocalDir = Dir(filepath.FromSlash("/usr/local"))
	ret.EventBrokerDir = *ret.UserConfigDir.AddToPath(subdir)
	ret.EventBrokerLogDir = *ret.EventBrokerDir.AddToPath(defaultLogBaseDir)
	ret.EventBrokerEtcDir = *ret.EventBrokerDir.AddToPath(defaultEtcBaseDir)
	ret.EventBrokerWorkingDir = ret.EventBrokerDir
	//ret.EventBrokerDir = Dir(filepath.FromSlash(fmt.Sprintf("%s/dist/eventbroker", ret.UserConfigDir)))

	//ret.ChannelsDir = Dir(filepath.FromSlash(fmt.Sprintf("%s", ret.EventBrokerDir)))
	//ret.MqttClientDir = Dir(filepath.FromSlash(fmt.Sprintf("%s", ret.EventBrokerDir)))

	return &ret
}


func (d *Dir) AddToPath(dir ...string) *Dir {

	var ret Dir
	var da []string

	da = append(da, string(*d))
	da = append(da, dir...)

	ret = Dir(filepath.FromSlash(strings.Join(da, "/")))

	return &ret
}


func (d *Dir) AddFileToPath(format string, fn ...interface{}) *File {

	var ret File
	var da []string

	da = append(da, string(*d))
	da = append(da, fmt.Sprintf(format, fn...))

	ret = File(filepath.FromSlash(strings.Join(da, "/")))

	return &ret
}


func (f *File) FileExists() error {

	var err error

	if f == nil {
		err = errors.New("File is nil")
		return err
	}

	_, err = os.Stat(f.String())
	if os.IsNotExist(err) {
		//fmt.Printf("Not exists PATH: '%s'\n", f.String())
	}

	return err
}


func (f *File) FileDelete() error {

	var err error

	_, err = os.Stat(f.String())
	if os.IsNotExist(err) {
		return err

	} else {
		err = os.Remove(f.String())
	}

	return err
}


func (d *Dir) DirExists() error {

	var err error

	if d == nil {
		err = errors.New("Dir is nil")
		return err
	}

	_, err = os.Stat(d.String())
	if os.IsNotExist(err) {
		//fmt.Printf("Not exists PATH: '%s'\n", d.String())
	}

	return err
}


func (d *Dir) CreateIfNotExists() (created bool, err error) {

	if d.DirExists() != nil {
		//fmt.Printf("CreateDirIfNotExists PATH: '%s'\n", d.String())
		err = os.MkdirAll(d.String(), os.ModePerm)
		if err == nil {
			created = true
		}
	}

	return created, err
}


func (d *Dirs) Append(dir ...string) *Dirs {

	var ret Dirs
	if d != nil {
		ret = *d
	}

	for _, s := range dir {
		ret = append(ret, Dir(s))
	}

	return &ret
}

//noinspection GoUnusedExportedFunction
func NewPath() *Paths {

	var ret Paths

	return &ret
}


func (p *Paths) AppendFile(file ...string) *Paths {

	var ret Paths
	if p != nil {
		ret = *p
	}

	for _, s := range file {
		ret = append(ret, *Split(s))
	}

	return &ret
}


func (p *Paths) AppendDir(dir ...string) *Paths {

	var ret Paths
	if p != nil {
		ret = *p
	}

	for _, s := range dir {
		if s == "" {
			continue
		}

		ret = append(ret, Path{Dir: Dir(s), File: ""})
	}

	return &ret
}


func (p *BasePaths) IsNil() error {
	var err error

	for range only.Once {
		if p == nil {
			err = errors.New("basepaths is nil")
			break
		}
	}

	return err
}


func (p *BasePaths) CreateIfNotExists() error {
	var err error

	for range only.Once {
		_, err = p.EventBrokerDir.CreateIfNotExists()
		if err != nil {
			break
		}

		_, err = p.EventBrokerEtcDir.CreateIfNotExists()
		if err != nil {
			break
		}

		_, err = p.EventBrokerLogDir.CreateIfNotExists()
		if err != nil {
			break
		}

		_, err = p.EventBrokerWorkingDir.CreateIfNotExists()
		if err != nil {
			break
		}
	}

	return err
}


func (p *Paths) CreateIfNotExists() (err error) {

	for _, p := range *p {
		if p.Dir.String() == "" {
			continue
		}

		_, err = p.Dir.CreateIfNotExists()
		if err != nil {
			break
		}
	}

	return err
}


func (p *Path) CreateIfNotExists() (created bool, err error) {

	created, err = p.Dir.CreateIfNotExists()
	if err != nil {
		fmt.Printf("CreateFileIfNotExists PATH: '%s'\n", p.String())
		err = os.MkdirAll(p.Dir.String(), os.ModePerm)
		created = true
	}

	return created, err
}


//func (me *Path) DirName() (created bool, err error) {
//
//	created, err = me.CreateIfNotExists()
//	if err != nil {
//		fmt.Printf("CreateFileIfNotExists PATH: '%s'\n", me.String())
//		err = os.MkdirAll(me.String(), os.ModePerm)
//		created = true
//	}
//
//	return created, err
//}


func (d *Dir) String() string {

	return string(*d)
}


func (f *File) String() string {

	return string(*f)
}


func (p *Path) String() string {

	return filepath.FromSlash(p.Dir.String() + "/"+ p.File.String())
}


func (p *Path) Abs() string {

	return filepath.FromSlash(p.Dir.String() + "/"+ p.File.String())
}


func Split(fn string) *Path {

	var pn Path

	d, f := filepath.Split(fn)
	pn.Dir = Dir(d)
	pn.File = File(f)

	return &pn
}

