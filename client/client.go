package client

import (
	"time"

	"github.com/crazyprograms/callpull"
)

type IDocumentWhere interface {
}
type IClient interface {
	GetConfiguration() string
	BeginTransaction() (string, error)
	CommitTransaction(TransactionUID string) error
	RollbackTransaction(TransactionUID string) error
	Listen(Name string, TimeoutWait time.Duration) (Param map[string]interface{}, Result chan callpull.Result, errResult error)
	Call(Name string, Params map[string]interface{}, TimeoutWait time.Duration) (callpull.Result, error)
	GetDocumentsPoles(TransactionUID string, DocumentType string, poles []string, wheres []IDocumentWhere) (map[string]map[string]interface{}, error)
	NewDocument(TransactionUID string, DocumentType string, poles map[string]interface{}) (string, error)
	SetDocumentPoles(TransactionUID string, DocumentUID string, poles map[string]interface{}) error
}
