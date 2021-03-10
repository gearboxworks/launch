// Code generated by github.com/newclarity/PackageReflect DO NOT EDIT.

package cmd

import "reflect"

var Types = map[string]reflect.Type{
	"LaunchArgs": reflect.TypeOf((*LaunchArgs)(nil)).Elem(),
	"TypeLaunchArgs": reflect.TypeOf((*TypeLaunchArgs)(nil)).Elem(),
}

var Functions = map[string]reflect.Value{
	"CheckReturns": reflect.ValueOf(CheckReturns),
	"DeterminePath": reflect.ValueOf(DeterminePath),
	"Execute": reflect.ValueOf(Execute),
	"GetGearboxDir": reflect.ValueOf(GetGearboxDir),
	"New": reflect.ValueOf(New),
	"SetCmd": reflect.ValueOf(SetCmd),
}

var Variables = map[string]reflect.Value{
	"Cmd": reflect.ValueOf(&Cmd),
	"CmdScribe": reflect.ValueOf(&CmdScribe),
	"CmdSelfUpdate": reflect.ValueOf(&CmdSelfUpdate),
	"CobraHelp": reflect.ValueOf(&CobraHelp),
	"ConfigFile": reflect.ValueOf(&ConfigFile),
}

var Consts = map[string]reflect.Value{
	"DefaultJsonFile": reflect.ValueOf(DefaultJsonFile),
	"DefaultJsonString": reflect.ValueOf(DefaultJsonString),
	"DefaultTemplateFile": reflect.ValueOf(DefaultTemplateFile),
	"DefaultTemplateString": reflect.ValueOf(DefaultTemplateString),
}

