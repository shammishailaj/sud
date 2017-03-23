package corebase

import (
	"errors"
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

func (w *DocumentWhereContainPole) GetDocumentWhereType() string { return "ContainPole" }

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

func (w *DocumentWhereNotContainPole) GetDocumentWhereType() string { return "NotContainPole" }

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

func (w *DocumentWhereLimit) GetDocumentWhereType() string { return "Limit" }

func (dw *DocumentWhereLimit) Load(Poles map[string]interface{}) error {
	dw.Skip = Poles["DocumentWhere.Limit.Skip"].(int64)
	dw.Count = Poles["DocumentWhere.Limit.Count"].(int64)
	return nil
}

type DocumentWhereOrder struct {
	PoleName string
	ASC      bool
}

func (w *DocumentWhereOrder) GetDocumentWhereType() string { return "Order" }

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

func (w *DocumentWhereSelectDocument) GetDocumentWhereType() string { return "SelectDocument" }

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

func (w *DocumentWhereCompare) GetDocumentWhereType() string { return "Compare" }

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
