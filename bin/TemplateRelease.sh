#!/bin/bash

DIR="$(dirname $0)"
CMD="$1"
VERSION="$2"

shift
shift
DESCRIPTION="$@"


help() {
cat<<EOF

$(basename $0)
	Updates the docker-template GitHub repository with either a new or updated release.

$(basename $0) create [version] [description] - Creates a new release on GitHub.
$(basename $0) update [version] [description] - Updates an existing release on GitHub.

if [version] isn't specified, then...
- Prompt the user for a new version.
- Do nothing.
EOF
}

check() {
	echo "# Gearbox[${GB_GITREPO}]: Checking if release v${VERSION} exists."
	${DIR}/github-release info \
		-u gearboxworks \
		-r "${GB_GITREPO}" \
		-t "${VERSION}" >& /dev/null
	RETURN="$?"
	# 0 - Exists.
	# 1 - Doesn't exist.
}

create() {
	echo "# Gearbox[${GB_GITREPO}]: Creating release v${VERSION} on GitHub."
	${DIR}/github-release release \
		--user "gearboxworks" \
		--repo "${GB_GITREPO}" \
		--tag "${VERSION}" \
		--name "Release ${VERSION}" \
		--description "${DESCRIPTION}"
}

update() {
	echo "# Gearbox[${GB_GITREPO}]: Updating release v${VERSION} on GitHub."
	${DIR}/github-release edit \
		--user "gearboxworks" \
		--repo "${GB_GITREPO}" \
		--tag "${VERSION}" \
		--name "Release ${VERSION}" \
		--description "${DESCRIPTION}"
}

upload() {
	FILES="Makefile TEMPLATE bin"
	echo "# Gearbox[${GB_GITREPO}]: Creating tarball release from files: ${FILES}"
	tar zcf docker-template.tgz ${FILES}

	echo "# Gearbox[${GB_GITREPO}]: Uploading tarball release v${VERSION} to GitHub."
	${DIR}/github-release upload \
		--user "gearboxworks" \
		--repo "${GB_GITREPO}" \
		--tag "${VERSION}" \
		--name "docker-template.tgz" \
		--label "docker-template.tgz" \
		-R \
		-f docker-template.tgz

	rm -f docker-template.tgz
}

find_last_version() {
	${DIR}/github-release info -u gearboxworks -r "${GB_GITREPO}" -j > /tmp/z
	LAST_VERSION="$(${DIR}/JsonToConfig -json /tmp/z -template-string '{{ with index .Json.Releases 0 }}{{ .tag_name }}{{ end }}')"
}

get_description() {
	GET_VERSION="$1"; export GET_VERSION
	${DIR}/github-release info -u gearboxworks -r "${GB_GITREPO}" -j > /tmp/z
	LAST_DESCRIPTION="$(${DIR}/JsonToConfig -json /tmp/z -template-string '{{ range $k,$v := .Json.Releases }}{{ if eq .tag_name $.Env.GET_VERSION }}{{ .body }}{{ end }}{{ end }}')"
}


################################################################################
GB_GITURL="$(git config --get remote.origin.url)"; export GB_GITURL
if [ "${GB_GITURL}" == "" ]
then
	GB_GITREPO=""; export GB_GITREPO
else
	GB_GITREPO="$(basename -s .git ${GB_GITURL})"; export GB_GITREPO
fi

if [ "${GB_GITREPO}" != "docker-template" ]
then
	./bin/TemplateUpdate.sh
fi


################################################################################
if [ "${CMD}" == "" ]
then
	echo "# Gearbox[${GB_GITREPO}]: Doing nothing."
	help

	find_last_version
	echo "# Gearbox[${GB_GITREPO}]: Last version: ${LAST_VERSION}"

	get_description ${LAST_VERSION}
	echo "# Gearbox[${GB_GITREPO}]: Last description: \"${LAST_DESCRIPTION}\""

	exit 1
fi

if [ "${VERSION}" == "" ]
then
	echo "# Gearbox[${GB_GITREPO}]: No release version specified."

	find_last_version
	echo "# Gearbox[${GB_GITREPO}]: Last version in repo is ${LAST_VERSION}"
	echo -n "Enter a release version: "
	read VERSION

	if [ "${VERSION}" == "" ]
	then
		echo "# Gearbox[${GB_GITREPO}]: No version entered? OK doing nothing."
		exit 1
	fi

#	if [ "${VERSION}" == "${LAST_VERSION}" ]
#	then
#	fi
fi

if [ "${DESCRIPTION}" == "" ]
then
	get_description ${VERSION}
	if [ "${LAST_DESCRIPTION}" == "" ]
	then
		DESCRIPTION="Release ${VERSION}"
	else
		DESCRIPTION="${LAST_DESCRIPTION}"
	fi
fi


if [ "${GB_GITREPO}" == "" ]
then
	echo "# Gearbox[${GB_GITREPO}]: GB_GITREPO isn't set... Strange..."
	echo "# Gearbox[${GB_GITREPO}]: Abandoning GitHub release changes."
	exit 1
fi

if [ "${GITHUB_USER}" == "" ]
then
	echo "# Gearbox[${GB_GITREPO}]: GITHUB_USER needs to be set to ${CMD} a release."
	echo "# Gearbox[${GB_GITREPO}]: Abandoning GitHub release changes."
	exit 1
fi

if [ "${GITHUB_TOKEN}" == "" ]
then
	echo "# Gearbox[${GB_GITREPO}]: GITHUB_TOKEN needs to be set to ${CMD} a release."
	echo "# Gearbox[${GB_GITREPO}]: Abandoning GitHub release changes."
	exit 1
fi


echo "################################################################################"
echo "# Gearbox[${GB_GITREPO}]: Creating release ${VERSION} on GitHub."
echo "# Gearbox[${GB_GITREPO}]: Description: \"${DESCRIPTION}.\""

echo "# Gearbox[${GB_GITREPO}]: Pushing repo to GitHub."
git add . && git commit -a -m "${DESCRIPTION}" && git push

export RETURN
case "${CMD}" in
	'create')
		check
		if [ "${RETURN}" == "0" ]
		then
			echo "# Gearbox[${GB_GITREPO}]: Release v${VERSION} already exists. Abandoning create."
			exit 1
		fi

		create

		if [ "${GB_GITREPO}" == "docker-template" ]
		then
			upload
		fi

		echo "# Gearbox[${GB_GITREPO}]: Release v${VERSION} OK."
		;;

	'update')
		check
		if [ "${RETURN}" != "0" ]
		then
			echo "# Gearbox[${GB_GITREPO}]: No release v${VERSION} found."
			exit 1
		fi

		if [ "${GB_GITREPO}" == "docker-template" ]
		then
			upload
		fi

		update

		echo "# Gearbox[${GB_GITREPO}]: Release v${VERSION} OK."
		;;

	'delete')
		;;
esac

