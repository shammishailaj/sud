package storage

type PoleInfo struct {
	ConfigurationName string
	DocumentType      string
	PoleName          string
	PoleType          string
	Title             string
	New               bool
	Edit              bool
	Remove            bool
	Default           Object
	IndexType         string
	Checker           IPoleChecker
}

func (pi *PoleInfo) GetConfigurationName() string { return pi.ConfigurationName }
func (pi *PoleInfo) GetDocumentType() string      { return pi.DocumentType }
func (pi *PoleInfo) GetPoleName() string          { return pi.PoleName }
func (pi *PoleInfo) GetPoleType() string          { return pi.PoleType }
func (pi *PoleInfo) GetTitle() string             { return pi.Title }
func (pi *PoleInfo) GetNew() bool                 { return pi.New }
func (pi *PoleInfo) GetEdit() bool                { return pi.Edit }
func (pi *PoleInfo) GetRemove() bool              { return pi.Remove }
func (pi *PoleInfo) GetDefault() Object           { return pi.Default }
func (pi *PoleInfo) GetIndexType() string         { return pi.IndexType }
func (pi *PoleInfo) GetChecker() IPoleChecker     { return pi.Checker }

var _pi IPoleInfo = (IPoleInfo)(&PoleInfo{})
