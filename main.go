package main

import (
	"launch/cmd"
	"launch/only"
	"launch/ux"
	"os"
)

func init() {
	_ = ux.Open()
}

func main() {
	var state ux.State

	for range only.Once {
		state = cmd.Execute()
		if state.IsError() {
			break
		}
		state = cmd.GetState()

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

	exit := 0
	if state.IsError() {
		exit = 1
	}
	if state.IsWarning() {
		exit = 2
	}
	state.Print()
	ux.Close()
	os.Exit(exit)
}
