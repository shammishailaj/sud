package core

import "github.com/crazyprograms/sud/corebase"

type TypeInfo struct {
	ConfigurationName string
	RecordType        string
	Title             string
	AccessNew         string
	AccessRead        string
	AccessSave        string
	AccessType        string
}

func (ti *TypeInfo) GetConfigurationName() string { return ti.ConfigurationName }
func (ti *TypeInfo) GetRecordType() string        { return ti.RecordType }
func (ti *TypeInfo) GetTitle() string             { return ti.Title }
func (ti *TypeInfo) GetAccessType() string        { return ti.AccessType }
func (ti *TypeInfo) GetAccessNew() string         { return ti.AccessNew }
func (ti *TypeInfo) GetAccessRead() string        { return ti.AccessRead }
func (ti *TypeInfo) GetAccessSave() string        { return ti.AccessSave }

var _ti corebase.ITypeInfo = (corebase.ITypeInfo)(&TypeInfo{})
