package core

import (
	"errors"
	"strings"
)

var baseConfiguration = map[string]*Configuration{}

func InitAddBaseConfiguration(Name string, conf *Configuration) {
	baseConfiguration[Name] = conf
}

type Configuration struct {
	polesInfo               map[string]map[string]IPoleInfo
	typesInfo               map[string]ITypeInfo
	callsInfo               map[string]ICallInfo
	dependConfigurationName []string
}

func (conf *Configuration) GetCallInfo(CallName string) (ICallInfo, error) {
	if ci, ok := conf.callsInfo[CallName]; ok {
		return ci, nil
	}
	return nil, errors.New("call not found: " + CallName)
}
func (conf *Configuration) GetPoleInfo(DocumentType string, PoleName string) (IPoleInfo, error) {
	if m, ok := conf.polesInfo[DocumentType]; ok {
		if pi, ok := m[PoleName]; ok {
			return pi, nil
		}
	}
	return nil, errors.New("pole not found: " + PoleName)
}
func (conf *Configuration) GetPolesInfo(DocumentType string, Poles []string) map[string]IPoleInfo {
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
func (conf *Configuration) AddDependConfiguration(DependConfigurationName string) {
	conf.dependConfigurationName = append(conf.dependConfigurationName, DependConfigurationName)
}
func (conf *Configuration) AddCall(ConfigurationName string, Name string, PullName string, Call bool, Listen bool, Title string) {
	conf.callsInfo[Name] = &CallInfo{ConfigurationName: ConfigurationName, Name: Name, PullName: PullName, Call: Call, Listen: Listen, Title: Title}
}
func (conf *Configuration) AddType(ConfigurationName string, DocumentType string, New bool, Read bool, Save bool, Title string) {
	conf.typesInfo[DocumentType] = &TypeInfo{ConfigurationName: ConfigurationName, DocumentType: DocumentType, New: New, Read: Read, Save: Save, Title: Title}
}
func (conf *Configuration) AddPole(ConfigurationName string, DocumentType string, PoleName string, PoleType string, Default Object, IndexType string, Checker IPoleChecker, New bool, Edit bool, Title string) {
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
func NewConfiguration() *Configuration {
	return &Configuration{
		polesInfo:               make(map[string]map[string]IPoleInfo),
		typesInfo:               make(map[string]ITypeInfo),
		callsInfo:               make(map[string]ICallInfo),
		dependConfigurationName: []string{},
	}
}
func (core *Core) loadDefaultsConfiguration() {
	for Name, Conf := range baseConfiguration {
		core.addBaseConfiguration(Name, Conf)
	}
}
