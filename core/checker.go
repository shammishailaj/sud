package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/crazyprograms/sud/corebase"
)

type PoleCheckerStringValue struct {
	List        map[string]bool
	AllowNull   bool
	MaxLen      bool
	MaxLenValue int64
}

func (pcsv *PoleCheckerStringValue) CheckPoleValue(Value interface{}) error {
	if pcsv.AllowNull && corebase.IsNull(Value) {
		return nil
	}
	var ok bool
	var ValueString string
	if ValueString, ok = Value.(string); !ok {
		return &corebase.Error{ErrorType: corebase.ErrorTypeInfo, Action: "CheckPoleValue", Name: "PoleCheckerStringValue", Info: "Value type not string"}
	}
	if len(pcsv.List) > 0 {
		if _, ok := pcsv.List[ValueString]; !ok {
			return &corebase.Error{ErrorType: corebase.ErrorTypeInfo, Action: "CheckPoleValue", Name: "PoleCheckerStringValue", Info: "Value " + ValueString + "  note contain list " + fmt.Sprintln(pcsv.List)}
		}
	}
	return nil
}
func (pcsv *PoleCheckerStringValue) Load(Poles map[string]interface{}) {
	if CheckerStringValueAllowNull, ok := Poles["Configuration.PoleInfo.CheckerStringValueAllowNull"]; ok && !corebase.IsNull(CheckerStringValueAllowNull) {
		if v, ok := CheckerStringValueAllowNull.(string); ok {
			pcsv.AllowNull = (v == "True")
		}
	}
	if CheckerStringValueList, ok := Poles["Configuration.PoleInfo.CheckerStringValueList"]; ok && !corebase.IsNull(CheckerStringValueList) {
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

func (pciv *PoleCheckerInt64Value) CheckPoleValue(Value interface{}) error {
	var ok bool
	var v int64
	if v, ok = Value.(int64); !ok {
		return &corebase.Error{ErrorType: corebase.ErrorTypeInfo, Action: "CheckPoleValue", Name: "PoleCheckerInt64Value", Info: "Value type not int64"}
	}
	if len(pciv.List) > 0 {
		if _, ok := pciv.List[v]; !ok {
			return &corebase.Error{ErrorType: corebase.ErrorTypeInfo, Action: "CheckPoleValue", Name: "PoleCheckerInt64Value", Info: fmt.Sprintln("Value ", v, "  note contain list ", pciv.List)}
		}
	}
	if pciv.Min && v < pciv.MinValue {
		return &corebase.Error{ErrorType: corebase.ErrorTypeInfo, Action: "CheckPoleValue", Name: "PoleCheckerInt64Value", Info: fmt.Sprintln("min", v, "<", pciv.MinValue)}
	}
	if pciv.Max && v > pciv.MaxValue {
		return &corebase.Error{ErrorType: corebase.ErrorTypeInfo, Action: "CheckPoleValue", Name: "PoleCheckerInt64Value", Info: fmt.Sprintln("max", v, ">", pciv.MaxValue)}
	}
	return nil
}
func (pciv *PoleCheckerInt64Value) Load(Poles map[string]interface{}) {

	if CheckerInt64ValueMin, ok1 := Poles["Configuration.PoleInfo.CheckerInt64ValueMin"]; ok1 {
		if v, ok := CheckerInt64ValueMin.(int64); ok {
			pciv.Min = true
			pciv.MinValue = v
		}
	}
	if CheckerInt64ValueMax, ok1 := Poles["Configuration.PoleInfo.CheckerInt64ValueMax"]; ok1 {
		if v, ok := CheckerInt64ValueMax.(int64); ok {
			pciv.Max = true
			pciv.MaxValue = v
		}
	}
	if CheckerInt64ValueList, ok1 := Poles["Configuration.PoleInfo.CheckerInt64ValueList"]; ok1 {
		if list, ok := CheckerInt64ValueList.(string); ok {
			for _, s := range strings.Split(list, ",") {
				if v, err := strconv.ParseInt(s, 10, 64); err == nil {
					pciv.List[v] = true
				}
			}
		}
	}
}
