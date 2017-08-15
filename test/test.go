package test

import (
	"time"

	"github.com/crazyprograms/sud/callpull"
	"github.com/crazyprograms/sud/core"
	"github.com/crazyprograms/sud/corebase"
)

func testStd(cr *core.Core, Name string, Param map[string]interface{}, timeOutWait time.Duration, Access corebase.IAccess) (callpull.Result, error) {
	return callpull.Result{Result: "Test Ok"}, nil
}

func InitModule(c *core.Core) error {
	return core.AddStdCall("TestStd", testStd)
}
