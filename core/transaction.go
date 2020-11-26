package core

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/shammishailaj/sud/corebase"
)

type transaction struct {
	core *Core
	tx   *sql.Tx
}

func (t *transaction) Commit() error {
	return t.tx.Commit()
}
func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}
func (t *transaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}
func (t *transaction) Prepare(query string) (*sql.Stmt, error) {
	return t.tx.Prepare(query)
}
func (t *transaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(query, args...)
}
func (t *transaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(query, args...)
}
func (tx *transaction) GetRecordsPoles(ConfigurationName string, RecordType string, poles []string, wheres []corebase.IRecordWhere) (map[string]map[string]interface{}, error) {
	var err error
	var config *Configuration
	if config, err = tx.core.LoadConfiguration(ConfigurationName); err != nil {
		return nil, err
	}
	var RecordUID string
	state := NewQueryState()
	state.AddPoleSQL(`"Record"."__RecordUID"`, &RecordUID)
	for poleName, pi := range config.GetPolesInfo(RecordType, poles) {
		state.AddPole(poleName, pi)
	}
	state.AddWhere(`"Record"."RecordType"=` + state.AddParam(RecordType))
	processWheres(config, RecordType, state, wheres, tx)
	var SQLTop string
	SQLPoles := state.GenSQLPoles()
	SQLTables := state.GenSQLTables(`"Record"`)
	if SQLTop, err = state.GenSQLTop(); err != nil {
		return nil, err
	}
	SQLWheres := state.GenSQLWheres()
	SQLOrders := state.GenSQLOrder()
	SQL := "SELECT " + SQLTop + SQLPoles + SQLTables + SQLWheres + SQLOrders
	params := state.GetParams()
	rows, err := tx.Query(SQL, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	values := state.GetPoleValueArray()
	Records := map[string]map[string]interface{}{}
	for rows.Next() {
		if err = rows.Scan(values...); err != nil {
			return nil, err
		}
		doc := map[string]interface{}{}
		if err = state.SetRecordPoles(doc, values); err != nil {
			return nil, err
		}
		Records[RecordUID] = doc
	}
	return Records, nil
}
func (tx *transaction) SetRecordPoles(ConfigurationName string, RecordUID string, Poles map[string]interface{}) error {
	var err error
	var ok bool
	var pi corebase.IPoleInfo
	var tl []*PoleTableInfo
	var RecordType string
	var Readonly bool
	var DeleteRecord bool
	if err = tx.QueryRow(`SELECT "RecordType", "Readonly", "DeleteRecord" FROM "Record" WHERE "__RecordUID"=$1`, RecordUID).Scan(&RecordType, &Readonly, &DeleteRecord); err != nil {
		return err
	}

	var config *Configuration
	if config, err = tx.core.LoadConfiguration(ConfigurationName); err != nil {
		return err
	}
	TablePole := map[string][]*PoleTableInfo{}

	for pole, value := range Poles {
		if pi, err = config.GetPoleInfo(RecordType, pole); err != nil {
			return err
		}
		checker := pi.GetChecker()
		if checker != nil {
			if err = checker.CheckPoleValue(value); err != nil {
				return err
			}
		}
		pti := &PoleTableInfo{}
		pti.FromPoleInfo(pi)
		if tl, ok = TablePole[pti.TableName]; !ok {
			tl = []*PoleTableInfo{}
			TablePole[pti.TableName] = tl
		}
		TablePole[pti.TableName] = append(tl, pti)
	}

	for tableName, tl := range TablePole {
		var count int64
		if err = tx.QueryRow(`SELECT COUNT(*) FROM "`+tableName+`" WHERE "__RecordUID"=$1`, RecordUID).Scan(&count); err != nil {
			return err
		}

		if count == 0 {
			if _, err = tx.Exec(`INSERT INTO "`+tableName+`"("__RecordUID") VALUES ($1)`, RecordUID); err != nil {
				return err
			}
		}
		num := len(tl)
		poles := make([]string, num, num)
		values := make([]interface{}, num+1, num+1)
		for i := 0; i < num; i++ {
			poles[i] = `"` + tl[i].PoleName + `"= $` + strconv.Itoa(i+1)
			values[i] = Poles[tl[i].PoleInfo.GetPoleName()]
		}
		values[num] = RecordUID
		if _, err = tx.Exec(`UPDATE "`+tableName+`" SET `+strings.Join(poles, ", ")+` WHERE "__RecordUID"=$`+strconv.Itoa(num+1), values...); err != nil {
			return err
		}

	}
	return nil
}
func (tx *transaction) NewRecord(ConfigurationName string, RecordType string, Poles map[string]interface{}) (string, error) {
	var err error
	var config *Configuration
	if config, err = tx.core.LoadConfiguration(ConfigurationName); err != nil {
		return "", err
	}
	var ti corebase.ITypeInfo
	if ti, err = config.GetTypeInfo(RecordType); err != nil {
		return "", err
	}
	if !ti.GetNew() {
		return "", errors.New("new record. access denied: " + RecordType)
	}
	DouceumentUID := corebase.NewUUID().String()
	_, err = tx.Exec(`INSERT INTO "Record"("__RecordUID", "RecordType", "Readonly", "DeleteRecord") VALUES ($1,$2,$3,$4)`, DouceumentUID, RecordType, false, false)
	if err != nil {
		return "", err
	}
	err = tx.SetRecordPoles(ConfigurationName, DouceumentUID, Poles)
	return DouceumentUID, err
}
