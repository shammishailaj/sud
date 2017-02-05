package core

import (
	"container/list"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/crazyprograms/callpull"
	uuid "github.com/satori/go.uuid"
)

type Server struct {
	databaseName          string
	database              *sql.DB
	lockTrancactions      *sync.Mutex
	trancactions          map[string]*transaction
	lockBaseConfiguration *sync.Mutex
	baseConfiguration     map[string]*Configuration
	lockFullConfiguration *sync.Mutex
	fullConfiguration     map[string]*Configuration
	lockUsers             *sync.Mutex
	users                 map[string]IUser
	close                 bool
	genUID                chan string
	listenpulls           map[string]IListenPull
	callpulls             map[string]ICallPull
}

func (server *Server) addBaseConfiguration(ConfigurationName string, conf *Configuration) bool {
	server.lockBaseConfiguration.Lock()
	defer server.lockBaseConfiguration.Unlock()
	_, ok := server.baseConfiguration[ConfigurationName]
	if ok {
		return false
	}
	server.baseConfiguration[ConfigurationName] = conf
	return true
}
func (server *Server) addFullConfiguration(ConfigurationName string, conf *Configuration) bool {
	server.lockFullConfiguration.Lock()
	defer server.lockFullConfiguration.Unlock()
	_, ok := server.fullConfiguration[ConfigurationName]
	if ok {
		return false
	}
	server.fullConfiguration[ConfigurationName] = conf
	return true
}

func (server *Server) getUser(UserName string) IUser {
	user, ok := server.users[UserName]
	if ok {
		return user
	}
	return nil
}
func (server *Server) loadBaseConfiguration(ConfigurationName string) *Configuration {
	server.lockBaseConfiguration.Lock()
	defer server.lockBaseConfiguration.Unlock()
	if conf, ok := server.baseConfiguration[ConfigurationName]; ok {
		return conf
	}
	return nil
}
func (server *Server) getConfiguration(ConfigurationName string) *Configuration {
	server.lockFullConfiguration.Lock()
	defer server.lockFullConfiguration.Unlock()
	conf, ok := server.fullConfiguration[ConfigurationName]
	if ok {
		return conf
	}
	return nil
}
func (server *Server) LoadConfiguration(ConfigurationName string) (*Configuration, error) {
	var conf *Configuration
	if conf = server.getConfiguration(ConfigurationName); conf != nil {
		return conf, nil
	}
	conf = NewConfiguration()
	loadConfiguration := make(map[string]*Configuration)
	depend := &list.List{}
	var addDepend func(ConfigurationName string) error
	addDepend = func(ConfigurationName string) error {
		var lconf *Configuration
		var ok bool
		if lconf, ok = loadConfiguration[ConfigurationName]; !ok {
			lconf = server.loadBaseConfiguration(ConfigurationName)
			if lconf == nil {
				return errors.New("configuration not found: " + ConfigurationName)
			}
			loadConfiguration[ConfigurationName] = lconf
			for i := 0; i < len(lconf.dependConfigurationName); i++ {
				if err := addDepend(lconf.dependConfigurationName[i]); err != nil {
					return err
				}
			}
			fmt.Println(lconf)
			depend.PushBack(lconf)
		}
		return nil
	}
	if err := addDepend(ConfigurationName); err != nil {
		return nil, err
	}
	if depend.Front() == nil {
		return nil, errors.New("Configuration not found: " + ConfigurationName)
	}
	for e := depend.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Configuration)
		for DocumentType, typeInfo := range c.typesInfo {
			conf.typesInfo[DocumentType] = typeInfo
		}
		for CallName, callInfo := range c.callsInfo {
			conf.callsInfo[CallName] = callInfo
		}
		for DocumentType, poleTypeMap := range c.polesInfo {
			ptm, ok := conf.polesInfo[DocumentType]
			if !ok {
				ptm = make(map[string]IPoleInfo)
				conf.polesInfo[DocumentType] = ptm
			}
			for PoleName, poleInfo := range poleTypeMap {
				ptm[PoleName] = poleInfo
			}
		}
	}
	server.addFullConfiguration(ConfigurationName, conf)
	return conf, nil
}

