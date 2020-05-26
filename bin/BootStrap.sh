#!/bin/bash

GB_GITREPO="$1"; export GB_GITREPO
GB_JSON="$2"

################################################################################
GB_GITURL="$(git config --get remote.origin.url)"; export GB_GITURL
if [ "${GB_GITURL}" != "" ]
then
	GB_GITREPO="$(basename -s .git ${GB_GITURL})"
fi

if [ "${GB_GITREPO}" == "docker-template" ]
then
	echo "################################### WARNING ###################################"
	echo "# This command can not ever be run from the docker-template GitHub repository."
	echo "################################### WARNING ###################################"
	exit 1
fi


echo "################################################################################"
echo "# Gearbox[bootstrap]: Create a new docker image repository from docker-template."
if [ "${GB_GITREPO}" == "" ]
then
	echo -n "Enter name of new repository: "
	read GB_GITREPO
	if [ "${GB_GITREPO}" == "" ]
	then
		echo "# Gearbox[bootstrap]: Nothing entered. Doing nothing."
		exit 1
	fi
fi


################################################################################
echo "# Gearbox[bootstrap]: Pulling docker-template from GitHub."
git clone https://github.com/gearboxworks/docker-template "${GB_GITREPO}"
if [ ! -d "${GB_GITREPO}" ]
then
	echo "# Gearbox[bootstrap]: Problems creating directory ${GB_GITREPO}."
	exit 1
fi
rm -rf "${GB_GITREPO}/.git"


################################################################################
if [ "${GB_JSON}" != "" ]
then
	echo "# Gearbox[bootstrap]: JSON file provided. Running 'make init' for you."
	cp "${GB_JSON}" "${GB_GITREPO}/${GB_GITREPO}.json"
	cd "${GB_GITREPO}"
	make init
	echo "# Gearbox[bootstrap]: Completed OK."
else
	echo "# Gearbox[bootstrap]: No JSON file provided."
	echo "# Gearbox[bootstrap]: Create one using 'gearbox-TEMPLATE.json' as a reference."
	echo "# Gearbox[bootstrap]: Then run 'make init' manually."
fi

