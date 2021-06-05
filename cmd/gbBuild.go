package cmd

import (
	"fmt"
	"github.com/gearboxworks/scribeHelpers/toolGear"
	"github.com/gearboxworks/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
	"strings"
)

// https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
// https://www.nearform.com/blog/building-docker-images-in-go/
// https://stackoverflow.com/questions/46878793/golang-docker-api-reports-invalid-argument-while-hitting-imagebuild
// https://www.loginradius.com/blog/async/build-push-docker-images-golang/
//

// ******************************************************************************** //
var gbBuildCmd = &cobra.Command{
	Use: "build",
	//Aliases:				[]string{"show"},
	Short:                 ux.SprintfBlue("Create %s %s", defaults.LanguageAppName, defaults.LanguageContainerPluralName),
	Long:                  ux.SprintfBlue("Create %s %s.", defaults.LanguageAppName, defaults.LanguageContainerPluralName),
	Example:               ux.SprintfWhite("launch create"),
	DisableFlagParsing:    false,
	DisableFlagsInUseLine: false,
	Run:                   gbBuildFunc,
	Args:                  cobra.RangeArgs(0, 2),
}

//gbBuildFunc takes a pointer to cobra.command and
//string arguments and output command state object
func gbBuildFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, false)
		if Cmd.State.IsNotOk() {
			break
		}

		//for _, c := range cmd.Commands() {
		//	for k, v := range c.Annotations {
		//		fmt.Printf(`%s - "%s": "%s"`,
		//			c.Use,
		//			k,
		//			v,
		//			)
		//		fmt.Printf("\n")
		//	}
		//}

		switch {
		case len(args) == 0:
			_ = cmd.Help()
		}
	}
}

// ******************************************************************************** //
var gbBuildCreateCmd = &cobra.Command{
	Use:                fmt.Sprintf("create <%s name>", defaults.LanguageContainerName),
	SuggestFor:         []string{"compile", "generate"},
	Short:              ux.SprintfBlue("Build a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:               ux.SprintfBlue("Allows building of arbitrary containers as a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:            ux.SprintfWhite("launch build create golang"),
	DisableFlagParsing: false,
	Run:                gbBuildCreateFunc,
	Args:               cobra.ExactArgs(1),
}

//goland:noinspection GoUnusedParameter
//gbBuildCreateFunc takes a pointer to cobra.command and
//string arguments and create command state object
func gbBuildCreateFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbBuildCreateFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}
func (ga *LaunchArgs) gbBuildCreateFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var gearJsonFile string
		gearJsonFile = "gearbox.json"
		ga.State = ga.Gears.CreateImage(gearJsonFile, defaults.Organization, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}

