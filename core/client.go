package core

import (
	"sync"
	"time"

	"github.com/crazyprograms/sud/callpull"
	iclient "github.com/crazyprograms/sud/client"
	"github.com/crazyprograms/sud/corebase"
	_ "github.com/lib/pq"
)

type resultInfo struct {
	ResultUID string
	Result    chan callpull.Result
	Access    corebase.IAccess
}
type coreClient struct {
	core              *Core
	user              corebase.IUser
	configurationName string
	transactions      map[string]bool
	lockResult        sync.Mutex
	result            map[string]*resultInfo
}

func (core *Core) NewClient(Login string, Password string, ConfigurationName string) (iclient.IClient, error) {
	user := core.getUser(Login)
	if user == nil {
		return nil, &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "NewClient", Name: Login}
	}
	if !user.GetCheckPassword(Password) {
		return nil, &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "NewClient", Name: Login}
	}
	//configuration := core.LoadConfiguration(ConfigurationName)
	return &coreClient{user: user, configurationName: ConfigurationName, transactions: map[string]bool{}, core: core, result: map[string]*resultInfo{}}, nil
}

func (client *coreClient) GetConfiguration() string { return client.configurationName }
func (client *coreClient) BeginTransaction() (string, error) {
	if TransactionUID, err := client.core.BeginTransaction(); err != nil {
		return "", err
	} else {
		client.transactions[TransactionUID] = true
		return TransactionUID, nil
	}
}
func (client *coreClient) CommitTransaction(TransactionUID string) error {
	delete(client.transactions, TransactionUID)
	return client.core.CommitTransaction(TransactionUID)
}
func (client *coreClient) RollbackTransaction(TransactionUID string) error {
	delete(client.transactions, TransactionUID)
	return client.core.RollbackTransaction(TransactionUID)
}
func (client *coreClient) newResult() *resultInfo {
	result := &resultInfo{}
	result.ResultUID = <-client.core.genUID
	client.lockResult.Lock()
	client.result[result.ResultUID] = result
	client.lockResult.Unlock()
	return result
}
func (client *coreClient) getResultAccess(AccessResultUID string) (corebase.IAccess, error) {
	if AccessResultUID == "" {
		return client.user, nil
	}
	client.lockResult.Lock()
	result, ok := client.result[AccessResultUID]
	client.lockResult.Unlock()
	if !ok {
		return nil, &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "ResultAccess", Name: AccessResultUID}
	}
	return AccessUnion(client.user, result.Access), nil
}
func (client *coreClient) freeResult(ResultUID string) *resultInfo {
	client.lockResult.Lock()
	if result, ok := client.result[ResultUID]; ok {
		delete(client.result, ResultUID)
		return result
	}
	client.lockResult.Unlock()
	return nil
}
func (client *coreClient) Listen(Name string, TimeoutWait time.Duration) (Param map[string]interface{}, ResultUID string, errResult error) {
	var cResult chan callpull.Result
	var Access corebase.IAccess
	Param, Access, cResult, errResult = client.core.Listen(client.configurationName, Name, TimeoutWait, client.user)
	if errResult == nil {
		result := client.newResult()
		result.Result = cResult
		result.Access = Access
		ResultUID = result.ResultUID
		//result.Access
	}
	return
}
func (client *coreClient) ListenResult(ResultUID string, Result interface{}, ResultError error) (err error) {
	result := client.freeResult(ResultUID)
	if result == nil {
		err = &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "ListenResult", Name: ResultUID}
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = &corebase.Error{ErrorType: corebase.ErrorTimeout, Action: "ListenResult", Name: ResultUID}
		}
	}()
	err = nil
	result.Result <- callpull.Result{Result: Result, Error: ResultError}
	return
}
func (client *coreClient) Call(Name string, Params map[string]interface{}, TimeoutWait time.Duration, AccessRusultUID string) (callpull.Result, error) {
	var err error
	var Access corebase.IAccess
	if Access, err = client.getResultAccess(AccessRusultUID); err != nil {
		return callpull.Result{Result: nil}, err
	}
	return client.core.Call(client.configurationName, Name, Params, TimeoutWait, Access)
}
func (client *coreClient) GetRecordsPoles(TransactionUID string, RecordType string, poles []string, wheres []corebase.IRecordWhere, AccessRusultUID string) (map[corebase.UUID]map[string]interface{}, error) {
	var err error
	var Access corebase.IAccess
	if Access, err = client.getResultAccess(AccessRusultUID); err != nil {
		return nil, err
	}
	return client.core.GetRecordsPoles(TransactionUID, client.configurationName, RecordType, poles, wheres, Access)
}
func (client *coreClient) SetRecordPoles(TransactionUID string, RecordUID corebase.UUID, poles map[string]interface{}, AccessRusultUID string) error {
	var err error
	var Access corebase.IAccess
	if Access, err = client.getResultAccess(AccessRusultUID); err != nil {
		return err
	}
	return client.core.SetRecordPoles(TransactionUID, client.configurationName, RecordUID, poles, Access)
}
func (client *coreClient) NewRecord(TransactionUID string, RecordType string, Poles map[string]interface{}, AccessRusultUID string) (corebase.UUID, error) {
	var err error
	var Access corebase.IAccess
	if Access, err = client.getResultAccess(AccessRusultUID); err != nil {
		return corebase.NullUUID, err
	}
	return client.core.NewRecord(TransactionUID, client.configurationName, RecordType, Poles, Access)
}
func (client *coreClient) GetRecordAccess(TransactionUID string, RecordUID corebase.UUID, AccessRusultUID string) (string, error) {
	var err error
	var Access corebase.IAccess
	if Access, err = client.getResultAccess(AccessRusultUID); err != nil {
		return "", err
	}
	return client.core.GetRecordAccess(TransactionUID, client.configurationName, RecordUID, Access)
}
func (client *coreClient) SetRecordAccess(TransactionUID string, RecordUID corebase.UUID, NewAccess string, AccessRusultUID string) error {
	var err error
	var Access corebase.IAccess
	if Access, err = client.getResultAccess(AccessRusultUID); err != nil {
		return err
	}
	return client.core.SetRecordAccess(TransactionUID, client.configurationName, RecordUID, NewAccess, Access)
}

/**/
