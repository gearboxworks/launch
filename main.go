package main

import (
	"flag"
	"launch/cmd"
	"launch/defaults"
	"launch/gear"
	"launch/only"
	"launch/ospaths"
	"launch/ux"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

var Version = "1.4.2"

func init() {
	_ = ux.Open()
}

func main() {
	//var response dockerClient.Response
	var state ux.State

	for range only.Once {
		var err error
		var args *Args
		args, err = ProcessArgs()
		if err != nil {
			ux.PrintError(err)
			args.Help()
			break
		}

		var g *gear.Gear
		g, state = gear.NewGear(Debug)
		if state.IsError() {
			break
		}

		// @TODO - testing.
		//ux.PrintfOk("OK - it works\n")
		//ux.PrintfWarning("WARNING - it may or may not have worked... sort of\n")
		//ux.PrintfError("ERROR - oops, no beer in the fridge\n")
		//err = errors.New("what? no beer!")
		//ux.PrintError(err)
		//_ = ux.Draw2()
		//ux.Draw3()
		//_ = ux.Draw4()
		//_ = ux.Draw5()


		if *args.List {
			state = g.Docker.ImageList(*args.ContainerName)
			if state.IsError() {
				break
			}

			state = g.Docker.ContainerList(*args.ContainerName)
			if state.IsError() {
				break
			}

			state = g.Docker.NetworkList(defaults.GearboxNetwork)
			break
		}


		if *args.ListContainers {
			state = g.Docker.ContainerList(*args.ContainerName)
			break
		}


		if *args.ListImages {
			state = g.Docker.ImageList(*args.ContainerName)
			break
		}


		if *args.ContainerName == "" {
			state.SetError("no container specified")
			args.Help()
			break
		}

		var found bool
		found, state = g.Docker.FindContainer(*args.ContainerName, "")
		if state.IsError() {
			break
		}
		state.ClearAll()

		state = g.Docker.NetworkCreate(defaults.GearboxNetwork)
		if state.IsError() {
			break
		}
		state.ClearAll()

		// Stop a container.
		if *args.ContainerStop {
			if found {
				ux.Printf("Stopping gear '%s': ", *args.ContainerName)
				state = g.Docker.Container.Stop()
				if state.IsError() {
					ux.PrintfRed("error stopping - %s\n", state.Error)
				} else if state.IsExited() {
					ux.PrintfGreen("OK\n")
				} else {
					ux.PrintfYellow("cannot be stopped\n")
				}
			} else {
				ux.PrintfWarning("Gear '%s' doesn't exist.\n", *args.ContainerName)
			}
			break
		}

		// Remove a container.
		if *args.ContainerRemove {
			if found {
				ux.Printf("Stopping gear '%s': ", *args.ContainerName)
				state = g.Docker.Container.Stop()
				if state.IsError() {
					ux.PrintfRed("error stopping - %s\n", state.Error)
					break
				} else if state.IsExited() {
					ux.PrintfGreen("OK\n")
				} else {
					ux.PrintfYellow("cannot be stopped\n")
					break
				}

				ux.Printf("Removing gear '%s': ", *args.ContainerName)
				state = g.Docker.Container.Remove()
				if state.IsError() {
					ux.PrintfRed("error removing - %s\n", state.Error)
				} else {
					ux.PrintfGreen("OK\n")
				}
			} else {
				ux.PrintfWarning("Gear '%s' doesn't exist.\n", *args.ContainerName)
			}
			break
		}

		// Remove an image.
		if *args.ImageRemove {
			if found {
				ux.Printf("Stopping gear '%s': ", *args.ContainerName)
				state = g.Docker.Container.Stop()
				if state.IsError() {
					ux.PrintfRed("error stopping - %s\n", state.Error)
					break
				} else if state.IsExited() {
					ux.PrintfGreen("OK\n")
				} else {
					ux.PrintfYellow("cannot be stopped\n")
					break
				}

				ux.Printf("Removing gear '%s': ", *args.ContainerName)
				state = g.Docker.Container.Remove()
				if state.IsError() {
					ux.PrintfRed("error removing - %s\n", state.Error)
				} else {
					ux.PrintfGreen("OK\n")
				}
			}

			var ok bool
			ok, state = g.Docker.FindImage(*args.ContainerName, *args.ContainerVersion)
			if state.IsError() {
				state.SetWarning("Gear image '%s' doesn't exist.\n", *args.ContainerName)
				break
			}
			if !ok {
				ux.PrintfYellow("Gear image '%s' doesn't exist.\n", *args.ContainerName)
				break
			}

			ux.Printf("Removing gear image '%s': ", *args.ContainerName)
			state = g.Docker.Image.Remove()
			if state.IsError() {
				ux.PrintfRed("error removing - '%s'\n", state.Error)
			} else {
				ux.PrintfGreen("OK\n")
			}

			break
		}

		// Default - run a shell.
		if !found {
			// state = g.ContainerCreate("golang", "", "/Users/mick/Documents/GitHub/containers/docker-golang")
			state = g.Docker.Container.ContainerCreate(*args.ContainerName, "", *args.DockerMount)
			if state.IsError() {
				break
			}
		}

		state = g.Docker.Container.Start()
		if state.IsError() {
			break
		}
		if !state.IsRunning() {
			state.SetError("container not started")
			break
		}

		state = g.Docker.ContainerSsh(*args.Shell, !*args.StatusLine, flag.Args()...)
		break
	}


	exit := 0
	if state.IsError() {
		exit = 1
	}
	if state.IsWarning() {
		exit = 2
	}
	state.Print()
	ux.Close()
	os.Exit(exit)
}

func (me *Args) Help() {
	for range only.Once {
		exe := path.Base(os.Args[0])
		exe = strings.TrimSuffix(exe, "-Darwin")
		exe = strings.TrimSuffix(exe, "-Linux")
		exe = strings.TrimSuffix(exe, "-Windows")

		ux.PrintfCyan("\nGearbox: %s v%s:\n", exe, Version)
		ux.PrintfWhite("\tLaunch an interactive container within the Gearbox environment.\n")
		ux.PrintfWhite("\n")

		// fmt.Printf("\n")
		// fmt.Printf("-%s \n\t%s - %s\n", me.DockerHost.Name, me.DockerHost.Usage, me.DockerHost.DefValue)

		flag.PrintDefaults()
		ux.PrintfWhite("\n")
		ux.PrintfWhite("Examples:\n")
		ux.PrintfWhite("Run 'ls -l' within a terminus container.\n")
		ux.PrintfGreen("\t%s -gb-name terminus -gb-shell -- ls -l\n", exe)

		ux.PrintfWhite("Run an interactive shell within a terminus container.\n")
		ux.PrintfGreen("\t%s -gb-name terminus -gb-shell\n", exe)

		ux.PrintfWhite("Run 'terminus' command within a terminus container.\n")
		ux.PrintfGreen("\t%s -gb-name terminus\n", exe)

		ux.PrintfWhite("Run 'terminus auth:login' within a terminus container.\n")
		ux.PrintfGreen("\t%s -gb-name terminus auth:login\n", exe)

		ux.PrintfWhite("\n")
		ux.PrintfWhite("If %s is symlinked to 'terminus', then you can drop the '-gb-name terminus' ...\n", exe)

		ux.PrintfWhite("Run 'ls -l' within a terminus container.\n")
		ux.PrintfGreen("\tterminus -gb-shell -- ls -l\n")

		ux.PrintfWhite("Run an interactive shell within a terminus container.\n")
		ux.PrintfGreen("\tterminus -gb-shell\n")

		ux.PrintfWhite("Run 'terminus' command within a terminus container.\n")
		ux.PrintfGreen("\tterminus\n")

		ux.PrintfWhite("Run 'terminus auth:login' within a terminus container.\n")
		ux.PrintfGreen("\tterminus auth:login\n")

		ux.PrintfWhite("\t\n")
	}
}

// func HelpVariables() {
// 	for range only.Once {
// 		fmt.Printf("Keys accessible within your template file:\n")
// 		fmt.Printf("\t{{ .Json }} - Your JSON file will appear here.\n")
// 		fmt.Printf("\n")
// 		fmt.Printf("\t{{ .Env }} - A map containing the runtime environment.\n")
// 		fmt.Printf("\n")
// 		fmt.Printf("\t{{ .ExecName }} - Executable used to produce the resulting file.\n")
// 		fmt.Printf("\t{{ .ExecVersion }} - Version of executable used to produce the resulting file.\n")
// 		fmt.Printf("\n")
// 		fmt.Printf("\t{{ .CreationDate }} - Creation date of resulting file.\n")
// 		fmt.Printf("\t{{ .CreationEpoch }} - Creation date, (unix epoch), of resulting file.\n")
// 		fmt.Printf("\t{{ .CreationInfo }} - More creation information.\n")
// 		fmt.Printf("\t{{ .CreationWarning }} - Generic 'DO NOT EDIT' warning.\n")
// 		fmt.Printf("\n")
// 		fmt.Printf("\t{{ .TemplateFile.Dir }} - template file absolute directory.\n")
// 		fmt.Printf("\t{{ .TemplateFile.Name }} - template filename.\n")
// 		fmt.Printf("\t{{ .TemplateFile.CreationDate }} - template file creation date.\n")
// 		fmt.Printf("\t{{ .TemplateFile.CreationEpoch }} - template file creation date, (unix epoch).\n")
// 		fmt.Printf("\n")
// 		fmt.Printf("\t{{ .JsonFile.Dir }} - json file absolute directory.\n")
// 		fmt.Printf("\t{{ .JsonFile.Name }} - json filename.\n")
// 		fmt.Printf("\t{{ .JsonFile.CreationDate }} - json file creation date.\n")
// 		fmt.Printf("\t{{ .JsonFile.CreationEpoch }} - json file creation date, (unix epoch).\n")
// 		fmt.Printf("\n")
// 	}
// }
//
// func HelpFunctions() {
// 	for range only.Once {
// 		fmt.Printf("Functions accessible within your template file:\n")
// 		fmt.Printf("\t{{ isInt $value }} - is $value an integer?\n")
// 		fmt.Printf("\t{{ isString $value }} - is $value a string?\n")
// 		fmt.Printf("\t{{ isSlice $value }} - is $value a slice?\n")
// 		fmt.Printf("\t{{ isArray $value }} - is $value an array?\n")
// 		fmt.Printf("\t{{ isMap $value }} - is $value a map?\n")
// 		fmt.Printf("\n")
// 		fmt.Printf("\t{{ ToUpper $value }} - uppercase $value.\n")
// 		fmt.Printf("\t{{ ToLower $value }} - lowercase $value.\n")
// 		fmt.Printf("\t{{ ToString $value }} - convert $value to a string.\n")
// 		fmt.Printf("\t{{ FindInMap $map $value }} - find $value in $map and return reference.\n")
// 		fmt.Printf("\t{{ ReadFile $file }} - read in $file and print verbatim. \n")
// 		fmt.Printf("\n")
// 		fmt.Printf("See http://masterminds.github.io/sprig/ for additional functions...\n")
// 		fmt.Printf("\n")
// 	}
// }
//
// func HelpExamples() {
// 	for range only.Once {
// 		fmt.Printf("Examples:\n")
// 		fmt.Printf("# Print out .dir key from config.json\n")
// 		fmt.Printf("\tJsonConfig -json config.json -template-string '{{ .Json.dir }}'\n")
//
// 		fmt.Printf("# Process Dockerfile.tmpl file and output to STDOUT.\n")
// 		fmt.Printf("\tJsonConfig -json config.json -template DockerFile.tmpl\n")
//
// 		fmt.Printf("# Process Dockerfile.tmpl file and output to Dockerfile.\n")
// 		fmt.Printf("\tJsonConfig -json config.json -template DockerFile.tmpl -out Dockerfile\n")
//
// 		fmt.Printf("# Process nginx.conf.tmpl file, output to nginx.conf and remove nginx.conf.tmpl afterwards.\n")
// 		fmt.Printf("\tJsonConfig -json config.json -create nginx.conf\n")
//
// 		fmt.Printf("# Process setup.sh.tmpl file, output to setup.sh and execute as a shell script.\n")
// 		fmt.Printf("\tJsonConfig -json config.json -template setup.sh -shell\n")
//
// 		fmt.Printf("# Process setup.sh.tmpl file, output to setup.sh, execute as a shell script and remove afterwards.\n")
// 		fmt.Printf("\tJsonConfig -json config.json -create setup.sh.tmpl -shell\n")
// 	}
// }

var Debug bool

type Args struct {
	Debug            *bool

	DockerHost       *string
	DockerPort       *string
	DockerDaemon     *url.URL

	//AltCommand       *bool
	DockerMount      *string
	ContainerName    *string
	ContainerVersion *string
	Shell            *bool
	StatusLine       *bool

	List             *bool
	ListImages       *bool
	ListContainers   *bool

	ContainerStop    *bool
	ContainerRemove  *bool
	ImageRemove      *bool
	ImageBuild       *bool
}

type Hargs struct {
	DockerHost       *flag.Flag
	DockerPort       *flag.Flag
	DockerDaemon     *url.URL

	Command          *flag.Flag
	DockerMount      *flag.Flag
	ContainerName    *flag.Flag
	ContainerVersion *flag.Flag
	Shell            *flag.Flag
	StatusLine       *flag.Flag

	List             *flag.Flag
	ListImages       *flag.Flag
	ListContainers   *flag.Flag

	ContainerStop    *flag.Flag
	ContainerRemove  *flag.Flag
	ImageRemove      *flag.Flag
	ImageBuild       *flag.Flag
}

type boolFlag struct {
	set   bool
	value bool
}

func (sf *boolFlag) Set(x bool) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *boolFlag) String() bool {
	return sf.value
}

