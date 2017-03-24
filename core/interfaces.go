package core

import (
	"time"

	"github.com/crazyprograms/callpull"
)

type IUser interface {
	GetUserName() string
	GetCheckPassword(Password string) bool
	CheckAccess(Access string) bool
}

type ICallPull interface {
	Call(Name string, Param map[string]interface{}, timeoutWait time.Duration) (callpull.Result, error)
}
type IListenPull interface {
	Listen(Name string, timeoutWait time.Duration) (Param map[string]interface{}, Result chan callpull.Result, err error)
}

type ICallListenPull interface {
	ICallPull
	IListenPull
}
