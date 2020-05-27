#!/bin/sh

BINARYREPO="$(tools/getBinaryRepo.sh)"
USER="$(echo "${BINARYREPO}" | awk -F/ '{print$1}')"
REPO="$(echo "${BINARYREPO}" | awk -F/ '{print$1}')"

VERSION="$(tools/getVersion.sh)"

set -x
github-release release \
	--user "${USER}" \
	--repo "${GB_GITREPO}" \
	--tag "${VERSION}" \
	--name "Release ${VERSION}" \
	--description "${DESCRIPTION}"

