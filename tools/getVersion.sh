#!/bin/sh
awk '/BinaryVersion/{gsub("\"", ""); print$4}' defaults/version.go
