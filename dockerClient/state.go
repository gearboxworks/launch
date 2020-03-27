package dockerClient

import "github.com/docker/docker/client"

type State struct {
	Docker string
	Error error
	Ok bool
}
func (me *State) IsRunning() bool {
	var ok bool
	if me.Docker == "running" {
		ok = true
	}
	return ok
}

func (me *State) IsPaused() bool {
	var ok bool
	if me.Docker == "paused" {
		ok = true
	}
	return ok
}

func (me *State) IsCreated() bool {
	var ok bool
	if me.Docker == "created" {
		ok = true
	}
	return ok
}

func (me *State) IsRestarting() bool {
	var ok bool
	if me.Docker == "restarting" {
		ok = true
	}
	return ok
}

func (me *State) IsRemoving() bool {
	var ok bool
	if me.Docker == "removing" {
		ok = true
	}
	return ok
}

func (me *State) IsExited() bool {
	var ok bool
	if me.Docker == "exited" {
		ok = true
	}
	return ok
}

func (me *State) IsDead() bool {
	var ok bool
	if me.Docker == "dead" {
		ok = true
	}
	return ok
}

// "created", "running", "paused", "restarting", "removing", "exited", or "dead"
func (me *State) SetState(s string) {
	me.Docker = s
}

func (me *State) IsError(c *client.Client) bool {
	var ok bool

	// if me.Error != nil {
	//
	// 	fmt.Printf("Gearbox error: %s\n", state.Error)
	// }

	return ok
}