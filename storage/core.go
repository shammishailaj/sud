package storage

import (
	"time"

	"errors"

	"github.com/crazyprograms/sud/core"
)

func stdGetStream(cr *core.Core, Name string, Param map[string]interface{}, timeOutWait time.Duration) (interface{}, error) {
	var ok bool
	var err error
	var HashI interface{}
	var Hash string
	var TransactionUIDI interface{}
	var TransactionUID string
	if HashI, ok = Param["Hash"]; !ok {
		return nil, errors.New("not found parameter Storage")
	}
	if Hash, ok = HashI.(string); !ok {
		return nil, errors.New("parameter Storage is not string")
	}
	if TransactionUIDI, ok = Param["TransactionUID"]; !ok {
		TransactionUID, err = cr.BeginTransaction()
		if err != nil {
			return nil, err
		}
		defer cr.CommitTransaction(TransactionUID)
	} else {
		if TransactionUID, ok = TransactionUIDI.(string); !ok {
			return nil, errors.New("parameter Storage is not string")
		}
	}
	var docs []core.IDocument
	if docs, err = cr.GetDocuments(TransactionUID, "Storage", "Storage.Stream", []string{"Storage.Stream.*"}, []core.IDocumentWhere{&core.DocumentWhereCompare{PoleName: "Storage.Stream.Hash", Operation: "Equally", Value: core.NewObject(Hash)}}); err != nil {
		return nil, err
	}
	if len(docs) != 1 {
		return nil, errors.New("stream not found: " + Hash)
	}
	var Storage string
	if Storage, err = docs[0].GetPole("Storage.Stream.Storage").String(); err != nil {
		return nil, err
	}
	return cr.Call("Storage."+Storage, "Storage."+Storage+".GetStream", Param, timeOutWait)
}
func stdSetStream(cr *core.Core, Name string, Param map[string]interface{}, timeOutWait time.Duration) (interface{}, error) {
	var ok bool
	var StorageI interface{}
	var Storage string
	if StorageI, ok = Param["Storage"]; !ok {
		return nil, errors.New("not found parameter Storage")
	}
	if Storage, ok = StorageI.(string); !ok {
		return nil, errors.New("parameter Storage is not string")
	}
	return cr.Call("Storage."+Storage, "Storage."+Storage+".SetStream", Param, timeOutWait)
}
func init() {
	initConfiguration()
	core.AddStdCall("Storage.GetStream", stdGetSetStream)
	core.AddStdCall("Storage.SetStream", stdSetStream)
}
func initConfiguration() {
	conf := core.NewConfiguration()
	conf.AddCall("Storage", "Storage.GetStream", "std", true, false, "")
	conf.AddCall("Storage", "Storage.SetStream", "std", true, false, "")

	conf.AddType("Storage", "Storage.Stream", true, true, true, "Поток")
	conf.AddPole("Storage", "Storage.Stream", "Storage.Stream.Hash", "StringValue", core.NewObject(nil), "Unique", &core.PoleCheckerStringValue{}, true, true, "Хеш потока")
	conf.AddPole("Storage", "Storage.Stream", "Storage.Stream.Size", "Int64Value", core.NewObject(nil), "None", &core.PoleCheckerInt64Value{}, true, true, "Размер в байтах")
	conf.AddPole("Storage", "Storage.Stream", "Storage.Stream.Storage", "StringValue", core.NewObject(nil), "None", &core.PoleCheckerStringValue{}, true, true, "Имя хранилища в котором хранится")
	core.InitAddBaseConfiguration("Storage", conf)
}
func initDefaultStorageConfiguration() {
	conf := core.NewConfiguration()
	conf.AddCall("Storage.Default", "Storage.Default.GetStream", "async", true, true, "")
	conf.AddCall("Storage.Default", "Storage.Default.SetStream", "async", true, true, "")
	conf.AddDependConfiguration("Storage")
	core.InitAddBaseConfiguration("Storage.Default", conf)
}
