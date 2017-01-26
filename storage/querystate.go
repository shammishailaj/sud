package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type queryState struct {
	Skip        int
	Count       int
	num         int
	tables      map[string]string
	polesSQL    []string
	polesSQLink []interface{}
	poles       map[string]*ptiInfo
	orders      []string
	wheres      []string
	params      []interface{}
}
type ptiInfo struct {
	PoleTableInfo
	SQL string
	N   int
}

func NewQueryState() *queryState {
	return &queryState{tables: map[string]string{}, polesSQL: make([]string, 0, 10), poles: map[string]*ptiInfo{}, params: make([]interface{}, 0, 10), wheres: make([]string, 0, 10), orders: make([]string, 0, 3)}
}

func (state *queryState) AddPole(poleName string, pi IPoleInfo) {
	pti := &ptiInfo{}
	pti.FromPoleInfo(pi)
	pti.SQL = "[dbo].[" + pti.TableName + "].[" + pti.PoleName + "]"
	pti.N = -1
	state.num++
	state.poles[pi.GetPoleName()] = pti
	state.AddTable(pti.TableName, " LEFT JOIN [dbo].["+pti.TableName+"] ON [dbo].["+pti.TableName+"].[__DocumentUID] = [dbo].[Document].[__DocumentUID]")
}
func (state *queryState) AddPoleSQL(poleSQL string, link interface{}) int {
	state.polesSQL = append(state.polesSQL, poleSQL)
	state.polesSQLink = append(state.polesSQLink, link)
	return len(state.polesSQL) - 1
}
func (state *queryState) AddTable(tableName string, SQL string) {
	if _, ok := state.tables[tableName]; !ok {
		state.tables[tableName] = SQL
	}
}
func (state *queryState) AddParam(value interface{}) string {
	state.params = append(state.params, value)
	return "$" + strconv.Itoa(len(state.params))
}
func (state *queryState) AddParamObject(value Object) (string, error) {
	if value.IsNull() {
		return state.AddParam(nil), nil
	}
	switch value.Type() {
	case "BooleanValue":
		if v, err := value.Boolean(); err != nil {
			return "", err
		} else {
			return state.AddParam(v), nil
		}
	case "StringValue":
		if v, err := value.String(); err != nil {
			return "", err
		} else {
			return state.AddParam(v), nil
		}
	case "Int64Value":
		if v, err := value.Int64(); err != nil {
			return "", err
		} else {
			return state.AddParam(v), nil
		}
	case "DateValue":
		if v, err := value.Date(); err != nil {
			return "", err
		} else {
			return state.AddParam(v), nil
		}
	case "DateTimeValue":
		if v, err := value.DateTime(); err != nil {
			return "", err
		} else {
			return state.AddParam(v), nil
		}

	case "DocumentLinkValue":
		dl, err := value.DocumentLink()
		if err != nil {
			return "", err
		}
		if len(dl) != 16 {
			return "", errors.New("convert type error DocumentLinkValue")
		}
		return state.AddParam(dl[0:16]), nil
	default:
		return "", errors.New("convert type error " + value.Type())
	}
}
func (state *queryState) AddOrder(order string) {
	state.orders = append(state.orders, order)
}
func (state *queryState) AddWhere(where string) {
	state.wheres = append(state.wheres, where)
}
func (state *queryState) GetParams() []interface{} {
	return state.params
}

