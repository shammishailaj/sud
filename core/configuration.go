package core

import (
	"strings"

	"github.com/crazyprograms/sud/corebase"
)

/*var baseConfiguration = map[string]*Configuration{}

func InitAddBaseConfiguration(Name string, conf *Configuration) {
	baseConfiguration[Name] = conf
}*/

type Configuration struct {
	polesInfo               map[string]map[string]corebase.IPoleInfo
	typesInfo               map[string]corebase.ITypeInfo
	callsInfo               map[string]corebase.ICallInfo
	dependConfigurationName []string
	accessList              []string
}

func (conf *Configuration) AddAccess(Access string) {
	for _, a := range conf.accessList {
		if a == Access {
			return
		}
	}
	conf.accessList = append(conf.accessList, Access)
}
func (conf *Configuration) GetAccessList() []string {
	return conf.accessList
}

func (conf *Configuration) CheckAccess(Access corebase.IAccess) bool {
	for _, a := range conf.accessList {
		if !Access.CheckAccess(a) {
			return false
		}
	}
	return true
}

func (conf *Configuration) GetCallInfo(CallName string) (corebase.ICallInfo, error) {
	if ci, ok := conf.callsInfo[CallName]; ok {
		return ci, nil
	}
	return nil, &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "GetCallInfo", Name: CallName}
}
func (conf *Configuration) GetTypeInfo(RecordType string) (corebase.ITypeInfo, error) {
	if m, ok := conf.typesInfo[RecordType]; ok {
		return m, nil
	}
	return nil, &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "GetTypeInfo", Name: RecordType}
}
func (conf *Configuration) GetPoleInfo(RecordType string, PoleName string) (corebase.IPoleInfo, error) {
	if m, ok := conf.polesInfo[RecordType]; ok {
		if pi, ok := m[PoleName]; ok {
			return pi, nil
		}
	}
	return nil, &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "GetPoleInfo", Name: PoleName}
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
func (conf *Configuration) AddCall(info CallInfo) {
	conf.callsInfo[info.Name] = &info
}

/*func (conf *Configuration) AddCall(ConfigurationName string, Name string, PullName string, AccessCall bool, AccessListen bool, Title string) {
	conf.callsInfo[Name] = &CallInfo{ConfigurationName: ConfigurationName, Name: Name, PullName: PullName, AccessCall: AccessCall, AccessListen: AccessListen, Title: Title}
}*/
func (conf *Configuration) AddType(info TypeInfo) {
	conf.typesInfo[info.RecordType] = &info
}

/*func (conf *Configuration) AddType(ConfigurationName string, RecordType string, AccessType string, New string, Read string, Save string, Title string) {
	conf.typesInfo[RecordType] = &TypeInfo{ConfigurationName: ConfigurationName, RecordType: RecordType, AccessType: AccessType, New: New, Read: Read, Save: Save, Title: Title}
}*/
func (conf *Configuration) AddPole(info PoleInfo) {
	_, ok := conf.polesInfo[info.RecordType]
	if !ok {
		conf.polesInfo[info.RecordType] = make(map[string]corebase.IPoleInfo)
	}
	conf.polesInfo[info.RecordType][info.PoleName] = &info
}

/*func (conf *Configuration) AddPole(ConfigurationName string, RecordType string, PoleName string, PoleType string, Default corebase.Object, IndexType string, Checker corebase.IPoleChecker, Read string, Write string, Title string) {
	_, ok := conf.polesInfo[RecordType]
	if !ok {
		conf.polesInfo[RecordType] = make(map[string]corebase.IPoleInfo)
	}
	conf.polesInfo[RecordType][PoleName] = &PoleInfo{
		ConfigurationName: ConfigurationName,
		RecordType:        RecordType,
		PoleName:          PoleName,
		PoleType:          PoleType,
		Default:           Default,
		IndexType:         IndexType,
		Checker:           Checker,
		Read:              Read,
		Write:             Write,
		Title:             Title,
	}
}*/
func NewConfiguration(AccessList []string) *Configuration {
	return &Configuration{
		polesInfo:               make(map[string]map[string]corebase.IPoleInfo),
		typesInfo:               make(map[string]corebase.ITypeInfo),
		callsInfo:               make(map[string]corebase.ICallInfo),
		dependConfigurationName: []string{},
		accessList:              AccessList,
	}
}
