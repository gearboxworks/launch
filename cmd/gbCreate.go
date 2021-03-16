package cmd

import (
	"github.com/newclarity/scribeHelpers/ux"
	"github.com/spf13/cobra"
	"launch/defaults"
)


func gbCreateFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		switch {
			case len(args) == 0:
				_ = cmd.Help()
		}
	}
}

func gbBuildFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbBuildFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}

func gbBuildCleanFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbCleanFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}

func gbPublishFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbPublishFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}

func gbSaveFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		ux.PrintfWarning("Command not yet implemented.\n")
	}
}

func gbLoadFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		ux.PrintfWarning("Command not yet implemented.\n")
	}
}

func gbUnitTestFunc(cmd *cobra.Command, args []string) {
	for range onlyOnce {
		var ga LaunchArgs

		Cmd.State = ga.ProcessArgs(rootCmd, args)
		if Cmd.State.IsNotOk() {
			break
		}

		Cmd.State = ga.gbUnitTestFunc()
		if Cmd.State.IsNotOk() {
			break
		}
	}
}


func (ga *LaunchArgs) gbBuildFunc() *ux.State {
	if state := ux.IfNilReturnError(ga); state.IsError() {
		return state
	}

	for range onlyOnce {
		ga.State = ga.Gears.CreateImage(ga.Name, ga.Version)
	}

	if !ga.Quiet {
		ga.State.PrintResponse()
	}
	return ga.State
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
		ga.State = ga.Gears.Selected.Container.Stop()
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
		ga.State = ga.Gears.Selected.ContainerSsh(true, ga.SshStatus, ga.Mount, ga.Args)

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
