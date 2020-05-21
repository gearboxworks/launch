package main

import (
	"launch/cmd"
	"launch/only"
	"launch/ux"
	"os"
)

func init() {
	_ = ux.Open("Gearbox: ")
}

func main() {
	state := ux.NewState(false)

	for range only.Once {
		state = cmd.Execute()
		if state.IsError() {
			break
		}
		//state = cmd.GetState()

		// @TODO - testing.
		//ux.PrintfOk("OK - it works\n")
		//ux.PrintfWarning("WARNING - it may or may not have worked... sort of\n")
		//ux.PrintfError("ERROR - oops, no beer in the fridge\n")
		//err = errors.New("what? no beer!")
		//ux.PrintError(err)
		//_ = ux.Draw2()
		//ux.Draw3()
		//_ = ux.Draw4()
		//_ = ux.Draw5()
	}

	if state.IsError() {
		state.SetExitCode(1)
	}

	//if state.IsWarning() {
	//	exit = 2
	//}

	if state.IsNotOk() {
		state.PrintResponse()
	}

	ux.Close()
	os.Exit(state.ExitCode)
}

// @TODO - Add '--strict' flag to fail early, (fail as soon as something isn't right EG: don't create or install).

// @TODO - Add 'wait for' code to get around delays in an image/container appearing within docker just after a creation.

// @TODO - Add '--update' flag to 'version' command OR add 'self-update' command.

// @TODO - Considering setting '-m' flag automatically when running as a symlink.

// @TODO - When piping, don't print ANSI codes.

// @TODO - Add LABEL showing source repo for container/image.

// @TODO - '-t' flag should probably start the container as an 'exec' rather than a full SSH.
// @TODO -		- Will still have to perform an SSHFS mount when using the '-m' flag.

// @TODO - Consider changing 'list' to 'show' OR 'info'.
// @TODO - Add several sub-commands to 'list' or have them as flags.
// @TODO -		'--short'	OR 'list short'		- produce basic info.
// @TODO -		'--long'	OR 'list long'		- produce longer info, (default).
// @TODO -		'--inspect'	OR 'list inspect'	- inspect container/image.
// @TODO -		'--logs'	OR 'list logs'		- show logs for container/image.
// @TODO -		'--ports'	OR 'list ports'		- show defined ports within container/image.
// @TODO -		'--repo'	OR 'list repo'		- show github src repo for container/image.

// @TODO - Add several sub-commands to 'build' or have them as flags.
// @TODO -		'--release'	OR 'build release'	- clean, build, test, push to dockerhub and clean again in one hit.
// @TODO -		'--test'	OR 'build test'		- move 'launch test' to 'launch build test'.
// @TODO -		'--release' OR 'build release'	- .
// @TODO -		'--release' OR 'build release'	- .
// @TODO -		'--release' OR 'build release'	- .

// @TODO - Add several sub-commands to 'uninstall' or have them as flags.
// @TODO -		'--image'	OR 'uninstall image'	- remove the image as well as the container.

// @TODO - Add several sub-commands to 'version' or have them as flags.
// @TODO -		'--update'	OR 'version update'	- update the current launch binary.

// @TODO - Needs to be a better way to determine correct image/container validity. The current method works, but is a little fragile.
