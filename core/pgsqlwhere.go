package core

import (
	"errors"
	"strconv"

	"github.com/crazyprograms/sud/corebase"
)

func pwRecordWhereLimit(conf *Configuration, RecordType string, state *queryState, w *corebase.RecordWhereLimit, tx *transaction) error {
	if w.Skip > 0 {
		state.Skip = int(w.Skip)
	}
	if w.Count > 0 {
		state.Count = int(w.Count)
	}
	return nil
}
func pwRecordWhereOrder(conf *Configuration, RecordType string, state *queryState, w *corebase.RecordWhereOrder, tx *transaction) error {
	var err error
	var info corebase.IPoleInfo
	if info, err = conf.GetPoleInfo(RecordType, w.PoleName); err != nil {
		return err
	}
	pti := PoleTableInfo{}
	pti.FromPoleInfo(info)
	state.AddTable(pti.TableName, `JOIN "`+pti.TableName+`" ON ("`+pti.TableName+`"."__RecordUID" = "Record"."__RecordUID")`)
	if w.ASC {
		state.AddOrder(`"` + pti.TableName + `"."` + pti.PoleName + `" ASC`)
	} else {
		state.AddOrder(`"` + pti.TableName + `"."` + pti.PoleName + `" DESC`)
	}
	return nil
}
func pwRecordWhereContainPole(conf *Configuration, RecordType string, state *queryState, w *corebase.RecordWhereContainPole, tx *transaction) error {
	var err error
	var info corebase.IPoleInfo
	if info, err = conf.GetPoleInfo(RecordType, w.PoleName); err != nil {
		return err
	}
	pti := PoleTableInfo{}
	pti.FromPoleInfo(info)
	Table := `table` + strconv.Itoa(len(state.tables))
	WithSQL := ``
	state.AddTable(`Contain_`+info.GetPoleName(), ` LEFT JOIN "`+pti.TableName+`" AS "`+Table+`" `+WithSQL+` ON ("`+Table+`"."__RecordUID" = "Record"."__RecordUID")`)
	state.AddWhere(`"` + Table + `"."` + pti.PoleName + `" IS NOT NULL`)
	return nil
}
func pwRecordWhereNotContainPole(conf *Configuration, RecordType string, state *queryState, w *corebase.RecordWhereNotContainPole, tx *transaction) error {
	var err error
	var info corebase.IPoleInfo
	if info, err = conf.GetPoleInfo(RecordType, w.PoleName); err != nil {
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
		state.AddPoleSQL(`"`+Table+`"."__RecordUID" AS "`+w.InTableName+`"`, d)
	}
	state.AddTable(`Contain_`+info.GetPoleName(), ` LEFT JOIN "`+pti.TableName+`" AS "`+Table+`" `+WithSQL+` ON ("`+Table+`"."__RecordUID" = "Record"."__RecordUID")`)
	state.AddWhere(`"` + Table + `"."` + pti.PoleName + `" IS NULL`)
	return nil
}
func pwRecordWhereCompare(conf *Configuration, RecordType string, state *queryState, w *corebase.RecordWhereCompare, tx *transaction) error {
	var err error
	var info corebase.IPoleInfo
	if info, err = conf.GetPoleInfo(RecordType, w.PoleName); err != nil {
		return err
	}
	pti := PoleTableInfo{}
	pti.FromPoleInfo(info)
	state.AddTable(pti.TableName, `JOIN "`+pti.TableName+`" ON ("`+pti.TableName+`"."__RecordUID" = "Record"."__RecordUID")`)
	//String VarName = State.AddParam(where.getValue(m_Configuration));
	Value := w.Value
	switch w.Operation {
	case `Equally`:
		{
			if corebase.IsNull(Value) {
				state.AddWhere(`"` + pti.TableName + `"."` + pti.PoleName + `" IS NULL `)
			} else {
				state.AddWhere(`"` + pti.TableName + `"."` + pti.PoleName + `" = ` + state.AddParam(w.Value))
			}
		}
	case `Not_Equally`:
		{
			if corebase.IsNull(Value) {
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
func processWheres(conf *Configuration, RecordType string, state *queryState, wheres []corebase.IRecordWhere, tx *transaction) error {
	for _, where := range wheres {
		switch w := where.(type) {
		case *corebase.RecordWhereLimit:
			return pwRecordWhereLimit(conf, RecordType, state, w, tx)
		case *corebase.RecordWhereOrder:
			return pwRecordWhereOrder(conf, RecordType, state, w, tx)
		case *corebase.RecordWhereContainPole:
			return pwRecordWhereContainPole(conf, RecordType, state, w, tx)
		case *corebase.RecordWhereNotContainPole:
			return pwRecordWhereNotContainPole(conf, RecordType, state, w, tx)
		case *corebase.RecordWhereCompare:
			return pwRecordWhereCompare(conf, RecordType, state, w, tx)
		}
	}
	return nil
}
