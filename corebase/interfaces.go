package corebase

type IPoleChecker interface {
	CheckPoleValue(Value interface{}) error
	Load(poles map[string]interface{}) error
	Save() map[string]interface{}
}

type IRecordWhere interface {
	Load(Poles map[string]interface{}) error
	Save() (string, map[string]interface{}, error)
}

type ITypeInfo interface {
	GetConfigurationName() string
	GetRecordType() string
	GetTitle() string
	GetAccessType() string
	GetAccessNew() string
	GetAccessRead() string
	GetAccessSave() string
}

type IPoleInfo interface {
	GetConfigurationName() string
	GetRecordType() string
	GetPoleName() string
	GetPoleType() string
	GetTitle() string
	GetAccessRead() string
	GetAccessWrite() string
	GetDefault() interface{}
	GetIndexType() string
	GetChecker() IPoleChecker
}
type ICallInfo interface {
	GetConfigurationName() string
	GetTitle() string
	GetName() string
	GetPullName() string
	GetAccessCall() string
	GetAccessListen() string
}

type IAccess interface {
	CheckAccess(Access string) bool
	Users() []IUser
}

type IUser interface {
	IAccess
	GetLogin() string
	GetCheckPassword(Password string) bool
}
