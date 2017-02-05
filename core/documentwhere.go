package core

type DocumentWhereContainPole struct {
	PoleName string
}

func (dw *DocumentWhereContainPole) Load(doc IDocument) error {
	var err error
	if dw.PoleName, err = doc.GetPole("DocumentWhere.ContainPole.PoleName").String(); err != nil {
		return err
	}
	return nil
}

type DocumentWhereNotContainPole struct {
	PoleName    string
	InTableName string
	TabLock     bool
}

func (dw *DocumentWhereNotContainPole) Load(doc IDocument) error {
	var err error
	if dw.PoleName, err = doc.GetPole("DocumentWhere.NotContainPole.PoleName").String(); err != nil {
		return err
	}
	if dw.InTableName, err = doc.GetPole("DocumentWhere.NotContainPole.InTableName").String(); err != nil {
		return err
	}
	if doc.GetPole("DocumentWhere.NotContainPole.TabLock").IsNull() {
		if dw.TabLock, err = doc.GetPole("DocumentWhere.NotContainPole.TabLock").Boolean(); err != nil {
			return err
		}
	} else {
		dw.TabLock = false
	}
	return nil
}

type DocumentWhereLimit struct {
	Skip  int64
	Count int64
}

func (dw *DocumentWhereLimit) Load(doc IDocument) error {
	var err error
	if dw.Skip, err = doc.GetPole("DocumentWhere.Limit.Skip").Int64(); err != nil {
		return err
	}
	if dw.Count, err = doc.GetPole("DocumentWhere.Limit.Count").Int64(); err != nil {
		return err
	}
	return nil
}

type DocumentWhereOrder struct {
	PoleName string
	ASC      bool
}

func (dw *DocumentWhereOrder) Load(doc IDocument) error {
	var err error
	if dw.PoleName, err = doc.GetPole("DocumentWhere.Order.PoleName").String(); err != nil {
		return err
	}
	var ASC string
	if ASC, err = doc.GetPole("DocumentWhere.Order.ASC").String(); err != nil {
		dw.ASC = ASC == "ASC"
		return err
	}
	return nil
}

type DocumentWhereSelectDocument struct {
	DocumentUID UUID
}

func (dw *DocumentWhereSelectDocument) Load(doc IDocument) error {
	var err error
	if dw.DocumentUID, err = doc.GetPole("DocumentWhere.SelectDocument.DocumentUID").DocumentLink(); err != nil {
		return err
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
	var err error
	if dw.PoleName, err = doc.GetPole("DocumentWhere.Compare.PoleName").String(); err != nil {
		return err
	}
	if dw.Operation, err = doc.GetPole("DocumentWhere.Compare.Operation").String(); err != nil {
		return err
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
	if ext := doc.GetPole("DocumentWhere.Compare.Extension"); !ext.IsNull() {
		if dw.ExtensionName, err = ext.String(); err != nil {
			return err
		}
	}
	return nil
}
