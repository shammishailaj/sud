package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/publicsuffix"

	"bytes"

	"fmt"

	"strings"

	"github.com/crazyprograms/callpull"
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
	fmt.Println("send:", string(outBuff))
	if r, err = client.http.Post(client.url+Name, "application/json", bytes.NewReader(outBuff)); err != nil {
		return err
	}
	if _, err = inBuff.ReadFrom(r.Body); err != nil {
		return err
	}
	fmt.Println("recv:", string(inBuff.Bytes()))
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
	if err = client.httpJsonClient("/json/login", jsonLogin{Login: login, Password: login, ConfigurationName: configurationName}, &result); err != nil {
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
func (client *Client) Listen(Name string, TimeoutWait time.Duration) (Param map[string]interface{}, Result chan callpull.Result, errResult error) {
	var err error
	result := jsonListenResult{}
	if err = client.httpJsonClient("/json/listen", &jsonListen{Name: Name, TimeoutWait: int(TimeoutWait.Nanoseconds() / 1000)}, &result); err != nil {
		return nil, nil, err
	}
	if result.Error != "" {
		return nil, nil, errors.New(result.Error)
	}
	var p map[string]interface{}
	if p, err = jsonUnPackMap(*result.Param); err != nil {
		return nil, nil, err
	}
	ResultChan := make(chan callpull.Result)
	go client.listenResult(ResultChan, TimeoutWait, result.ResultUID)
	return p, ResultChan, nil
}
func (client *Client) Call(Name string, Params map[string]interface{}, TimeoutWait time.Duration) (callpull.Result, error) {
	var err error
	var m map[string]*jsonParam
	if m, err = jsonPackMap(Params); err != nil {
		return callpull.Result{}, err
	}
	result := jsonCallResult{}
	if err = client.httpJsonClient("/json/call", &jsonCall{Name: Name, Params: &m, TimeoutWait: int(TimeoutWait / time.Millisecond)}, &result); err != nil {
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
func (client *Client) GetDocumentsPoles(TransactionUID string, DocumentType string, poles []string, wheres []corebase.IDocumentWhere) (map[string]map[string]interface{}, error) {
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
			if n[0] != "DocumentWhere" || n[1] != WhereType {
				return nil, errors.New("where param name error " + name)
			}
			p[strings.ToLower(n[2][0:1])+n[2][1:]] = value
		}
		if w[i], err = jsonPackMap(p); err != nil {
			return nil, err
		}
	}
	result := jsonGetDocumentPolesResult{}
	if err = client.httpJsonClient("/json/getdocumentpoles", &jsonGetDocumentPoles{TransactionUID: TransactionUID, DocumentType: DocumentType, Poles: poles, Wheres: w}, &result); err != nil {
		return nil, err
	}
	if result.Error != "" {
		return nil, errors.New(result.Error)
	}
	Documents := map[string]map[string]interface{}{}
	if result.Documents != nil {
		for DocumentUID, Poles := range *result.Documents {
			if Documents[DocumentUID], err = jsonUnPackMap(Poles); err != nil {
				return nil, err
			}
		}
	}
	return Documents, nil
}
func (client *Client) NewDocument(TransactionUID string, DocumentType string, poles map[string]interface{}) (string, error) {
	var err error
	var m map[string]*jsonParam
	if m, err = jsonPackMap(poles); err != nil {
		return "", err
	}
	result := jsonNewDocumentResult{}
	if err = client.httpJsonClient("/json/newdocument", &jsonNewDocument{TransactionUID: TransactionUID, DocumentType: DocumentType, Poles: &m}, &result); err != nil {
		return "", err
	}
	if result.Error != "" {
		return "", errors.New(result.Error)
	}
	return result.DocumentUID, nil
}
func (client *Client) SetDocumentPoles(TransactionUID string, DocumentUID string, poles map[string]interface{}) error {
	var err error
	var m map[string]*jsonParam
	if m, err = jsonPackMap(poles); err != nil {
		return err
	}
	result := jsonSetDocumentPolesResult{}
	if err = client.httpJsonClient("/json/setdocumentpoles", &jsonSetDocumentPoles{TransactionUID: TransactionUID, DocumentUID: DocumentUID, Poles: &m}, &result); err != nil {
		return err
	}
	if result.Error != "" {
		return errors.New(result.Error)
	}
	return nil
}

var c client.IClient = &Client{}
