package defaults

import "time"


const (
	BrandName = "Gearbox"
	Organization = "gearboxworks"
	DefaultTimeout = time.Second * 20
	DefaultProject = "/home/gearbox/projects/default"
	GearboxNetwork = "gearboxnet"
	DefaultUnitTestCmd = "/etc/gearbox/unit-tests/run.sh"
	DefaultCommandName = "default"

	DefaultMinTimeout = 2
	DefaultMaxTimeout = 30

	DefaultPathNone = "none"
	DefaultPathCwd = "cwd"
	DefaultPathHome = "home"
	DefaultPathEmpty = ""
	DefaultProvider = "docker"

	EnvPrefix = "LAUNCH"

	LanguageAppName = "Launch"
	LanguageContainerName = "Container"
	LanguageImageName = "Container image"
	LanguageContainerPluralName = "Containers"
	LanguageImagePluralName = "Container images"

	//LanguageAppName = "Gearbox"
	//LanguageContainerName = "Gear"
	//LanguageImageName = "Gear image"
)
