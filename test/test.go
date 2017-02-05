package test

import (
	"time"

	"github.com/crazyprograms/sud/core"
)

func testStd(cr *core.Core, Name string, Param map[string]interface{}, timeOutWait time.Duration) (interface{}, error) {
	return "Test Ok", nil
}
func init() {
	initDocumentConfiguration()
	core.AddStdCall("TestStd", testStd)
}
func initDocumentConfiguration() {
	conf := core.NewConfiguration()
	conf.AddType("Test", "Test", true, true, true, "Тест")
	conf.AddPole("Test", "Test", "Test.Test1", "StringValue", core.NewObject(nil), "Index", &core.PoleCheckerStringValue{}, true, true, "Поле для тестирования")
	conf.AddCall("Test", "TestStd", "std", true, false, "")
	conf.AddCall("Test", "TestAsync", "async", true, true, "")
	core.InitAddBaseConfiguration("Test", conf)
}
