package storage

import (
	"time"

	"errors"

	"github.com/crazyprograms/sud/callpull"
	"github.com/crazyprograms/sud/core"
	"github.com/crazyprograms/sud/corebase"
	"github.com/crazyprograms/sud/sortex"
)

func stdGetStream(cr *core.Core, Name string, Param map[string]interface{}, timeOutWait time.Duration, Access corebase.IAccess) (callpull.Result, error) {
	var ok bool
	var err error
	var HashI interface{}
	var Hash string
	var TransactionUIDI interface{}
	var TransactionUID string
	if HashI, ok = Param["Hash"]; !ok {
		return callpull.Result{Result: nil, Error: errors.New("not found parameter Storage")}, nil
	}
	if Hash, ok = HashI.(string); !ok {
		return callpull.Result{Result: nil, Error: errors.New("parameter Storage is not string")}, nil
	}
	if TransactionUIDI, ok = Param["TransactionUID"]; !ok {
		TransactionUID, err = cr.BeginTransaction()
		if err != nil {
			return callpull.Result{Result: nil, Error: err}, nil
		}
		defer cr.CommitTransaction(TransactionUID)
		Param["TransactionUID"] = TransactionUID
	} else {
		if TransactionUID, ok = TransactionUIDI.(string); !ok {
			return callpull.Result{Result: nil, Error: errors.New("parameter Storage is not string")}, nil
		}
	}
	var docs map[corebase.UUID]map[string]interface{}
	if docs, err = cr.GetRecordsPoles(TransactionUID, "Storage", "Storage.Stream", []string{"Storage.Stream.*"}, []corebase.IRecordWhere{&corebase.RecordWhereCompare{PoleName: "Storage.Stream.Hash", Operation: "Equally", Value: Hash}}, Access); err != nil {
		return callpull.Result{Result: nil, Error: err}, nil
	}
	if len(docs) == 0 {
		return callpull.Result{Result: nil, Error: errors.New("stream not found: " + Hash)}, nil
	}
	Storages := make([]string, len(docs), len(docs))
	Prioritets := map[string]int64{}
	n := 0
	for _, doc := range docs {
		var StorageP string
		if StorageP, ok = doc["Storage.Stream.Storage"].(string); !ok {
			return callpull.Result{Result: nil, Error: errors.New("not found Storage.Stream.Storage ")}, nil
		}
		var PriorityP int64
		if PriorityP, ok = doc["Storage.Stream.Priority"].(int64); !ok {
			return callpull.Result{Result: nil, Error: errors.New("not found Storage.Stream.Priority ")}, nil
		}
		Prioritets[StorageP] = PriorityP
		Storages[n] = StorageP
		n++
	}
	sortex.SortStrings(Storages, func(i, j int) bool { return Prioritets[Storages[i]] < Prioritets[Storages[j]] })
	var r callpull.Result
	for i := 0; i < len(Storages); i++ {
		r, err = cr.Call("Storage."+Storages[i], "Storage."+Storages[i]+".GetStream", Param, timeOutWait, Access)
		if err == nil {
			return r, nil
		}
	}
	return r, err
}
func stdSetStream(cr *core.Core, Name string, Param map[string]interface{}, timeOutWait time.Duration, Access corebase.IAccess) (callpull.Result, error) {
	var ok bool
	var StorageI interface{}
	var Storage string
	if StorageI, ok = Param["Storage"]; !ok {
		return callpull.Result{Result: nil, Error: errors.New("not found parameter Storage")}, nil
	}
	if Storage, ok = StorageI.(string); !ok {
		return callpull.Result{Result: nil, Error: errors.New("parameter Storage is not string")}, nil
	}
	return cr.Call("Storage."+Storage, "Storage."+Storage+".SetStream", Param, timeOutWait, Access)
}

func initConfiguration(c *core.Core) bool {
	conf := core.NewConfiguration([]string{"Storage"})
	conf.AddCall(core.CallInfo{ConfigurationName: "Storage", Name: "Storage.GetStream", PullName: "std", AccessCall: "Storage", AccessListen: "", Title: "Получить поток"})
	conf.AddCall(core.CallInfo{ConfigurationName: "Storage", Name: "Storage.SetStream", PullName: "std", AccessCall: "Storage", AccessListen: "", Title: "Передать поток"})

	conf.AddType(core.TypeInfo{ConfigurationName: "Storage", RecordType: "Storage.Stream", AccessNew: "StorageEngine", AccessRead: "Storage", AccessSave: "StorageEngine", Title: "Поток"})
	conf.AddPole(core.PoleInfo{ConfigurationName: "Storage", RecordType: "Storage.Stream", PoleName: "Storage.Stream.Hash", PoleType: "StringValue", Default: corebase.NULL, IndexType: "Index", Checker: &core.PoleCheckerStringValue{}, AccessRead: "Storage", AccessWrite: "StorageEngine", Title: "Хеш потока"})
	conf.AddPole(core.PoleInfo{ConfigurationName: "Storage", RecordType: "Storage.Stream", PoleName: "Storage.Stream.Size", PoleType: "Int64Value", Default: corebase.NULL, IndexType: "None", Checker: &core.PoleCheckerInt64Value{}, AccessRead: "Storage", AccessWrite: "StorageEngine", Title: "Размер в байтах"})
	conf.AddPole(core.PoleInfo{ConfigurationName: "Storage", RecordType: "Storage.Stream", PoleName: "Storage.Stream.Priority", PoleType: "Int64Value", Default: corebase.NULL, IndexType: "None", Checker: &core.PoleCheckerInt64Value{}, AccessRead: "Storage", AccessWrite: "StorageEngine", Title: "Приоритет"})
	conf.AddPole(core.PoleInfo{ConfigurationName: "Storage", RecordType: "Storage.Stream", PoleName: "Storage.Stream.Storage", PoleType: "StringValue", Default: corebase.NULL, IndexType: "None", Checker: &core.PoleCheckerStringValue{}, AccessRead: "Storage", AccessWrite: "StorageEngine", Title: "Имя хранилища в котором хранится"})
	return c.AddBaseConfiguration("Storage", conf)
}
func initDefaultStorageConfiguration(c *core.Core) bool {
	conf := core.NewConfiguration([]string{"StorageEngine"})
	conf.AddCall(core.CallInfo{ConfigurationName: "Storage.Default", Name: "Storage.Default.GetStream", PullName: "async", AccessCall: "Storage", AccessListen: "StorageEngine", Title: ""})
	conf.AddCall(core.CallInfo{ConfigurationName: "Storage.Default", Name: "Storage.Default.SetStream", PullName: "async", AccessCall: "Storage", AccessListen: "StorageEngine", Title: ""})
	conf.AddDependConfiguration("Storage")
	return c.AddBaseConfiguration("Storage.Default", conf)
}
func InitModule(c *core.Core) error {
	var err error
	initConfiguration(c)
	initDefaultStorageConfiguration(c)
	if err = core.AddStdCall("Storage.GetStream", stdGetStream); err != nil {
		return err
	}
	if err = core.AddStdCall("Storage.SetStream", stdSetStream); err != nil {
		return err
	}
	return nil
}