func ProcessArgs() (*Args, error) {
	var err error
	var args Args

	for range only.Once {
		//var hargs Hargs

		//foo := ospaths.Split("C:\\\\Users\\\\mick\\\\Documents\\\\launch")
		//foo := ospaths.Split("C:\\\\Users\\\\mick\\\\Documents\\\\launch-Darwin")
		//foo := ospaths.Split("./bin/launch-Darwin")
		foo := ospaths.Split(os.Args[0])
		exe := foo.File.String()
		//fmt.Printf("F2: %s %s\n", foo.File.String(), foo.Dir.String())
		var ok bool
		ok, err = regexp.MatchString("^" + defaults.BinaryName, exe)
		if !ok {
			break
		}

		if ok {
			exe = ""
		}

		err = cmd.Execute()
		if err != nil {
			break
		}

		break

		//help_all := flag.Bool("gb-help", false, "Show all help.")
		//
		//args.Debug = flag.Bool("gb-debug", false, "DEBUG")
		//
		//args.DockerHost = flag.String("gb-host", "", "Specify an alternative Docker host.")
		//hargs.DockerHost = flag.Lookup("gb-host")
		//
		//args.DockerPort = flag.String("gb-docker-port", "", "Specify an alternative Docker port.")
		//hargs.DockerPort = flag.Lookup("gb-docker-port")
		//
		//args.ContainerName = flag.String("gb-name", exe, "Specify container name.")
		//hargs.ContainerName = flag.Lookup("gb-name")
		//
		//args.ContainerVersion = flag.String("gb-version", "latest", "Specify container version.")
		//hargs.ContainerVersion = flag.Lookup("gb-version")
		//
		//args.Shell = flag.Bool("gb-shell", false, "Run a shell instead of the default container command.")
		//hargs.Shell = flag.Lookup("gb-shell")
		//
		//args.StatusLine = flag.Bool("gb-status", false, "Include a Gearbox status line within the container shell.")
		//hargs.StatusLine = flag.Lookup("gb-status")
		//
		//args.List = flag.Bool("gb-list", false, "List all images and containers.")
		//hargs.List = flag.Lookup("gb-list")
		//
		//args.ListImages = flag.Bool("gb-images", false, "List all images downloaded.")
		//hargs.ListImages = flag.Lookup("gb-images")
		//
		//args.ListContainers = flag.Bool("gb-containers", false, "List all containers created.")
		//hargs.ListContainers = flag.Lookup("gb-containers")
		//
		//args.ContainerStop = flag.Bool("gb-stop", false, "Stop a created container.")
		//hargs.ContainerStop = flag.Lookup("gb-stop")
		//
		//args.ContainerRemove = flag.Bool("gb-remove", false, "Remove a created container.")
		//hargs.ContainerRemove = flag.Lookup("gb-remove")
		//
		//args.ImageRemove = flag.Bool("gb-clean", false, "Remove downloaded image.")
		//hargs.ImageRemove = flag.Lookup("gb-clean")
		//
		//args.ImageBuild = flag.Bool("gb-build", false, "Build an image.")
		//hargs.ImageBuild = flag.Lookup("gb-build")
		//
		//args.DockerMount = flag.String("gb-project", "", "Specify a project mount point.")
		//hargs.DockerMount = flag.Lookup("gb-project")
		//
		//flag.Parse()
		//
		//Debug = *args.Debug

		//if (*args.DockerHost != "") && (*args.DockerPort != "") {
		//	args.DockerDaemon, err = client.ParseHostURL(fmt.Sprintf("tcp://%s:%s", *args.DockerHost, *args.DockerPort))
		//	if err != nil {
		//		break
		//	}
		//
		//	err = os.Setenv("DOCKER_HOST", args.DockerDaemon.String())
		//}
		//
		//ok, _ = regexp.MatchString("^" + defaults.BinaryName, *args.ContainerName)
		//if ok {
		//	*args.ContainerName = ""
		//}
		//
		//// Show help.
		////if *help_all {
		////	args.Help()
		////	os.Exit(0)
		////}
		//
		//if *args.ListImages {
		//	break
		//}
		//
		//if *args.List {
		//	break
		//}
		//
		//if *args.ListContainers {
		//	break
		//}
		//
		//if *args.ContainerStop {
		//	break
		//}
		//
		//if *args.ContainerRemove {
		//	break
		//}
		//
		//if *args.ImageRemove {
		//	break
		//}
		//
		//if *args.Shell {
		//	break
		//}
		//
		// @TODO Need to figure this logic out.
		//args.Help()
		//os.Exit(0)
	}

	return &args, err
}

