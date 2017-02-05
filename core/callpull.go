package core

import "time"
import "errors"

type FCall func(core *Core, Name string, Param map[string]interface{}, timeOutWait time.Duration) (interface{}, error)
type LocalCallPull struct {
	core  *Core
	calls map[string]FCall
}

func (cp *LocalCallPull) Call(Name string, Param map[string]interface{}, timeOutWait time.Duration) (interface{}, error) {
	if _, ok := cp.calls[Name]; !ok {
		return nil, errors.New("call " + Name + " not exists")
	}
	return cp.calls[Name](cp.core, Name, Param, timeOutWait)
}

var strCalls = map[string]FCall{}

func GetStdPull(core *Core) ICallPull {
	return &LocalCallPull{calls: strCalls, core: core}
}
func AddStdCall(Name string, Call FCall) error {
	if _, ok := strCalls[Name]; ok {
		return errors.New("call " + Name + " already exists")
	}
	strCalls[Name] = Call
	return nil

}
