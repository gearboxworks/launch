package ux

import (
	"errors"
	"fmt"
)

type State struct {
	Error error
	Warning error
	Ok error
	String string
}

func (me *State) Print() {
	switch {
		case me.Error != nil:
			PrintfError("%s", me.Error)
		case me.Warning != nil:
			PrintfWarning("%s", me.Warning)
		case me.Ok != nil:
			PrintfOk("%s", me.Ok)
	}
}


func (me *State) IsError() bool {
	var ok bool

	if me.Error != nil {
		ok = true
	}

	return ok
}

func (me *State) IsWarning() bool {
	var ok bool

	if me.Warning != nil {
		ok = true
	}

	return ok
}

func (me *State) IsOk() bool {
	var ok bool

	if me.Ok != nil {
		ok = true
	}

	return ok
}


func (me *State) SetError(format string, args ...interface{}) {
	me.Ok = nil
	me.Warning = nil
	me.Error = errors.New(fmt.Sprintf(format, args...))
}

func (me *State) SetWarning(format string, args ...interface{}) {
	me.Ok = nil
	me.Warning = errors.New(fmt.Sprintf(format, args...))
	me.Error = nil
}

func (me *State) SetOk(format string, args ...interface{}) {
	me.Ok = errors.New(fmt.Sprintf(format, args...))
	me.Warning = nil
	me.Error = nil
}

func (me *State) ClearError() {
	me.Error = nil
}

func (me *State) ClearAll() {
	me.Ok = nil
	me.Warning = nil
	me.Error = nil
}


func (me *State) IsRunning() bool {
	var ok bool
	if me.String == StateRunning {
		ok = true
	}
	return ok
}

func (me *State) IsPaused() bool {
	var ok bool
	if me.String == StatePaused {
		ok = true
	}
	return ok
}

func (me *State) IsCreated() bool {
	var ok bool
	if me.String == StateCreated {
		ok = true
	}
	return ok
}

func (me *State) IsRestarting() bool {
	var ok bool
	if me.String == StateRestarting {
		ok = true
	}
	return ok
}

func (me *State) IsRemoving() bool {
	var ok bool
	if me.String == StateRemoving {
		ok = true
	}
	return ok
}

func (me *State) IsExited() bool {
	var ok bool
	if me.String == StateExited {
		ok = true
	}
	return ok
}

func (me *State) IsDead() bool {
	var ok bool
	if me.String == StateDead {
		ok = true
	}
	return ok
}

// "created", "running", "paused", "restarting", "removing", "exited", or "dead"
func (me *State) SetString(s string) {
	me.String = s
}
