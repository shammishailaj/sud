package storage

import (
	"database/sql"
	"errors"
	"strings"

	"log"

	"sync"

	"container/list"

	"fmt"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type IQuery interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
type transaction struct {
	server *Server
	tx     *sql.Tx
}

func (t *transaction) Commit() error {
	return t.tx.Commit()
}
func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}
func (t *transaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}
func (t *transaction) Prepare(query string) (*sql.Stmt, error) {
	return t.tx.Prepare(query)
}
func (t *transaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(query, args...)
}
func (t *transaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(query, args...)
}

type Server struct {
	databaseName          string
	database              *sql.DB
	lockTrancactions      *sync.Mutex
	trancactions          map[string]*transaction
	lockBaseConfiguration *sync.Mutex
	baseConfiguration     map[string]*configuration
	lockFullConfiguration *sync.Mutex
	fullConfiguration     map[string]*configuration
	lockUsers             *sync.Mutex
	users                 map[string]IUser
	close                 bool
	genUID                chan string
}
type Client struct {
	server            *Server
	user              IUser
	configurationName string
}
type configuration struct {
	polesInfo               map[string]map[string]IPoleInfo
	typesInfo               map[string]ITypeInfo
	callsInfo               map[string]ICallInfo
	dependConfigurationName []string
}

