package corebase

type IPoleChecker interface {
	CheckPoleValue(Value interface{}) error
	Load(poles map[string]interface{})
}

type IRecordWhere interface {
	Load(Poles map[string]interface{}) error
	Save() (string, map[string]interface{}, error)
}

type ITypeInfo interface {
	GetConfigurationName() string
	GetRecordType() string
	GetTitle() string
	GetNew() bool
	GetRead() bool
	GetSave() bool
}

type IPoleInfo interface {
	GetConfigurationName() string
	GetRecordType() string
	GetPoleName() string
	GetPoleType() string
	GetTitle() string
	GetNew() bool
	GetEdit() bool
	GetRemove() bool
	GetDefault() Object
	GetIndexType() string
	GetChecker() IPoleChecker
}
type ICallInfo interface {
	GetConfigurationName() string
	GetTitle() string
	GetName() string
	GetCall() bool
	GetPullName() string
	GetListen() bool
}
