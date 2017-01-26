package storage

import (
	"strconv"
	"strings"
)

type PoleCheckerStringValue struct {
	List        map[string]bool
	AllowNull   bool
	MaxLen      bool
	MaxLenValue int64
}

func (pcsv *PoleCheckerStringValue) CheckPoleValue(Value Object) bool {
	if pcsv.AllowNull && Value.IsNull() {
		return true
	}
	if Value.Type() != "String" {
		return false
	}
	if len(pcsv.List) > 0 {
		if v, err := Value.String(); err == nil {
			if _, ok := pcsv.List[v]; !ok {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
func (pcsv *PoleCheckerStringValue) Load(doc IDocument) {
	if !doc.GetPole("Configuration.PoleInfo.CheckerStringValueAllowNull").IsNull() {
		if v, err := doc.GetPole("Configuration.PoleInfo.CheckerStringValueAllowNull").String(); err != nil {
			pcsv.AllowNull = (v == "True")
		}
	}
	if !doc.GetPole("Configuration.PoleInfo.CheckerStringValueList").IsNull() {
		if v, err := doc.GetPole("Configuration.PoleInfo.CheckerStringValueList").String(); err != nil {
			for _, s := range strings.Split(v, ",") {
				pcsv.List[s] = true
			}
		}
	}
}

type PoleCheckerInt64Value struct {
	Min       bool
	MinValue  int64
	Max       bool
	MaxValue  int64
	List      map[int64]bool
	AllowNull bool
}

func (pciv *PoleCheckerInt64Value) CheckPoleValue(Value Object) bool {
	v, err := Value.Int64()
	if err != nil {
		return false
	}
	if len(pciv.List) > 0 {
		if _, ok := pciv.List[v]; !ok {
			return false
		}
	}
	if pciv.Min && v < pciv.MinValue {
		return false
	}
	if pciv.Max && v > pciv.MaxValue {
		return false
	}
	return true

}
func (pciv *PoleCheckerInt64Value) Load(doc IDocument) {
	if v, err := doc.GetPole("Configuration.PoleInfo.CheckerInt64ValueMin").Int64(); err != nil {
		pciv.Min = true
		pciv.MinValue = v
	}
	if v, err := doc.GetPole("Configuration.PoleInfo.CheckerInt64ValueMax").Int64(); err != nil {
		pciv.Max = true
		pciv.MaxValue = v
	}
	if !doc.GetPole("Configuration.PoleInfo.CheckerInt64ValueList").IsNull() {
		if v, err := doc.GetPole("Configuration.PoleInfo.CheckerInt64ValueList").String(); err != nil {
			for _, s := range strings.Split(v, ",") {
				if v, err := strconv.ParseInt(s, 10, 64); err == nil {
					pciv.List[v] = true
				}
			}
		}
	}
}