func (conf *configuration) getPoleInfo(DucumentType string, PoleName string) (IPoleInfo, error) {
	if m, ok := conf.polesInfo[DucumentType]; ok {
		if pi, ok := m[PoleName]; ok {
			return pi, nil
		}
	}
	return nil, errors.New("pole not found: " + PoleName)
}
func (conf *configuration) getPolesInfo(DocumentType string, Poles []string) map[string]IPoleInfo {
	if m, ok := conf.polesInfo[DocumentType]; ok {
		polesInfo := map[string]IPoleInfo{}
		if len(Poles) != 0 {
			for _, pole := range Poles {
				p := strings.Split(pole, ".")
				// Document.Table.Name.
				if p[len(p)-1] == "" {
					for poleName, info := range m {
						if strings.HasPrefix(poleName, pole) {
							polesInfo[poleName] = info
						}
					}
				} else if p[len(p)-1] == "*" {
					p[len(p)-1] = ""
					pole = strings.Join(p, ".")
					for poleName, info := range m {
						if strings.HasPrefix(poleName, pole) {
							polesInfo[poleName] = info
						}
					}
				} else {
					for poleName, info := range m {
						if poleName == pole {
							polesInfo[poleName] = info
						}
					}
				}
			}
		} else {
			for name, info := range m {
				polesInfo[name] = info
			}
		}
		return polesInfo
	}
	return map[string]IPoleInfo{}
}
func (conf *configuration) addType(ConfigurationName string, DocumentType string, New bool, Read bool, Save bool, Title string) {
	conf.typesInfo[DocumentType] = &TypeInfo{ConfigurationName: ConfigurationName, DocumentType: DocumentType, New: New, Read: Read, Save: Save, Title: Title}
}
func (conf *configuration) addPole(ConfigurationName string, DocumentType string, PoleName string, PoleType string, Default Object, IndexType string, Checker IPoleChecker, New bool, Edit bool, Title string) {
	_, ok := conf.polesInfo[DocumentType]
	if !ok {
		conf.polesInfo[DocumentType] = make(map[string]IPoleInfo)
	}
	conf.polesInfo[DocumentType][PoleName] = &PoleInfo{
		ConfigurationName: ConfigurationName,
		DocumentType:      DocumentType,
		PoleName:          PoleName,
		PoleType:          PoleType,
		Default:           Default,
		IndexType:         IndexType,
		Checker:           &PoleCheckerStringValue{},
		New:               New,
		Edit:              Edit,
		Title:             Title,
	}
}
func (server *Server) addBaseConfiguration(ConfigurationName string, conf *configuration) bool {
	server.lockBaseConfiguration.Lock()
	defer server.lockBaseConfiguration.Unlock()
	_, ok := server.baseConfiguration[ConfigurationName]
	if ok {
		return false
	}
	server.baseConfiguration[ConfigurationName] = conf
	return true
}
func (server *Server) addFullConfiguration(ConfigurationName string, conf *configuration) bool {
	server.lockFullConfiguration.Lock()
	defer server.lockFullConfiguration.Unlock()
	_, ok := server.fullConfiguration[ConfigurationName]
	if ok {
		return false
	}
	server.fullConfiguration[ConfigurationName] = conf
	return true
}
func (server *Server) newConfiguration() *configuration {

	return &configuration{
		polesInfo:               make(map[string]map[string]IPoleInfo),
		typesInfo:               make(map[string]ITypeInfo),
		callsInfo:               make(map[string]ICallInfo),
		dependConfigurationName: []string{},
	}
}
func (server *Server) loadDefaultsConfiguration() {
	loadConfConfiguration(server)
	loadConfEditConfiguration(server)
	loadDocumentConfiguration(server)
	//newConfiguration(ConfigurationName)
}
func (server *Server) getUser(UserName string) IUser {
	user, ok := server.users[UserName]
	if ok {
		return user
	}
	return nil
}
func (server *Server) loadBaseConfiguration(ConfigurationName string) *configuration {
	server.lockBaseConfiguration.Lock()
	defer server.lockBaseConfiguration.Unlock()
	if conf, ok := server.baseConfiguration[ConfigurationName]; ok {
		return conf
	}
	return nil
}
func (server *Server) getConfiguration(ConfigurationName string) *configuration {
	server.lockFullConfiguration.Lock()
	defer server.lockFullConfiguration.Unlock()
	conf, ok := server.fullConfiguration[ConfigurationName]
	if ok {
		return conf
	}
	return nil
}
func (server *Server) LoadConfiguration(ConfigurationName string) (*configuration, error) {
	var conf *configuration
	if conf = server.getConfiguration(ConfigurationName); conf != nil {
		return conf, nil
	}
	conf = server.newConfiguration()
	loadConfiguration := make(map[string]*configuration)
	depend := &list.List{}
	var addDepend func(ConfigurationName string)
	addDepend = func(ConfigurationName string) {
		var lconf *configuration
		var ok bool
		if lconf, ok = loadConfiguration[ConfigurationName]; !ok {
			lconf = server.loadBaseConfiguration(ConfigurationName)
			loadConfiguration[ConfigurationName] = lconf
			for i := 0; i < len(lconf.dependConfigurationName); i++ {
				addDepend(lconf.dependConfigurationName[i])
			}
			fmt.Println(lconf)
			depend.PushBack(lconf)
		}
	}
	addDepend(ConfigurationName)
	if depend.Front() == nil {
		return nil, errors.New("configuration not found: " + ConfigurationName)
	}
	for e := depend.Front(); e != nil; e = e.Next() {
		c := e.Value.(*configuration)
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
	server := &Server{
		databaseName:          DatabaseName,
		genUID:                make(chan string),
		lockTrancactions:      &sync.Mutex{},
		trancactions:          make(map[string]*transaction),
		lockBaseConfiguration: &sync.Mutex{},
		baseConfiguration:     make(map[string]*configuration),
		lockFullConfiguration: &sync.Mutex{},
		fullConfiguration:     make(map[string]*configuration),
		lockUsers:             &sync.Mutex{},
		users:                 make(map[string]IUser),
		close:                 false,
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
func (server *Server) NewClient(Login string, Password string, ConfigurationName string) *Client {
	user := server.getUser(Login)
	if !user.GetCheckPassword(Password) {
		return nil
	}
	//configuration := server.LoadConfiguration(ConfigurationName)
	return &Client{user: user, configurationName: ConfigurationName, server: server}
}

/**/
