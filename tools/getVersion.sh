#!/bin/sh
VERSION="$(awk '/BinaryVersion/{gsub("\"", ""); print$4}' defaults/version.go)"
if [ -z "${VERSION}" ]
then
	VERSION="$(awk '/BinaryVersion/{gsub("\"", ""); print$3}' defaults/version.go)"
fi
echo "${VERSION}"
