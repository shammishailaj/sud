package core

import "github.com/shammishailaj/sud/corebase"

type TypeInfo struct {
	ConfigurationName string
	RecordType      string
	Title             string
	New               bool
	Read              bool
	Save              bool
}

func (ti *TypeInfo) GetConfigurationName() string { return ti.ConfigurationName }
func (ti *TypeInfo) GetRecordType() string      { return ti.RecordType }
func (ti *TypeInfo) GetTitle() string             { return ti.Title }
func (ti *TypeInfo) GetNew() bool                 { return ti.New }
func (ti *TypeInfo) GetRead() bool                { return ti.Read }
func (ti *TypeInfo) GetSave() bool                { return ti.Save }

var _ti corebase.ITypeInfo = (corebase.ITypeInfo)(&TypeInfo{})
