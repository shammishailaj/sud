package httpserver

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

type jsonErrorReturn struct {
	Error string `json:"error,omitempty"`
}
type jsonLogin struct {
	Login             string `json:"login"`
	Password          string `json:"password"`
	ConfigurationName string `json:"configurationName"`
}

type jsonLoginResult struct {
	Login bool   `json:"login,omitempty"`
	Error string `json:"error,omitempty"`
}

type jsonListen struct {
	Name        string `json:"name"`
	TimeoutWait int    `json:"timeoutWait"` //в милесекундах
}
type jsonListenResult struct {
	Param     *map[string]*jsonParam `json:"param,omitempty"`
	ResultUID string                 `json:"resultUID,omitempty"`
	Error     string                 `json:"error,omitempty"`
}
type jsonListenReturn struct {
	ResultUID string     `json:"resultUID,omitempty"`
	Result    *jsonParam `json:"result,omitempty"`
	Error     string     `json:"error,omitempty"`
}
type jsonListenReturnResult struct {
	Error string `json:"error,omitempty"`
}
type jsonBeginTransactionResult struct {
	TransactionUID string `json:"transactionUID,omitempty"`
	Error          string `json:"error,omitempty"`
}

type jsonCommitTransaction struct {
	TransactionUID string `json:"transactionUID"`
}

type jsonCommitTransactionResult struct {
	Error string `json:"error,omitempty"`
}

type jsonRollbackTransaction struct {
	TransactionUID string `json:"transactionUID"`
}

type jsonRollbackTransactionResult struct {
	Error string `json:"error,omitempty"`
}

type jsonCall struct {
	Name        string                 `json:"name"`
	Params      *map[string]*jsonParam `json:"params,omitempty"`
	TimeoutWait int                    `json:"timeoutWait"` //в милесекундах
}
type jsonCallResult struct {
	Result *jsonParam `json:"result"`
	Error  string     `json:"error,omitempty"`
}
type jsonSetDocumentPoles struct {
	TransactionUID string                 `json:"transactionUID"`
	DocumentUID    string                 `json:"documentUID"`
	Poles          *map[string]*jsonParam `json:"poles,omitempty"`
}
type jsonSetDocumentPolesResult struct {
	Error string `json:"error,omitempty"`
}
type jsonParam struct {
	String *string                `json:"string,omitempty"`
	Int    *int64                 `json:"int,omitempty"`
	Time   *time.Time             `json:"time,omitempty"`
	Bytes  *string                `json:"bytes,omitempty"`
	Map    *map[string]*jsonParam `json:"map,omitempty"`
	List   *[]*jsonParam          `json:"list,omitempty"`
}

func jsonPackMap(Params map[string]interface{}) (map[string]*jsonParam, error) {
	m := map[string]*jsonParam{}
	for name, value := range Params {
		var err error
		if m[name], err = jsonPack(value); err != nil {
			return nil, err
		}
	}
	return m, nil
}
func jsonPackList(Params []interface{}) ([]*jsonParam, error) {
	m := make([]*jsonParam, len(Params), len(Params))
	for name, value := range Params {
		var err error
		if m[name], err = jsonPack(value); err != nil {
			return nil, err
		}
	}
	return m, nil
}
func jsonUnPackMap(Params map[string]*jsonParam) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	for name, value := range Params {
		var err error
		if m[name], err = jsonUnPack(value); err != nil {
			return nil, err
		}
	}
	return m, nil
}
func jsonUnPackList(Params []*jsonParam) ([]interface{}, error) {
	m := make([]interface{}, len(Params), len(Params))
	for name, value := range Params {
		var err error
		if m[name], err = jsonUnPack(value); err != nil {
			return nil, err
		}
	}
	return m, nil
}
func jsonPack(Param interface{}) (*jsonParam, error) {
	var err error
	switch v := Param.(type) {
	case string:
		return &jsonParam{String: &v}, nil
	case int64:
		return &jsonParam{Int: &v}, nil
	case int:
		v1 := int64(v)
		return &jsonParam{Int: &v1}, nil
	case time.Time:
		return &jsonParam{Time: &v}, nil
	case []byte:
		s := base64.StdEncoding.EncodeToString(v)
		return &jsonParam{Bytes: &s}, nil
	case map[string]interface{}:
		var m map[string]*jsonParam
		if m, err = jsonPackMap(v); err != nil {
			return nil, err
		}
		return &jsonParam{Map: &m}, nil
	case []interface{}:
		var m []*jsonParam
		if m, err = jsonPackList(v); err != nil {
			return nil, err
		}
		return &jsonParam{List: &m}, nil
	default:
		return nil, errors.New("json pack error " + fmt.Sprintln(v))
	}
}
func jsonUnPack(Param *jsonParam) (interface{}, error) {
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
	if Param.Time != nil {
		return *Param.Time, nil
	}
	if Param.Bytes != nil {
		return base64.StdEncoding.DecodeString(*Param.Bytes)
	}
	if Param.Map != nil {
		m := map[string]interface{}{}
		if m, err = jsonUnPackMap(*Param.Map); err != nil {
			return nil, err
		}
		return m, nil
	}
	if Param.List != nil {
		m := make([]interface{}, len(*Param.List), len(*Param.List))
		if m, err = jsonUnPackList(*Param.List); err != nil {
			return nil, err
		}
		return m, nil
	}
	return nil, errors.New("json unpack error")
}
