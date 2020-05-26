#!/bin/bash

DIR="$(dirname $0)"
CMD="$1"
VERSION="$2"

help() {
cat<<EOF

$(basename $0)
	Updates docker-template files from GitHub with either a specified version or latest.

$(basename $0) [version] - Updates TEMPLATE files.

EOF
}


################################################################################
GB_GITURL="$(git config --get remote.origin.url)"; export GB_GITURL
if [ "${GB_GITURL}" == "" ]
then
	GB_GITREPO=""; export GB_GITREPO
else
	GB_GITREPO="$(basename -s .git ${GB_GITURL})"; export GB_GITREPO
fi

if [ "${GB_GITREPO}" == "docker-template" ]
then
	echo "################################### WARNING ###################################"
	echo "# This command can not ever be run from the docker-template GitHub repository."
	echo "################################### WARNING ###################################"
	exit 1
fi

if [ "${VERSION}" == "" ]
then
	VERSION="-l"
else
	VERSION="-t ${VERSION}"
fi


echo "################################################################################"
echo "# Gearbox[docker-template]: Updating docker-template files from GitHub."

${DIR}/github-release download \
	--user "gearboxworks" \
	--repo "docker-template" \
	${VERSION} \
	--name docker-template.tgz

if [ -f docker-template.tgz ]
then
	echo "# Gearbox[docker-template]: Extracting."
	mkdir docker-template/
	tar zxf docker-template.tgz -C docker-template/
	rsync -HvaxP --delete docker-template/* .
	if [ ! -d .github/workflows ]
	then
		mkdir -p .github/workflows
	fi
	cp TEMPLATE/release.yml .github/workflows/
	rm -rf docker-template.tgz docker-template/
	echo "# Gearbox[docker-template]: Done."

	${DIR}/JsonToConfig -template ./TEMPLATE/README.md.tmpl -json gearbox.json -out README.md
else
	echo "# Gearbox[docker-template]: Cannot find docker-template repository."
fi

make init
