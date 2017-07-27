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

func initRecordConfiguration(c *core.Core) {
	conf := core.NewConfiguration([]string{"Default"})
	conf.AddType(core.TypeInfo{ConfigurationName: "Test", RecordType: "Test", AccessType: "Free", AccessNew: "Default", AccessRead: "Default", AccessSave: "Default", Title: "Тест"})
	conf.AddPole(core.PoleInfo{ConfigurationName: "Test", RecordType: "Test", PoleName: "Test.Test1", PoleType: "StringValue", Default: corebase.NULL, IndexType: "Index", Checker: &core.PoleCheckerStringValue{}, AccessRead: "Default", AccessWrite: "Default", Title: "Поле для тестирования"})
	conf.AddCall(core.CallInfo{ConfigurationName: "Test", Name: "TestStd", PullName: "std", AccessCall: "Default", AccessListen: "", Title: ""})
	conf.AddCall(core.CallInfo{ConfigurationName: "Test", Name: "TestAsync", PullName: "async", AccessCall: "Default", AccessListen: "Default", Title: ""})
	c.AddBaseConfiguration("Test", conf)
}

func InitModule(c *core.Core) error {
	initRecordConfiguration(c)
	return core.AddStdCall("TestStd", testStd)
}
