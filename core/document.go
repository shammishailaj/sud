package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type document struct {
	configuration  *Configuration
	documentUID    UUID
	documentType   string
	readOnly       bool
	deleteDocument bool
	editpoles      map[string]bool
	poles          map[string]Object
}

func (doc *document) GetDocumentUID() UUID       { return doc.documentUID }
func (doc *document) GetDocumentType() string    { return doc.documentType }
func (doc *document) GetReadOnly() bool          { return doc.readOnly }
func (doc *document) GetDeleteDocument() bool    { return doc.deleteDocument }
func (doc *document) GetPole(name string) Object { return doc.poles[name] }
func (doc *document) GetPoleValue(name string) interface{} {
	obj := doc.poles[name]
	return (&obj).Get()
}
func (doc *document) SetDocumentType(documenttype string) { doc.documentType = documenttype }
func (doc *document) SetReadOnly(readonly bool)           { doc.readOnly = readonly }
func (doc *document) SetDeleteDocument(delete bool)       { doc.deleteDocument = delete }
func (doc *document) SetPoleValue(name string, value interface{}) error {
	return doc.SetPole(name, NewObject(value))
}
func (doc *document) SetPole(name string, value Object) error {
	var err error
	info, err := doc.configuration.GetPoleInfo(doc.documentType, name)
	if err != nil {
		return err
	}
	checker := info.GetChecker()
	if checker != nil {
		if err = checker.CheckPoleValue(value); err != nil {
			return err
		}
	}
	doc.editpoles[name] = true
	doc.poles[name] = value
	return nil
}
func (doc *document) GetPoleNames() []string {
	poles := make([]string, len(doc.poles))
	n := 0
	for polename := range doc.poles {
		poles[n] = polename
		n++
	}
	return poles
}
func (doc *document) GetConfiguration() *Configuration { return doc.configuration }
func (tx *transaction) NewDocument(ConfigurationName string, DocumentType string) (*document, error) {
	var err error
	var config *Configuration
	if config, err = tx.core.LoadConfiguration(ConfigurationName); err != nil {
		return nil, err
	}
	doc := &document{configuration: config, documentType: DocumentType, poles: map[string]Object{}, editpoles: map[string]bool{}}
	doc.documentUID = NewUUID()
	_, err = tx.Exec(`INSERT INTO "Document"("__DocumentUID", "DocumentType", "Readonly", "DeleteDocument") VALUES ($1,$2,$3,$4)`, doc.documentUID, DocumentType, false, false)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

/*
type IDocument interface {
	GetDocumentUID() uuid.UUID
	GetDocumentType() string
	GetReadOnly() bool
	GetDeleteDocument() bool
	GetPole(name string) Object
	SetDocumentUID(UID uuid.UUID)
	SetDocumentType(documenttype string)
	SetReadOnly(readonly bool)
	SetDeleteDocument(delete bool)
	SetPole(name string, value Object)
	GetPoleNames() bool
	GetConfiguration() *Configuration
}*/
func printValue(pval *interface{}) {
	switch v := (*pval).(type) {
	case nil:
		fmt.Print("NULL")
	case bool:
		if v {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
	case []byte:
		fmt.Print(string(v))
	case time.Time:
		fmt.Print(v.Format("2006-01-02 15:04:05.999"))
	default:
		fmt.Print(v)
	}
}

func (tx *transaction) processWheres(conf *Configuration, DocumentType string, state *queryState, wheres []IDocumentWhere) error {
	var where IDocumentWhere
	var info IPoleInfo
	var err error
	for where = range wheres {
		switch w := where.(type) {
		case DocumentWhereLimit:
			if w.Skip > 0 {
				state.Skip = int(w.Skip)
			}
			if w.Count > 0 {
				state.Count = int(w.Count)
			}
		case DocumentWhereOrder:

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
		case DocumentWhereContainPole:
			if info, err = conf.GetPoleInfo(DocumentType, w.PoleName); err != nil {
				return err
			}
			pti := PoleTableInfo{}
			pti.FromPoleInfo(info)
			Table := `table` + strconv.Itoa(len(state.tables))
			WithSQL := ``
			state.AddTable(`Contain_`+info.GetPoleName(), ` LEFT JOIN "`+pti.TableName+`" AS "`+Table+`" `+WithSQL+` ON ("`+Table+`"."__DocumentUID" = "Document"."__DocumentUID")`)
			state.AddWhere(`"` + Table + `"."` + pti.PoleName + `" IS NOT NULL`)

		case DocumentWhereNotContainPole:
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
		case DocumentWhereCompare:
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
					if Value.IsNull() {
						state.AddWhere(`"` + pti.TableName + `"."` + pti.PoleName + `" IS NULL `)
					} else {
						state.AddWhere(`"` + pti.TableName + `"."` + pti.PoleName + `" = ` + state.AddParam(w.Value))
					}
				}
			case `Not_Equally`:
				{
					if Value.IsNull() {
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
		}
	}
	return nil
}
func (core *Core) GetDocuments(TransactionUID string, ConfigurationName string, DocumentType string, poles []string, wheres []IDocumentWhere) ([]IDocument, error) {
	var err error
	var tx *transaction
	if tx, err = core.getTransaction(TransactionUID); err != nil {
		return nil, err
	}
	return tx.GetDocuments(ConfigurationName, DocumentType, poles, wheres)
}
func (tx *transaction) GetDocuments(ConfigurationName string, DocumentType string, poles []string, wheres []IDocumentWhere) ([]IDocument, error) {
	var err error
	var config *Configuration
	if config, err = tx.core.LoadConfiguration(ConfigurationName); err != nil {
		return nil, err
	}
	var DocumentUID []byte
	state := NewQueryState()
	state.AddPoleSQL(`"Document"."__DocumentUID"`, &DocumentUID)
	for poleName, pi := range config.GetPolesInfo(DocumentType, poles) {
		state.AddPole(poleName, pi)
	}
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
	documents := make([]IDocument, 0, 16)
	for rows.Next() {
		if err = rows.Scan(values...); err != nil {
			return nil, err
		}
		doc := &document{configuration: config, documentType: DocumentType, poles: map[string]Object{}, editpoles: map[string]bool{}}
		doc.documentUID = UUID(DocumentUID[0:16])
		state.SetDocumentPoles(doc, values)
		documents = append(documents, doc)
	}
	//e, _ := rows.Columns()
	//fmt.Println(SQL, e)
	return documents, nil
}

func (core *Core) SaveDocument(TransactionUID string, doc *document) error {
	var err error
	var ok bool
	var pi IPoleInfo
	var tl []*PoleTableInfo
	var tx *transaction
	if tx, err = core.getTransaction(TransactionUID); err != nil {
		return err
	}
	TablePole := map[string][]*PoleTableInfo{}
	for pole := range doc.editpoles {
		if pi, err = doc.configuration.GetPoleInfo(doc.documentType, pole); err != nil {
			return err
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
		var count int
		tx.QueryRow(`SELECT COUNT(*) FROM "`+tableName+`" WHERE "__DocumentUID"=$1`, doc.documentUID).Scan(&count)
		if count == 0 {
			if _, err = tx.Exec(`INSERT INTO "`+tableName+`"("__DocumentUID") VALUES ($1)`, doc.documentUID); err != nil {
				return err
			}
		}
		num := len(tl)
		poles := make([]string, num, num)
		values := make([]interface{}, num+1, num+1)
		for i := 0; i < num; i++ {
			poles[i] = `"` + tl[i].PoleName + `"= $` + strconv.Itoa(i+1)
			values[i] = doc.poles[tl[i].PoleInfo.GetPoleName()]
		}
		values[num] = doc.documentUID
		if _, err = tx.Exec(`UPDATE "`+tableName+`" SET `+strings.Join(poles, ", ")+` WHERE "__DocumentUID"=$`+strconv.Itoa(num+1), values...); err != nil {
			return err
		}

	}
	return nil
}
