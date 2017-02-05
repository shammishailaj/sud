package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type PoleCheckerStringValue struct {
	List        map[string]bool
	AllowNull   bool
	MaxLen      bool
	MaxLenValue int64
}

func (pcsv *PoleCheckerStringValue) CheckPoleValue(Value Object) error {
	if pcsv.AllowNull && IsNull(Value) {
		return nil
	}
	var ok bool
	var ValueString string
	if ValueString, ok = Value.(string); ok {
		return errors.New("Value type not string")
	}
	if len(pcsv.List) > 0 {
		if _, ok := pcsv.List[ValueString]; !ok {
			return errors.New("Value " + ValueString + "  note contain list " + fmt.Sprintln(pcsv.List))
		}
	}
	return nil
}
func (pcsv *PoleCheckerStringValue) Load(doc IDocument) {
	CheckerStringValueAllowNull := doc.GetPole("Configuration.PoleInfo.CheckerStringValueAllowNull")
	CheckerStringValueList := doc.GetPole("Configuration.PoleInfo.CheckerStringValueList")
	if !IsNull(CheckerStringValueAllowNull) {
		if v, ok := CheckerStringValueAllowNull.(string); ok {
			pcsv.AllowNull = (v == "True")
		}
	}
	if !IsNull(CheckerStringValueList) {
		if v, ok := CheckerStringValueList.(string); ok {
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

func (pciv *PoleCheckerInt64Value) CheckPoleValue(Value Object) error {
	var ok bool
	var v int64
	if v, ok = Value.(int64); !ok {
		return errors.New("Value type not int64")
	}
	if len(pciv.List) > 0 {
		if _, ok := pciv.List[v]; !ok {
			return errors.New(fmt.Sprintln("Value ", v, "  note contain list ", pciv.List))
		}
	}
	if pciv.Min && v < pciv.MinValue {
		return errors.New(fmt.Sprintln("min", v, "<", pciv.MinValue))
	}
	if pciv.Max && v > pciv.MaxValue {
		return errors.New(fmt.Sprintln("max", v, "<", pciv.MaxValue))
	}
	return nil
}
func (pciv *PoleCheckerInt64Value) Load(doc IDocument) {
	CheckerInt64ValueMin := doc.GetPole("Configuration.PoleInfo.CheckerInt64ValueMin")
	CheckerInt64ValueMax := doc.GetPole("Configuration.PoleInfo.CheckerInt64ValueMax")
	CheckerInt64ValueList := doc.GetPole("Configuration.PoleInfo.CheckerInt64ValueList")
	if v, ok := CheckerInt64ValueMin.(int64); ok {
		pciv.Min = true
		pciv.MinValue = v
	}
	if v, ok := CheckerInt64ValueMax.(int64); ok {
		pciv.Max = true
		pciv.MaxValue = v
	}
	if list, ok := CheckerInt64ValueList.(string); ok {
		for _, s := range strings.Split(list, ",") {
			if v, err := strconv.ParseInt(s, 10, 64); err == nil {
				pciv.List[v] = true
			}
		}
	}
}
