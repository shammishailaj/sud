package test

import (
	"time"

	"github.com/crazyprograms/callpull"
	"github.com/shammishailaj/sud/core"
	"github.com/shammishailaj/sud/corebase"
)

func testStd(cr *core.Core, Name string, Param map[string]interface{}, timeOutWait time.Duration) (callpull.Result, error) {
	return callpull.Result{Result: "Test Ok"}, nil
}
func init() {
	initRecordConfiguration()
	core.AddStdCall("TestStd", testStd)
}
func initRecordConfiguration() {
	conf := core.NewConfiguration()
	conf.AddType("Test", "Test", true, true, true, "Тест")
	conf.AddPole("Test", "Test", "Test.Test1", "StringValue", corebase.NULL, "Index", &core.PoleCheckerStringValue{}, true, true, "Поле для тестирования")
	conf.AddCall("Test", "TestStd", "std", true, false, "")
	conf.AddCall("Test", "TestAsync", "async", true, true, "")
	core.InitAddBaseConfiguration("Test", conf)
}
