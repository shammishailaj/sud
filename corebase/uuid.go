package corebase

import (
	"database/sql/driver"
	"errors"

	uuid "github.com/satori/go.uuid"
)

type UUID string

const NullUUID UUID = "00000000-0000-0000-0000-000000000000"

/*struct {
	uuid.UUID
}*/

func NewUUID() UUID {
	return UUID(uuid.NewV4().String())
}
func (u UUID) String() string {
	return string(u)
}

func (u *UUID) Value() (driver.Value, error) {
	if u == nil {
		return nil, nil
	}
	return string(*u), nil
}
func (u *UUID) Scan(value interface{}) error {

	if bytes, ok := value.([]byte); ok && len(bytes) == 16 {
		var id uuid.UUID
		if err := id.UnmarshalBinary(bytes); err != nil {
			return err
		}
		*u = UUID(id.String())
		return nil
	} else if bytes, ok := value.([]byte); ok && len(bytes) == 36 {
		var id uuid.UUID
		if err := id.UnmarshalText(bytes); err != nil {
			return err
		}
		*u = UUID(id.String())
		return nil
	}
	if strings, ok := value.(string); ok {
		var id uuid.UUID
		if err := id.UnmarshalText(([]byte)(strings)); err != nil {
			return err
		}
		*u = UUID(id.String())
		return nil
	}
	return errors.New("convert error UUID")
}
