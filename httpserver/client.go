package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/publicsuffix"

	"bytes"

	"strings"

	"github.com/crazyprograms/sud/callpull"
	"github.com/crazyprograms/sud/client"
	"github.com/crazyprograms/sud/corebase"
)

type Client struct {
	url  string
	jar  http.CookieJar
	http http.Client
}

func (client *Client) httpJsonClient(Name string, InParam interface{}, OutParam interface{}) error {
	var err error
	var inBuff bytes.Buffer
	var outBuff []byte
	if outBuff, err = json.Marshal(InParam); err != nil {
		return err
	}
	var r *http.Response
	fmt.Println("send("+client.url+Name+"):", string(outBuff))
	if r, err = client.http.Post(client.url+Name, "application/json", bytes.NewReader(outBuff)); err != nil {
		return err
	}
	if _, err = inBuff.ReadFrom(r.Body); err != nil {
		return err
	}
	fmt.Println("recv("+client.url+Name+"):", string(inBuff.Bytes()))
	if err = json.Unmarshal(inBuff.Bytes(), OutParam); err != nil {
		return err
	}
	return nil
}
func NewClient(url, login, password, configurationName string) (*Client, error) {
	var err error
	var jar http.CookieJar
	if jar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List}); err != nil {
		return nil, err
	}
	client := &Client{url: url, jar: jar, http: http.Client{Jar: jar}}
	result := jsonLoginResult{}
	if err = client.httpJsonClient("/json/login", jsonLogin{Login: login, Password: password, ConfigurationName: configurationName}, &result); err != nil {
		return nil, err
	}
	if !result.Login {
		return nil, errors.New(result.Error)
	}
	return client, nil
}
func (client *Client) GetConfiguration() string {
	var err error
	result := jsonGetConfigurationResult{}
	if err = client.httpJsonClient("/json/getconfiguration", struct{}{}, &result); err != nil {
		return ""
	}
	if result.Error != "" {
		return ""
	}
	return result.Configuration
}
func (client *Client) BeginTransaction() (string, error) {
	var err error
	result := jsonBeginTransactionResult{}
	if err = client.httpJsonClient("/json/begintransaction", struct{}{}, &result); err != nil {
		return "", err
	}
	if result.Error != "" {
		return "", errors.New(result.Error)
	}
	return result.TransactionUID, nil
}
func (client *Client) CommitTransaction(TransactionUID string) error {
	var err error
	result := jsonCommitTransactionResult{}
	if err = client.httpJsonClient("/json/committransaction", &jsonCommitTransaction{TransactionUID: TransactionUID}, &result); err != nil {
		return err
	}
	if result.Error != "" {
		return errors.New(result.Error)
	}
	return nil
}
func (client *Client) RollbackTransaction(TransactionUID string) error {
	var err error
	result := jsonRollbackTransactionResult{}
	if err = client.httpJsonClient("/json/rollbacktransaction", &jsonRollbackTransaction{TransactionUID: TransactionUID}, &result); err != nil {
		return err
	}
	if result.Error != "" {
		return errors.New(result.Error)
	}
	return nil
}
func (client *Client) listenResult(Result chan callpull.Result, TimeoutWait time.Duration, ResultUID string) error {
	var err error
	select {
	case r := <-Result:
		var resultParam *jsonParam
		if resultParam, err = jsonPack(r.Result); err != nil {
			return err
		}
		var resultError string
		if r.Error != nil {
			resultError = r.Error.Error()
		}
		var result jsonListenReturnResult
		if err = client.httpJsonClient("/json/listenreturn", &jsonListenReturn{ResultUID: ResultUID, Result: resultParam, Error: resultError}, &result); err != nil {
			return err
		}
		if result.Error != "" {
			return errors.New(result.Error)
		}
		return nil
	case <-time.After(TimeoutWait + time.Second):
		return errors.New("timeout error")
	}
}

