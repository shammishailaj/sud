package core

import (
	"time"

	"github.com/crazyprograms/sud/callpull"
	"github.com/crazyprograms/sud/corebase"
)

type FCall func(core *Core, Name string, Param map[string]interface{}, timeOutWait time.Duration, Access corebase.IAccess) (callpull.Result, error)
type LocalCallPull struct {
	core  *Core
	calls map[string]FCall
}

func (cp *LocalCallPull) Call(Name string, Param map[string]interface{}, timeOutWait time.Duration, Access corebase.IAccess) (callpull.Result, error) {
	if _, ok := cp.calls[Name]; !ok {
		return callpull.Result{Result: nil}, &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "Call", Name: Name}
	}
	return cp.calls[Name](cp.core, Name, Param, timeOutWait, Access)
}

var strCalls = map[string]FCall{}

func GetStdPull(core *Core) ICallPull {
	return &LocalCallPull{calls: strCalls, core: core}
}
func AddStdCall(Name string, Call FCall) error {
	if _, ok := strCalls[Name]; ok {
		return &corebase.Error{ErrorType: corebase.ErrorTypeAlreadyExists, Action: "AddStdCall", Name: Name}
	}
	strCalls[Name] = Call
	return nil

}
