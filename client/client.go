package client

import (
	"time"

	"github.com/crazyprograms/callpull"
	"github.com/crazyprograms/sud/corebase"
)

type IClient interface {
	GetConfiguration() string
	BeginTransaction() (string, error)
	CommitTransaction(TransactionUID string) error
	RollbackTransaction(TransactionUID string) error
	Listen(Name string, TimeoutWait time.Duration) (Param map[string]interface{}, Result chan callpull.Result, errResult error)
	Call(Name string, Params map[string]interface{}, TimeoutWait time.Duration) (callpull.Result, error)
	GetRecordsPoles(TransactionUID string, RecordType string, poles []string, wheres []corebase.IRecordWhere) (map[string]map[string]interface{}, error)
	NewRecord(TransactionUID string, RecordType string, poles map[string]interface{}) (string, error)
	SetRecordPoles(TransactionUID string, RecordUID string, poles map[string]interface{}) error
}
