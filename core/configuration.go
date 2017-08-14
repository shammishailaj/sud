package core

import (
	"strings"

	"github.com/crazyprograms/sud/corebase"
	"github.com/crazyprograms/sud/jparam"
	"github.com/crazyprograms/sud/structures"
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
func (conf *Configuration) GetDependConfigurationName() []string {
	return conf.dependConfigurationName
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
func (conf *Configuration) GetCallNames() []string {
	list := make([]string, len(conf.callsInfo))
	i := 0
	for name := range conf.callsInfo {
		list[i] = name
		i++
	}
	return list
}
func (conf *Configuration) GetCallInfo(CallName string) (corebase.ICallInfo, error) {
	if ci, ok := conf.callsInfo[CallName]; ok {
		return ci, nil
	}
	return nil, &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "GetCallInfo", Name: CallName}
}
func (conf *Configuration) GetTypeNames() []string {
	list := make([]string, len(conf.typesInfo))
	i := 0
	for name := range conf.typesInfo {
		list[i] = name
		i++
	}
	return list
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
func (conf *Configuration) GetPolesRecordTypes() []string {
	list := make([]string, len(conf.polesInfo))
	i := 0
	for name := range conf.polesInfo {
		list[i] = name
		i++
	}
	return list
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
func (conf *Configuration) Save() (map[string]interface{}, error) {
	confInfo := make(map[string]interface{})
	structures.MapSetStringList(confInfo, "dependConfigurationName", conf.GetDependConfigurationName())
	structures.MapSetStringList(confInfo, "accessList", conf.GetAccessList())
	callsInfo := make(map[string]interface{})
	for _, callName := range conf.GetCallNames() {
		if info, err := conf.GetCallInfo(callName); err == nil {
			callInfo := make(map[string]interface{})
			callInfo["configurationName"] = info.GetConfigurationName()
			callInfo["name"] = info.GetName()
			callInfo["pullName"] = info.GetPullName()
			callInfo["accessCall"] = info.GetAccessCall()
			callInfo["accessListen"] = info.GetAccessListen()
			callInfo["title"] = info.GetTitle()
			callsInfo[callName] = callInfo
		}
	}
	confInfo["callInfo"] = callsInfo
	typesInfo := make(map[string]interface{})
	for _, typeName := range conf.GetTypeNames() {
		if info, err := conf.GetTypeInfo(typeName); err == nil {
			typeInfo := make(map[string]interface{})
			typeInfo["configurationName"] = info.GetConfigurationName()
			typeInfo["recordType"] = info.GetRecordType()
			typeInfo["accessType"] = info.GetAccessType()
			typeInfo["accessNew"] = info.GetAccessNew()
			typeInfo["accessRead"] = info.GetAccessRead()
			typeInfo["accessSave"] = info.GetAccessSave()
			typeInfo["title"] = info.GetTitle()
			typesInfo[typeName] = typeInfo
		}
	}
	confInfo["typeInfo"] = typesInfo
	polesInfo := make(map[string]interface{})
	for _, recordType := range conf.GetPolesRecordTypes() {
		polesRT := make(map[string]interface{})
		for _, info := range conf.GetPolesInfo(recordType, []string{}) {
			poleInfo := make(map[string]interface{})
			poleInfo["configurationName"] = info.GetConfigurationName()
			poleInfo["poleName"] = info.GetPoleName()
			poleInfo["poleType"] = info.GetPoleType()
			poleInfo["recordType"] = info.GetRecordType()
			poleInfo["indexType"] = info.GetIndexType()
			poleInfo["default"] = info.GetDefault()
			poleInfo["checker"] = SaveChecker(info.GetChecker())
			poleInfo["accessRead"] = info.GetAccessRead()
			poleInfo["accessWrite"] = info.GetAccessWrite()
			poleInfo["title"] = info.GetTitle()
			polesRT[info.GetPoleName()] = poleInfo
		}
		polesInfo[recordType] = polesRT
	}
	confInfo["polesInfo"] = polesInfo
	return confInfo, nil
}
func (conf *Configuration) SaveJson() ([]byte, error) {
	var err error
	var data interface{}
	if data, err = conf.Save(); err != nil {
		return nil, err
	}
	return jparam.ToJson(data)
}
func (conf *Configuration) Load(confInfo map[string]interface{}) error {
	if list, ok := structures.MapGetStringList(confInfo, "dependConfigurationName"); ok {
		for _, dconf := range list {
			conf.AddDependConfiguration(dconf)
		}
	}
	if list, ok := structures.MapGetStringList(confInfo, "accessList"); ok {
		for _, a := range list {
			conf.AddAccess(a)
		}
	}
	if item, ok := confInfo["callInfo"]; ok {
		if callInfo, ok := item.(map[string]interface{}); ok {
			for callName, info := range callInfo {
				if infoExt, ok := info.(map[string]interface{}); ok {
					var ci CallInfo
					var configurationNameOK, nameOK, pullNameOK, accessCallOK, accessListenOK, titleOK bool
					ci.ConfigurationName, configurationNameOK = structures.MapGetString(infoExt, "configurationName")
					ci.Name, nameOK = structures.MapGetString(infoExt, "name")
					ci.PullName, pullNameOK = structures.MapGetString(infoExt, "pullName")
					ci.AccessCall, accessCallOK = structures.MapGetString(infoExt, "accessCall")
					ci.AccessListen, accessListenOK = structures.MapGetString(infoExt, "accessListen")
					ci.Title, titleOK = structures.MapGetString(infoExt, "title")
					if callName != ci.Name {
						return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "callName", Name: "conf stucture error"}
					}
					if configurationNameOK && nameOK && pullNameOK && accessCallOK && accessListenOK && titleOK {
						conf.AddCall(ci)
					} else {
						return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "configurationName, name, pullName, accessCall, accessListen, title", Name: "conf stucture error"}
					}
				} else {
					return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "callInfo", Name: "conf sutucture error"}
				}
			}
		}
	}
	if item, ok := confInfo["typeInfo"]; ok {
		if typeInfo, ok := item.(map[string]interface{}); ok {
			for typeName, info := range typeInfo {
				if infoExt, ok := info.(map[string]interface{}); ok {
					var ti TypeInfo
					var configurationNameOK, recordTypeOK, accessTypeOK, accessNewOK, accessReadOK, accessSaveOK, titleOK bool
					ti.ConfigurationName, configurationNameOK = structures.MapGetString(infoExt, "configurationName")
					ti.RecordType, recordTypeOK = structures.MapGetString(infoExt, "recordType")
					ti.AccessType, accessTypeOK = structures.MapGetString(infoExt, "accessType")
					ti.AccessNew, accessNewOK = structures.MapGetString(infoExt, "accessNew")
					ti.AccessRead, accessReadOK = structures.MapGetString(infoExt, "accessRead")
					ti.AccessSave, accessSaveOK = structures.MapGetString(infoExt, "accessSave")
					ti.Title, titleOK = structures.MapGetString(infoExt, "title")
					if typeName != ti.RecordType {
						return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "typeName", Name: "conf stucture error"}
					}
					if configurationNameOK && recordTypeOK && accessTypeOK && accessNewOK && accessReadOK && accessSaveOK && titleOK {
						conf.AddType(ti)
					} else {
						return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "configurationName, recordType, accessType, accessNew, accessRead, accessSave, title", Name: "conf stucture error"}
					}
				} else {
					return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "callInfo", Name: "conf sutucture error"}
				}
			}
		}
	}
	if item, ok := confInfo["polesInfo"]; ok {
		if polesInfo, ok := item.(map[string]interface{}); ok {
			for poleName, info := range polesInfo {
				if infoExt, ok := info.(map[string]interface{}); ok {
					var pi PoleInfo
					var configurationNameOK, poleNameOK, poleTypeOK, recordTypeOK, indexTypeOK, defaultOK, checkerOK, accessReadOK, accessWriteOK, titleOK bool
					pi.ConfigurationName, configurationNameOK = structures.MapGetString(infoExt, "configurationName")
					pi.PoleName, poleNameOK = structures.MapGetString(infoExt, "poleName")
					pi.PoleType, poleTypeOK = structures.MapGetString(infoExt, "poleType")
					pi.RecordType, recordTypeOK = structures.MapGetString(infoExt, "recordType")
					pi.IndexType, indexTypeOK = structures.MapGetString(infoExt, "indexType")
					pi.Default, defaultOK = structures.MapGetValue(infoExt, "default")
					var Checker interface{}
					var err error
					Checker, checkerOK = structures.MapGetValue(infoExt, "checker")
					if CheckerM, ok := Checker.(map[string]interface{}); ok {
						if pi.Checker, err = LoadChecker(CheckerM); err != nil {
							return nil
						}
					}
					pi.AccessRead, accessReadOK = structures.MapGetString(infoExt, "accessRead")
					pi.AccessWrite, accessWriteOK = structures.MapGetString(infoExt, "accessWrite")
					pi.Title, titleOK = structures.MapGetString(infoExt, "title")
					if poleName != pi.PoleName {
						return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "poleName", Name: "conf stucture error"}
					}
					if configurationNameOK && poleNameOK && poleTypeOK && recordTypeOK && indexTypeOK && defaultOK && checkerOK && accessReadOK && accessWriteOK && titleOK {
						conf.AddPole(pi)
					} else {
						return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "configurationName, poleName, poleType, recordType, indexType, default, checker, accessRead, accessWrite, title", Name: "conf stucture error"}
					}
				} else {
					return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "callInfo", Name: "conf sutucture error"}
				}
			}
		}
	}
	return nil
}
func (conf *Configuration) LoadJson(data []byte) error {
	var err error
	var params interface{}
	if params, err = jparam.FromJson(data); err != nil {
		return err
	}
	var paramsM map[string]interface{}
	switch params.(type) {
	case map[string]interface{}:
		paramsM = params.(map[string]interface{})
	default:
		return &corebase.Error{Action: "Configuration:LoadJson", ErrorType: corebase.ErrorFormat, Info: "{map:{...}}", Name: "json sutucture error"}
	}
	return conf.Load(paramsM)
}
