#!/bin/bash

# WARNING: This file is SOURCED. Don't add in any "exit", otherwise your shell will exit.
export ARCH GB_BINFILE GB_ARCHDIR GB_BINDIR GB_BASEDIR GB_JSONFILE GB_VERSIONS GB_VERSION GITBIN GB_GITURL GB_GITREPO

ARCH="$(uname -s)"
case ${ARCH} in
	'Linux')
		LOG_ARGS='-t 10'
		;;
	*)
		LOG_ARGS='-r -t 10'
		;;
esac

GB_BINFILE="$(./bin/${ARCH}/JsonToConfig -json-string '{}' -template-string '{{ .ExecName }}')"
GB_ARCHDIR="$(./bin/${ARCH}/JsonToConfig -json-string '{}' -template-string '{{ .DirPath }}')"
GB_BINDIR="$(dirname "$GB_ARCHDIR")"
GB_BASEDIR="$(dirname "$GB_BINDIR")"
GB_JSONFILE="${GB_BASEDIR}/gearbox.json"

LAUNCHBIN="${GB_BINDIR}/${ARCH}/launch"

if [ -f "${GB_JSONFILE}" ]
then
	GB_VERSIONS="$(${GB_BINFILE} -json ${GB_JSONFILE} -template-string '{{ range $version, $value := .Json.versions }}{{ if ne $version "" }}{{ $version }}{{ end }} {{ end }}')"
	GB_VERSIONS="$(echo ${GB_VERSIONS})"	# Easily remove CR

	GB_IMAGENAME="$(${GB_BINFILE} -json ${GB_JSONFILE} -template-string '{{ .Json.meta.organization }}/{{ .Json.meta.name }}')"

	GB_NAME="$(${GB_BINFILE} -json ${GB_JSONFILE} -template-string '{{ .Json.meta.name }}')"
fi

GB_BASE="$(${GB_BINFILE} -json gearbox.json -template-string '{{ .Json.build.base }}')"

GITBIN="$(which git)"
GB_GITURL="$(${GITBIN} config --get remote.origin.url)"
if [ "${GB_GITURL}" == "" ]
then
	GB_GITREPO=""
else
	GB_GITREPO="$(basename -s .git ${GB_GITURL})"
fi

. ${GB_BINDIR}/_Colors.sh


################################################################################
_getVersions() {
	if [ ! -f "${GB_JSONFILE}" ]
	then
		c_err "Can't find JSON file: ${GB_JSONFILE}"
		return 0
	fi

	if [ "${GB_VERSIONS}" == "" ]
	then
		c_err "No versions found"
		return 0
	fi

	if [ "${GB_GITREPO}" == "docker-template" ]
	then
		c_warn "Cannot run this command from the docker-template repository."
		c_warn "IF THIS IS A COPY of that repo, then..."
		c_warn "	1. Remove the .git directory."
		c_warn "	2. Run the \"make init\" command."
		c_warn ""
		unset GB_VERSIONS GB_GITURL GB_GITREPO

	else
		case $1 in
			'all')
				;;
			'')
				c_warn "No versions specified."
				c_info "Versions available:"
				_listVersions
				unset GB_VERSIONS
				return 0
				;;
			*)
				GB_VERSIONS="$@"
				;;
		esac
	fi

	return 1
}


################################################################################
_listVersions() {
	echo "	all - All versions"
	${GB_BINFILE} -json ${GB_JSONFILE} -template-string '{{ range $version, $value := .Json.versions }}{{ if ne $version "" }}\t{{ $version }} - {{ $.Json.meta.organization }}/{{ $.Json.meta.name }}:{{ $version }}\n{{ end }}{{ end }}'
	echo ""
}


################################################################################
gb_getenv() {
	GB_VERDIR="${GB_BASEDIR}/versions/$1"
	export GB_VERDIR

	if [ -f "TEMPLATE/version/.env.tmpl" ]
	then
		${GB_BINFILE} -json "${GB_JSONFILE}" -template "TEMPLATE/version/.env.tmpl" -out "${GB_VERDIR}/.env"
	fi

	if [ -f "${GB_VERDIR}/.env" ]
	then
		. "${GB_VERDIR}/.env"
	fi
}


