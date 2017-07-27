package core

import "github.com/crazyprograms/sud/corebase"

type CallInfo struct {
	ConfigurationName string
	Title             string
	Name              string
	PullName          string
	AccessCall        string
	AccessListen      string
}

func (ci *CallInfo) GetConfigurationName() string { return ci.ConfigurationName }
func (ci *CallInfo) GetTitle() string             { return ci.Title }
func (ci *CallInfo) GetName() string              { return ci.Name }
func (ci *CallInfo) GetPullName() string          { return ci.PullName }
func (ci *CallInfo) GetAccessCall() string        { return ci.AccessCall }
func (ci *CallInfo) GetAccessListen() string      { return ci.AccessListen }

var _ci corebase.ICallInfo = (corebase.ICallInfo)(&CallInfo{})
