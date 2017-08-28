package core

import (
	"time"

	"github.com/crazyprograms/sud/callpull"
	"github.com/crazyprograms/sud/corebase"
)

func stdConfigurationList(cr *Core, Name string, Param map[string]interface{}, timeOutWait time.Duration, Access corebase.IAccess) (callpull.Result, error) {
	configurations := cr.Configurator().GetConfigurations()
	listConf := make(map[string]interface{})
	for name := range configurations {
		var err error
		var confInfo map[string]interface{}
		if confInfo, err = configurations[name].Save(); err != nil {
			return callpull.Result{Error: err, Result: nil}, nil
		}
		listConf[name] = confInfo
	}
	return callpull.Result{Error: nil, Result: listConf}, nil
}

func InitStdModule(c *Core) bool {
	if AddStdCall("std.Configuration.Get", stdConfigurationList) != nil {
		return false
	}
	return true
}