func (state *queryState) GenSQLTop() (string, error) {
	if len(state.orders) != 0 {
		return "", nil
	}
	if state.Skip == 0 && state.Count == 0 {
		return "", nil
	}
	if state.Skip == 0 {
		return " TOP " + strconv.Itoa(state.Count) + " ", nil
	}
	return "", errors.New("Skip items can only be ordered in the table")
}
func (state *queryState) GenSQLOrder() string {
	if len(state.orders) == 0 {
		return ""
	}
	OffsetSQL := ""
	if state.Count == 0 && state.Skip > 0 {
		OffsetSQL += " OFFSET " + strconv.Itoa(state.Skip) + " ROWS"
	} else if state.Count > 0 && state.Skip > 0 {
		OffsetSQL += " OFFSET " + strconv.Itoa(state.Skip) + " ROWS FETCH NEXT " + strconv.Itoa(state.Count) + " ROWS ONLY"
	}
	return " ORDER BY " + strings.Join(state.orders, " ") + OffsetSQL
}
func (state *queryState) GenSQLTables(baseTable string) string {
	SQLTables := baseTable
	for _, SQL := range state.tables {
		SQLTables += SQL
	}
	return " FROM " + SQLTables
}
func (state *queryState) GenSQLWheres() string {
	SQLWheres := ""
	if len(state.wheres) == 0 {
		return ""
	}
	for _, SQL := range state.wheres {
		SQLWheres += SQL
	}
	return " WHERE " + SQLWheres
}
func (state *queryState) GenSQLPoles() string {
	poles := make([]string, len(state.polesSQL)+len(state.poles))
	copy(poles[0:len(state.polesSQL)], state.polesSQL)
	N := len(state.polesSQL)
	for _, pti := range state.poles {
		poles[N] = pti.SQL
		pti.N = N
		N++
	}
	state.num = N
	return strings.Join(poles, ", ")
}

func (state *queryState) GetPoleValueArray() [](interface{}) {
	values := make([](interface{}), state.num)
	copy(values[0:len(state.polesSQLink)], state.polesSQLink)
	NumTime := 0
	NumBoolean := 0
	NumString := 0
	NumInt64 := 0
	for _, pti := range state.poles {
		switch pti.PoleInfo.GetPoleType() {
		case "BooleanValue":
			NumBoolean++
		case "StringValue":
			NumString++
		case "DateValue":
			NumTime++
		case "Int64Value":
			NumInt64++
		default:
			fmt.Println(pti.PoleInfo.GetPoleType())
			panic("")
		}
	}
	ValueTime := make([]*time.Time, NumTime)
	ValueBool := make([]sql.NullBool, NumString)
	ValueStrings := make([]sql.NullString, NumString)
	ValueInt64 := make([]sql.NullInt64, NumInt64)
	NumTime = 0
	NumBoolean = 0
	NumString = 0
	NumInt64 = 0
	for _, pti := range state.poles {
		switch pti.PoleInfo.GetPoleType() {
		case "BooleanValue":
			values[pti.N] = &ValueBool[NumBoolean]
			NumBoolean++
		case "StringValue":
			values[pti.N] = &ValueStrings[NumString]
			NumString++
		case "DateValue":
			values[pti.N] = &ValueTime[NumTime]
			NumTime++
		case "Int64Value":
			values[pti.N] = &ValueInt64[NumInt64]
			NumInt64++
		default:
			fmt.Println(pti.PoleInfo.GetPoleType())
			panic("")
		}
	}
	return values
}
func (state *queryState) SetDocumentPoles(doc *document, values [](interface{})) error {
	for poleName, pti := range state.poles {
		o := Object{}
		switch pti.PoleInfo.GetPoleType() {
		case "BooleanValue":
			v, ok := values[pti.N].(*sql.NullBool)
			if !ok {
				return errors.New("pole read error " + poleName)
			}
			if !v.Valid {
				o.SetNull()
			} else {
				o.SetBoolean(v.Bool)
			}
		case "StringValue":
			v, ok := values[pti.N].(*sql.NullString)
			if !ok {
				return errors.New("pole read error " + poleName)
			}
			if !v.Valid {
				o.SetNull()
			} else {
				o.SetString(v.String)
			}
		case "DateValue":
			v, ok := values[pti.N].(**time.Time)
			if !ok {
				return errors.New("pole read error " + poleName)
			}
			if *v == nil {
				o.SetNull()
			} else {
				o.SetDate(*(*v))
			}
			//o.SetDate(**v)
			//fmt.Println(v, *v, v == nil, *v == nil, ok)
			/*if !v.Valid */
		case "Int64Value":
			v, ok := values[pti.N].(*sql.NullInt64)
			if !ok {
				return errors.New("pole read error " + poleName)
			}
			if !v.Valid {
				o.SetNull()
			} else {
				o.SetInt64(v.Int64)
			}
		default:
			fmt.Println(pti.PoleInfo.GetPoleType())
			panic("")
		}
		doc.poles[poleName] = o
	}
	return nil
}