// ******************************************************************************** //
var gbBuildCleanCmd = &cobra.Command{
	Use:                fmt.Sprintf("clean <%s name> [%s version]", defaults.LanguageContainerName, defaults.LanguageContainerName),
	SuggestFor:         []string{},
	Short:              ux.SprintfBlue("Remove a %s %s build", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:               ux.SprintfBlue("Remove a %s %s build.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:            ux.SprintfWhite("launch build clean golang"),
	DisableFlagParsing: false,
	Run:                gbBuildCleanFunc,
	Args:               cobra.RangeArgs(1, 2),
}

//goland:noinspection GoUnusedParameter
//gbBuildCleanFunc takes a pointer to cobra.command and
//string arguments and clean command state object
func gbBuildCleanFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbBuildCleanFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}
func (ga *LaunchArgs) gbBuildCleanFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var containers map[string]*toolGear.Gear
		containers, ga.State = ga.Gears.FindContainers(ga.Name)

		foo := ga.Gears.SelectedVersions()
		foo.GetVersions()

		for _, c := range containers {
			ux.PrintflnGreen("Build clean %s '%s:%s'", defaults.LanguageContainerName, c.Container.Name, c.Container.Version)
			//fmt.Printf("C: %s\n", c.GearConfig.Meta.String())
			ga.Gears.Selected = c

			ga.State = c.Stop()
			if ga.State.IsError() {
				continue
			}
			if !ga.State.IsExited() {
				ga.State.SetWarning("%s '%s:%s' didn't stop", defaults.LanguageContainerName, c.Container.Name, c.Container.Version)
				continue
			}

			ga.State = c.Remove()
			if ga.State.IsError() {
				ga.State.SetWarning("%s '%s:%s' didn't remove", defaults.LanguageContainerName, c.Container.Name, c.Container.Version)
				continue
			}

			ux.PrintflnGreen("Build clean %s '%s:%s'", defaults.LanguageImageName, c.Container.Name, c.Container.Version)
			ga.State = c.ImageRemove()
			if ga.State.IsError() {
				ga.State.SetWarning("%s '%s:%s' didn't remove", defaults.LanguageImageName, c.Image.Name, c.Image.Version)
				continue
			}

			//if ga.State.IsExited() {
			//	ga.State.SetOk("%s '%s:%s' stopped OK", defaults.LanguageContainerName, ga.Name, ga.Version)
			//	continue
			//}
			//if ga.State.IsCreated() {
			//	ga.State.SetOk("%s '%s:%s' stopped OK", defaults.LanguageContainerName, ga.Name, ga.Version)
			//	continue
			//}

		}

		//ga.State = ga.gbUninstallFunc()
		//if ga.State.IsError() {
		//	break
		//}
		//
		//var found bool
		//ga.State = ga.Gears.FindImage(ga.Name, ga.Version)
		//if ga.State.IsError() {
		//	break
		//}
		//found = ga.Gears.State.GetResponseAsBool()
		//if !found {
		//	ga.State.SetOk("%s '%s:%s' already removed.", defaults.LanguageImageName, ga.Name, ga.Version)
		//	ga.State.SetOutput("")
		//	break
		//}
		//ga.State.Clear()
		//
		//if !ga.Quiet {
		//	ux.PrintflnNormal("Removing %s '%s:%s': ", defaults.LanguageContainerName, ga.Name, ga.Version)
		//}
		//ga.State = ga.Gears.SelectedImageRemove()
		//if ga.State.IsError() {
		//	ga.State.SetError("%s '%s:%s' remove error - %s", defaults.LanguageImageName, ga.Name, ga.Version, ga.State.GetError())
		//	break
		//}
		//
		//if ga.State.IsOk() {
		//	if ga.Temporary {
		//		ga.State.Clear()
		//		break
		//	}
		//
		//	ga.State.SetOk("%s '%s:%s' removed OK", defaults.LanguageImageName, ga.Name, ga.Version)
		//	ga.State.SetOutput("")
		//	break
		//}
		//
		//ga.State.SetWarning("%s '%s:%s' cannot be removed", defaults.LanguageImageName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}

// ******************************************************************************** //
var gbBuildStartCmd = &cobra.Command{
	Use:                fmt.Sprintf("start <%s name>", defaults.LanguageContainerName),
	Short:              ux.SprintfBlue("Start a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:               ux.SprintfBlue("Start a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:            ux.SprintfWhite("launch build start golang"),
	DisableFlagParsing: false,
	Run:                gbBuildStartFunc,
	Args:               cobra.ExactArgs(1),
}

//goland:noinspection GoUnusedParameter
//gbBuildCleanFunc takes a pointer to cobra.command and
//string arguments and start command state object
func gbBuildStartFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}
		//Cmd.SetDebug(ga.Debug)

		Cmd.State = ga.gbStartFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}
