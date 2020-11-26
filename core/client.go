package core

import (
	"errors"
	"time"

	"github.com/crazyprograms/callpull"
	"github.com/shammishailaj/sud/corebase"
	_ "github.com/lib/pq"
)

type Client struct {
	core              *Core
	user              IUser
	configurationName string
	transactions      map[string]bool
}

func (core *Core) NewClient(Login string, Password string, ConfigurationName string) (*Client, error) {
	user := core.getUser(Login)
	if user == nil {
		return nil, errors.New("Login error")
	}
	if !user.GetCheckPassword(Password) {
		return nil, errors.New("Login error")
	}
	//configuration := core.LoadConfiguration(ConfigurationName)
	return &Client{user: user, configurationName: ConfigurationName, transactions: map[string]bool{}, core: core}, nil
}

func (client *Client) GetConfiguration() string { return client.configurationName }
func (client *Client) BeginTransaction() (string, error) {
	if TransactionUID, err := client.core.BeginTransaction(); err != nil {
		return "", err
	} else {
		client.transactions[TransactionUID] = true
		return TransactionUID, nil
	}
}
func (client *Client) CommitTransaction(TransactionUID string) error {
	delete(client.transactions, TransactionUID)
	return client.core.CommitTransaction(TransactionUID)
}
func (client *Client) RollbackTransaction(TransactionUID string) error {
	delete(client.transactions, TransactionUID)
	return client.core.RollbackTransaction(TransactionUID)
}
func (client *Client) Listen(Name string, TimeoutWait time.Duration) (Param map[string]interface{}, Result chan callpull.Result, errResult error) {
	return client.core.Listen(client.configurationName, Name, TimeoutWait)
}
func (client *Client) Call(Name string, Params map[string]interface{}, TimeoutWait time.Duration) (callpull.Result, error) {
	return client.core.Call(client.configurationName, Name, Params, TimeoutWait)
}
func (client *Client) GetRecordsPoles(TransactionUID string, RecordType string, poles []string, wheres []corebase.IRecordWhere) (map[string]map[string]interface{}, error) {
	return client.core.GetRecordsPoles(TransactionUID, client.configurationName, RecordType, poles, wheres)
}
func (client *Client) SetRecordPoles(TransactionUID string, RecordUID string, poles map[string]interface{}) error {
	return client.core.SetRecordPoles(TransactionUID, client.configurationName, RecordUID, poles)
}
func (client *Client) NewRecord(TransactionUID string, RecordType string, Poles map[string]interface{}) (string, error) {
	return client.core.NewRecord(TransactionUID, client.configurationName, RecordType, Poles)
}

/**/
