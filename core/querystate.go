package core

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/crazyprograms/sud/corebase"
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

func (state *queryState) AddPole(poleName string, pi corebase.IPoleInfo) {
	pti := &ptiInfo{}
	pti.FromPoleInfo(pi)
	pti.SQL = `"` + pti.TableName + `"."` + pti.PoleName + `"`
	pti.N = -1
	state.num++
	state.poles[pi.GetPoleName()] = pti
	state.AddTable(pti.TableName, `LEFT JOIN "`+pti.TableName+`" ON "`+pti.TableName+`"."__RecordUID" = "Record"."__RecordUID"`)
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
func (state *queryState) AddParamObject(value corebase.Object) (string, error) {
	if corebase.IsNull(value) {
		return state.AddParam(nil), nil
	}
	switch v := value.(type) {
	case bool:
		return state.AddParam(v), nil
	case string:
		return state.AddParam(v), nil
	case int64:
		return state.AddParam(v), nil
	case time.Time:
		return state.AddParam(v), nil
	case corebase.UUID:
		return state.AddParam(v), nil
	default:
		return "", errors.New("type not support")
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
	if len(state.wheres) == 0 {
		return ""
	}
	return " WHERE (" + strings.Join(state.wheres, ") AND (") + ")"
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
			panic(pti.PoleInfo.GetPoleType())
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
			panic(pti.PoleInfo.GetPoleType())
		}
	}
	return values
}
func (state *queryState) SetRecordPoles(doc map[string]interface{}, values [](interface{})) error {
	for poleName, pti := range state.poles {
		var o corebase.Object
		switch pti.PoleInfo.GetPoleType() {
		case "BooleanValue":
			v, ok := values[pti.N].(*sql.NullBool)
			if !ok {
				return errors.New("pole read error " + poleName)
			}
			if !v.Valid {
				o = corebase.NULL
			} else {
				o = v.Bool
			}
		case "StringValue":
			v, ok := values[pti.N].(*sql.NullString)
			if !ok {
				return errors.New("pole read error " + poleName)
			}
			if !v.Valid {
				o = corebase.NULL
			} else {
				o = v.String
			}
		case "DateValue":
			v, ok := values[pti.N].(**time.Time)
			if !ok {
				return errors.New("pole read error " + poleName)
			}
			if *v == nil {
				o = corebase.NULL
			} else {
				o = *(*v)
			}
		case "Int64Value":
			v, ok := values[pti.N].(*sql.NullInt64)
			if !ok {
				return errors.New("pole read error " + poleName)
			}
			if !v.Valid {
				o = corebase.NULL
			} else {
				o = v.Int64
			}
		default:
			panic(pti.PoleInfo.GetPoleType())
		}
		doc[poleName] = o
	}
	return nil
}