func (client *Client) Listen(Name string, TimeoutWait time.Duration) (Param map[string]interface{}, ResultUID string, errResult error) {
	var err error
	result := jsonListenResult{}
	if err = client.httpJsonClient("/json/listen", &jsonListen{Name: Name, TimeoutWait: int(TimeoutWait.Nanoseconds() / 1000)}, &result); err != nil {
		return nil, "", err
	}
	if result.Error != "" {
		return nil, "", errors.New(result.Error)
	}
	var p map[string]interface{}
	if p, err = jsonUnPackMap(*result.Param); err != nil {
		return nil, "", err
	}
	return p, result.ResultUID, nil
}
func (client *Client) ListenResult(ResultUID string, Result interface{}, ResultError error) error {
	var err error
	result := jsonListenReturnResult{}
	var ResultPack *jsonParam
	if ResultPack, err = jsonPack(Result); err != nil {
		return err
	}
	if err = client.httpJsonClient("/json/listenreturn", &jsonListenReturn{ResultUID: ResultUID, Result: ResultPack}, &result); err != nil {
		return err
	}
	if result.Error != "" {
		return errors.New(result.Error)
	}
	return nil
}
func (client *Client) Call(Name string, Params map[string]interface{}, TimeoutWait time.Duration, AccessResultUID string) (callpull.Result, error) {
	var err error
	var m map[string]*jsonParam
	if m, err = jsonPackMap(Params); err != nil {
		return callpull.Result{}, err
	}
	result := jsonCallResult{}
	if err = client.httpJsonClient("/json/call", &jsonCall{Name: Name, Params: &m, TimeoutWait: int(TimeoutWait / time.Millisecond), AccessResultUID: AccessResultUID}, &result); err != nil {
		return callpull.Result{}, err
	}
	var Result interface{}
	if Result, err = jsonUnPack(result.Result); err != nil {
		return callpull.Result{}, err
	}
	var callpullError error
	if result.Error != "" {
		callpullError = errors.New(result.Error)
	}
	return callpull.Result{Result: Result, Error: callpullError}, nil
}
func (client *Client) GetRecordsPoles(TransactionUID string, RecordType string, poles []string, wheres []corebase.IRecordWhere, AccessResultUID string) (map[corebase.UUID]map[string]interface{}, error) {
	var err error
	w := make([]map[string]*jsonParam, len(wheres), len(wheres))
	for i, where := range wheres {
		p := map[string]interface{}{}
		var WhereType string
		var Params map[string]interface{}
		WhereType, Params, err = where.Save()
		p["whereType"] = WhereType
		for name, value := range Params {
			n := strings.Split(name, ".")
			if n[0] != "RecordWhere" || n[1] != WhereType {
				return nil, errors.New("where param name error " + name)
			}
			p[strings.ToLower(n[2][0:1])+n[2][1:]] = value
		}
		if w[i], err = jsonPackMap(p); err != nil {
			return nil, err
		}
	}
	result := jsonGetRecordPolesResult{}
	if err = client.httpJsonClient("/json/getrecordpoles", &jsonGetRecordPoles{TransactionUID: TransactionUID, RecordType: RecordType, Poles: poles, Wheres: w, AccessResultUID: AccessResultUID}, &result); err != nil {
		return nil, err
	}
	if result.Error != "" {
		return nil, errors.New(result.Error)
	}
	Records := map[corebase.UUID]map[string]interface{}{}
	if result.Records != nil {
		for RecordUID, Poles := range *result.Records {
			if Records[RecordUID], err = jsonUnPackMap(Poles); err != nil {
				return nil, err
			}
		}
	}
	return Records, nil
}
func (client *Client) NewRecord(TransactionUID string, RecordType string, poles map[string]interface{}, AccessResultUID string) (corebase.UUID, error) {
	var err error
	var m map[string]*jsonParam
	if m, err = jsonPackMap(poles); err != nil {
		return "", err
	}
	result := jsonNewRecordResult{}
	if err = client.httpJsonClient("/json/newrecord", &jsonNewRecord{TransactionUID: TransactionUID, RecordType: RecordType, Poles: &m, AccessResultUID: AccessResultUID}, &result); err != nil {
		return "", err
	}
	if result.Error != "" {
		return "", errors.New(result.Error)
	}
	return result.RecordUID, nil
}
func (client *Client) SetRecordPoles(TransactionUID string, RecordUID corebase.UUID, poles map[string]interface{}, AccessResultUID string) error {
	var err error
	var m map[string]*jsonParam
	if m, err = jsonPackMap(poles); err != nil {
		return err
	}
	result := jsonSetRecordPolesResult{}
	if err = client.httpJsonClient("/json/setrecordpoles", &jsonSetRecordPoles{TransactionUID: TransactionUID, RecordUID: RecordUID, Poles: &m, AccessResultUID: AccessResultUID}, &result); err != nil {
		return err
	}
	if result.Error != "" {
		return errors.New(result.Error)
	}
	return nil
}

func (client *Client) GetRecordAccess(TransactionUID string, RecordUID corebase.UUID, AccessResultUID string) (string, error) {
	var err error
	result := jsonGetRecordAccessResult{}
	if err = client.httpJsonClient("/json/getrecordaccess", &jsonGetRecordAccess{TransactionUID: TransactionUID, RecordUID: RecordUID, AccessResultUID: AccessResultUID}, &result); err != nil {
		return "", err
	}
	if result.Error != "" {
		return "", errors.New(result.Error)
	}
	return result.Access, nil
}
func (client *Client) SetRecordAccess(TransactionUID string, RecordUID corebase.UUID, NewAccess string, AccessResultUID string) error {
	var err error
	result := jsonSetRecordAccessResult{}
	if err = client.httpJsonClient("/json/setrecordaccess", &jsonSetRecordAccess{TransactionUID: TransactionUID, RecordUID: RecordUID, NewAccess: NewAccess, AccessResultUID: AccessResultUID}, &result); err != nil {
		return err
	}
	if result.Error != "" {
		return errors.New(result.Error)
	}
	return nil
}

var c client.IClient = &Client{}
