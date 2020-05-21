package ux

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)


type StateGetter interface {
	Print()
	IsError() bool
	IsWarning() bool
	IsOk() bool
	SetError(format string, args ...interface{})
	SetWarning(format string, args ...interface{})
	SetOk(format string, args ...interface{})
	ClearError()
	ClearAll()
	IsRunning() bool
	IsPaused() bool
	IsCreated() bool
	IsRestarting() bool
	IsRemoving() bool
	IsExited() bool
	IsDead() bool
	SetString(s string)
}

type State struct {
	prefix      string
	prefixArray []string
	_Package    string
	_Function   string

	_Fatal      error
	_Error      error
	_Warning    error
	_Ok         error
	_Debug      error
	ExitCode    int
	debug       RuntimeDebug

	RunState    string

	Output      string
	_Separator  string
	OutputArray []string
	Response    interface{}
}

type RuntimeDebug struct {
	Enabled  bool
	File     string
	Line     int
	Function string
}


const DefaultSeparator = "\n"


func NewState(debugMode bool) *State {
	me := State{}
	me.Clear()
	me.debug.Enabled = debugMode

	return &me
}

func (p *State) EnsureNotNil() *State {
	for range OnlyOnce {
		if p == nil {
			p = NewState(false)
		}
		p.Clear()
	}
	return p
}

func EnsureStateNotNil(p *State) *State {
	for range OnlyOnce {
		if p == nil {
			p = NewState(false)
		}
		p.Clear()
	}
	return p
}

func IfNilReturnError(ref interface{}) *State {
	if ref == nil {
		s := NewState(true)
		s._Fatal = errors.New("SW ERROR")
		s.ExitCode = 255
		return s
	}

	state := SearchStructureForUxState(ref)
	if state == nil {
		state = NewState(false)
	}
	return state
	//return ref.(*State)
}

// Search a given structure for the State object and return it's pointer.
func SearchStructureForUxState(m interface{}) *State {
	var state *State

	s := reflect.ValueOf(m).Elem()
	typeOfT := s.Type()
	//fmt.Println("t=", m)
	for i := 0; i < s.NumField(); i++ {
		if typeOfT.Field(i).Name == "State" {
			state = s.Field(i).Interface().(*State)
			break
		}
	}

	return state
}

func (p *State) Clear() {
	if p != nil {
		p._Debug = nil
		p._Fatal = nil
		p._Error = nil
		p._Warning = nil
		p._Ok = errors.New("")
		p.ExitCode = 0

		p.Output = ""
		p._Separator = DefaultSeparator
		p.OutputArray = []string{}
		p.Response = nil
	} else {
		panic(p)
	}
}


func (p *State) GetPrefix() string {
	return p.prefix
}
func (p *State) GetPackage() string {
	return p._Package
}
func (p *State) GetFunction() string {
	return p._Function
}
func (p *State) SetPackage(s string) {
	if s == "" {
		// Discover package name.
		//pc, file, no, ok := runtime.Caller(1)
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			//s = file + ":" + string(no)
			details := runtime.FuncForPC(pc)
			s = filepath.Base(details.Name())
			sa := strings.Split(s, ".")
			if len(sa) > 0 {
				s = sa[0]
			}
		}
	}

	p._Package = s
	if p._Function == "" {
		p.prefix = p._Package
	} else {
		p.prefix = p._Package + "." + p._Function + "()"
		p.prefixArray = append(p.prefixArray, p.prefix)
	}
}
func (p *State) SetFunction(s string) {
	if s == "" {
		// Discover function name.
		//pc, file, no, ok := runtime.Caller(1)
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			//s = file + ":" + string(no)
			details := runtime.FuncForPC(pc)
			foo := details.Name()
			s = filepath.Base(foo)
			sa := strings.Split(s, ".")
			switch {
				case len(sa) > 2:
					s = sa[2]
				case len(sa) > 1:
					s = sa[1]
				case len(sa) > 0:
					s = sa[0]
			}
		}
	}

	p._Function = s
	if p._Package == "" {
		p.prefix = p._Function + "()"
	} else {
		p.prefix = p._Package + "." + p._Function + "()"
	}

	p.prefixArray = append(p.prefixArray, p.prefix)
}
func (p *State) SetFunctionCaller() {
	var s string
	// Discover function name.
	pc, _, _, ok := runtime.Caller(2)
	if ok {
		//s = file + ":" + string(no)
		details := runtime.FuncForPC(pc)
		s = filepath.Base(details.Name())
		sa := strings.Split(s, ".")
		if len(sa) > 0 {
			s = sa[1]
		}
	}

	p.SetFunction(s)
}


func (p *State) GetState() *bool {
	var b bool
	return &b
}
func (s *State) SetState(p *State) {
	if s == nil {
		s = NewState(true)
		s._Fatal = errors.New("SW ERROR")
		s.ExitCode = 255
		return
	}
	s._Error =      p._Error
	s._Warning =    p._Warning
	s._Ok =         p._Ok
	s._Debug =      p._Debug
	s.ExitCode =    p.ExitCode
	s.Output =      p.Output
	s.OutputArray = p.OutputArray
	s.Response =    p.Response
	s.RunState =    p.RunState
}


