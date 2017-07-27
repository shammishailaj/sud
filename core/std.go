package core

import (
	"time"

	"github.com/crazyprograms/sud/callpull"
	"github.com/crazyprograms/sud/corebase"
)

func stdConfigurationList(cr *Core, Name string, Param map[string]interface{}, timeOutWait time.Duration, Access corebase.IAccess) (callpull.Result, error) {
	configurations := cr.GetConfiguration()
	list := make([]interface{}, len(configurations))
	i := 0
	for name := range configurations {
		list[i] = name
		i++
	}
	return callpull.Result{Error: nil, Result: list}, nil
}

func InitStdModule(c *Core) bool {
	conf := NewConfiguration([]string{"std"})
	conf.AddCall(CallInfo{ConfigurationName: "std", Name: "std.Configuration.List", PullName: "std", AccessCall: "Configuration", AccessListen: "", Title: "Список конфигураций"})
	if AddStdCall("std.Configuration.List", stdConfigurationList) != nil {
		return false
	}
	return c.AddBaseConfiguration("std", conf)
}
