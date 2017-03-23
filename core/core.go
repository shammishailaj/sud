package core

import (
	"container/list"
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/crazyprograms/callpull"
	"github.com/crazyprograms/sud/corebase"
	uuid "github.com/satori/go.uuid"
)

type Core struct {
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

func (core *Core) addBaseConfiguration(ConfigurationName string, conf *Configuration) bool {
	core.lockBaseConfiguration.Lock()
	defer core.lockBaseConfiguration.Unlock()
	_, ok := core.baseConfiguration[ConfigurationName]
	if ok {
		return false
	}
	core.baseConfiguration[ConfigurationName] = conf
	return true
}
func (core *Core) addFullConfiguration(ConfigurationName string, conf *Configuration) bool {
	core.lockFullConfiguration.Lock()
	defer core.lockFullConfiguration.Unlock()
	_, ok := core.fullConfiguration[ConfigurationName]
	if ok {
		return false
	}
	core.fullConfiguration[ConfigurationName] = conf
	return true
}

func (core *Core) getUser(UserName string) IUser {
	user, ok := core.users[UserName]
	if ok {
		return user
	}
	return nil
}
func (core *Core) loadBaseConfiguration(ConfigurationName string) *Configuration {
	core.lockBaseConfiguration.Lock()
	defer core.lockBaseConfiguration.Unlock()
	if conf, ok := core.baseConfiguration[ConfigurationName]; ok {
		return conf
	}
	return nil
}
func (core *Core) getConfiguration(ConfigurationName string) *Configuration {
	core.lockFullConfiguration.Lock()
	defer core.lockFullConfiguration.Unlock()
	conf, ok := core.fullConfiguration[ConfigurationName]
	if ok {
		return conf
	}
	return nil
}
func (core *Core) LoadConfiguration(ConfigurationName string) (*Configuration, error) {
	var conf *Configuration
	if conf = core.getConfiguration(ConfigurationName); conf != nil {
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
			lconf = core.loadBaseConfiguration(ConfigurationName)
			if lconf == nil {
				return errors.New("configuration not found: " + ConfigurationName)
			}
			loadConfiguration[ConfigurationName] = lconf
			for i := 0; i < len(lconf.dependConfigurationName); i++ {
				if err := addDepend(lconf.dependConfigurationName[i]); err != nil {
					return err
				}
			}
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
				ptm = make(map[string]corebase.IPoleInfo)
				conf.polesInfo[DocumentType] = ptm
			}
			for PoleName, poleInfo := range poleTypeMap {
				ptm[PoleName] = poleInfo
			}
		}
	}
	core.addFullConfiguration(ConfigurationName, conf)
	return conf, nil
}

