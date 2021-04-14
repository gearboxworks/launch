## About `launch`

`launch` is a tool specifically designed to simplify running Docker containers.

Launch provides three (3) important functional areas, without any Docker container learning curve:

- Allows control over Gearbox Docker containers: stop, start, create, remove.
- Build, update, modify and release Docker images.
- Acts as a proxy for interactive commands within a Gearbox Docker container.

It also provides a functional SSH daemon for connecting remotely as well as a standard set of common tools and utilities.


### Download launch
`launch` is currently in beta testing and is included along with all Gearbox Docker repos.
Once out of beta, it will be included within the Gearbox installation package.

For now, simply download the standalone `launch` binary for your O/S.
- [Mac OSX 64bit](https://github.com/gearboxworks/docker-template/raw/master/bin/Darwin/launch)
- [Linux 64bit](https://github.com/gearboxworks/docker-template/raw/master/bin/Linux/launch)
- [Windows 64bit](https://github.com/gearboxworks/docker-template/raw/master/bin/Windows/launch)


## Usage

 Description:
 	launch - Gearbox gear launcher.

 Usage:

    launch [flags] [command] <gear name>


## Top level commands

Where [command] is one of the below commands.

Use launch help [command] for more information about a command.

### Help

 	help		- Help about any command

### Building & creation
 
 	build		- Build a Gearbox gear
 	export		- Save state of a Gearbox gear
 	import		- Load a Gearbox gear
 	publish		- Publish a Gearbox gear
 	test		- Execute unit tests in Gearbox gear

### Install & uninstall

 	install		- Install a Gearbox gear
 	uninstall	- Uninstall a Gearbox gear
 	reinstall	- Update a Gearbox gear
 	list		- List a Gearbox gear

### Start & stop

 	start		- Start a Gearbox gear
 	stop		- Stop a Gearbox gear

### Running

 	run		- Run default Gearbox gear command
 	shell		- Execute shell in Gearbox gear


## Flags

### General flags

   -b, --completion        Generate BASH completion script.
   -v, --version           Display version of launch

### Provider flags

       --provider string   Set virtual provider (default "docker")
       --host string       Set virtual provider host.
       --port string       Set virtual provider port.

### Runtime flags

       --config string     Config file. (default "$HOME/.gearbox/launch.json")
   -d, --debug             Debug mode.
   -q, --quiet             Silence all Gearbox messsages.
   -n, --no-create         Don't create container.
   -t, --temporary         Temporary container - remove after running command.

### Interactive shell flags

   -s, --status            Show shell status line.

### Mounting flags

   -m, --mount string      Mount arbitrary directory via SSHFS. (default "none")
   -p, --project string    Mount project directory. (default "none")


## Running launch

### Examples

There are many ways to call launch, either directly or indirectly.
Additionally, all host environment variables will be imported into the container seamlessly.
This allows a devloper to try multiple versions of software as though they were installed locally.

If a container is missing, it will be downloaded and created. Multiple versions can co-exist.

Install, create, and start the gearbox-base Gearbox container.

`./launch install gearbox-base`

Create, and start the gearbox-base Gearbox container. Run a shell.

`./launch shell gearbox-base`

Create, and start the gearbox-base Gearbox container with version alpine-3.4 and run a shell.

`./launch shell gearbox-base:alpine-3.4`

`./launch shell gearbox-base:alpine-3.4 ls -l`

`./launch shell gearbox-base:alpine-3.4 ps -eaf`


### Available commands
If gearbox-base is symlinked to `launch`, then the Gearbox container will be determined automatically and the default command will be run.
All available commands for a Gearbox container will be automatically symlinked upon installation.

`./gearbox-base`

Running gearbox-base Gearbox container default command. If a container has a default interactive command, arguments can be supplied without specifying that command.

`./gearbox-base -flag1 -flag2 variable`

`./launch gearbox-base:alpine-3.4 -flag1 -flag2 variable`

Gearbox containers may have multiple executables that can be run. The gearbox-base Gearbox container has the following available commands:
- The default command will execute `` within the container.


### Remote connection
ssh - All [Gearbox](https://github.com/gearboxworks/) containers have a running SSH daemon. So you can connect remotely.
To show what ports are exported to the host, use the following command.

`./launch list gearbox-base`

