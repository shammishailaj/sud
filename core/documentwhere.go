package core

import (
	"errors"
)

type DocumentWhereContainPole struct {
	PoleName string
}

func (dw *DocumentWhereContainPole) Load(doc IDocument) error {
	var ok bool
	if dw.PoleName, ok = doc.GetPole("DocumentWhere.ContainPole.PoleName").(string); !ok {
		return errors.New("not found DocumentWhere.ContainPole.PoleName")
	}
	return nil
}

type DocumentWhereNotContainPole struct {
	PoleName    string
	InTableName string
	TabLock     bool
}

func (dw *DocumentWhereNotContainPole) Load(doc IDocument) error {
	var ok bool
	if dw.PoleName, ok = doc.GetPole("DocumentWhere.NotContainPole.PoleName").(string); !ok {
		return errors.New("not found DocumentWhere.NotContainPole.PoleName")
	}
	if dw.InTableName, ok = doc.GetPole("DocumentWhere.NotContainPole.InTableName").(string); !ok {
		return errors.New("not found DocumentWhere.NotContainPole.InTableName")
	}
	if dw.TabLock, ok = doc.GetPole("DocumentWhere.NotContainPole.TabLock").(bool); !ok {
		dw.TabLock = false
	}
	return nil
}

type DocumentWhereLimit struct {
	Skip  int64
	Count int64
}

func (dw *DocumentWhereLimit) Load(doc IDocument) error {
	dw.Skip = doc.GetPole("DocumentWhere.Limit.Skip").(int64)
	dw.Count = doc.GetPole("DocumentWhere.Limit.Count").(int64)
	return nil
}

type DocumentWhereOrder struct {
	PoleName string
	ASC      bool
}

func (dw *DocumentWhereOrder) Load(doc IDocument) error {
	var ok bool
	if dw.PoleName, ok = doc.GetPole("DocumentWhere.Order.PoleName").(string); !ok {
		return errors.New("not found DocumentWhere.Order.PoleName")
	}
	var ASC string
	if ASC, ok = doc.GetPole("DocumentWhere.Order.ASC").(string); ok {
		dw.ASC = ASC == "ASC"
	}
	return nil
}

type DocumentWhereSelectDocument struct {
	DocumentUID UUID
}

func (dw *DocumentWhereSelectDocument) Load(doc IDocument) error {
	var ok bool
	var UID string
	if UID, ok = doc.GetPole("DocumentWhere.Order.ASC").(string); ok {
		if err := dw.DocumentUID.Scan(UID); err != nil {
			return err
		}
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

func (dw *DocumentWhereCompare) Load(doc IDocument) error {
	var ok bool
	if dw.PoleName, ok = doc.GetPole("DocumentWhere.Compare.PoleName").(string); !ok {
		return errors.New("not found DocumentWhere.Compare.PoleName")
	}
	if dw.Operation, ok = doc.GetPole("DocumentWhere.Compare.Operation").(string); !ok {
		return errors.New("not found DocumentWhere.Compare.Operation")
	}
	poles := map[string]bool{
		"DocumentWhere.Compare.StringValue":       true,
		"DocumentWhere.Compare.Int64Value":        true,
		"DocumentWhere.Compare.DateValue":         true,
		"DocumentWhere.Compare.DateTimeValue":     true,
		"DocumentWhere.Compare.DocumentLinkValue": true,
	}
	for _, p := range doc.GetPoleNames() {
		if _, ok := poles[p]; ok {
			dw.Value = doc.GetPole(p)
			break
		}
	}
	if ext := doc.GetPole("DocumentWhere.Compare.Extension"); !IsNull(ext) {
		dw.ExtensionName = ext.(string)
	}
	return nil
}
