package core

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
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
func (tx *transaction) GetDocumentsPoles(ConfigurationName string, DocumentType string, poles []string, wheres []IDocumentWhere) (map[string]map[string]interface{}, error) {
	var err error
	var config *Configuration
	if config, err = tx.core.LoadConfiguration(ConfigurationName); err != nil {
		return nil, err
	}
	var DocumentUID string
	state := NewQueryState()
	state.AddPoleSQL(`"Document"."__DocumentUID"`, &DocumentUID)
	for poleName, pi := range config.GetPolesInfo(DocumentType, poles) {
		state.AddPole(poleName, pi)
	}
	state.AddWhere(`"Document"."DocumentType"=` + state.AddParam(DocumentType))
	processWheres(config, DocumentType, state, wheres, tx)
	var SQLTop string
	SQLPoles := state.GenSQLPoles()
	SQLTables := state.GenSQLTables(`"Document"`)
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
	Documents := map[string]map[string]interface{}{}
	for rows.Next() {
		if err = rows.Scan(values...); err != nil {
			return nil, err
		}
		doc := map[string]interface{}{}
		if err = state.SetDocumentPoles(doc, values); err != nil {
			return nil, err
		}
		Documents[DocumentUID] = doc
	}
	return Documents, nil
}
func (tx *transaction) SetDocumentPoles(ConfigurationName string, DocumentUID string, Poles map[string]interface{}) error {
	var err error
	var ok bool
	var pi IPoleInfo
	var tl []*PoleTableInfo
	var DocumentType string
	var Readonly bool
	var DeleteDocument bool
	if err = tx.QueryRow(`SELECT "DocumentType", "Readonly", "DeleteDocument" FROM "Document" WHERE "__DocumentUID"=$1`, DocumentUID).Scan(&DocumentType, &Readonly, &DeleteDocument); err != nil {
		return err
	}

	var config *Configuration
	if config, err = tx.core.LoadConfiguration(ConfigurationName); err != nil {
		return err
	}
	TablePole := map[string][]*PoleTableInfo{}

	for pole, value := range Poles {
		if pi, err = config.GetPoleInfo(DocumentType, pole); err != nil {
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
		if err = tx.QueryRow(`SELECT COUNT(*) FROM "`+tableName+`" WHERE "__DocumentUID"=$1`, DocumentUID).Scan(&count); err != nil {
			return err
		}

		if count == 0 {
			if _, err = tx.Exec(`INSERT INTO "`+tableName+`"("__DocumentUID") VALUES ($1)`, DocumentUID); err != nil {
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
		values[num] = DocumentUID
		if _, err = tx.Exec(`UPDATE "`+tableName+`" SET `+strings.Join(poles, ", ")+` WHERE "__DocumentUID"=$`+strconv.Itoa(num+1), values...); err != nil {
			return err
		}

	}
	return nil
}
func (tx *transaction) NewDocument(ConfigurationName string, DocumentType string, Poles map[string]interface{}) (string, error) {
	var err error
	var config *Configuration
	if config, err = tx.core.LoadConfiguration(ConfigurationName); err != nil {
		return "", err
	}
	var ti ITypeInfo
	if ti, err = config.GetTypeInfo(DocumentType); err != nil {
		return "", err
	}
	if !ti.GetNew() {
		return "", errors.New("new document. access denied: " + DocumentType)
	}
	DouceumentUID := NewUUID().String()
	_, err = tx.Exec(`INSERT INTO "Document"("__DocumentUID", "DocumentType", "Readonly", "DeleteDocument") VALUES ($1,$2,$3,$4)`, DouceumentUID, DocumentType, false, false)
	if err != nil {
		return "", err
	}
	err = tx.SetDocumentPoles(ConfigurationName, DouceumentUID, Poles)
	return DouceumentUID, err
}
