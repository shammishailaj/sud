package core

type TypeInfo struct {
	ConfigurationName string
	DocumentType      string
	Title             string
	New               bool
	Read              bool
	Save              bool
}

func (ti *TypeInfo) GetConfigurationName() string { return ti.ConfigurationName }
func (ti *TypeInfo) GetDocumentType() string      { return ti.DocumentType }
func (ti *TypeInfo) GetTitle() string             { return ti.Title }
func (ti *TypeInfo) GetNew() bool                 { return ti.New }
func (ti *TypeInfo) GetRead() bool                { return ti.Read }
func (ti *TypeInfo) GetSave() bool                { return ti.Save }

var _ti ITypeInfo = (ITypeInfo)(&TypeInfo{})