func (p *State) Sprint() string {
	var ret string

	e := ""
	if p.ExitCode != 0 {
		e = fmt.Sprintf("Exit(%d) - ", p.ExitCode)
	}

	pa := ""
	if len(p.prefixArray) > 0 {
		pa = fmt.Sprintf("[%s] - ", p.prefixArray[0])
	}

	ou := ""
	if p.Output != "" {
		ou = "\n" + SprintfOk("%s ", p.Output)
	}

	switch {
		case p._Error != nil:
			ret = SprintfError("ERROR: %s%s%s%s", pa, e, p._Error, ou)

		case p._Warning != nil:
			ret = SprintfWarning("WARNING: %s%s%s%s", pa, e, p._Warning, ou)

		case p._Ok != nil:
			ret = SprintfOk("%s%s", p._Ok, ou)

		case p.debug.Enabled:
			if p._Debug != nil {
				ret = SprintfDebug("%s%s", p._Debug, ou)
			}
	}

	return ret
}
func (p *State) SprintResponse() string {
	return p.Sprint()
}
func (p *State) PrintResponse() {
	_, _ = fmt.Fprintf(os.Stderr, p.Sprint() + "\n")
}
func (p *State) SprintError() string {
	var ret string

	for range OnlyOnce {
		if p._Ok != nil {
			// If we have an OK response.
			break
		}

		ret = p.Sprint()
	}

	return ret
}


func (p *State) IsError() bool {
	var ok bool

	if p == nil {
		//fmt.Printf("DUH\n")
		ok = true
	} else if p._Error != nil {
		ok = true
	}

	return ok
}

func (p *State) IsWarning() bool {
	var ok bool

	if p._Warning != nil {
		ok = true
	}

	return ok
}

func (p *State) IsOk() bool {
	var ok bool

	if p._Ok != nil {
		ok = true
	}

	return ok
}
func (p *State) IsNotOk() bool {
	ok := true

	for range OnlyOnce {
		if p._Warning != nil {
			break
		}
		if p._Error != nil {
			break
		}
		ok = false
	}

	return ok
}

func (p *State) SetExitCode(e int) {
	if p == nil {
		return
	}
	p.ExitCode = e
}
func (p *State) GetExitCode() int {
	return p.ExitCode
}


func (p *State) SetError(error ...interface{}) {
	for range OnlyOnce {
		if p == nil {
			panic(p)
			break
		}
		p.debug.fetchRuntimeDebug(2)

		p._Ok = nil
		p._Warning = nil

		if len(error) == 0 {
			p._Error = errors.New("ERROR")
			break
		}

		if error[0] == nil {
			p._Error = nil
			break
		}

		debugPrefix := ""
		if p.debug.Enabled {
			debugPrefix = SprintfCyan("%s:%d [%s] - ", p.debug.File, p.debug.Line, p.debug.Function)
		}
		p._Error = errors.New(debugPrefix + _Sprintf(error...))
		if p.debug.Enabled {
			p.PrintResponse()
		}
	}
}
func (p *State) GetError() error {
	return p._Error
}


func (p *State) SetWarning(warning ...interface{}) {
	for range OnlyOnce {
		if p == nil {
			panic(p)
			break
		}
		p.debug.fetchRuntimeDebug(2)

		p._Ok = nil
		p._Error = nil

		if len(warning) == 0 {
			p._Warning = errors.New("WARNING")
			break
		}

		if warning[0] == nil {
			p._Warning = nil
			break
		}

		debugPrefix := ""
		if p.debug.Enabled {
			debugPrefix = SprintfCyan("%s:%d [%s] - ", p.debug.File, p.debug.Line, p.debug.Function)
		}
		p._Warning = errors.New(debugPrefix + _Sprintf(warning...))
		if p.debug.Enabled {
			p.PrintResponse()
		}
	}
}
func (p *State) GetWarning() error {
	return p._Warning
}


func (p *State) SetOk(msg ...interface{}) {
	for range OnlyOnce {
		if p == nil {
			panic(p)
			break
		}
		p.debug.fetchRuntimeDebug(2)

		p._Error = nil
		p._Warning = nil
		p.ExitCode = 0

		if len(msg) == 0 {
			p._Ok = errors.New("")
			break
		}

		if msg[0] == nil {
			p._Ok = errors.New("")
			break
		}

		debugPrefix := ""
		if p.debug.Enabled {
			debugPrefix = SprintfCyan("%s:%d [%s] - ", p.debug.File, p.debug.Line, p.debug.Function)
		}
		p._Ok = errors.New(debugPrefix + _Sprintf(msg...))
		if p.debug.Enabled {
			p.PrintResponse()
		}
	}
}
func (p *State) GetOk() error {
	return p._Ok
}


func (p *State) ClearError() {
	p._Error = nil
}


func (p *State) IsRunning() bool {
	var ok bool
	if p.RunState == StateRunning {
		ok = true
	}
	return ok
}