func (ga *LaunchArgs) gbBuildStartFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var found bool
		found, ga.State = ga.Gears.FindContainer(ga.Name, ga.Version)
		if ga.State.IsNotOk() {
			break
		}
		if !found {
			ga.State.SetError("%s not found '%s:%s'.", defaults.LanguageContainerName, ga.Name, ga.Version)
			break
		}

		if ga.State.IsRunning() {
			ga.State.SetOk("%s '%s:%s' already started.", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		if !ga.Quiet {
			ux.PrintflnNormal("Starting %s '%s:%s': ", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.Gears.SelectedStart()
		if ga.State.IsError() {
			if strings.Contains(ga.State.GetError().Error(), "address already in use") {
				//ux.PrintflnRed("Error: There are ports already used.")
				//saved := gear.State.GetError()
				ga.Gears.Selected.ListImagePorts()
				ga.State.SetError("Error: There are ports already used.")
				break
			}
			ga.State.SetError("%s '%s:%s' start error - %s", defaults.LanguageContainerName, ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsRunning() {
			ga.State.SetOk("%s '%s:%s' started OK", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetError("%s '%s:%s' cannot be started", defaults.LanguageContainerName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}

// ******************************************************************************** //
var gbBuildStopCmd = &cobra.Command{
	Use:                fmt.Sprintf("stop <%s name>", defaults.LanguageContainerName),
	Short:              ux.SprintfBlue("Stop a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:               ux.SprintfBlue("Stop a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:            ux.SprintfWhite("launch build stop golang"),
	DisableFlagParsing: false,
	Run:                gbStopFunc,
	Args:               cobra.ExactArgs(1),
}

// ******************************************************************************** //
var gbPublishCmd = &cobra.Command{
	Use:                fmt.Sprintf("publish <%s name>", defaults.LanguageContainerName),
	SuggestFor:         []string{"upload"},
	Short:              ux.SprintfBlue("Publish a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:               ux.SprintfBlue("Publish a %s %s to GitHub or DockerHub.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:            ux.SprintfWhite("launch build publish golang"),
	DisableFlagParsing: false,
	Run:                gbPublishFunc,
	Args:               cobra.ExactArgs(1),
}

//goland:noinspection GoUnusedParameter
//gbPublishFunc takes a pointer to cobra.command and
//string arguments and publish command state object
func gbPublishFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbPublishFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}
func (ga *LaunchArgs) gbPublishFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		var found bool
		found, ga.State = ga.Gears.FindContainer(ga.Name, ga.Version)
		if ga.State.IsError() {
			break
		}
		if !found {
			break
		}
		if ga.State.IsExited() {
			ga.State.SetOk("%s '%s:%s' already stopped.", defaults.LanguageImageName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		if !ga.Quiet {
			ux.PrintflnNormal("Stopping %s '%s:%s': ", defaults.LanguageContainerName, ga.Name, ga.Version)
		}
		ga.State = ga.Gears.SelectedStop()
		if ga.State.IsError() {
			ga.State.SetError("%s '%s:%s' stop error - %s", defaults.LanguageContainerName, ga.Name, ga.Version, ga.State.GetError())
			break
		}

		if ga.State.IsExited() {
			ga.State.SetOk("%s '%s:%s' stopped OK", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		if ga.State.IsCreated() {
			ga.State.SetOk("%s '%s:%s' stopped OK", defaults.LanguageContainerName, ga.Name, ga.Version)
			ga.State.SetOutput("")
			break
		}

		ga.State.SetWarning("%s '%s:%s' cannot be stopped", defaults.LanguageContainerName, ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}

// ******************************************************************************** //
var gbSaveCmd = &cobra.Command{
	Use:                fmt.Sprintf("export <%s name>", defaults.LanguageContainerName),
	SuggestFor:         []string{"save"},
	Short:              ux.SprintfBlue("Save state of a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:               ux.SprintfBlue("Save state of a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:            ux.SprintfWhite("launch save golang"),
	DisableFlagParsing: false,
	Run:                gbSaveFunc,
	Args:               cobra.ExactArgs(1),
}

//goland:noinspection GoUnusedParameter
//gbSaveFunc takes a pointer to cobra.command and
//string arguments and save command state object
func gbSaveFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}

		ux.PrintfWarning("Command not yet implemented.\n")
	}
}

// ******************************************************************************** //
var gbLoadCmd = &cobra.Command{
	Use:                fmt.Sprintf("import <%s name>", defaults.LanguageContainerName),
	SuggestFor:         []string{"load"},
	Short:              ux.SprintfBlue("Load a %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:               ux.SprintfBlue("Load a %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:            ux.SprintfWhite("launch load golang"),
	DisableFlagParsing: false,
	Run:                gbLoadFunc,
	Args:               cobra.ExactArgs(1),
}

//goland:noinspection GoUnusedParameter
//gbLoadFunc takes a pointer to cobra.command and
//string arguments and load command state object
func gbLoadFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}

		ux.PrintfWarning("Command not yet implemented.\n")
	}
}

// ******************************************************************************** //
var gbUnitTestCmd = &cobra.Command{
	Use:                   fmt.Sprintf("test <%s name>", defaults.LanguageContainerName),
	Short:                 ux.SprintfBlue("Execute unit tests in %s %s", defaults.LanguageAppName, defaults.LanguageContainerName),
	Long:                  ux.SprintfBlue("Execute unit tests in %s %s.", defaults.LanguageAppName, defaults.LanguageContainerName),
	Example:               ux.SprintfWhite("launch build test terminus"),
	DisableFlagParsing:    true,
	DisableFlagsInUseLine: true,
	Run:                   gbUnitTestFunc,
	Args:                  cobra.MinimumNArgs(1),
}

//goland:noinspection GoUnusedParameter
//gbUnitTestFunc takes a pointer to cobra.command and
//string arguments and unit test command state object
func gbUnitTestFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args, true)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbUnitTestFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}
func (ga *LaunchArgs) gbUnitTestFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.gbStartFunc()
		if !ga.State.IsRunning() {
			ga.State.SetError("%s not started", defaults.LanguageContainerName)
			break
		}

		ga.Args = []string{defaults.DefaultUnitTestCmd}
		ga.State = ga.Gears.SelectedSsh(true, ga.SshStatus, ga.Mount, ga.Args)

		if ga.Temporary {
			ga.State = ga.gbUninstallFunc()
		}
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
}

/*
gb_build() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Building image for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		gb_getenv ${GB_VERSION}
		gb_getdockerfile ${GB_VERSION}


		# LOGFILE="${GB_VERDIR}/logs/$(date +'%Y%m%d-%H%M%S').log"
		LOGFILE="${GB_VERDIR}/logs/build.log"
		if [ ! -d "${GB_VERDIR}/logs/" ]
		then
			mkdir -p "${GB_VERDIR}/logs"
		fi

		if [ "${GB_REF}" == "base" ]
		then
			DOCKER_ARGS="--squash"
			p_info "${GB_IMAGENAME}:${GB_VERSION}" "This is a base container."

		elif [ "${GB_REF}" != "" ]
		then
			DOCKER_ARGS=""
			p_info "${GB_IMAGENAME}:${GB_VERSION}" "Pull ref container."
			docker pull "${GB_REF}"
			if [ "${GB_RUN}" == "" ]
			then
				p_info "${GB_IMAGENAME}:${GB_VERSION}" "Query ref container."
				GEARBOX_ENTRYPOINT="$(docker inspect --format '{{ with }}{{ else }}{{ with .ContainerConfig.Entrypoint}}{{ index . 0 }}{{ end }}' "${GB_REF}")"
				export GEARBOX_ENTRYPOINT
				GEARBOX_ENTRYPOINT_ARGS="$(docker inspect --format '{{ join .ContainerConfig.Entrypoint " " }}' "${GB_REF}")"
				export GEARBOX_ENTRYPOINT_ARGS
			else
				GEARBOX_ENTRYPOINT="${GB_RUN}"
				export GEARBOX_ENTRYPOINT
				GEARBOX_ENTRYPOINT_ARGS="${GB_ARGS}"
				export GEARBOX_ENTRYPOINT_ARGS
			fi
		fi

		p_info "${GB_IMAGENAME}:${GB_VERSION}" "Building container."
		if [ "${GITHUB_ACTIONS}" == "" ]
		then
			script ${LOG_ARGS} ${LOGFILE} \
				docker build -t ${GB_IMAGENAME}:${GB_VERSION} -f ${GB_DOCKERFILE} --build-arg GEARBOX_ENTRYPOINT --build-arg GEARBOX_ENTRYPOINT_ARGS ${DOCKER_ARGS} .
			p_info "${GB_IMAGENAME}:${GB_VERSION}" "Log file saved to \"${LOGFILE}\""
		else
			docker build -t ${GB_IMAGENAME}:${GB_VERSION} -f ${GB_DOCKERFILE} --build-arg GEARBOX_ENTRYPOINT --build-arg GEARBOX_ENTRYPOINT_ARGS ${DOCKER_ARGS} .
		fi

		if [ "${GB_MAJORVERSION}" != "" ]
		then
			docker tag ${GB_IMAGENAME}:${GB_VERSION} ${GB_IMAGENAME}:${GB_MAJORVERSION}
		fi

		if [ "${GB_LATEST}" == "true" ]
		then
			docker tag ${GB_IMAGENAME}:${GB_VERSION} ${GB_IMAGENAME}:latest
		fi
	done

	return ${EXIT}
}
*/