//getConfiguration
func NewCore(DatabaseName string, ConnectionString string) (*Core, error) {
	var err error
	core := &Core{
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
		listenpulls:           map[string]IListenPull{},
		callpulls:             map[string]ICallPull{},
	}
	async := callpull.NewCallPull()
	core.listenpulls["async"] = async
	core.callpulls["async"] = async
	core.callpulls["std"] = GetStdPull(core)
	core.database, err = sql.Open("postgres", ConnectionString)
	if err != nil {
		return nil, err
	}
	err = core.database.Ping()
	if err != nil {
		core.database.Close()
		log.Fatalln(err)
		return nil, err
	}
	go core.gogenUID()
	core.loadDefaultsConfiguration()
	core.users["Test"] = &User{UserName: "Test", HashPassword: GenHashPassword("Test"), Access: map[string]bool{"CheckConfiguration": true}}
	return core, nil
}
func (core *Core) gogenUID() {
	for !core.close {
		core.genUID <- uuid.NewV4().String()
	}
}
func (core *Core) getTransaction(TransactionUID string) (*transaction, error) {
	core.lockTrancactions.Lock()
	defer core.lockTrancactions.Unlock()
	tx, ok := core.trancactions[TransactionUID]
	if ok {
		return tx, nil
	}
	return nil, errors.New("transaction not found: " + TransactionUID)
}
func (core *Core) BeginTransaction() (string, error) {
	uid := <-core.genUID
	tx, err := core.database.Begin()
	if err != nil {
		return "", nil
	}
	core.lockTrancactions.Lock()
	defer core.lockTrancactions.Unlock()
	core.trancactions[uid] = &transaction{tx: tx, core: core}
	return uid, nil
}
func (core *Core) RollbackTransaction(TransactionUID string) error {
	core.lockTrancactions.Lock()
	defer core.lockTrancactions.Unlock()
	t, ok := core.trancactions[TransactionUID]
	if ok {
		defer delete(core.trancactions, TransactionUID)
	} else {
		return errors.New("transaction not found")
	}
	err := t.tx.Rollback()
	if err != nil {
		log.Println(err)
	}
	return err
}
func (core *Core) CommitTransaction(TransactionUID string) error {
	core.lockTrancactions.Lock()
	defer core.lockTrancactions.Unlock()
	t, ok := core.trancactions[TransactionUID]
	if ok {
		defer delete(core.trancactions, TransactionUID)
	} else {
		return errors.New("transaction not found")
	}
	err := t.tx.Commit()
	if err != nil {
		log.Println(err)
	}
	return err
}
func (core *Core) Listen(ConfigurationName string, Name string, TimeoutWait time.Duration) (Param map[string]interface{}, Result chan callpull.Result, errResult error) {
	var err error
	var ok bool
	var config *Configuration
	var callinfo corebase.ICallInfo
	if config, err = core.LoadConfiguration(ConfigurationName); err != nil {
		return nil, nil, err
	}
	if callinfo, err = config.GetCallInfo(Name); err != nil {
		return nil, nil, err
	}
	pullName := callinfo.GetPullName()
	var listenpull IListenPull
	if listenpull, ok = core.listenpulls[pullName]; !ok {
		return nil, nil, errors.New("listen pull not found: " + pullName)
	}
	return listenpull.Listen(Name, TimeoutWait)
}
func (core *Core) Call(ConfigurationName string, Name string, Params map[string]interface{}, TimeoutWait time.Duration) (callpull.Result, error) {
	var err error
	var ok bool
	var config *Configuration
	var callinfo corebase.ICallInfo
	if config, err = core.LoadConfiguration(ConfigurationName); err != nil {
		return callpull.Result{Result: nil}, err
	}
	if callinfo, err = config.GetCallInfo(Name); err != nil {
		return callpull.Result{Result: nil}, err
	}
	pullName := callinfo.GetPullName()
	var cp ICallPull
	if cp, ok = core.callpulls[pullName]; !ok {
		return callpull.Result{Result: nil}, errors.New("call pull not found: " + pullName)
	}
	return cp.Call(Name, Params, TimeoutWait)
}
func (core *Core) GetDocumentsPoles(TransactionUID string, ConfigurationName string, DocumentType string, poles []string, wheres []corebase.IDocumentWhere) (map[string]map[string]interface{}, error) {
	var err error
	var tx *transaction
	if tx, err = core.getTransaction(TransactionUID); err != nil {
		return nil, err
	}
	return tx.GetDocumentsPoles(ConfigurationName, DocumentType, poles, wheres)
}
func (core *Core) SetDocumentPoles(TransactionUID string, ConfigurationName string, DocumentUID string, poles map[string]interface{}) error {
	var err error
	var tx *transaction
	if tx, err = core.getTransaction(TransactionUID); err != nil {
		return err
	}
	return tx.SetDocumentPoles(ConfigurationName, DocumentUID, poles)
}
func (core *Core) NewDocument(TransactionUID string, ConfigurationName string, DocumentType string, Poles map[string]interface{}) (string, error) {
	var err error
	var tx *transaction
	if tx, err = core.getTransaction(TransactionUID); err != nil {
		return "", err
	}
	return tx.NewDocument(ConfigurationName, DocumentType, Poles)
}
