package corebase

import (
	"database/sql/driver"
	"errors"

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
