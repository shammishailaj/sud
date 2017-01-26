package storage

import (
	"errors"
	"time"
)

type Object struct {
	value       interface{}
	currentType string
}

func NewObject(value interface{}) Object {
	o := Object{}
	o.Set(value)
	return o
}
func (o *Object) Get() interface{} {
	return o.value
}
func (o *Object) Set(value interface{}) {
	switch v := value.(type) {
	case nil:
		o.SetNull()
	case string:
		o.SetString(v)
	case int:
		o.SetInt64(int64(v))
	case int32:
		o.SetInt64(int64(v))
	case int64:
		o.SetInt64(v)
	case time.Time:
		o.SetDateTime(v)
	case UUID:
		o.SetDocumentLink(v)
	default:
		panic("Object type error")
	}
}

func (o *Object) SetNull() {
	o.value = nil
	o.currentType = ""
}
func (o Object) IsNull() bool {
	return o.currentType == ""
}
func (o *Object) SetString(value string) {
	o.value = value
	o.currentType = "String"
}

func (o Object) String() (string, error) {
	if o.currentType != "String" {
		return "", errors.New("type error")
	}
	return o.value.(string), nil
}

func (o *Object) SetInt64(value int64) {
	o.value = value
	o.currentType = "Int64"
}
func (o Object) Int64() (int64, error) {
	if o.currentType != "Int64" {
		return 0, errors.New("type error")
	}
	return o.value.(int64), nil
}

func (o *Object) SetBoolean(value bool) {
	o.value = value
	o.currentType = "Boolean"
}
func (o Object) Boolean() (bool, error) {
	if o.currentType != "Boolean" {
		return false, errors.New("type error")
	}
	return o.value.(bool), nil
}
func (o *Object) SetDateTime(value time.Time) {
	o.value = value
	o.currentType = "DateTime"
}
func (o Object) DateTime() (time.Time, error) {
	if o.currentType != "DateTime" {
		return time.Time{}, errors.New("type error")
	}
	return o.value.(time.Time), nil
}
func (o *Object) SetDate(value time.Time) {
	o.value = value
	o.currentType = "Date"
}
func (o Object) Date() (time.Time, error) {
	if o.currentType != "Date" {
		return time.Time{}, errors.New("type error")
	}
	return o.value.(time.Time), nil
}
func (o *Object) SetDocumentLink(value UUID) {
	o.value = value
	o.currentType = "DocumentLink"
}

func (o Object) DocumentLink() (UUID, error) {
	if o.currentType != "DocumentLink" {
		return nil, errors.New("type error")
	}
	return o.value.(UUID), nil
}

func (o Object) Type() string {
	return o.currentType
}
