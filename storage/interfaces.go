package storage

type UUID []byte
type IPoleChecker interface {
	CheckPoleValue(Value Object) bool
	Load(doc IDocument)
}
type IDocument interface {
	GetDocumentUID() UUID
	GetDocumentType() string
	GetReadOnly() bool
	GetDeleteDocument() bool
	GetPole(name string) Object
	SetDocumentType(documenttype string)
	SetReadOnly(readonly bool)
	SetDeleteDocument(delete bool)
	SetPole(name string, value Object)
	GetPoleNames() []string
	GetConfiguration() *configuration
}
type ITypeInfo interface {
	GetConfigurationName() string
	GetDocumentType() string
	GetTitle() string
	GetNew() bool
	GetRead() bool
	GetSave() bool
}
type IUser interface {
	GetUserName() string
	GetCheckPassword(Password string) bool
	CheckAccess(Access string) bool
}

type IPoleInfo interface {
	GetConfigurationName() string
	GetDocumentType() string
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
	GetCall() bool
	GetListen() bool
}
type IDocumentWhere interface {
}
type loadConfiguration func(configuration Configuration, State loadState) error
