package core

import (
	"errors"
	"strconv"

	"github.com/crazyprograms/sud/client"
)

func pwDocumentWhereLimit(conf *Configuration, DocumentType string, state *queryState, w *DocumentWhereLimit, tx *transaction) error {
	if w.Skip > 0 {
		state.Skip = int(w.Skip)
	}
	if w.Count > 0 {
		state.Count = int(w.Count)
	}
	return nil
}
func pwDocumentWhereOrder(conf *Configuration, DocumentType string, state *queryState, w *DocumentWhereOrder, tx *transaction) error {
	var err error
	var info IPoleInfo
	if info, err = conf.GetPoleInfo(DocumentType, w.PoleName); err != nil {
		return err
	}
	pti := PoleTableInfo{}
	pti.FromPoleInfo(info)
	state.AddTable(pti.TableName, `JOIN "`+pti.TableName+`" ON ("`+pti.TableName+`"."__DocumentUID" = "Document"."__DocumentUID")`)
	if w.ASC {
		state.AddOrder(`"` + pti.TableName + `"."` + pti.PoleName + `" ASC`)
	} else {
		state.AddOrder(`"` + pti.TableName + `"."` + pti.PoleName + `" DESC`)
	}
	return nil
}
func pwDocumentWhereContainPole(conf *Configuration, DocumentType string, state *queryState, w *DocumentWhereContainPole, tx *transaction) error {
	var err error
	var info IPoleInfo
	if info, err = conf.GetPoleInfo(DocumentType, w.PoleName); err != nil {
		return err
	}
	pti := PoleTableInfo{}
	pti.FromPoleInfo(info)
	Table := `table` + strconv.Itoa(len(state.tables))
	WithSQL := ``
	state.AddTable(`Contain_`+info.GetPoleName(), ` LEFT JOIN "`+pti.TableName+`" AS "`+Table+`" `+WithSQL+` ON ("`+Table+`"."__DocumentUID" = "Document"."__DocumentUID")`)
	state.AddWhere(`"` + Table + `"."` + pti.PoleName + `" IS NOT NULL`)
	return nil
}
func pwDocumentWhereNotContainPole(conf *Configuration, DocumentType string, state *queryState, w *DocumentWhereNotContainPole, tx *transaction) error {
	var err error
	var info IPoleInfo
	if info, err = conf.GetPoleInfo(DocumentType, w.PoleName); err != nil {
		return err
	}
	pti := PoleTableInfo{}
	pti.FromPoleInfo(info)
	Table := `table` + strconv.Itoa(len(state.tables))
	WithSQL := ``
	if w.TabLock {
		WithSQL = `WITH (UPDLOCK)`
	}
	if w.InTableName != `` {
		d := new([]byte)
		state.AddPoleSQL(`"`+Table+`"."__DocumentUID" AS "`+w.InTableName+`"`, d)
	}
	state.AddTable(`Contain_`+info.GetPoleName(), ` LEFT JOIN "`+pti.TableName+`" AS "`+Table+`" `+WithSQL+` ON ("`+Table+`"."__DocumentUID" = "Document"."__DocumentUID")`)
	state.AddWhere(`"` + Table + `"."` + pti.PoleName + `" IS NULL`)
	return nil
}
func pwDocumentWhereCompare(conf *Configuration, DocumentType string, state *queryState, w *DocumentWhereCompare, tx *transaction) error {
	var err error
	var info IPoleInfo
	if info, err = conf.GetPoleInfo(DocumentType, w.PoleName); err != nil {
		return err
	}
	pti := PoleTableInfo{}
	pti.FromPoleInfo(info)
	state.AddTable(pti.TableName, `JOIN "`+pti.TableName+`" ON ("`+pti.TableName+`"."__DocumentUID" = "Document"."__DocumentUID")`)
	//String VarName = State.AddParam(where.getValue(m_Configuration));
	Value := w.Value
	switch w.Operation {
	case `Equally`:
		{
			if IsNull(Value) {
				state.AddWhere(`"` + pti.TableName + `"."` + pti.PoleName + `" IS NULL `)
			} else {
				state.AddWhere(`"` + pti.TableName + `"."` + pti.PoleName + `" = ` + state.AddParam(w.Value))
			}
		}
	case `Not_Equally`:
		{
			if IsNull(Value) {
				state.AddWhere(`"` + pti.TableName + `"."` + pti.PoleName + `" IS NOT NULL `)
			} else {
				state.AddWhere(`"` + pti.TableName + `"."` + pti.PoleName + `" <> ` + state.AddParam(w.Value))
			}
		}

	case `Less`:
		state.AddWhere(`"` + pti.TableName + `"."` + pti.PoleName + `" < ` + state.AddParam(w.Value))

	case `More`:
		state.AddWhere(`"` + pti.TableName + `"."` + pti.PoleName + `" > ` + state.AddParam(w.Value))

	case `NotLess`:
		state.AddWhere(`"` + pti.TableName + `"."` + pti.PoleName + `" >= ` + state.AddParam(w.Value))
	case `NotMore`:
		state.AddWhere(`"` + pti.TableName + `"."` + pti.PoleName + `" <= ` + state.AddParam(w.Value))
	default:
		return errors.New(w.Operation + ` not implemented`)
	}
	return nil
}
func processWheres(conf *Configuration, DocumentType string, state *queryState, wheres []client.IDocumentWhere, tx *transaction) error {
	for _, where := range wheres {
		switch w := where.(type) {
		case DocumentWhereLimit:
			return pwDocumentWhereLimit(conf, DocumentType, state, &w, tx)
		case *DocumentWhereLimit:
			return pwDocumentWhereLimit(conf, DocumentType, state, w, tx)
		case DocumentWhereOrder:
			return pwDocumentWhereOrder(conf, DocumentType, state, &w, tx)
		case *DocumentWhereOrder:
			return pwDocumentWhereOrder(conf, DocumentType, state, w, tx)
		case DocumentWhereContainPole:
			return pwDocumentWhereContainPole(conf, DocumentType, state, &w, tx)
		case *DocumentWhereContainPole:
			return pwDocumentWhereContainPole(conf, DocumentType, state, w, tx)
		case DocumentWhereNotContainPole:
			return pwDocumentWhereNotContainPole(conf, DocumentType, state, &w, tx)
		case *DocumentWhereNotContainPole:
			return pwDocumentWhereNotContainPole(conf, DocumentType, state, w, tx)
		case DocumentWhereCompare:
			return pwDocumentWhereCompare(conf, DocumentType, state, &w, tx)
		case *DocumentWhereCompare:
			return pwDocumentWhereCompare(conf, DocumentType, state, w, tx)
		}
	}
	return nil
}
