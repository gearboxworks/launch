![Gearbox](https://gearboxworks.github.io/assets/images/gearbox-logo.png)


# [Gearbox](https://github.com/gearboxworks/) Docker container template
This is the main template repository for creating Docker containers for [Gearbox](https://github.com/gearboxworks/).

Use it to create a Gearbox container from scratch or wrap an already existing container to be used within Gearbox.


## Repository Info
GitHub commit: ![commit-date](https://img.shields.io/github/last-commit/gearboxworks/docker-template?style=flat-square)

GitHub release(latest): ![last-release-date](https://img.shields.io/github/release-date/gearboxworks/docker-template) ![last-release-date](https://img.shields.io/github/v/tag/gearboxworks/docker-template?sort=semver)

All supported Docker containers: [here](gears/)

[![GoDoc Badge][]][GoDoc]


* * *

## Simple setup, (using BootStrap.sh).

### 1 Fetch this repo.
Open a terminal window and fetch using wget.

`wget https://github.com/gearboxworks/docker-template/raw/master/bin/BootStrap.sh`

Or download [from here](bin/BootStrap.sh)

### 2 Run BootStrap.sh script.
`bash ./BootStrap.sh` - Let the script prompt you:

`bash ./BootStrap.sh my-container` - Specify a directory. You will have to create your JSON file and run `make init` afterwards.

`bash ./BootStrap.sh my-container MyContainer.json` - Specify a directory and JSON file - this will run the final `make init` for you.


* * *

## Setup via GitHub.

### 1 Fetch this repo.

`git clone https://github.com/gearboxworks/docker-template`

`cd docker-template`

`rm -rf .git`

### 2 Create JSON file.

Check out the gearbox-TEMPLATE.json file for an example.

### 3 Initialize repo based off JSON file.

`make init`

This will generate a `build` directory as well as an version directories that may be defined in the JSON file.

At this point you can remove the `TEMPLATE` directory and src JSON file. Or leave them for later updates.

### 4 Update scripts, directories and any other files under rootfs.

	- 4.7.6/DockerfileRuntime
	- 4.7.6/gearbox.json
	- 4.7.6/logs
	- 4.7.6/logs/.keep
	- build/build-adminer.apks
	- build/build-adminer.env
	- build/build-adminer.sh
	- build/gearbox-adminer.json
	- build/rootfs
	- build/rootfs/etc
	- build/rootfs/etc/gearbox
	- build/rootfs/etc/gearbox/services
	- build/rootfs/etc/gearbox/services/adminer
	- build/rootfs/etc/gearbox/services/adminer/finish
	- build/rootfs/etc/gearbox/services/adminer/run
	- build/rootfs/etc/gearbox/unit-tests
	- build/rootfs/etc/gearbox/unit-tests/adminer
	- build/rootfs/etc/gearbox/unit-tests/adminer/01-base.sh


* * *

## Using the docker-template tools.

### Creating docker images.

`make build-4.7.6` - Build a specific version.

`make build-all` - Build all versions.


### Show build logs.

`make logs-4.7.6` - Show logs from a specific version.

`make logs-all` - Show logs from all versions.


### Show image and container information.

`make info-4.7.6` - Info on a specific version.

`make info-all` - Info on all versions.


### Cleanup all docker images and containers.

`make clean-4.7.6` - Cleanup a specific version.

`make clean-all` - Cleanup all versions.


### Test a build.

`make test-4.7.6` - Test a specific version.

`make test-all` - Test all versions.


### Maybe run a shell.

`make shell-4.7.6` - Shell into a specific version.

`make shell-all` - Shell into all versions.


### Push repos.

`make push-4.7.6` - Push specific version to GitHub and DockerHub.

`make push-all` - Push all versions to GitHub and DockerHub.


### Run a full release cycle.

`make release-4.7.6` - Release a specific version.

`make release-all` - Release all versions.

This will perform the following actions:
- make clean
- make build
- make test
- make push

If any step fails in the sequence, the process will be terminated.


* * *

## Additional docker-template tools.

### Adding a new version.
To add a new version of a docker image.
- Create a JSON file - Either a new one, or copying an older version and editing with updated version details.
- `./bin/CreateVersion.sh my-new-version.json` - CreateVersion.sh will populate the version directory.

### Updating docker-template.
New releases of the docker-template repository can be pulled using:

`./bin/TemplateUpdate.sh`

This will pull the latest release and update the TEMPLATE, bin and Makefile directories.

Alternatively, you can specifying a version different to "latest".

`./bin/TemplateUpdate.sh 1.0.1`


* * *

### Updating this template container.

*NOTE: For gearboxworks staff only...*

```
TemplateRelease.sh
    Updates the docker-template GitHub repository with either a new or updated release.

TemplateRelease.sh create [version] - Creates a new release on GitHub.
TemplateRelease.sh update [version] - Updates an existing release on GitHub.

if [version] isn't specified, then...
- Prompt the user for a new version.
- Do nothing.
```
