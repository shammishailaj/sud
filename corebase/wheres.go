package corebase

import (
	"errors"
	"time"
)

var recordWheres = map[string]func() IRecordWhere{
	"ContainPole": func() IRecordWhere {
		return &RecordWhereContainPole{}
	},
	"NotContainPole": func() IRecordWhere {
		return &RecordWhereNotContainPole{}
	},
	"Limit": func() IRecordWhere {
		return &RecordWhereLimit{}
	},
	"Order": func() IRecordWhere {
		return &RecordWhereOrder{}
	},
	"SelectRecord": func() IRecordWhere {
		return &RecordWhereSelectRecord{}
	},
	"Compare": func() IRecordWhere {
		return &RecordWhereCompare{}
	},
}

func InitAddRecordWhere(RecordWhereType string, newWhere func() IRecordWhere) {
	recordWheres[RecordWhereType] = newWhere
}
func NewRecordWhere(RecordWhereType string) (IRecordWhere, error) {
	var n func() IRecordWhere
	var ok bool
	if n, ok = recordWheres[RecordWhereType]; !ok {
		return nil, errors.New("RecordWhere " + RecordWhereType + " not found")
	}
	return n(), nil
}

type RecordWhereContainPole struct {
	PoleName string
}

func (dw *RecordWhereContainPole) Save() (string, map[string]interface{}, error) {
	return "ContainPole", map[string]interface{}{"RecordWhere.ContainPole.PoleName": dw.PoleName}, nil
}

func (dw *RecordWhereContainPole) Load(Poles map[string]interface{}) error {
	var ok bool
	if dw.PoleName, ok = Poles["RecordWhere.ContainPole.PoleName"].(string); !ok {
		return errors.New("not found RecordWhere.ContainPole.PoleName")
	}
	return nil
}

type RecordWhereNotContainPole struct {
	PoleName    string
	InTableName string
	TabLock     bool
}

func (dw *RecordWhereNotContainPole) Save() (string, map[string]interface{}, error) {
	p := map[string]interface{}{"RecordWhere.NotContainPole.PoleName": dw.PoleName, "RecordWhere.NotContainPole.InTableName": dw.InTableName}
	if dw.TabLock == true {
		p["RecordWhere.NotContainPole.TabLock"] = dw.TabLock
	}
	return "NotContainPole", p, nil
}

func (dw *RecordWhereNotContainPole) Load(Poles map[string]interface{}) error {
	var ok bool
	if dw.PoleName, ok = Poles["RecordWhere.NotContainPole.PoleName"].(string); !ok {
		return errors.New("not found RecordWhere.NotContainPole.PoleName")
	}
	if dw.InTableName, ok = Poles["RecordWhere.NotContainPole.InTableName"].(string); !ok {
		return errors.New("not found RecordWhere.NotContainPole.InTableName")
	}
	if dw.TabLock, ok = Poles["RecordWhere.NotContainPole.TabLock"].(bool); !ok {
		dw.TabLock = false
	}
	return nil
}

type RecordWhereLimit struct {
	Skip  int64
	Count int64
}

func (dw *RecordWhereLimit) Save() (string, map[string]interface{}, error) {
	return "Limit", map[string]interface{}{"RecordWhere.Limit.Skip": dw.Skip, "RecordWhere.Limit.Count": dw.Count}, nil
}

func (dw *RecordWhereLimit) Load(Poles map[string]interface{}) error {
	dw.Skip = Poles["RecordWhere.Limit.Skip"].(int64)
	dw.Count = Poles["RecordWhere.Limit.Count"].(int64)
	return nil
}

type RecordWhereOrder struct {
	PoleName string
	ASC      bool
}

func (dw *RecordWhereOrder) Save() (string, map[string]interface{}, error) {
	p := map[string]interface{}{"RecordWhere.Order.PoleName": dw.PoleName}
	if dw.ASC == true {
		p["RecordWhere.Order.ASC"] = "ASC"
	}
	return "Order", p, nil
}

func (dw *RecordWhereOrder) Load(Poles map[string]interface{}) error {
	var ok bool
	if dw.PoleName, ok = Poles["RecordWhere.Order.PoleName"].(string); !ok {
		return errors.New("not found RecordWhere.Order.PoleName")
	}
	var ASC string
	if ASC, ok = Poles["RecordWhere.Order.ASC"].(string); ok {
		dw.ASC = ASC == "ASC"
	}
	return nil
}

type RecordWhereSelectRecord struct {
	RecordUID UUID
}

func (dw *RecordWhereSelectRecord) Save() (string, map[string]interface{}, error) {
	return "SelectRecord", map[string]interface{}{"RecordWhere.SelectRecord.RecordUID": dw.RecordUID}, nil
}

func (dw *RecordWhereSelectRecord) Load(Poles map[string]interface{}) error {
	var ok bool
	var UID UUID
	if UID, ok = Poles["RecordWhere.SelectRecord.RecordUID"].(UUID); ok {
		dw.RecordUID = UID
	}
	return nil
}

type RecordWhereCompare struct {
	PoleName      string
	Operation     string
	ExtensionName string
	//dCallFunction Extension = null;
	Value interface{}
}

func (dw *RecordWhereCompare) Save() (string, map[string]interface{}, error) {
	p := map[string]interface{}{"RecordWhere.Compare.PoleName": dw.PoleName, "RecordWhere.Compare.Operation": dw.Operation}
	switch v := dw.Value.(type) {
	case string:
		p["RecordWhere.Compare.StringValue"] = v
	case int64:
		p["RecordWhere.Compare.Int64Value"] = v
	case time.Time:
		p["RecordWhere.Compare.DateTimeValue"] = v
	}
	if dw.ExtensionName != "" {
		p["RecordWhere.Compare.Extension"] = dw.ExtensionName
	}
	return "Compare", p, nil
}

func (dw *RecordWhereCompare) Load(Poles map[string]interface{}) error {
	var ok bool
	if dw.PoleName, ok = Poles["RecordWhere.Compare.PoleName"].(string); !ok {
		return errors.New("not found RecordWhere.Compare.PoleName")
	}
	if dw.Operation, ok = Poles["RecordWhere.Compare.Operation"].(string); !ok {
		return errors.New("not found RecordWhere.Compare.Operation")
	}
	poles := map[string]bool{
		"RecordWhere.Compare.StringValue":     true,
		"RecordWhere.Compare.Int64Value":      true,
		"RecordWhere.Compare.DateValue":       true,
		"RecordWhere.Compare.DateTimeValue":   true,
		"RecordWhere.Compare.RecordLinkValue": true,
	}
	for p, value := range Poles {
		if _, ok := poles[p]; ok {
			dw.Value = value
			break
		}
	}
	if ext := Poles["RecordWhere.Compare.Extension"]; !IsNull(ext) {
		dw.ExtensionName = ext.(string)
	}
	return nil
}
