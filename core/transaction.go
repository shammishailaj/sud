package core

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/crazyprograms/sud/corebase"
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
func (tx *transaction) GetRecordsPoles(ConfigurationName string, RecordType string, poles []string, wheres []corebase.IRecordWhere, Access corebase.IAccess) (map[corebase.UUID]map[string]interface{}, error) {
	var err error
	var config *Configuration
	if config, err = tx.core.LoadConfiguration(ConfigurationName, Access); err != nil {
		return nil, err
	}
	var TypeInfo corebase.ITypeInfo
	if TypeInfo, err = config.GetTypeInfo(RecordType); err != nil {
		return nil, err
	}
	FreeAccess := TypeInfo.GetAccessType() == "Free"
	CheckAccess := TypeInfo.GetAccessType() == "Check"
	if !Access.CheckAccess(TypeInfo.GetAccessRead()) {
		return nil, &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "GetRecordsPoles", Name: RecordType}
	}
	var RecordUID corebase.UUID
	var RecordAccess string
	state := NewQueryState()
	state.AddPoleSQL(`"Record"."__RecordUID"`, &RecordUID)
	state.AddPoleSQL(`"Record"."RecordAccess"`, &RecordAccess)
	for poleName, pi := range config.GetPolesInfo(RecordType, poles) {
		if Access.CheckAccess(pi.GetAccessRead()) {
			state.AddPole(poleName, pi)
		}
	}
	state.AddWhere(`"Record"."RecordType"=` + state.AddParam(RecordType))
	if CheckAccess {
		for an, a := range Access.Users() {
			tablename := `"Access_` + strconv.Itoa(an) + `"`
			state.AddTable(tablename, `INNER JOIN Access AS `+tablename+` ON `+tablename+`."Access"="Record"."RecordAccess"`)
			state.AddWhere(tablename + `."Login"=` + state.AddParam(a))
		}
	}
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
	Records := map[corebase.UUID]map[string]interface{}{}
	for rows.Next() {
		if err = rows.Scan(values...); err != nil {
			return nil, err
		}
		if FreeAccess || Access.CheckAccess(RecordAccess) {
			doc := map[string]interface{}{}
			if err = state.SetRecordPoles(doc, values); err != nil {
				return nil, err
			}
			Records[RecordUID] = doc
		}
	}
	return Records, nil
}
func (tx *transaction) SetRecordPoles(ConfigurationName string, RecordUID corebase.UUID, Poles map[string]interface{}, Access corebase.IAccess) error {
	var err error
	var ok bool
	var pi corebase.IPoleInfo
	var tl []*PoleTableInfo
	var RecordType string
	var RecordAccess string
	if err = tx.QueryRow(`SELECT "RecordType", "RecordAccess" FROM "Record" WHERE "__RecordUID"=$1`, RecordUID).Scan(&RecordType, &RecordAccess); err != nil {
		return err
	}

	var config *Configuration
	if config, err = tx.core.LoadConfiguration(ConfigurationName, Access); err != nil {
		return err
	}
	var TypeInfo corebase.ITypeInfo
	if TypeInfo, err = config.GetTypeInfo(RecordType); err != nil {
		return err
	}
	FreeAccess := TypeInfo.GetAccessType() == "Free"
	if !Access.CheckAccess(TypeInfo.GetAccessSave()) {
		return &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "SetRecordPoles:RecordType", Name: RecordType, RecordUID: RecordUID.String()}
	}
	if !FreeAccess {
		if !Access.CheckAccess(RecordAccess) {
			return &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "SetRecordPoles:RecordAccess", Name: RecordType, RecordUID: RecordUID.String()}
		}
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
func (tx *transaction) NewRecord(ConfigurationName string, RecordType string, Poles map[string]interface{}, Access corebase.IAccess) (corebase.UUID, error) {
	var err error
	var config *Configuration
	if config, err = tx.core.LoadConfiguration(ConfigurationName, Access); err != nil {
		return "", err
	}
	var TypeInfo corebase.ITypeInfo
	if TypeInfo, err = config.GetTypeInfo(RecordType); err != nil {
		return "", err
	}
	if !Access.CheckAccess(TypeInfo.GetAccessNew()) {
		return "", &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "NewRecord:RecordType", Name: RecordType}
	}
	RecordUID := corebase.NewUUID()
	_, err = tx.Exec(`INSERT INTO "Record"("__RecordUID", "RecordType", "RecordAccess") VALUES ($1,$2,$3,$4)`, RecordUID, "Default")
	if err != nil {
		return "", err
	}
	err = tx.SetRecordPoles(ConfigurationName, RecordUID, Poles, Access)
	return RecordUID, err
}
func (tx *transaction) GetRecordAccess(ConfigurationName string, RecordUID corebase.UUID, Access corebase.IAccess) (string, error) {
	var err error
	var RecordType string
	var RecordAccess string
	if err = tx.QueryRow(`SELECT "RecordType", "RecordAccess" FROM "Record" WHERE "__RecordUID"=$1`, RecordUID).Scan(&RecordType, &RecordAccess); err != nil {
		return "", err
	}

	var config *Configuration
	if config, err = tx.core.LoadConfiguration(ConfigurationName, Access); err != nil {
		return "", err
	}
	var TypeInfo corebase.ITypeInfo
	if TypeInfo, err = config.GetTypeInfo(RecordType); err != nil {
		return "", err
	}
	FreeAccess := TypeInfo.GetAccessType() == "Free"
	if !Access.CheckAccess(TypeInfo.GetAccessRead()) {
		return "", &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "GetRecordAccess:RecordType", Name: RecordType, RecordUID: RecordUID.String()}
	}
	if !FreeAccess {
		if !Access.CheckAccess(RecordAccess) {
			return "", &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "GetRecordAccess:RecordAccess", Name: RecordType, RecordUID: RecordUID.String()}
		}
	}
	return RecordAccess, nil
}
func (tx *transaction) SetRecordAccess(ConfigurationName string, RecordUID corebase.UUID, NewAccess string, Access corebase.IAccess) error {
	if !Access.CheckAccess(NewAccess) {
		return &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "SetRecordAccess:Access", Name: NewAccess, RecordUID: RecordUID.String()}
	}
	var err error
	var RecordType string
	var RecordAccess string
	if err = tx.QueryRow(`SELECT "RecordType", "RecordAccess" FROM "Record" WHERE "__RecordUID"=$1`, RecordUID).Scan(&RecordType, &RecordAccess); err != nil {
		return err
	}

	var config *Configuration
	if config, err = tx.core.LoadConfiguration(ConfigurationName, Access); err != nil {
		return err
	}
	var TypeInfo corebase.ITypeInfo
	if TypeInfo, err = config.GetTypeInfo(RecordType); err != nil {
		return err
	}
	FreeAccess := TypeInfo.GetAccessType() == "Free"
	if !Access.CheckAccess(TypeInfo.GetAccessSave()) {
		return &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "SetRecordAccess:RecordType", Name: RecordType, RecordUID: RecordUID.String()}
	}
	if !FreeAccess {
		if !Access.CheckAccess(RecordAccess) {
			return &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "SetRecordAccess:RecordAccess", Name: RecordType, RecordUID: RecordUID.String()}
		}
	}
	if _, err = tx.Exec(`UPDATE "Record" SET "RecordAccess"=$2 WHERE "__RecordUID"=$1`, RecordUID, NewAccess); err != nil {
		return err
	}
	return nil
}
