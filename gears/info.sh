#!/bin/bash

export D

cat <<EOF
# Docker auto-release build status.
The following represents the status of all auto-release builds for Gearbox docker images.


| Repository | Last Commit | Release Date | Release Version | Release State | Docker Latest |
| ---------- | ----------- | ------------ | --------------- | ------------- | ------------- |
EOF

for D in docker-*
do
	if [ ! -d ${D} ]
	then
		continue
	fi

	if [ "${D}" == "docker-test" ]
	then
		continue
	fi

	perl -e '
BEGIN {
	use Env qw(D);
	$N=$D;
	$N=~s/^docker-//;
}

$U = sprintf("[%s](https://github.com/gearboxworks/%s/)", $D, $D);
$last_commit = sprintf("![commit-date](https://img.shields.io/github/last-commit/gearboxworks/%s?style=flat-square)", $D);
$release_date = sprintf("![release-date](https://img.shields.io/github/release-date/gearboxworks/%s)", $D);
$release_version = sprintf("![release-date](https://img.shields.io/github/v/tag/gearboxworks/%s?sort=semver)", $D);
$release_state = sprintf("![release-state](https://github.com/gearboxworks/%s/workflows/release/badge.svg?event=release)", $D);
$workflow = sprintf("[%s](https://github.com/gearboxworks/%s/actions?query=workflow%3Arelease)", $release_state, $D);
$docker_latest_by_semver = sprintf("![docker-latest](https://img.shields.io/docker/v/gearboxworks/%s?sort=semver)", $N);
# $docker_size = sprintf("![Docker Image Size (tag)](https://img.shields.io/docker/image-size/gearboxworks/adminer/4.7.6)", $N);
$docker_latest_by_date = sprintf("![docker-latest](https://img.shields.io/docker/v/gearboxworks/%s)", $N);

if ($D eq "docker-template") {
	$docker_latest_by_semver = $workflow = "N/A";
} elsif ($D eq "docker-repo") {
	exit();
}

printf("| %s | %s | %s | %s | %s | %s |\n",
	$U,
	$last_commit,
	$release_date,
	$release_version,
	$workflow,
	$docker_latest_by_semver
	);
'
done

