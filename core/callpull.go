package core

import "time"
import "errors"

type FCall func(Name string, Param map[string]interface{}, timeOutWait time.Duration) (interface{}, error)
type LocalCallPull struct {
	calls map[string]FCall
}

func (cp *LocalCallPull) Call(Name string, Param map[string]interface{}, timeOutWait time.Duration) (interface{}, error) {
	if _, ok := cp.calls[Name]; !ok {
		return nil, errors.New("call " + Name + " not exists")
	}
	return cp.calls[Name](Name, Param, timeOutWait)
}
func (cp *LocalCallPull) AddCall(Name string, Call FCall) error {
	if _, ok := cp.calls[Name]; ok {
		return errors.New("call " + Name + " already exists")
	}
	cp.calls[Name] = Call
	return nil
}

var StdCallPull *LocalCallPull = &LocalCallPull{calls: map[string]FCall{}}
