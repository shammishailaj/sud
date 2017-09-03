package core

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/crazyprograms/sud/callpull"
	"github.com/crazyprograms/sud/corebase"
	uuid "github.com/satori/go.uuid"
)

type Core struct {
	databaseName     string
	database         *sql.DB
	lockTrancactions *sync.Mutex
	trancactions     map[string]*transaction
	configurator     *Configurator
	lockUsers        *sync.RWMutex
	users            map[string]corebase.IUser
	close            bool
	genUID           chan string
	listenpulls      map[string]IListenPull
	callpulls        map[string]ICallPull
}

func (c *Core) Configurator() *Configurator { return c.configurator }

/*
func (core *Core) GetConfiguration() map[string]*Configuration {
	l := make(map[string]*Configuration, len(core.baseConfiguration))
	for key, value := range core.baseConfiguration {
		l[key] = value
	}
	return l
}
func (core *Core) AddBaseConfiguration(ConfigurationName string, conf *Configuration) bool {
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
}/**/

func (core *Core) getUser(Login string) corebase.IUser {
	core.lockUsers.RLock()
	defer core.lockUsers.RUnlock()
	user, ok := core.users[Login]
	if ok {
		return user
	}
	return nil
}

/*
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
func (core *Core) LoadConfiguration(ConfigurationName string, Access corebase.IAccess) (*Configuration, error) {
	var conf *Configuration
	if conf = core.getConfiguration(ConfigurationName); conf != nil {
		if !conf.CheckAccess(Access) {
			return nil, &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "LoadConfiguration", Name: ConfigurationName}
		}
		return conf, nil
	}
	AccessConfiguration := map[string]bool{}
	loadConfiguration := map[string]*Configuration{}
	depend := &list.List{}
	var addDepend func(ConfigurationName string) error
	addDepend = func(ConfigurationName string) error {
		var lconf *Configuration
		var ok bool
		if lconf, ok = loadConfiguration[ConfigurationName]; !ok {
			lconf = core.loadBaseConfiguration(ConfigurationName)
			if lconf == nil {
				return &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "LoadConfiguration", Name: ConfigurationName}
			}
			for _, a := range lconf.GetAccessList() {
				AccessConfiguration[a] = true
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
		return nil, &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "LoadConfiguration", Name: ConfigurationName}
	}
	AccessList := make([]string, len(AccessConfiguration), len(AccessConfiguration))
	ANum := 0
	for a, _ := range AccessConfiguration {
		AccessList[ANum] = a
		ANum++
	}
	conf = NewConfiguration(AccessList)
	for e := depend.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Configuration)
		for RecordType, typeInfo := range c.typesInfo {
			conf.typesInfo[RecordType] = typeInfo
		}
		for CallName, callInfo := range c.callsInfo {
			conf.callsInfo[CallName] = callInfo
		}
		for RecordType, poleTypeMap := range c.polesInfo {
			ptm, ok := conf.polesInfo[RecordType]
			if !ok {
				ptm = make(map[string]corebase.IPoleInfo)
				conf.polesInfo[RecordType] = ptm
			}
			for PoleName, poleInfo := range poleTypeMap {
				ptm[PoleName] = poleInfo
			}
		}
	}
	core.addFullConfiguration(ConfigurationName, conf)
	return conf, nil
}
*/
func (core *Core) Congigurator() *Configurator {
	return core.configurator
}

