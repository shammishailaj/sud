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
	name                    string
	polesInfo               map[string]map[string]corebase.IPoleInfo
	typesInfo               map[string]corebase.ITypeInfo
	callsInfo               map[string]corebase.ICallInfo
	accessInfo              map[string]corebase.IAccessInfo
	dependConfigurationName []string
	accessLoad              []string
}

func (conf *Configuration) Name() string {
	return conf.name
}
func (conf *Configuration) AddAccessLoad(Access string) {
	for _, a := range conf.accessLoad {
		if a == Access {
			return
		}
	}
	conf.accessLoad = append(conf.accessLoad, Access)
}
func (conf *Configuration) DependConfigurationName() []string {
	return conf.dependConfigurationName
}
func (conf *Configuration) AccessLoad() []string {
	return conf.accessLoad
}

func (conf *Configuration) CheckAccessLoad(Access corebase.IAccess) bool {
	for _, a := range conf.accessLoad {
		if !Access.CheckAccess(a) {
			return false
		}
	}
	return true
}
func (conf *Configuration) CallNames() []string {
	list := make([]string, len(conf.callsInfo))
	i := 0
	for name := range conf.callsInfo {
		list[i] = name
		i++
	}
	return list
}
func (conf *Configuration) CallInfo(CallName string) (corebase.ICallInfo, error) {
	if ci, ok := conf.callsInfo[CallName]; ok {
		return ci, nil
	}
	return nil, &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "Get CallInfo", Name: CallName}
}
func (conf *Configuration) TypeNames() []string {
	list := make([]string, len(conf.typesInfo))
	i := 0
	for name := range conf.typesInfo {
		list[i] = name
		i++
	}
	return list
}
func (conf *Configuration) TypeInfo(RecordType string) (corebase.ITypeInfo, error) {
	if m, ok := conf.typesInfo[RecordType]; ok {
		return m, nil
	}
	return nil, &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "Get TypeInfo", Name: RecordType}
}
func (conf *Configuration) PoleInfo(RecordType string, PoleName string) (corebase.IPoleInfo, error) {
	if m, ok := conf.polesInfo[RecordType]; ok {
		if pi, ok := m[PoleName]; ok {
			return pi, nil
		}
	}
	return nil, &corebase.Error{ErrorType: corebase.ErrorTypeNotFound, Action: "Get PoleInfo", Name: PoleName}
}
func (conf *Configuration) PolesRecordTypes() []string {
	list := make([]string, len(conf.polesInfo))
	i := 0
	for name := range conf.polesInfo {
		list[i] = name
		i++
	}
	return list
}
func (conf *Configuration) PolesInfo(RecordType string, Poles []string) map[string]corebase.IPoleInfo {
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
func (conf *Configuration) AddAccess(info AccessInfo) {
	conf.accessInfo[info.Name] = &info
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
func NewConfiguration(Name string, AccessLoad []string) *Configuration {
	return &Configuration{
		name:                    Name,
		accessInfo:              make(map[string]corebase.IAccessInfo),
		polesInfo:               make(map[string]map[string]corebase.IPoleInfo),
		typesInfo:               make(map[string]corebase.ITypeInfo),
		callsInfo:               make(map[string]corebase.ICallInfo),
		dependConfigurationName: []string{},
		accessLoad:              AccessLoad,
	}
}
func (conf *Configuration) Save() (map[string]interface{}, error) {
	confInfo := make(map[string]interface{})
	structures.MapSetStringList(confInfo, "dependConfigurationName", conf.DependConfigurationName())
	structures.MapSetStringList(confInfo, "accessLoad", conf.AccessLoad())
	callsInfo := make(map[string]interface{})
	for _, callName := range conf.CallNames() {
		if info, err := conf.CallInfo(callName); err == nil {
			callInfo := make(map[string]interface{})
			callInfo["pullName"] = info.GetPullName()
			callInfo["accessCall"] = info.GetAccessCall()
			callInfo["accessListen"] = info.GetAccessListen()
			callInfo["title"] = info.GetTitle()
			callsInfo[callName] = callInfo
		}
	}
	confInfo["callInfo"] = callsInfo
	typesInfo := make(map[string]interface{})
	for _, typeName := range conf.TypeNames() {
		if info, err := conf.TypeInfo(typeName); err == nil {
			typeInfo := make(map[string]interface{})
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
	for _, recordType := range conf.PolesRecordTypes() {
		polesRT := make(map[string]interface{})
		for _, info := range conf.PolesInfo(recordType, []string{}) {
			poleInfo := make(map[string]interface{})
			poleInfo["poleType"] = info.GetPoleType()
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
	if list, ok := structures.MapGetStringList(confInfo, "accessLoad"); ok {
		for _, a := range list {
			conf.AddAccessLoad(a)
		}
	}
	if item, ok := confInfo["accessInfo"]; ok {
		if accessInfo, ok := item.(map[string]interface{}); ok {
			for accessName, info := range accessInfo {
				if infoExt, ok := info.(map[string]interface{}); ok {
					var ai AccessInfo
					var nameOK, titleOK bool
					ai.Name, nameOK = structures.MapGetString(infoExt, "name")
					ai.Title, titleOK = structures.MapGetString(infoExt, "title")
					if nameOK && accessName != ai.Name {
						return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "accessName", Name: "conf stucture error"}
					}
					if !nameOK {
						ai.Name = accessName
					}
					if titleOK {
						conf.AddAccess(ai)
					} else {
						return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "title", Name: "conf stucture error"}
					}
				}
			}
		}
	}
	if item, ok := confInfo["callInfo"]; ok {
		if callInfo, ok := item.(map[string]interface{}); ok {
			for callName, info := range callInfo {
				if infoExt, ok := info.(map[string]interface{}); ok {
					var ci CallInfo
					var pullNameOK, accessCallOK, accessListenOK, titleOK bool
					ci.ConfigurationName = conf.Name()
					ci.Name = callName
					ci.PullName, pullNameOK = structures.MapGetString(infoExt, "pullName")
					ci.AccessCall, accessCallOK = structures.MapGetString(infoExt, "accessCall")
					ci.AccessListen, accessListenOK = structures.MapGetString(infoExt, "accessListen")
					ci.Title, titleOK = structures.MapGetString(infoExt, "title")
					if pullNameOK && accessCallOK && accessListenOK && titleOK {
						conf.AddCall(ci)
					} else {
						return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "configurationName, pullName, accessCall, accessListen, title", Name: "conf stucture error"}
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
					var accessTypeOK, accessNewOK, accessReadOK, accessSaveOK, titleOK bool
					ti.ConfigurationName = conf.Name()
					ti.RecordType = typeName
					ti.AccessType, accessTypeOK = structures.MapGetString(infoExt, "accessType")
					ti.AccessNew, accessNewOK = structures.MapGetString(infoExt, "accessNew")
					ti.AccessRead, accessReadOK = structures.MapGetString(infoExt, "accessRead")
					ti.AccessSave, accessSaveOK = structures.MapGetString(infoExt, "accessSave")
					ti.Title, titleOK = structures.MapGetString(infoExt, "title")
					if accessTypeOK && accessNewOK && accessReadOK && accessSaveOK && titleOK {
						conf.AddType(ti)
					} else {
						return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "accessType, accessNew, accessRead, accessSave, title", Name: "conf stucture error"}
					}
				} else {
					return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "callInfo", Name: "conf sutucture error"}
				}
			}
		}
	}
	if item, ok := confInfo["polesInfo"]; ok {
		if recordTypes, ok := item.(map[string]interface{}); ok {
			for recordTypeName, recordTypesInfo := range recordTypes {
				if polesInfo, ok := recordTypesInfo.(map[string]interface{}); ok {
					for poleName, info := range polesInfo {
						if infoExt, ok := info.(map[string]interface{}); ok {
							var pi PoleInfo
							var poleTypeOK, indexTypeOK, defaultOK, checkerOK, accessReadOK, accessWriteOK, titleOK bool
							pi.ConfigurationName = conf.Name()
							pi.PoleName = poleName
							pi.PoleType, poleTypeOK = structures.MapGetString(infoExt, "poleType")
							pi.RecordType = recordTypeName
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
							if poleTypeOK && indexTypeOK && defaultOK && checkerOK && accessReadOK && accessWriteOK && titleOK {
								conf.AddPole(pi)
							} else {
								return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "poleType, indexType, default, checker, accessRead, accessWrite, title", Name: "conf stucture error"}
							}
						} else {
							return &corebase.Error{Action: "Configuration:Load", ErrorType: corebase.ErrorFormat, Info: "callInfo", Name: "conf sutucture error"}
						}
					}
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
