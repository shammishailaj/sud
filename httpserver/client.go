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
	client := &Client{url: url, jar: jar, http: http.Client{Jar: jar, Timeout}}
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
	return ""

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
	return nil, nil
}
func (client *Client) NewDocument(TransactionUID string, DocumentType string, poles map[string]interface{}) (string, error) {
	return "", nil
}
func (client *Client) SetDocumentPoles(TransactionUID string, DocumentUID string, poles map[string]interface{}) error {
	var err error
	var r *http.Response
	var m map[string]*jsonParam
	if m, err = jsonPackMap(poles); err != nil {
		return err
	}
	j := jsonSetDocumentPoles{TransactionUID: TransactionUID, DocumentUID: DocumentUID, Poles: &m}
	var data []byte
	if data, err = json.Marshal(&j); err == nil {
		if r, err = client.http.Post(client.url+"/json/setdocument", "application/json", bytes.NewReader(data)); err != nil {
			return err
		}
		var b bytes.Buffer
		if _, err = b.ReadFrom(r.Body); err == nil {
			var result jsonSetDocumentPolesResult
			fmt.Println("recv:", string(b.Bytes()))
			if err = json.Unmarshal(b.Bytes(), &result); err != nil {
				return err
			}
			if result.Error != "" {
				return errors.New(result.Error)
			}
			return nil
		}
	}
	return err
}

var c client.IClient = &Client{}
