package storage

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type document struct {
	configuration  *configuration
	documentUID    UUID
	documentType   string
	readOnly       bool
	deleteDocument bool
	poles          map[string]Object
}

func (doc *document) GetDocumentUID() UUID                { return doc.documentUID }
func (doc *document) GetDocumentType() string             { return doc.documentType }
func (doc *document) GetReadOnly() bool                   { return doc.readOnly }
func (doc *document) GetDeleteDocument() bool             { return doc.deleteDocument }
func (doc *document) GetPole(name string) Object          { return doc.poles[name] }
func (doc *document) SetDocumentType(documenttype string) { doc.documentType = documenttype }
func (doc *document) SetReadOnly(readonly bool)           { doc.readOnly = readonly }
func (doc *document) SetDeleteDocument(delete bool)       { doc.deleteDocument = delete }
func (doc *document) SetPole(name string, value Object)   { doc.poles[name] = value }
func (doc *document) GetPoleNames() []string {
	poles := make([]string, len(doc.poles))
	n := 0
	for polename := range doc.poles {
		poles[n] = polename
		n++
	}
	return poles
}
func (doc *document) GetConfiguration() *configuration { return doc.configuration }
func NewDocument(configuration *configuration, DocumentType string) *document {
	return &document{configuration: configuration, poles: map[string]Object{}}
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

func (server *Server) processWheres(conf *configuration, DocumentType string, state *queryState, wheres []IDocumentWhere) error {
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

			if info, err = conf.getPoleInfo(DocumentType, w.PoleName); err != nil {
				return err
			}
			pti := PoleTableInfo{}
			pti.FromPoleInfo(info)
			state.AddTable(pti.TableName, "JOIN ["+pti.TableName+"] ON (["+pti.TableName+"].[__DocumentUID] = [Document].[__DocumentUID])")
			if w.ASC {
				state.AddOrder("[" + pti.TableName + "].[" + pti.PoleName + "] ASC")
			} else {
				state.AddOrder("[" + pti.TableName + "].[" + pti.PoleName + "] DESC")
			}
		case DocumentWhereContainPole:
			if info, err = conf.getPoleInfo(DocumentType, w.PoleName); err != nil {
				return err
			}
			pti := PoleTableInfo{}
			pti.FromPoleInfo(info)
			Table := "table" + strconv.Itoa(len(state.tables))
			WithSQL := ""
			state.AddTable("Contain_"+info.GetPoleName(), " LEFT JOIN ["+pti.TableName+"] AS ["+Table+"] "+WithSQL+" ON (["+Table+"].[__DocumentUID] = [Document].[__DocumentUID])")
			state.AddWhere("[" + Table + "].[" + pti.PoleName + "] IS NOT NULL")

		case DocumentWhereNotContainPole:
			if info, err = conf.getPoleInfo(DocumentType, w.PoleName); err != nil {
				return err
			}
			pti := PoleTableInfo{}
			pti.FromPoleInfo(info)
			Table := "table" + strconv.Itoa(len(state.tables))
			WithSQL := ""
			if w.TabLock {
				WithSQL = "WITH (UPDLOCK)"
			}
			if w.InTableName != "" {
				d := new([]byte)
				state.AddPoleSQL("["+Table+"].[__DocumentUID] AS ["+w.InTableName+"]", d)
			}
			state.AddTable("Contain_"+info.GetPoleName(), " LEFT JOIN ["+pti.TableName+"] AS ["+Table+"] "+WithSQL+" ON (["+Table+"].[__DocumentUID] = [Document].[__DocumentUID])")
			state.AddWhere("[" + Table + "].[" + pti.PoleName + "] IS NULL")
		case DocumentWhereCompare:
			if info, err = conf.getPoleInfo(DocumentType, w.PoleName); err != nil {
				return err
			}
			pti := PoleTableInfo{}
			pti.FromPoleInfo(info)
			state.AddTable(pti.TableName, "JOIN ["+pti.TableName+"] ON (["+pti.TableName+"].[__DocumentUID] = [Document].[__DocumentUID])")
			//String VarName = State.AddParam(where.getValue(m_Configuration));
			Value := w.Value
			switch w.Operation {
			case "Equally":
				{
					if Value.IsNull() {
						state.AddWhere("[" + pti.TableName + "].[" + pti.PoleName + "] IS NULL ")
					} else {
						state.AddWhere("[" + pti.TableName + "].[" + pti.PoleName + "] = " + state.AddParam(w.Value))
					}
				}
			case "Not_Equally":
				{
					if Value.IsNull() {
						state.AddWhere("[" + pti.TableName + "].[" + pti.PoleName + "] IS NOT NULL ")
					} else {
						state.AddWhere("[" + pti.TableName + "].[" + pti.PoleName + "] <> " + state.AddParam(w.Value))
					}
				}

			case "Less":
				state.AddWhere("[" + pti.TableName + "].[" + pti.PoleName + "] < " + state.AddParam(w.Value))

			case "More":
				state.AddWhere("[" + pti.TableName + "].[" + pti.PoleName + "] > " + state.AddParam(w.Value))

			case "NotLess":
				state.AddWhere("[" + pti.TableName + "].[" + pti.PoleName + "] >= " + state.AddParam(w.Value))
			case "NotMore":
				state.AddWhere("[" + pti.TableName + "].[" + pti.PoleName + "] <= " + state.AddParam(w.Value))
			default:
				return errors.New(w.Operation + " not implemented")
			}
		}
	}
	return nil
}
func (server *Server) GetDocuments(TransactionUID string, ConfigurationName string, DocumentType string, poles []string, wheres []IDocumentWhere) ([]IDocument, error) {
	var err error
	//var ok bool
	var tx *transaction
	var config *configuration
	if tx, err = server.getTransaction(TransactionUID); err != nil {
		return nil, err
	}
	if config, err = server.LoadConfiguration(ConfigurationName); err != nil {
		return nil, err
	}
	var DocumentUID []byte
	state := NewQueryState()
	state.AddPoleSQL("[dbo].[Document].[__DocumentUID]", &DocumentUID)
	for poleName, pi := range config.getPolesInfo(DocumentType, poles) {
		state.AddPole(poleName, pi)
	}
	var SQLTop string
	SQLPoles := state.GenSQLPoles()
	SQLTables := state.GenSQLTables("[dbo].[Document]")
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
		doc := NewDocument(config, DocumentType)
		doc.documentUID = UUID(DocumentUID[0:16])
		state.SetDocumentPoles(doc, values)
		documents = append(documents, doc)
	}
	//e, _ := rows.Columns()
	//fmt.Println(SQL, e)
	return documents, nil
}
