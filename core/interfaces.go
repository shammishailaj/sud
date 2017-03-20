package core

import (
	"database/sql/driver"
	"time"

	"errors"

	"github.com/crazyprograms/callpull"
	uuid "github.com/satori/go.uuid"
)

type UUID struct {
	uuid.UUID
}

func NewUUID() UUID {
	return UUID{UUID: uuid.NewV4()}
}
func (u *UUID) Value() (driver.Value, error) {
	if u == nil {
		return nil, nil
	}
	return u.String(), nil
}
func (u *UUID) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if ok && len(bytes) == 16 {
		u.UnmarshalBinary(bytes)
		return nil
	}
	strings, ok := value.(string)
	if ok {
		if err := u.UnmarshalText(([]byte)(strings)); err != nil {
			return err
		}
	}
	return errors.New("convert error UUID")
}

type IPoleChecker interface {
	CheckPoleValue(Value interface{}) error
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
	SetPole(name string, value Object) error
	GetPoleNames() []string
	GetConfiguration() *Configuration
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
	GetName() string
	GetCall() bool
	GetPullName() string
	GetListen() bool
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
