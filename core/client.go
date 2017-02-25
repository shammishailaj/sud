package core

import (
	"time"

	"github.com/crazyprograms/callpull"
	_ "github.com/lib/pq"
)

type Client struct {
	core              *Core
	user              IUser
	configurationName string
	transactions      map[string]bool
}

func (core *Core) NewClient(Login string, Password string, ConfigurationName string) *Client {
	user := core.getUser(Login)
	if !user.GetCheckPassword(Password) {
		return nil
	}
	//configuration := core.LoadConfiguration(ConfigurationName)
	return &Client{user: user, configurationName: ConfigurationName, transactions: map[string]bool{}, core: core}
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
func (client *Client) GetDocumentsPoles(TransactionUID string, DocumentType string, poles []string, wheres []IDocumentWhere) (map[string]map[string]interface{}, error) {
	return client.core.GetDocumentsPoles(TransactionUID, client.configurationName, DocumentType, poles, wheres)
}
func (client *Client) SetDocumentPoles(TransactionUID string, DocumentUID string, poles map[string]interface{}) error {
	return client.core.SetDocumentPoles(TransactionUID, client.configurationName, DocumentUID, poles)
}
func (client *Client) NewDocument(TransactionUID string, DocumentType string, Poles map[string]interface{}) (string, error) {
	return client.core.NewDocument(TransactionUID, client.configurationName, DocumentType, Poles)
}

/**/
