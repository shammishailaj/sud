package storage

import (
	"database/sql/driver"
	"encoding/hex"

	"errors"

	uuid "github.com/satori/go.uuid"
)

type UUID []byte

func NewUUID() UUID {
	return uuid.NewV4().Bytes()
}
func (u UUID) String() string {
	if u == nil {
		return "00000000-0000-0000-0000-000000000000"
	}
	const dash byte = '-'
	buf := make([]byte, 36)
	hex.Encode(buf[0:8], u[0:4])
	buf[8] = dash
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = dash
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = dash
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = dash
	hex.Encode(buf[24:], u[10:])
	return string(buf)
}
func (u UUID) Value() (driver.Value, error) {
	if u == nil {
		return nil, nil
	}
	return u.String(), nil
}
func (u *UUID) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if ok && len(bytes) == 16 {
		*u = bytes
		return nil
	}
	return errors.New("convert error UUID")

}

type IPoleChecker interface {
	CheckPoleValue(Value Object) error
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