################################################################################
gb_getdockerfile() {
	GB_VERDIR="${GB_BASEDIR}/versions/$1"
	export GB_VERDIR

	if [ -f "TEMPLATE/version/DockerfileRuntime.tmpl" ]
	then
		${GB_BINFILE} -json "${GB_JSONFILE}" -template "TEMPLATE/version/DockerfileRuntime.tmpl" -out "${GB_VERDIR}/DockerfileRuntime"
	fi
}


################################################################################
gb_checkImage() {
	STATE="$(docker image ls -q "$1")"
	if [ "${STATE}" == "" ]
	then
		# Not created.
		STATE="MISSING"
	else
		STATE="PRESENT"
	fi
}


################################################################################
gb_checkContainer() {
	STATE="$(docker container ls -q -a -f name="^$1")"
	if [ "${STATE}" == "" ]
	then
		# Not created.
		STATE="MISSING"
		return
	fi

	STATE="$(docker container ls -q -f name="^$1")"
	if [ "${STATE}" == "" ]
	then
		# Not created.
		STATE="STOPPED"
		return
	fi

	STATE="STARTED"
}


################################################################################
gb_checknetwork() {
	STATE="$(docker network ls -qf "name=gearboxnet")"
	if [ "${STATE}" == "" ]
	then
		# Create network
		echo "Creating network"
		docker network create --subnet 172.42.0.0/24 gearboxnet
	fi
}


################################################################################
gb_init() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "Initializing repo."

	gb_create-build ${GB_JSONFILE}
	gb_create-version ${GB_JSONFILE}
	# ${DIR}/JsonToConfig-$(uname -s) -json "${GB_JSONFILE}" -template TEMPLATE/README.md.tmpl -out README.md

	return 0
}


################################################################################
gb_create-build() {
	if _getVersions $@
	then
		return 1
	fi

	if [ -d build ]
	then
		p_ok "${FUNCNAME[0]}" "Updating build directory."
	else
		p_ok "${FUNCNAME[0]}" "Creating build directory."
		cp -i TEMPLATE/build.sh.tmpl .
		${GB_BINFILE} -json ${GB_JSONFILE} -create build.sh.tmpl -shell
		rm -f build.sh.tmpl build.sh
	fi

	${GB_BINFILE} -template ./TEMPLATE/README.md.tmpl -json ${GB_JSONFILE} -out README.md

	cp ./TEMPLATE/Makefile .
	if [ "${GB_BASE}" == "true" ]
	then
		cp "${GB_JSONFILE}" build/
	else
		cp "${GB_JSONFILE}" "build/gearbox-${GB_NAME}.json"
	fi

	return 0
}


################################################################################
gb_create-version() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "Creating/updating version directory for versions: ${GB_VERSIONS}"

	for GB_VERSION in ${GB_VERSIONS}
	do
		if [ -d "${GB_BASEDIR}/versions/${GB_VERSION}" ]
		then
			gb_getenv ${GB_VERSION}
			p_info "${FUNCNAME[0]}" "Updating version directory \"${GB_VERSION}\"."
			${GB_BINFILE} -json ${GB_JSONFILE} -template ./TEMPLATE/version/.env.tmpl -out "${GB_VERDIR}/.env"
			${GB_BINFILE} -json ${GB_JSONFILE} -template ./TEMPLATE/version/DockerfileRuntime.tmpl -out "${GB_VERDIR}/DockerfileRuntime"
			rm -f "${GB_VERDIR}/gearbox.json"

		else
			p_info "${FUNCNAME[0]}" "Creating version directory \"${GB_VERSION}\"."
			cp -i TEMPLATE/version.sh.tmpl .
			${GB_BINFILE} -json ${GB_JSONFILE} -create version.sh.tmpl -shell
			rm -f version.sh.tmpl version.sh "${GB_VERDIR}/gearbox.json"
		fi
	done

	${GB_BINFILE} -json ${GB_JSONFILE} -template ./TEMPLATE/README.md.tmpl -out README.md

	return 0
}