//getConfiguration
func NewCore(DatabaseName string, ConnectionString string) (*Core, error) {
	var err error
	core := &Core{
		databaseName:     DatabaseName,
		genUID:           make(chan string),
		lockTrancactions: &sync.Mutex{},
		trancactions:     make(map[string]*transaction),
		configurator:     NewConfigurator(),
		lockUsers:        &sync.RWMutex{},
		close:            false,
		listenpulls:      map[string]IListenPull{},
		callpulls:        map[string]ICallPull{},
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
	core.CreateDatabase()
	if err = core.LoadUsers(); err != nil {
		return nil, err
	}
	core.users["Test"] = &User{Login: "Test", HashPassword: GenHashPassword("Test"), Access: map[string]bool{"CheckConfiguration": true, "Storage": true, "Configuration": true}}
	core.users["Storage"] = &User{Login: "Storage", HashPassword: GenHashPassword("Test"), Access: map[string]bool{"StorageEngine": true, "Storage": true}}
	return core, nil
}
func (core *Core) LoadUsers() error {
	var err error
	core.lockUsers.Lock()
	defer core.lockUsers.Unlock()
	if core.users, err = getUsers(core.database); err != nil {
		return err
	}
	return nil
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
	return nil, &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "GetTransaction", Name: TransactionUID}
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
		return &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "RollbackTransaction", Name: TransactionUID}
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
		return &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "CommitTransaction", Name: TransactionUID}
	}
	err := t.tx.Commit()
	if err != nil {
		log.Println(err)
	}
	return err
}
func (core *Core) Listen(ConfigurationName string, Name string, TimeoutWait time.Duration, Access corebase.IAccess) (Param map[string]interface{}, AccessClient corebase.IAccess, Result chan callpull.Result, errResult error) {
	var err error
	var ok bool
	var config *Configuration
	var callinfo corebase.ICallInfo
	if config, err = core.configurator.GetConfiguration(ConfigurationName, Access); err != nil {
		return nil, nil, nil, err
	}
	if callinfo, err = config.CallInfo(Name); err != nil {
		return nil, nil, nil, err
	}
	if !Access.CheckAccess(callinfo.GetAccessListen()) {
		return nil, nil, nil, &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "Listen:Call", Name: Name}
	}
	pullName := callinfo.GetPullName()
	var listenpull IListenPull
	if listenpull, ok = core.listenpulls[pullName]; !ok {
		return nil, nil, nil, &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "Listen:Pull", Name: pullName}
	}
	return listenpull.Listen(Name, TimeoutWait)
}
func (core *Core) Call(ConfigurationName string, Name string, Params map[string]interface{}, TimeoutWait time.Duration, Access corebase.IAccess) (callpull.Result, error) {
	var err error
	var ok bool
	var config *Configuration
	var callinfo corebase.ICallInfo
	if config, err = core.configurator.GetConfiguration(ConfigurationName, Access); err != nil {
		return callpull.Result{Result: nil}, err
	}
	if callinfo, err = config.CallInfo(Name); err != nil {
		return callpull.Result{Result: nil}, err
	}
	if !Access.CheckAccess(callinfo.GetAccessCall()) {
		return callpull.Result{Result: nil}, &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "Call", Name: Name}
	}
	pullName := callinfo.GetPullName()
	var cp ICallPull
	if cp, ok = core.callpulls[pullName]; !ok {
		return callpull.Result{Result: nil}, &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "Call:Pull", Name: pullName}
	}
	return cp.Call(Name, Params, TimeoutWait, Access)
}
func (core *Core) GetRecordsPoles(TransactionUID string, ConfigurationName string, RecordType string, poles []string, wheres []corebase.IRecordWhere, Access corebase.IAccess) (map[corebase.UUID]map[string]interface{}, error) {
	var err error
	var tx *transaction
	if tx, err = core.getTransaction(TransactionUID); err != nil {
		return nil, err
	}
	return tx.GetRecordsPoles(ConfigurationName, RecordType, poles, wheres, Access)
}
func (core *Core) SetRecordPoles(TransactionUID string, ConfigurationName string, RecordUID corebase.UUID, poles map[string]interface{}, Access corebase.IAccess) error {
	var err error
	var tx *transaction
	if tx, err = core.getTransaction(TransactionUID); err != nil {
		return err
	}
	return tx.SetRecordPoles(ConfigurationName, RecordUID, poles, Access)
}
func (core *Core) NewRecord(TransactionUID string, ConfigurationName string, RecordType string, Poles map[string]interface{}, Access corebase.IAccess) (corebase.UUID, error) {
	var err error
	var tx *transaction
	if tx, err = core.getTransaction(TransactionUID); err != nil {
		return corebase.NullUUID, err
	}
	return tx.NewRecord(ConfigurationName, RecordType, Poles, Access)
}
func (core *Core) GetRecordAccess(TransactionUID string, ConfigurationName string, RecordUID corebase.UUID, Access corebase.IAccess) (string, error) {
	var err error
	var tx *transaction
	if tx, err = core.getTransaction(TransactionUID); err != nil {
		return "", err
	}
	return tx.GetRecordAccess(ConfigurationName, RecordUID, Access)
}
func (core *Core) SetRecordAccess(TransactionUID string, ConfigurationName string, RecordUID corebase.UUID, NewAccess string, Access corebase.IAccess) error {
	var err error
	var tx *transaction
	if tx, err = core.getTransaction(TransactionUID); err != nil {
		return err
	}
	return tx.SetRecordAccess(ConfigurationName, RecordUID, NewAccess, Access)
}
