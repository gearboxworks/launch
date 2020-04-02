package dockerClient
//
//import (
//	"github.com/docker/docker/client"
//)
//
////type Response struct {
//	//Docker string
//	//State ux.State
//	//Ok bool
////}
//
//type Response string
//
//func (me *Response) IsRunning() bool {
//	var ok bool
//	if *me == "running" {
//		ok = true
//	}
//	return ok
//}
//
//func (me *Response) IsPaused() bool {
//	var ok bool
//	if *me == "paused" {
//		ok = true
//	}
//	return ok
//}
//
//func (me *Response) IsCreated() bool {
//	var ok bool
//	if *me == "created" {
//		ok = true
//	}
//	return ok
//}
//
//func (me *Response) IsRestarting() bool {
//	var ok bool
//	if *me == "restarting" {
//		ok = true
//	}
//	return ok
//}
//
//func (me *Response) IsRemoving() bool {
//	var ok bool
//	if *me == "removing" {
//		ok = true
//	}
//	return ok
//}
//
//func (me *Response) IsExited() bool {
//	var ok bool
//	if *me == "exited" {
//		ok = true
//	}
//	return ok
//}
//
//func (me *Response) IsDead() bool {
//	var ok bool
//	if *me == "dead" {
//		ok = true
//	}
//	return ok
//}
//
//// "created", "running", "paused", "restarting", "removing", "exited", or "dead"
//func (me *Response) SetState(s string) {
//	*me = Response(s)
//}
//
//func (me *Response) IsError(c *client.Client) bool {
//	var ok bool
//
//	// if me.Error != nil {
//	//
//	// 	fmt.Printf("Gearbox error: %s\n", state.Error)
//	// }
//
//	return ok
//}