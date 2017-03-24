package corebase

import (
	"errors"
	"time"
)

var documentWheres = map[string]func() IDocumentWhere{
	"ContainPole": func() IDocumentWhere {
		return &DocumentWhereContainPole{}
	},
	"NotContainPole": func() IDocumentWhere {
		return &DocumentWhereNotContainPole{}
	},
	"Limit": func() IDocumentWhere {
		return &DocumentWhereLimit{}
	},
	"Order": func() IDocumentWhere {
		return &DocumentWhereOrder{}
	},
	"SelectDocument": func() IDocumentWhere {
		return &DocumentWhereSelectDocument{}
	},
	"Compare": func() IDocumentWhere {
		return &DocumentWhereCompare{}
	},
}

func InitAddDocumentWhere(DocumentWhereType string, newWhere func() IDocumentWhere) {
	documentWheres[DocumentWhereType] = newWhere
}
func NewDocumentWhere(DocumentWhereType string) (IDocumentWhere, error) {
	var n func() IDocumentWhere
	var ok bool
	if n, ok = documentWheres[DocumentWhereType]; !ok {
		return nil, errors.New("DocumentWhere " + DocumentWhereType + " not found")
	}
	return n(), nil
}

type DocumentWhereContainPole struct {
	PoleName string
}

func (dw *DocumentWhereContainPole) Save() (string, map[string]interface{}, error) {
	return "ContainPole", map[string]interface{}{"DocumentWhere.ContainPole.PoleName": dw.PoleName}, nil
}

func (dw *DocumentWhereContainPole) Load(Poles map[string]interface{}) error {
	var ok bool
	if dw.PoleName, ok = Poles["DocumentWhere.ContainPole.PoleName"].(string); !ok {
		return errors.New("not found DocumentWhere.ContainPole.PoleName")
	}
	return nil
}

type DocumentWhereNotContainPole struct {
	PoleName    string
	InTableName string
	TabLock     bool
}

func (dw *DocumentWhereNotContainPole) Save() (string, map[string]interface{}, error) {
	p := map[string]interface{}{"DocumentWhere.NotContainPole.PoleName": dw.PoleName, "DocumentWhere.NotContainPole.InTableName": dw.InTableName}
	if dw.TabLock == true {
		p["DocumentWhere.NotContainPole.TabLock"] = dw.TabLock
	}
	return "NotContainPole", p, nil
}

func (dw *DocumentWhereNotContainPole) Load(Poles map[string]interface{}) error {
	var ok bool
	if dw.PoleName, ok = Poles["DocumentWhere.NotContainPole.PoleName"].(string); !ok {
		return errors.New("not found DocumentWhere.NotContainPole.PoleName")
	}
	if dw.InTableName, ok = Poles["DocumentWhere.NotContainPole.InTableName"].(string); !ok {
		return errors.New("not found DocumentWhere.NotContainPole.InTableName")
	}
	if dw.TabLock, ok = Poles["DocumentWhere.NotContainPole.TabLock"].(bool); !ok {
		dw.TabLock = false
	}
	return nil
}

type DocumentWhereLimit struct {
	Skip  int64
	Count int64
}

func (dw *DocumentWhereLimit) Save() (string, map[string]interface{}, error) {
	return "Limit", map[string]interface{}{"DocumentWhere.Limit.Skip": dw.Skip, "DocumentWhere.Limit.Count": dw.Count}, nil
}

func (dw *DocumentWhereLimit) Load(Poles map[string]interface{}) error {
	dw.Skip = Poles["DocumentWhere.Limit.Skip"].(int64)
	dw.Count = Poles["DocumentWhere.Limit.Count"].(int64)
	return nil
}

type DocumentWhereOrder struct {
	PoleName string
	ASC      bool
}

func (dw *DocumentWhereOrder) Save() (string, map[string]interface{}, error) {
	p := map[string]interface{}{"DocumentWhere.Order.PoleName": dw.PoleName}
	if dw.ASC == true {
		p["DocumentWhere.Order.ASC"] = "ASC"
	}
	return "Order", p, nil
}

func (dw *DocumentWhereOrder) Load(Poles map[string]interface{}) error {
	var ok bool
	if dw.PoleName, ok = Poles["DocumentWhere.Order.PoleName"].(string); !ok {
		return errors.New("not found DocumentWhere.Order.PoleName")
	}
	var ASC string
	if ASC, ok = Poles["DocumentWhere.Order.ASC"].(string); ok {
		dw.ASC = ASC == "ASC"
	}
	return nil
}

type DocumentWhereSelectDocument struct {
	DocumentUID string
}

func (dw *DocumentWhereSelectDocument) Save() (string, map[string]interface{}, error) {
	return "SelectDocument", map[string]interface{}{"DocumentWhere.SelectDocument.DocumentUID": dw.DocumentUID}, nil
}

func (dw *DocumentWhereSelectDocument) Load(Poles map[string]interface{}) error {
	var ok bool
	var UID string
	if UID, ok = Poles["DocumentWhere.SelectDocument.DocumentUID"].(string); ok {
		dw.DocumentUID = UID
	}
	return nil
}

type DocumentWhereCompare struct {
	PoleName      string
	Operation     string
	ExtensionName string
	//dCallFunction Extension = null;
	Value Object
}

func (dw *DocumentWhereCompare) Save() (string, map[string]interface{}, error) {
	p := map[string]interface{}{"DocumentWhere.Compare.PoleName": dw.PoleName, "DocumentWhere.Compare.Operation": dw.Operation}
	switch v := dw.Value.(type) {
	case string:
		p["DocumentWhere.Compare.StringValue"] = v
	case int64:
		p["DocumentWhere.Compare.Int64Value"] = v
	case time.Time:
		p["DocumentWhere.Compare.DateTimeValue"] = v
	}
	if dw.ExtensionName != "" {
		p["DocumentWhere.Compare.Extension"] = dw.ExtensionName
	}
	return "Compare", p, nil
}

func (dw *DocumentWhereCompare) Load(Poles map[string]interface{}) error {
	var ok bool
	if dw.PoleName, ok = Poles["DocumentWhere.Compare.PoleName"].(string); !ok {
		return errors.New("not found DocumentWhere.Compare.PoleName")
	}
	if dw.Operation, ok = Poles["DocumentWhere.Compare.Operation"].(string); !ok {
		return errors.New("not found DocumentWhere.Compare.Operation")
	}
	poles := map[string]bool{
		"DocumentWhere.Compare.StringValue":       true,
		"DocumentWhere.Compare.Int64Value":        true,
		"DocumentWhere.Compare.DateValue":         true,
		"DocumentWhere.Compare.DateTimeValue":     true,
		"DocumentWhere.Compare.DocumentLinkValue": true,
	}
	for p, value := range Poles {
		if _, ok := poles[p]; ok {
			dw.Value = value
			break
		}
	}
	if ext := Poles["DocumentWhere.Compare.Extension"]; !IsNull(ext) {
		dw.ExtensionName = ext.(string)
	}
	return nil
}