func (p *State) IsPaused() bool {
	var ok bool
	if p.RunState == StatePaused {
		ok = true
	}
	return ok
}

func (p *State) IsCreated() bool {
	var ok bool
	if p.RunState == StateCreated {
		ok = true
	}
	return ok
}

func (p *State) IsRestarting() bool {
	var ok bool
	if p.RunState == StateRestarting {
		ok = true
	}
	return ok
}

func (p *State) IsRemoving() bool {
	var ok bool
	if p.RunState == StateRemoving {
		ok = true
	}
	return ok
}

func (p *State) IsExited() bool {
	var ok bool
	if p.RunState == StateExited {
		ok = true
	}
	return ok
}

func (p *State) IsDead() bool {
	var ok bool
	if p.RunState == StateDead {
		ok = true
	}
	return ok
}


func (p *State) ExitOnNotOk() string {
	if p.IsNotOk() {
		_, _ = fmt.Fprintf(os.Stderr, p.Sprint() + "\n")
		os.Exit(p.ExitCode)
	}
	return ""
}


func (p *State) ExitOnError() string {
	if p.IsWarning() {
		_, _ = fmt.Fprintf(os.Stderr, p.Sprint() + "\n")
	}
	if p.IsError() {
		_, _ = fmt.Fprintf(os.Stderr, p.Sprint() + "\n")
		os.Exit(p.ExitCode)
	}
	return ""
}


func (p *State) ExitOnWarning() string {
	if p.IsWarning() {
		_, _ = fmt.Fprintf(os.Stderr, p.Sprint() + "\n")
		os.Exit(p.ExitCode)
	}
	return ""
}


func (p *State) Exit(e int) string {
	p.ExitCode = e
	_, _ = fmt.Fprintf(os.Stdout, p.Sprint())
	os.Exit(p.ExitCode)
	return ""
}


func Exit(e int64, msg ...interface{}) string {
	ret := _Sprintf(msg...)
	if e == 0 {
		_, _ = fmt.Fprintf(os.Stdout, SprintfOk(ret))
	} else {
		_, _ = fmt.Fprintf(os.Stderr, SprintfError(ret))
	}
	os.Exit(int(e))
	return ""	// Will never get here.
}


func _Sprintf(msg ...interface{}) string {
	var ret string

	for range OnlyOnce {
		if len(msg) == 0 {
			break
		}

		value := reflect.ValueOf(msg[0])
		switch value.Kind() {
			case reflect.String:
				if len(msg) == 1 {
					ret = fmt.Sprintf(msg[0].(string))
				} else {
					ret = fmt.Sprintf(msg[0].(string), msg[1:]...)
				}

			default:
				if len(msg) == 1 {
					ret = fmt.Sprintf("%v", msg)
				} else {
					var es string
					for _, e := range msg {
						es += fmt.Sprintf("%v ", e)
					}
					es = strings.TrimSuffix(es, " ")
					ret = es
				}
		}

		//ret = fmt.Sprintf(msg[0].(string), msg[1:]...)
	}

	return ret
}


func (p *RuntimeDebug) fetchRuntimeDebug(level int) {
	for range OnlyOnce {
		if p == nil {
			break
		}
		if level == 0 {
			level = 1
		}

		// Discover package name.
		var ok bool
		var pc uintptr
		pc, p.File, p.Line, ok = runtime.Caller(level)
		if ok {
			details := runtime.FuncForPC(pc)
			p.Function = details.Name()
			//f, l := details.FileLine(pc)
			//fmt.Printf("%s:%d - %s:%d\n",
			//	p.File,
			//	p.Line,
			//	f,
			//	l,
			//	)
		}
		//fmt.Printf("DEBUG => %s:%d [%s]\n", p.File, p.Line, p.Function)
	}
}

func (p *State) DebugEnable() {
	for range OnlyOnce {
		if p == nil {
			break
		}
		p.debug.Enabled = true
	}
}
func (p *State) DebugDisable() {
	for range OnlyOnce {
		if p == nil {
			break
		}
		p.debug.Enabled = false
	}
}
func (p *State) DebugSet(d bool) {
	for range OnlyOnce {
		if p == nil {
			break
		}
		p.debug.Enabled = d
	}
}

func (p *State) SetDebug(msg ...interface{}) {
	for range OnlyOnce {
		if p == nil {
			break
		}
		p.debug.fetchRuntimeDebug(2)

		if len(msg) == 0 {
			p._Debug = errors.New("DEBUG")
			break
		}

		if msg[0] == nil {
			p._Debug = errors.New("DEBUG")
			break
		}

		debugPrefix := ""
		if p.debug.Enabled {
			debugPrefix = SprintfCyan("%s:%d [%s] - ", p.debug.File, p.debug.Line, p.debug.Function)
		}
		p._Debug = errors.New(debugPrefix + _Sprintf(msg...))
		if p.debug.Enabled {
			p.PrintResponse()
		}
	}
}
func (p *State) GetDebug() error {
	return p._Debug
}
