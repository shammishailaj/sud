/*Package jparam - module of pack and unpack data
compatible types
 int64, string, tie.Time, []byte, []interface{}, map[string]interface{} */
package jparam

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/crazyprograms/sud/corebase"
)

type JParam struct {
	Bool   *bool               `json:"bool,omitempty"`
	String *string             `json:"string,omitempty"`
	Int    *int64              `json:"int,omitempty"`
	Time   *time.Time          `json:"time,omitempty"`
	Bytes  *string             `json:"bytes,omitempty"`
	NULL   *string             `json:"null,omitempty"`
	Map    *map[string]*JParam `json:"map,omitempty"`
	List   *[]*JParam          `json:"list,omitempty"`
}

func PackMap(Params map[string]interface{}) (map[string]*JParam, error) {
	m := map[string]*JParam{}
	for name, value := range Params {
		var err error
		if m[name], err = Pack(value); err != nil {
			return nil, err
		}
	}
	return m, nil
}
func PackList(Params []interface{}) ([]*JParam, error) {
	m := make([]*JParam, len(Params), len(Params))
	for name, value := range Params {
		var err error
		if m[name], err = Pack(value); err != nil {
			return nil, err
		}
	}
	return m, nil
}
func UnPackMap(Params map[string]*JParam) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	for name, value := range Params {
		var err error
		if m[name], err = UnPack(value); err != nil {
			return nil, err
		}
	}
	return m, nil
}
func UnPackList(Params []*JParam) ([]interface{}, error) {
	m := make([]interface{}, len(Params), len(Params))
	for name, value := range Params {
		var err error
		if m[name], err = UnPack(value); err != nil {
			return nil, err
		}
	}
	return m, nil
}
func Pack(Param interface{}) (*JParam, error) {
	var err error
	switch v := Param.(type) {
	case bool:
		return &JParam{Bool: &v}, nil
	case string:
		return &JParam{String: &v}, nil
	case int64:
		return &JParam{Int: &v}, nil
	case int:
		v1 := int64(v)
		return &JParam{Int: &v1}, nil
	case time.Time:
		return &JParam{Time: &v}, nil
	case []byte:
		s := base64.StdEncoding.EncodeToString(v)
		return &JParam{Bytes: &s}, nil
	case map[string]interface{}:
		var m map[string]*JParam
		if m, err = PackMap(v); err != nil {
			return nil, err
		}
		return &JParam{Map: &m}, nil
	case []interface{}:
		var m []*JParam
		if m, err = PackList(v); err != nil {
			return nil, err
		}
		return &JParam{List: &m}, nil
	case corebase.TNULL:
		n := "NULL"
		return &JParam{NULL: &n}, nil
	default:
		return nil, errors.New("json pack error " + fmt.Sprintln(v))
	}
}
func UnPack(Param *JParam) (interface{}, error) {
	var err error
	if Param == nil {
		return nil, nil
	}
	if Param.String != nil {
		return *Param.String, nil
	}
	if Param.Int != nil {
		return *Param.Int, nil
	}
	if Param.Bool != nil {
		return *Param.Bool, nil
	}
	if Param.Time != nil {
		return *Param.Time, nil
	}
	if Param.Bytes != nil {
		return base64.StdEncoding.DecodeString(*Param.Bytes)
	}
	if Param.Map != nil {
		m := map[string]interface{}{}
		if m, err = UnPackMap(*Param.Map); err != nil {
			return nil, err
		}
		return m, nil
	}
	if Param.List != nil {
		m := make([]interface{}, len(*Param.List), len(*Param.List))
		if m, err = UnPackList(*Param.List); err != nil {
			return nil, err
		}
		return m, nil
	}
	if Param.NULL != nil {
		return corebase.NULL, nil
	}
	return nil, errors.New("json unpack error")
}
func ToJson(Param interface{}) ([]byte, error) {
	var err error
	var InParam *JParam
	if InParam, err = Pack(Param); err != nil {
		return nil, err
	}
	return json.Marshal(InParam)
}
func FromJson(str []byte) (interface{}, error) {
	var err error
	var OutParam JParam
	if err = json.Unmarshal(str, &OutParam); err != nil {
		return nil, err
	}
	return UnPack(&OutParam)
}
