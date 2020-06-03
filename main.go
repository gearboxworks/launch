package main

import (
	"launch/cmd"
	"github.com/newclarity/scribeHelpers/ux"
	"os"
)

func init() {
	_ = ux.Open("Gearbox: ")
}

func main() {
	state := cmd.Execute()
	ux.Close()
	os.Exit(state.ExitCode)
}


/*

@TODO - Add '--sshfs-host' OR '--mount-host' flag, (string). When set to:
@TODO -		'' - Connect via direct SSH to host, (host.docker.internal).
@TODO -		'docker' - Connect via direct SSH to host, (host.docker.internal). (Enabled automatically when Docker server is local.)
@TODO -		'tunnel' - Connect via SSH tunnel, (host.docker.internal through client SSH). (Enabled automatically when Docker server is remote.)
@TODO -		'*' - Connect via direct SSH to host, (specified hostname).

@TODO - Add '--strict' flag to fail early, (fail as soon as something isn't right EG: don't create or install).

@TODO - Add '--bin-path' flag to allow alternative path links.

@TODO - Change configuration to point to 'launch.json' next to the launch binary.
@TODO -		Use '$HOME/.gearbox/launch.json' as a fallback.

@TODO - Add LABEL showing source repo for container/image.

@TODO - '-t' flag should probably start the container as an 'exec' rather than a full SSH.
@TODO -		Will still have to perform an SSHFS mount when using the '-m' flag.

@TODO - Consider changing 'list' to 'show' OR 'info'.
@TODO - Add several sub-commands to 'list' or have them as flags.
@TODO -		'--short'	OR 'list short'		- produce basic info.
@TODO -		'--long'	OR 'list long'		- produce longer info, (default).
@TODO -		'--inspect'	OR 'list inspect'	- inspect container/image.
@TODO -		'--logs'	OR 'list logs'		- show logs for container/image.
@TODO -		'--ports'	OR 'list ports'		- show defined ports within container/image.
@TODO -		'--repo'	OR 'list repo'		- show github src repo for container/image.

@TODO - Add several sub-commands to 'build' or have them as flags.
@TODO -		'--release'	OR 'build release'	- clean, build, test, push to dockerhub and clean again in one hit.
@TODO -		'--test'	OR 'build test'		- move 'launch test' to 'launch build test'.
@TODO -		'--release' OR 'build release'	- .
@TODO -		'--release' OR 'build release'	- .
@TODO -		'--release' OR 'build release'	- .
*/

/*
# Completed issues:
@DONE - When piping, don't print ANSI codes.

@DONE - Needs to be a better way to determine correct image/container validity. The current method works, but is a little fragile.

@DONE - Add 'wait for' code to get around delays in an image/container appearing within docker just after a creation.
@DONE - Seeing errors: "Docker client error: context deadline exceeded"

@DONE - Considering setting '-m' flag automatically when running as a symlink.

@DONE - Add several sub-commands to 'uninstall' or have them as flags.
@DONE -		'--image'	OR 'uninstall image'	- remove the image as well as the container.

@DONE - Refactor JTC and include helperDocker under the JTC framework.

@DONE - Add several sub-commands to 'version' or have them as flags.
@DONE -		'--update'	OR 'version update'	- update the current launch binary.
@DONE - Add '--update' flag to 'version' command OR add 'self-update' command.

@DONE - Add '--tmp-dir' flag to provide alternative mount point for TmpDir.
@DONE -		- The TmpDir should mount up automatically on every container, (not via SSHFS).
@DONE -		- This dir should be either $HOME/tmp/, /tmp/ or override flag above.
*/
