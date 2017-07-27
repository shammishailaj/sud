package client

import (
	"time"

	"github.com/crazyprograms/sud/callpull"
	"github.com/crazyprograms/sud/corebase"
)

type IClient interface {
	GetConfiguration() string
	BeginTransaction() (string, error)
	CommitTransaction(TransactionUID string) error
	RollbackTransaction(TransactionUID string) error
	//Listen(Name string, TimeoutWait time.Duration) (Param map[string]interface{}, Result chan callpull.Result, errResult error)
	Listen(Name string, TimeoutWait time.Duration) (Param map[string]interface{}, ResultUID string, errResult error)
	ListenResult(ResultUID string, Result interface{}, ResultError error) error
	Call(Name string, Params map[string]interface{}, TimeoutWait time.Duration, AccessResultUID string) (callpull.Result, error)
	GetRecordsPoles(TransactionUID string, RecordType string, poles []string, wheres []corebase.IRecordWhere, AccessResultUID string) (map[corebase.UUID]map[string]interface{}, error)
	NewRecord(TransactionUID string, RecordType string, poles map[string]interface{}, AccessResultUID string) (corebase.UUID, error)
	SetRecordPoles(TransactionUID string, RecordUID corebase.UUID, poles map[string]interface{}, AccessResultUID string) error
	GetRecordAccess(TransactionUID string, RecordUID corebase.UUID, AccessResultUID string) (string, error)
	SetRecordAccess(TransactionUID string, RecordUID corebase.UUID, NewAccess string, AccessResultUID string) error
}
