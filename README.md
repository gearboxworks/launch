![Gearbox](docs/Gearbox-100x.png)


## About `launch`

`launch` is a tool specifically designed to simplify running Docker containers.

Launch provides three (4) important functional areas, without any Docker container learning curve:

- Allows for "install" of applications as though natively installed - even on a remote Docker server!
- Allows control over Gearbox Docker containers: stop, start, create, remove.
- Build, update, modify and release Gearbox Docker images.
- Acts as a proxy for interactive commands within a Gearbox Docker container.

It also provides a functional SSH daemon for connecting remotely as well as a standard set of common tools and utilities.


### Download launch
`launch` is currently in beta testing and is included along with all Gearbox Docker repos.
Once out of beta, it will be included within the Gearbox installation package.

For now, simply download the standalone `launch` binary for your O/S.
| Mac OSX | Linux | Linux | Windows | Windows | SBC | SBC |
| ------- | ----- | ----- | ------- | ------- | --- | --- |
| [![Mac OSX](docs/logos/64x64/mac.png) Mac OSX 64bit](https://github.com/gearboxworks/launch/releases/download/latest/launch-darwin_amd64.tar.gz) | [![Linux](docs/logos/64x64/linux.png) Linux 64bit](https://github.com/gearboxworks/launch/releases/download/latest/launch-linux_amd64.tar.gz) | [![Linux](docs/logos/64x64/linux.png) Linux 32bit](https://github.com/gearboxworks/launch/releases/download/latest/launch-linux_i386.tar.gz) | [![Windows](docs/logos/64x64/windows.png) Windows 64bit](https://github.com/gearboxworks/launch/releases/download/latest/launch-windows_amd64.tar.gz) | [![Windows](docs/logos/64x64/windows.png) Windows 32bit](https://github.com/gearboxworks/launch/releases/download/latest/launch-windows_i386.tar.gz) | [![Raspberry Pi](docs/logos/64x64/linux.png) Raspberry Pi 32bit](https://github.com/gearboxworks/launch/releases/download/latest/launch-linux_arm.tar.gz) | [![SBC](docs/logos/64x64/linux.png) SBC 64bit](https://github.com/gearboxworks/launch/releases/download/latest/launch-linux_arm64.tar.gz) |


---
## Usage

 Description:
 	launch - Gearbox gear launcher.

 Usage:

    launch [flags] [command] <gear name>


## Top level commands

Where `launch [command]` is one of the below commands.
 
### Basic commands

	version		- Guru		- launch - Self-manage this executable.
	search		- Manage	- Search for available Launch Container
	list		- Manage	- List a Launch Container
	run		- Execute	- Run default Launch Container command
	shell		- Execute	- Execute shell in Launch Container
 
### Manage

	search		- Manage	- Search for available Launch Container
	list		- Manage	- List a Launch Container
	install		- Manage	- Install a Launch Container
	uninstall	- Manage	- Uninstall a Launch Container
	update		- Manage	- Update a Launch Container
	clean		- Manage	- Completely uninstall a Launch Container
	log		- Manage	- Show logs of Launch Container
	start		- Manage	- Start a Launch Container
	stop		- Manage	- Stop a Launch Container

### Run

	run		- Execute	- Run default Launch Container command
	shell		- Execute	- Execute shell in Launch Container

### Help

 	help		- Help		- Help about any command.
	assist		- Help		- Show additional help
 	assist flags	- Help		- Show additional flags
	assist examples	- Help		- Show examples
	assist basic	- Help		- Show basic help
	assist advanced	- Help		- Show advanced help
	assist all	- Help		- Show all help
 

---
## Build commands

Where `launch build [command]` is one of the below commands.

### Building & creation
 
	create		- Build		- Build a Launch Container
	test		- Build		- Execute unit tests in Launch Container
	publish		- Build		- Publish a Launch Container
	clean		- Build		- Remove a Launch Container build
	export		- Build		- Save state of a Launch Container
	import		- Build		- Load a Launch Container 


---
## Guru commands

Where `launch [command]` is one of the below commands.

	scribe		- Guru		- scribe - The ultimate scripting toolkit.
	version		- Guru		- launch - Self-manage this executable. 
	completion	- Guru		- Generate BASH completion file
	
	
---
## Flags

### Provider flags

	--provider string	Set virtual provider (default "docker")
	--host string		Set virtual provider host.
	--port string		Set virtual provider port.

### Runtime flags

	-d, --debug		Debug mode.
	-n, --no-create		Don't create container.
	-q, --quiet		Silence all launch messages.
	-t, --temporary		Temporary container - remove after running command.
	--tmp string		Alternate TMP dir mount point. (default "none")
	-v, --version		Display version of launch

### Interactive shell flags

	-s, --status		Show shell status line.

### Mounting flags

	-m, --mount string	Mount arbitrary directory via SSHFS. (default "none")
	-p, --project string	Mount project directory. (default "none")


---
## Running launch

### GoLang

Building a simple GoLang application with multiple versions.

![GoLang](docs/launch-GoLang-Pause.gif)

## Imagemagick

Identify and resize a plain PNG file.

![Imagemagick](docs/launch-Imagemagick-Pause.gif)


### Further examples

There are many ways to call launch, either directly or indirectly.
Additionally, all host environment variables will be imported into the container seamlessly.
This allows a devloper to try multiple versions of software as though they were installed locally.

If a container is missing, it will be downloaded and created. Multiple versions can co-exist.

Install, create, and start the gearbox-base Gearbox container.

`launch install golang`

Create, and start the gearbox-base Gearbox container. Run a shell.

`launch shell golang`

Create, and start the gearbox-base Gearbox container with version alpine-3.4 and run a shell.

`launch shell golang:latest`

`launch shell golang:1.20 ls -l`

`launch shell golang:1.11.1 ps -eaf`


### Available commands
If gearbox-base is symlinked to `launch`, then the Gearbox container will be determined automatically and the default command will be run.
All available commands for a Gearbox container will be automatically symlinked upon installation.

`go`

Running gearbox-base Gearbox container default command. If a container has a default interactive command, arguments can be supplied without specifying that command.

`go -flag1 -flag2 variable`

`launch run golang -flag1 -flag2 variable`

Gearbox containers may have multiple executables that can be run. These will be automatically made available when you install a launch container.


### Remote Docker

You don't have to have Docker installed on your PC! As long as you set the DOCKER_HOST environment variable, then launch will use the remote Docker instance!


### Remote connection
ssh - All [Gearbox](https://github.com/gearboxworks/) containers have a running SSH daemon. So you can connect remotely.
To show what ports are exported to the host, use the following command.

`launch list ports`