//getConfiguration
func NewServer(DatabaseName string, ConnectionString string) (*Server, error) {
	var err error
	async := callpull.NewCallPull()
	server := &Server{
		databaseName:          DatabaseName,
		genUID:                make(chan string),
		lockTrancactions:      &sync.Mutex{},
		trancactions:          make(map[string]*transaction),
		lockBaseConfiguration: &sync.Mutex{},
		baseConfiguration:     make(map[string]*Configuration),
		lockFullConfiguration: &sync.Mutex{},
		fullConfiguration:     make(map[string]*Configuration),
		lockUsers:             &sync.Mutex{},
		users:                 make(map[string]IUser),
		close:                 false,
		listenpulls:           map[string]IListenPull{"async": async},
		callpulls:             map[string]ICallPull{"async": async, "std": StdCallPull},
	}
	server.database, err = sql.Open("postgres", ConnectionString)
	if err != nil {
		return nil, err
	}
	err = server.database.Ping()
	if err != nil {
		server.database.Close()
		log.Fatalln(err)
		return nil, err
	}
	go server.gogenUID()
	server.loadDefaultsConfiguration()
	server.users["Test"] = &User{UserName: "Test", HashPassword: GenHashPassword("Test"), Access: map[string]bool{"CheckConfiguration": true}}
	return server, nil
}
func (server *Server) gogenUID() {
	for !server.close {
		server.genUID <- uuid.NewV4().String()
	}
}
func (server *Server) getTransaction(TransactionUID string) (*transaction, error) {
	server.lockTrancactions.Lock()
	defer server.lockTrancactions.Unlock()
	tx, ok := server.trancactions[TransactionUID]
	if ok {
		return tx, nil
	}
	return nil, errors.New("transaction not found: " + TransactionUID)
}
func (server *Server) BeginTransaction() (string, error) {
	uid := <-server.genUID
	tx, err := server.database.Begin()
	if err != nil {
		return "", nil
	}
	server.lockTrancactions.Lock()
	defer server.lockTrancactions.Unlock()
	server.trancactions[uid] = &transaction{tx: tx, server: server}
	return uid, nil
}
func (server *Server) RollbackTransaction(TransactionUID string) {
	server.lockTrancactions.Lock()
	defer server.lockTrancactions.Unlock()
	t, ok := server.trancactions[TransactionUID]
	if ok {
		defer delete(server.trancactions, TransactionUID)
	} else {
		return
	}
	err := t.tx.Rollback()
	if err != nil {
		log.Println(err)
	}
}
func (server *Server) CommitTransaction(TransactionUID string) {
	server.lockTrancactions.Lock()
	defer server.lockTrancactions.Unlock()
	t, ok := server.trancactions[TransactionUID]
	if ok {
		defer delete(server.trancactions, TransactionUID)
	} else {
		return
	}
	err := t.tx.Commit()
	if err != nil {
		log.Println(err)
	}
}

func (server *Server) Listen(ConfigurationName string, Name string, TimeoutWait time.Duration) (Param map[string]interface{}, Result chan interface{}, errResult error) {
	var err error
	var ok bool
	var config *Configuration
	var callinfo ICallInfo
	if config, err = server.LoadConfiguration(ConfigurationName); err != nil {
		return nil, nil, err
	}
	if callinfo, err = config.GetCallInfo(Name); err != nil {
		return nil, nil, err
	}
	pullName := callinfo.GetPullName()
	var listenpull IListenPull
	if listenpull, ok = server.listenpulls[pullName]; !ok {
		return nil, nil, errors.New("listen pull not found: " + pullName)
	}
	return listenpull.Listen(Name, TimeoutWait)
}
func (server *Server) Call(ConfigurationName string, Name string, Params map[string]interface{}, TimeoutWait time.Duration) (interface{}, error) {
	var err error
	var ok bool
	var config *Configuration
	var callinfo ICallInfo
	if config, err = server.LoadConfiguration(ConfigurationName); err != nil {
		return nil, err
	}
	if callinfo, err = config.GetCallInfo(Name); err != nil {
		return nil, err
	}
	pullName := callinfo.GetPullName()
	var callpull ICallPull
	if callpull, ok = server.callpulls[pullName]; !ok {
		return nil, errors.New("call pull not found: " + pullName)
	}
	return callpull.Call(Name, Params, TimeoutWait)
}
