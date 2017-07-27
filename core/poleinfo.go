package core

import "github.com/crazyprograms/sud/corebase"

type PoleInfo struct {
	ConfigurationName string
	RecordType        string
	PoleName          string
	PoleType          string
	Title             string
	AccessRead        string
	AccessWrite       string
	Default           corebase.Object
	IndexType         string
	Checker           corebase.IPoleChecker
}

func (pi *PoleInfo) GetConfigurationName() string      { return pi.ConfigurationName }
func (pi *PoleInfo) GetRecordType() string             { return pi.RecordType }
func (pi *PoleInfo) GetPoleName() string               { return pi.PoleName }
func (pi *PoleInfo) GetPoleType() string               { return pi.PoleType }
func (pi *PoleInfo) GetTitle() string                  { return pi.Title }
func (pi *PoleInfo) GetAccessRead() string             { return pi.AccessRead }
func (pi *PoleInfo) GetAccessWrite() string            { return pi.AccessWrite }
func (pi *PoleInfo) GetDefault() corebase.Object       { return pi.Default }
func (pi *PoleInfo) GetIndexType() string              { return pi.IndexType }
func (pi *PoleInfo) GetChecker() corebase.IPoleChecker { return pi.Checker }

var _pi corebase.IPoleInfo = (corebase.IPoleInfo)(&PoleInfo{})
