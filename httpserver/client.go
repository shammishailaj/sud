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
)

type Client struct {
	url  string
	jar  http.CookieJar
	http http.Client
}

func NewClient(url, login, password, configurationName string) (*Client, error) {
	var err error
	var jar http.CookieJar
	if jar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List}); err != nil {
		return nil, err
	}
	client := &Client{url: url, jar: jar, http: http.Client{Jar: jar}}
	j := jsonLogin{Login: login, Password: login, ConfigurationName: configurationName}
	var data []byte
	if data, err = json.Marshal(&j); err == nil {
		fmt.Println("send:", string(data))
		var r *http.Response
		r, err = client.http.Post(client.url+"/json/login", "application/json", bytes.NewReader(data))
		var b bytes.Buffer
		if _, err = b.ReadFrom(r.Body); err == nil {
			var result jsonLoginResult
			fmt.Println("recv:", string(b.Bytes()))
			if err = json.Unmarshal(b.Bytes(), &result); err == nil {
				if result.Login {
					return client, nil
				} else {
					return nil, errors.New(result.Error)
				}
			}
		}
	}
	return nil, err
}
func (client *Client) GetConfiguration() string {
	return ""

}
func (client *Client) BeginTransaction() (string, error) {
	var err error
	var r *http.Response
	if r, err = client.http.Post(client.url+"/json/begin", "application/json", bytes.NewReader([]byte("{}"))); err != nil {
		return "", err
	}
	var b bytes.Buffer
	if _, err = b.ReadFrom(r.Body); err == nil {
		var result jsonBeginTransactionResult
		fmt.Println("recv:", string(b.Bytes()))
		if err = json.Unmarshal(b.Bytes(), &result); err == nil {
			return result.TransactionUID, nil
		}
	}
	return "", err
}
func (client *Client) CommitTransaction(TransactionUID string) error {
	var err error
	var r *http.Response
	j := jsonCommitTransaction{TransactionUID: TransactionUID}
	var data []byte
	if data, err = json.Marshal(&j); err == nil {
		if r, err = client.http.Post(client.url+"/json/begin", "application/json", bytes.NewReader(data)); err != nil {
			return err
		}
		var b bytes.Buffer
		if _, err = b.ReadFrom(r.Body); err == nil {
			var result jsonCommitTransactionResult
			fmt.Println("recv:", string(b.Bytes()))
			if err = json.Unmarshal(b.Bytes(), &result); result.Commit {
				return nil
			}
			err = errors.New(result.Error)
		}
	}
	return err
}
func (client *Client) RollbackTransaction(TransactionUID string) error {
	var err error
	var r *http.Response
	j := jsonRollbackTransaction{TransactionUID: TransactionUID}
	var data []byte
	if data, err = json.Marshal(&j); err == nil {
		if r, err = client.http.Post(client.url+"/json/begin", "application/json", bytes.NewReader(data)); err != nil {
			return err
		}
		var b bytes.Buffer
		if _, err = b.ReadFrom(r.Body); err == nil {
			var result jsonRollbackTransactionResult
			fmt.Println("recv:", string(b.Bytes()))
			if err = json.Unmarshal(b.Bytes(), &result); result.Rollback {
				return nil
			}
			err = errors.New(result.Error)
		}
	}
	return err
}
func (client *Client) listenResult(Result chan callpull.Result, TimeoutWait time.Duration, ResultUID string) error {
	var err error
	select {
	case r := <-Result:
		j := jsonListenReturn{ResultUID: ResultUID}
		if j.Result, err = jsonPack(r.Result); err != nil {
			return err
		}
		if r.Error != nil {
			j.Error = r.Error.Error()
		}
		var data []byte
		if data, err = json.Marshal(&j); err == nil {
			fmt.Println("send:", string(data))
			var r *http.Response
			r, err = client.http.Post(client.url+"/json/listenreturn", "application/json", bytes.NewReader(data))
			var b bytes.Buffer
			if _, err = b.ReadFrom(r.Body); err == nil {
				var result jsonListenReturnResult
				fmt.Println("recv:", string(b.Bytes()))
				if err = json.Unmarshal(b.Bytes(), &result); err == nil {
					if result.Error != "" {
						err = errors.New(result.Error)
					}
				}
			}
		}
		return err
	case <-time.After(TimeoutWait + time.Second):
		return errors.New("timeout error")
	}
}
func (client *Client) Listen(Name string, TimeoutWait time.Duration) (Param map[string]interface{}, Result chan callpull.Result, errResult error) {
	var err error
	j := jsonListen{Name: Name, TimeoutWait: int(TimeoutWait.Nanoseconds() / 1000)}
	var data []byte
	if data, err = json.Marshal(&j); err == nil {
		fmt.Println("send:", string(data))
		var r *http.Response
		if r, err = client.http.Post(client.url+"/json/listen", "application/json", bytes.NewReader(data)); err != nil {
			return nil, nil, err
		}
		var b bytes.Buffer
		if _, err = b.ReadFrom(r.Body); err == nil {
			var result jsonListenResult
			fmt.Println("recv:", string(b.Bytes()))
			if err = json.Unmarshal(b.Bytes(), &result); err == nil {
				if result.Error == "" {
					var p map[string]interface{}
					if p, err = jsonUnPackMap(*result.Param); err == nil {
						ResultChan := make(chan callpull.Result)
						go client.listenResult(ResultChan, TimeoutWait, result.ResultUID)
						return p, ResultChan, nil
					}
				} else {
					err = errors.New(result.Error)
				}
			}
		}
	}
	return nil, nil, err
}
func (client *Client) Call(Name string, Params map[string]interface{}, TimeoutWait time.Duration) (callpull.Result, error) {
	var err error
	var r *http.Response
	var m map[string]*jsonParam
	if m, err = jsonPackMap(Params); err != nil {
		return callpull.Result{}, err
	}
	j := jsonCall{Name: Name, Params: &m, TimeoutWait: int(TimeoutWait / time.Millisecond)}
	var data []byte
	if data, err = json.Marshal(&j); err == nil {
		if r, err = client.http.Post(client.url+"/json/call", "application/json", bytes.NewReader(data)); err != nil {
			return callpull.Result{}, err
		}
		var b bytes.Buffer
		if _, err = b.ReadFrom(r.Body); err == nil {
			var result jsonCallResult
			fmt.Println("recv:", string(b.Bytes()))
			if err = json.Unmarshal(b.Bytes(), &result); err != nil {
				return callpull.Result{}, err
			}
			var Result interface{}
			if Result, err = jsonUnPack(result.Result); err != nil {
				return callpull.Result{}, err
			}
			r := callpull.Result{Result: Result}
			if result.Error != "" {
				r.Error = errors.New(result.Error)
			}
			return r, nil
		}
	}
	return callpull.Result{}, nil
}
func (client *Client) GetDocumentsPoles(TransactionUID string, DocumentType string, poles []string, wheres []client.IDocumentWhere) (map[string]map[string]interface{}, error) {
	return nil, nil
}
func (client *Client) NewDocument(TransactionUID string, DocumentType string, poles map[string]interface{}) (string, error) {
	return "", nil
}
func (client *Client) SetDocumentPoles(TransactionUID string, DocumentUID string, poles map[string]interface{}) error {
	return nil
}

var c client.IClient = &Client{}
