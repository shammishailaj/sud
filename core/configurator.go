package core

import (
	"container/list"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/crazyprograms/sud/corebase"
)

type Configurator struct {
	configurationPath     string
	lockBaseConfiguration sync.RWMutex
	baseConfiguration     map[string]*Configuration
	lockFullConfiguration sync.RWMutex
	fullConfiguration     map[string]*Configuration
}

func (cr *Configurator) loadConfigurations(configDir string) error {
	configDir = path.Clean(configDir)

	var err error
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(configDir); err != nil {
		return err
	}
	for _, file := range files {
		n := file.Name()
		nf := path.Join(configDir, n)
		if file.IsDir() {
			cr.loadConfigurations(nf)
		} else {
			var data []byte
			if data, err = ioutil.ReadFile(nf); err != nil {
				return err
			}
			confName := n[0 : len(n)-len(path.Ext(n))]
			conf := NewConfiguration([]string{})
			if err = conf.LoadJson(data); err != nil {
				return err
			}
			cr.AddBaseConfiguration(confName, conf)
		}
	}
	return nil
}

func NewConfigurator() *Configurator {
	return &Configurator{
		//lockBaseConfiguration: &sync.RWMutex{},
		baseConfiguration: make(map[string]*Configuration),
		//lockFullConfiguration: &sync.RWMutex{},
		fullConfiguration: make(map[string]*Configuration),
	}
}
func (cr *Configurator) SetConfigurationPath(configurationPath string) {
	cr.configurationPath = configurationPath
}
func (cr *Configurator) ReloadConfiguration() error {
	cr.lockFullConfiguration.Lock()
	defer cr.lockFullConfiguration.Unlock()
	cr.lockBaseConfiguration.Lock()
	cr.baseConfiguration = make(map[string]*Configuration)
	cr.fullConfiguration = make(map[string]*Configuration)
	cr.lockBaseConfiguration.Unlock()
	return cr.loadConfigurations(cr.configurationPath)
}
func (cr *Configurator) GetConfigurations() map[string]*Configuration {
	cr.lockBaseConfiguration.RLock()
	l := make(map[string]*Configuration, len(cr.baseConfiguration))
	for key, value := range cr.baseConfiguration {
		l[key] = value
	}
	cr.lockBaseConfiguration.RUnlock()
	return l
}
func (cr *Configurator) AddBaseConfiguration(ConfigurationName string, conf *Configuration) bool {
	cr.lockBaseConfiguration.Lock()
	_, ok := cr.baseConfiguration[ConfigurationName]
	if !ok {
		cr.baseConfiguration[ConfigurationName] = conf
	}
	return !ok
	cr.lockBaseConfiguration.Unlock()
	return true
}
func (cr *Configurator) GetConfiguration(ConfigurationName string, Access corebase.IAccess) (*Configuration, error) {
	var conf *Configuration
	var ok bool
	var err error
	cr.lockFullConfiguration.RLock()
	conf, ok = cr.fullConfiguration[ConfigurationName]
	cr.lockFullConfiguration.RUnlock()
	if !ok {
		if conf, err = cr.loadConfiguration(ConfigurationName); err != nil {
			return nil, err
		}
	}
	if !conf.CheckAccess(Access) {
		return nil, &corebase.Error{ErrorType: corebase.ErrorTypeAccessIsDenied, Action: "LoadConfiguration", Name: ConfigurationName}
	}
	return conf, nil
}
func (cr *Configurator) loadConfiguration(ConfigurationName string) (*Configuration, error) {
	var conf *Configuration
	cr.lockFullConfiguration.Lock()
	defer cr.lockFullConfiguration.Unlock()
	AccessConfiguration := map[string]bool{}
	loadConfiguration := map[string]*Configuration{}
	depend := &list.List{}
	var addDepend func(ConfigurationName string) error
	addDepend = func(ConfigurationName string) error {
		var lconf *Configuration
		var ok bool
		if lconf, ok = loadConfiguration[ConfigurationName]; !ok {
			if lconf, ok = cr.baseConfiguration[ConfigurationName]; !ok {
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
	AccessList := make([]string, len(AccessConfiguration))
	ANum := 0
	for a := range AccessConfiguration {
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
	cr.fullConfiguration[ConfigurationName] = conf
	return conf, nil
}