//type FileInfo struct {
//	Dir string
//	Name string
//	CreationEpoch int64
//	CreationDate string
//}
//
//func (me *Environment) ToString() string {
//	var s string
//
//	for range only.Once {
//		s = fmt.Sprintf("%s", *me)
//	}
//
//	return s
//}
//
//func (me *FileInfo) getPaths(f string) error {
//	var err error
//
//	for range only.Once {
//		var abs string
//		abs, err = filepath.Abs(f)
//		if err != nil {
//			break
//		}
//
//		me.Dir = filepath.Dir(abs)
//		me.Name = filepath.Base(abs)
//
//		var fstat os.FileInfo
//		fstat, err = os.Stat(abs)
//		if os.IsNotExist(err) {
//			break
//		}
//
//		me.CreationEpoch = fstat.ModTime().Unix()
//		me.CreationDate = fstat.ModTime().Format("2006-01-02T15:04:05-0700")
//	}
//
//	return err
//}
//
//func fileToString(fileName string) ([]byte, error) {
//	var jsonString []byte
//	var err error
//
//	for range only.Once {
//		_, err = os.Stat(fileName)
//		if os.IsNotExist(err) {
//			break
//		}
//
//		jsonString, err = ioutil.ReadFile(fileName)
//		if err != nil {
//			break
//		}
//	}
//
//	return jsonString, err
//}
