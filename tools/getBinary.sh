#!/bin/sh
VERSION="$(awk '/BinaryName/{gsub("\"", ""); print$4}' defaults/version.go)"
if [ -z "${VERSION}" ]
then
	VERSION="$(awk '/BinaryName/{gsub("\"", ""); print$3}' defaults/version.go)"
fi
echo "${VERSION}"
