package storage

type CallInfo struct {
	ConfigurationName string
	Title             string
	Call              bool
	Listen            bool
}

func (ci *CallInfo) GetConfigurationName() string { return ci.ConfigurationName }
func (ci *CallInfo) GetTitle() string             { return ci.Title }
func (ci *CallInfo) GetCall() bool                { return ci.Call }
func (ci *CallInfo) GetListen() bool              { return ci.Listen }

var _ci ICallInfo = (ICallInfo)(&CallInfo{})
