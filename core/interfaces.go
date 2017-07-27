package core

import (
	"time"

	"github.com/crazyprograms/sud/callpull"
	"github.com/crazyprograms/sud/corebase"
)

type ICallPull interface {
	Call(Name string, Param map[string]interface{}, timeoutWait time.Duration, Access corebase.IAccess) (callpull.Result, error)
}
type IListenPull interface {
	Listen(Name string, timeoutWait time.Duration) (Param map[string]interface{}, Access corebase.IAccess, Result chan callpull.Result, err error)
}

type ICallListenPull interface {
	ICallPull
	IListenPull
}
