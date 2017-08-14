package core

import (
	"fmt"
	"reflect"

	"github.com/crazyprograms/sud/structures"

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
func (pcsv *PoleCheckerStringValue) Load(Poles map[string]interface{}) error {
	pcsv.AllowNull = false
	pcsv.List = map[string]bool{}
	if AllowNull, ok := Poles["allowNull"]; ok {
		if v, ok := AllowNull.(bool); ok {
			pcsv.AllowNull = v
		}
	}
	if list, ok := structures.MapGetStringList(Poles, "list"); ok {
		for _, item := range list {
			pcsv.List[item] = true
		}
	}
	return nil
}
func (pcsv *PoleCheckerStringValue) Save() map[string]interface{} {
	r := map[string]interface{}{}
	r["allowNull"] = pcsv.AllowNull
	if pcsv.List != nil {
		list := make([]string, len(pcsv.List))
		i := 0
		for v := range pcsv.List {
			list[i] = v
			i++
		}
		structures.MapSetStringList(r, "list", list)
	}
	return r
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
func (pciv *PoleCheckerInt64Value) Load(Poles map[string]interface{}) error {
	pciv.AllowNull = false
	pciv.List = map[int64]bool{}
	if AllowNull, ok := Poles["allowNull"]; ok {
		if v, ok := AllowNull.(bool); ok {
			pciv.AllowNull = v
		}
	}
	if v, ok := Poles["min"]; ok {
		if min, ok := v.(int64); ok {
			pciv.Min = true
			pciv.MinValue = min
		}
	}
	if v, ok := Poles["max"]; ok {
		if max, ok := v.(int64); ok {
			pciv.Max = true
			pciv.MaxValue = max
		}
	}
	if list, ok := structures.MapGetInt64List(Poles, "list"); ok {
		for _, item := range list {
			pciv.List[item] = true
		}
	}
	return nil
}
func (pciv *PoleCheckerInt64Value) Save() map[string]interface{} {
	r := map[string]interface{}{}
	r["allowNull"] = pciv.AllowNull
	if pciv.Min {
		r["min"] = pciv.MinValue
	}
	if pciv.Max {
		r["max"] = pciv.MaxValue
	}
	if pciv.List != nil {
		list := make([]int64, len(pciv.List))
		i := 0
		for v := range pciv.List {
			list[i] = v
			i++
		}
		structures.MapSetInt64List(r, "list", list)
	}
	return r
}
func LoadChecker(data map[string]interface{}) (corebase.IPoleChecker, error) {
	var checker corebase.IPoleChecker
	var v interface{}
	var vc map[string]interface{}
	var ok bool
	if v, ok = data["PoleCheckerStringValue"]; ok {
		checker = &PoleCheckerStringValue{}
	} else if v, ok = data["PoleCheckerInt64Value"]; ok {
		checker = &PoleCheckerInt64Value{}
	}
	if vc, ok = v.(map[string]interface{}); !ok {
		return nil, &corebase.Error{ErrorType: corebase.ErrorFormat, Action: "LoadChecker", Info: "structure error"}
	}
	if err := checker.Load(vc); err != nil {
		return nil, err
	}
	return checker, nil
}
func SaveChecker(checker corebase.IPoleChecker) map[string]interface{} {
	r := make(map[string]interface{})
	if checker == nil {
		return nil
	}
	name := reflect.TypeOf(checker).Elem().Name()
	r[name] = checker.Save()
	return r
}
