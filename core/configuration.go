package core

import (
	"errors"
	"strings"

	"github.com/shammishailaj/sud/corebase"
)

var baseConfiguration = map[string]*Configuration{}

func InitAddBaseConfiguration(Name string, conf *Configuration) {
	baseConfiguration[Name] = conf
}

type Configuration struct {
	polesInfo               map[string]map[string]corebase.IPoleInfo
	typesInfo               map[string]corebase.ITypeInfo
	callsInfo               map[string]corebase.ICallInfo
	dependConfigurationName []string
}

func (conf *Configuration) GetCallInfo(CallName string) (corebase.ICallInfo, error) {
	if ci, ok := conf.callsInfo[CallName]; ok {
		return ci, nil
	}
	return nil, errors.New("call not found: " + CallName)
}
func (conf *Configuration) GetTypeInfo(RecordType string) (corebase.ITypeInfo, error) {
	if m, ok := conf.typesInfo[RecordType]; ok {
		return m, nil
	}
	return nil, errors.New("type not found: " + RecordType)
}
func (conf *Configuration) GetPoleInfo(RecordType string, PoleName string) (corebase.IPoleInfo, error) {
	if m, ok := conf.polesInfo[RecordType]; ok {
		if pi, ok := m[PoleName]; ok {
			return pi, nil
		}
	}
	return nil, errors.New("pole not found: " + PoleName)
}
func (conf *Configuration) GetPolesInfo(RecordType string, Poles []string) map[string]corebase.IPoleInfo {
	if m, ok := conf.polesInfo[RecordType]; ok {
		polesInfo := map[string]corebase.IPoleInfo{}
		if len(Poles) != 0 {
			for _, pole := range Poles {
				p := strings.Split(pole, ".")
				// Record.Table.Name.
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
	return map[string]corebase.IPoleInfo{}
}
func (conf *Configuration) AddDependConfiguration(DependConfigurationName string) {
	conf.dependConfigurationName = append(conf.dependConfigurationName, DependConfigurationName)
}
func (conf *Configuration) AddCall(ConfigurationName string, Name string, PullName string, Call bool, Listen bool, Title string) {
	conf.callsInfo[Name] = &CallInfo{ConfigurationName: ConfigurationName, Name: Name, PullName: PullName, Call: Call, Listen: Listen, Title: Title}
}
func (conf *Configuration) AddType(ConfigurationName string, RecordType string, New bool, Read bool, Save bool, Title string) {
	conf.typesInfo[RecordType] = &TypeInfo{ConfigurationName: ConfigurationName, RecordType: RecordType, New: New, Read: Read, Save: Save, Title: Title}
}
func (conf *Configuration) AddPole(ConfigurationName string, RecordType string, PoleName string, PoleType string, Default corebase.Object, IndexType string, Checker corebase.IPoleChecker, New bool, Edit bool, Title string) {
	_, ok := conf.polesInfo[RecordType]
	if !ok {
		conf.polesInfo[RecordType] = make(map[string]corebase.IPoleInfo)
	}
	conf.polesInfo[RecordType][PoleName] = &PoleInfo{
		ConfigurationName: ConfigurationName,
		RecordType:      RecordType,
		PoleName:          PoleName,
		PoleType:          PoleType,
		Default:           Default,
		IndexType:         IndexType,
		Checker:           Checker,
		New:               New,
		Edit:              Edit,
		Title:             Title,
	}
}
func NewConfiguration() *Configuration {
	return &Configuration{
		polesInfo:               make(map[string]map[string]corebase.IPoleInfo),
		typesInfo:               make(map[string]corebase.ITypeInfo),
		callsInfo:               make(map[string]corebase.ICallInfo),
		dependConfigurationName: []string{},
	}
}
func (core *Core) loadDefaultsConfiguration() {
	for Name, Conf := range baseConfiguration {
		core.addBaseConfiguration(Name, Conf)
	}
}
