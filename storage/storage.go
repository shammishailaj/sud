package storage

import (
	"io"
	"io/ioutil"
	"time"

	"errors"

	"sync"

	"encoding/hex"

	"os"
	"path"

	"crypto/sha1"

	"github.com/crazyprograms/sud/callpull"
	"github.com/crazyprograms/sud/client"
	"github.com/crazyprograms/sud/corebase"
)

type Storage struct {
	lockRW       sync.RWMutex
	tmpDir       string
	root         string
	close        bool
	errGetStream error
	name         string
	client       client.IClient
}

const tempDir = "tmp"
const storageDir = "storage"

func (storage *Storage) getStream(TransactionUID string, Hash string) ([]byte, error) {
	p := path.Join(storage.root, storageDir, Hash[len(Hash)-2:], Hash[len(Hash)-4:len(Hash)-2], Hash)
	storage.lockRW.RLock()
	defer storage.lockRW.RUnlock()
	data, err := ioutil.ReadFile(p)
	return data, err
}
func (storage *Storage) loopGetStream() {
	//var err error
	var ok bool
	for !storage.close {
		Param, ResultUID, err := storage.client.Listen("Storage."+storage.name+".GetStream", time.Second*10)
		if err == callpull.ErrorTimeout {
			continue
		}
		if err != nil {
			storage.errGetStream = err
			return
		}
		var TransactionUID string
		var Hash string
		var value interface{}
		if value, ok = Param["TransactionUID"]; !ok {
			storage.client.ListenResult(ResultUID, callpull.Result{Result: nil, Error: errors.New("not found param TransactionUID")}, nil)
			continue
		}
		if TransactionUID, ok = value.(string); !ok {
			storage.client.ListenResult(ResultUID, callpull.Result{Result: nil, Error: errors.New("not string param TransactionUID")}, nil)
			continue
		}
		if value, ok = Param["Hash"]; !ok {
			storage.client.ListenResult(ResultUID, callpull.Result{Result: nil, Error: errors.New("not found param Hash")}, nil)
			continue
		}
		if Hash, ok = value.(string); !ok {
			storage.client.ListenResult(ResultUID, callpull.Result{Result: nil, Error: errors.New("not string param Hash")}, nil)
			continue
		}
		var r interface{}
		if r, err = storage.getStream(TransactionUID, Hash); err != nil {
			storage.client.ListenResult(ResultUID, callpull.Result{Result: nil, Error: err}, nil)
			continue
		}
		storage.client.ListenResult(ResultUID, callpull.Result{Result: r, Error: nil}, nil)
	}
}
func (storage *Storage) setStream(TransactionUID string, Stream []byte) (string, error) {
	var err error
	tmpfile, err := ioutil.TempFile(storage.tmpDir, "set")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name()) // clean up
	defer tmpfile.Close()
	hashStream := sha1.New()
	w := io.MultiWriter(tmpfile, hashStream)
	if _, err = w.Write(Stream); err != nil {
		return "", err
	}
	if err = tmpfile.Close(); err != nil {
		return "", err
	}
	var hash string
	hash = hex.EncodeToString(hashStream.Sum([]byte{}))
	p := path.Join(storage.root, storageDir, hash[len(hash)-2:], hash[len(hash)-4:len(hash)-2], hash)
	var docs map[corebase.UUID]map[string]interface{}

	if docs, err = storage.client.GetRecordsPoles(TransactionUID, "Storage.Stream", []string{"Storage.Stream.*"}, []corebase.IRecordWhere{
		&corebase.RecordWhereCompare{PoleName: "Storage.Stream.Hash", Operation: "Equally", Value: hash},
	}, ""); err != nil {
		return "", err
	}
	for _, Poles := range docs {
		if !corebase.IsNull(Poles["Storage.Stream.Storage"]) && Poles["Storage.Stream.Storage"].(string) == storage.name {
			return hash, nil
		}
	}

	storage.lockRW.Lock()
	defer storage.lockRW.Unlock()
	os.Rename(tmpfile.Name(), p)
	if _, err = storage.client.NewRecord(TransactionUID, "Storage.Stream", map[string]interface{}{
		"Storage.Stream.Hash":     hash,
		"Storage.Stream.Storage":  storage.name,
		"Storage.Stream.Priority": int64(0),
		"Storage.Stream.Size":     int64(len(Stream)),
	}, ""); err != nil {
		return "", err
	}
	return hash, nil
}

func (storage *Storage) loopSetStream() {
	//var err error
	var ok bool
	for !storage.close {
		Param, ResultUID, err := storage.client.Listen("Storage."+storage.name+".SetStream", time.Second*10)
		if err == callpull.ErrorTimeout {
			continue
		}
		if err != nil {
			storage.errGetStream = err
			return
		}
		var TransactionUID string
		var Stream []byte
		var value interface{}
		if value, ok = Param["TransactionUID"]; !ok {
			storage.client.ListenResult(ResultUID, callpull.Result{Result: "", Error: errors.New("not found param TransactionUID")}, nil)
			continue
		}
		if TransactionUID, ok = value.(string); !ok {
			storage.client.ListenResult(ResultUID, callpull.Result{Result: "", Error: errors.New("not string param TransactionUID")}, nil)
			continue
		}
		if value, ok = Param["Stream"]; !ok {
			storage.client.ListenResult(ResultUID, callpull.Result{Result: "", Error: errors.New("not found param Hash")}, nil)
			continue
		}
		if Stream, ok = value.([]byte); !ok {
			storage.client.ListenResult(ResultUID, callpull.Result{Result: "", Error: errors.New("not string param Hash")}, nil)
			continue
		}

		var Hash string
		if Hash, err = storage.setStream(TransactionUID, Stream); err != nil {
			storage.client.ListenResult(ResultUID, "", err)
			continue
		}
		storage.client.ListenResult(ResultUID, Hash, nil) /**/
	}
}
func (storage *Storage) start() {
	storage.close = false
	go storage.loopGetStream()
	go storage.loopSetStream()
}
func (storage *Storage) createStorageStructure() error {
	var err error
	storage.tmpDir = path.Join(storage.root, tempDir)
	os.RemoveAll(storage.tmpDir)
	if err = os.Mkdir(storage.tmpDir, 0777); err != nil {
		return err
	}
	err = os.Mkdir(path.Join(storage.root, storageDir), 0777)
	if !os.IsExist(err) {
		for ia := 0; ia <= 255; ia++ {
			a := hex.EncodeToString([]byte{((byte)(ia))})
			os.Mkdir(path.Join(storage.root, storageDir, a), 0777)
			for ib := 0; ib <= 255; ib++ {
				b := hex.EncodeToString([]byte{((byte)(ib))})
				os.Mkdir(path.Join(storage.root, storageDir, a, b), 0777)
			}
		}
	}
	return nil
}
func StartStorage(Name string, client client.IClient, root string) *Storage {
	s := &Storage{client: client, name: Name, root: root}
	s.createStorageStructure()
	s.start()
	return s
}
