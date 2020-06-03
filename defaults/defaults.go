package defaults

import "time"

const (
	BrandName = "Gearbox"
	Organization = "gearboxworks"
	Timeout = time.Second * 2
	DefaultProject = "/home/gearbox/projects/default"
	GearboxNetwork = "gearboxnet"
	DefaultUnitTestCmd = "/etc/gearbox/unit-tests/run.sh"
	DefaultCommandName = "default"

	DefaultPathNone = "none"
	DefaultPathCwd = "cwd"
	DefaultPathHome = "home"
	DefaultPathEmpty = ""
	DefaultProvider = "docker"

	LanguageContainerName = "Gear"
	LanguageImageName = "Gear image"
)