################################################################################
gb_clean() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Cleaning up for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		gb_getenv ${GB_VERSION}


		p_info "${GB_IMAGEVERSION}" "Removing logs."
		rm -f ${GB_VERDIR}/logs/*.log


		${LAUNCHBIN} uninstall "${GB_NAME}:${GB_VERSION}"
		RETURN="$?"
		if [ "${RETURN}" != "0" ]
		then
			p_err "${FUNCNAME[0]}" "Error exit code: ${RETURN}"
			EXIT="1"
		fi


		gb_checkImage ${GB_IMAGEMAJORVERSION}
		case ${STATE} in
			'PRESENT')
				p_info "${GB_IMAGEMAJORVERSION}" "Removing image."
				docker image rm -f ${GB_IMAGEMAJORVERSION}
				;;
			*)
				p_warn "${GB_IMAGEMAJORVERSION}" "Image already removed."
				;;
		esac


		gb_checkImage ${GB_IMAGEVERSION}
		case ${STATE} in
			'PRESENT')
				p_info "${GB_IMAGEVERSION}" "Removing image."
				docker image rm -f ${GB_IMAGEVERSION}
				;;
			*)
				p_warn "${GB_IMAGEVERSION}" "Image already removed."
				;;
		esac


		gb_checkImage ${GB_IMAGENAME}:latest
		case ${STATE} in
			'PRESENT')
				p_info "${GB_IMAGENAME}:latest" "Removing image."
				docker image rm -f ${GB_IMAGENAME}:latest
				;;
			*)
				p_warn "${GB_IMAGENAME}:latest" "Image already removed."
				;;
		esac
	done

	return ${EXIT}
}


################################################################################
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
			p_info "${GB_IMAGENAME}:${GB_VERSION}" "Query ref container."
			GEARBOX_ENTRYPOINT="$(docker inspect --format '{{ with .ContainerConfig.Entrypoint}} {{ index . 0 }}{{ end }}' "${GB_REF}")"
			export GEARBOX_ENTRYPOINT
			GEARBOX_ENTRYPOINT_ARGS="$(docker inspect --format '{{ join .ContainerConfig.Entrypoint " " }}' "${GB_REF}")"
			export GEARBOX_ENTRYPOINT_ARGS
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


################################################################################
gb_info() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Image and container info for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		gb_getenv ${GB_VERSION}

		p_info "${GB_IMAGEMAJORVERSION}" "List image."
		docker image ls ${GB_IMAGEMAJORVERSION}
		p_info "${GB_IMAGEVERSION}" "List image."
		docker image ls ${GB_IMAGEVERSION}

		echo "# Gearbox[${GB_CONTAINERMAJORVERSION}]: List container."
		docker container ls -f name="^${GB_CONTAINERMAJORVERSION}"
		p_info "${GB_CONTAINERVERSION}" "List container."
		docker container ls -f name="^${GB_CONTAINERVERSION}"
	done

	return ${EXIT}
}


################################################################################
gb_inspect() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Inspecting image and container for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		gb_getenv ${GB_VERSION}

		p_info "${GB_IMAGEMAJORVERSION}" "Inspect image."
		docker image inspect ${GB_IMAGEMAJORVERSION} 2>&1
		p_info "${GB_IMAGEVERSION}" "Inspect image."
		docker image inspect ${GB_IMAGEVERSION} 2>&1

		echo "# Gearbox[${GB_CONTAINERMAJORVERSION}]: Inspect container."
		docker container inspect name="^${GB_CONTAINERMAJORVERSION}" 2>&1
		p_info "${GB_CONTAINERVERSION}" "Inspect container."
		docker container inspect name="^${GB_CONTAINERVERSION}" 2>&1
	done

	return ${EXIT}
}


################################################################################
gb_list() {
	if _getVersions $@
	then
		return 1
	fi

	${LAUNCHBIN} list "${GB_NAME}:${GB_VERSION}"
	RETURN="$?"
	return ${RETURN}
}


################################################################################
gb_logs() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Showing build logs for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		gb_getenv ${GB_VERSION}

		if [ -f "${GB_VERDIR}/logs/build.log" ]
		then
			p_info "${GB_IMAGEMAJORVERSION}" "Showing logs."
			script -dp "${GB_VERDIR}/logs/build.log" | less -SinR
		else
			p_warn "${GB_IMAGEMAJORVERSION}" "No logs."
		fi
	done

	return ${EXIT}
}


################################################################################
gb_ports() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Showing ports for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		gb_getenv ${GB_VERSION}

		gb_checkContainer ${GB_CONTAINERVERSION}
		case ${STATE} in
			'STARTED')
				p_info "${GB_CONTAINERVERSION}" "Showing exposed container ports."
				docker port ${GB_CONTAINERVERSION}
				;;
			'STOPPED')
				p_info "${GB_CONTAINERVERSION}" "Container needs to be started."
				;;
			'MISSING')
				p_info "${GB_CONTAINERVERSION}" "Need to create container first."
				;;
			*)
				p_err "${GB_CONTAINERVERSION}" "Unknown state."
				EXIT="1"
				continue
				;;
		esac
	done

	return ${EXIT}
}


################################################################################
gb_dockerhub() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Pushing to DockerHub for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		gb_getenv ${GB_VERSION}

		p_info "${GB_IMAGEVERSION}" "Pushing image to DockerHub."
		docker push ${GB_IMAGEVERSION}

		if [ "${GB_IMAGEMAJORVERSION}" != "" ]
		then
			p_info "${GB_IMAGEMAJORVERSION}" "Pushing image to DockerHub."
			docker push ${GB_IMAGEMAJORVERSION}
		fi

		if [ "${GB_LATEST}" == "true" ]
		then
			p_info "${GB_IMAGENAME}:latest" "Pushing image to DockerHub."
			docker push "${GB_IMAGENAME}:latest"
		fi
	done

	if [ "${EXIT}" != "0" ]
	then
		p_err "${FUNCNAME[0]}" "Push FAILED for versions: ${GB_VERSIONS}"
	fi
	return ${EXIT}
}


################################################################################
gb_github() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Pushing to GitHub for repo."

	if [ "${GITHUB_ACTIONS}" != "" ]
	then
		echo "# Gearbox[${GB_GITREPO}]: Running from GitHub action - ignoring."
		return 1
	fi

	echo "# Gearbox[${GB_GITREPO}]: Pushing repo to GitHub."
	git commit -a -m "Latest push" && git push

	return 0
}


################################################################################
gb_push() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Pushing to GitHub and DockerHub for versions: ${GB_VERSIONS}"

	gb_dockerhub ${GB_VERSIONS}
	gb_github ${GB_VERSIONS}
	return 0
}


################################################################################
gb_release() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Releasing for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		p_ok "${FUNCNAME[0]}" "#### Build and release version ${GB_VERSION}"

		gb_clean ${GB_VERSION}

		gb_build ${GB_VERSION}
		RETURN="$?"
		if [ "${RETURN}" != "0" ]
		then
			p_err "${FUNCNAME[0]}" "Error exit code: ${RETURN}"
			EXIT="1"
			continue
		fi

		# Just after a build, the image won't be visible for a short while.
		sleep 2

		# gb_test ${GB_VERSION}
		${LAUNCHBIN} test "${GB_NAME}:${GB_VERSION}"
		RETURN="$?"
		if [ "${RETURN}" != "0" ]
		then
			p_err "${FUNCNAME[0]}" "Error exit code: ${RETURN}"
			EXIT="1"
			gb_clean ${GB_VERSION}
			continue
		fi

		gb_dockerhub ${GB_VERSION}
		RETURN="$?"
		if [ "${RETURN}" != "0" ]
		then
			p_err "${FUNCNAME[0]}" "Error exit code: ${RETURN}"
			gb_clean ${GB_VERSION}
			EXIT="1"
		fi

		gb_clean ${GB_VERSION}
	done

	if [ "${EXIT}" != "0" ]
	then
		p_err "${FUNCNAME[0]}" "Release FAILED for versions: ${GB_VERSIONS}"
	fi
	return ${EXIT}
}


################################################################################
gb_rm() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Removing container for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		gb_getenv ${GB_VERSION}

		${LAUNCHBIN} uninstall "${GB_NAME}:${GB_VERSION}"
		RETURN="$?"
		if [ "${RETURN}" != "0" ]
		then
			EXIT="1"
		fi
	done

	return ${EXIT}
}


################################################################################
gb_shell() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Running shell for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		gb_getenv ${GB_VERSION}

		${LAUNCHBIN} shell "${GB_NAME}:${GB_VERSION}"
		RETURN="$?"
		if [ "${RETURN}" != "0" ]
		then
			EXIT="1"
		fi
	done

	if [ "${EXIT}" != "0" ]
	then
		p_err "${FUNCNAME[0]}" "Shell FAILED for versions: ${GB_VERSIONS}"
	fi
	return ${EXIT}
}


################################################################################
gb_bash() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Running shell for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		gb_getenv ${GB_VERSION}

		${LAUNCHBIN} start "${GB_NAME}:${GB_VERSION}"
		RETURN="$?"
		if [ "${RETURN}" != "0" ]
		then
			EXIT="1"
			p_err "${GB_CONTAINERVERSION}" "Failed to start."
			continue
		fi

		p_info "${GB_CONTAINERVERSION}" "Entering container."
		docker exec -i -t ${GB_CONTAINERVERSION} /bin/bash -l
	done

	if [ "${EXIT}" != "0" ]
	then
		p_err "${FUNCNAME[0]}" "Shell FAILED for versions: ${GB_VERSIONS}"
	fi
	return 0
}


################################################################################
gb_start() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "#### Starting container for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		gb_getenv ${GB_VERSION}

		${LAUNCHBIN} start "${GB_NAME}:${GB_VERSION}"
		RETURN="$?"
		if [ "${RETURN}" != "0" ]
		then
			EXIT="1"
		fi
	done

	if [ "${EXIT}" != "0" ]
	then
		p_err "${FUNCNAME[0]}" "Start FAILED for versions: ${GB_VERSIONS}"
	fi
	return ${EXIT}
}


################################################################################
gb_stop() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "Stopping container for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		gb_getenv ${GB_VERSION}

		${LAUNCHBIN} stop "${GB_NAME}:${GB_VERSION}"
		RETURN="$?"
		if [ "${RETURN}" != "0" ]
		then
			EXIT="1"
		fi
	done

	if [ "${EXIT}" != "0" ]
	then
		p_err "${FUNCNAME[0]}" "Stop FAILED for versions: ${GB_VERSIONS}"
	fi
	return ${EXIT}
}


################################################################################
gb_test() {
	if _getVersions $@
	then
		return 1
	fi
	p_ok "${FUNCNAME[0]}" "Testing container for versions: ${GB_VERSIONS}"

	EXIT="0"
	for GB_VERSION in ${GB_VERSIONS}
	do
		gb_getenv ${GB_VERSION}

		${LAUNCHBIN} test "${GB_NAME}:${GB_VERSION}"
		RETURN="$?"
		if [ "${RETURN}" != "0" ]
		then
			EXIT="1"
		fi
	done

	if [ "${EXIT}" != "0" ]
	then
		p_err "${FUNCNAME[0]}" "Testing FAILED for versions: ${GB_VERSIONS}"
	fi
	return ${EXIT}
}


