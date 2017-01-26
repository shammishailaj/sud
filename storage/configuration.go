package storage

// Устарел см. client.go
import (
	"errors"
	"log"
)

var (
	DefaultConfiguration = map[string]func(configuration *Configuration, State *loadState) error{
	//BaseConfigurationName: func(configuration *Configuration, State *loadState) error { return baseConf.load(configuration, State) },
	}
)

type Configuration struct {
	polesInfo             map[string]map[string]IPoleInfo
	typesInfo             map[string]ITypeInfo
	loadConfigurationName map[string]bool
}
type loadState struct {
	PolesInfo             map[string]map[string]IPoleInfo
	TypesInfo             map[string]ITypeInfo
	LoadConfigurationName map[string]bool
	QueueConfiguration    []string
	ListDepend            []string
	/*//public IConnection connection;
	  public Queue<String> QueueConfiguration = new Queue<String>();
	  public List<String> ListDepend = new List<String>();
	  public Dictionary<String, IEnumerable<IDocument>> ConfigList = new Dictionary<String, IEnumerable<IDocument>>();
	  public Dictionary<String, Dictionary<String, PoleInfo>> PolesInfo = new Dictionary<String, Dictionary<String, Configuration.PoleInfo>>();
	  public Dictionary<String, DocumentTypeInfo> DocumentTypesInfo = new Dictionary<String, DocumentTypeInfo>();
	  public HashSet<String> LoadConfigurationName = new HashSet<String>();
	*/
}

func newLoadState() *loadState {
	return &loadState{
		PolesInfo:             map[string]map[string]IPoleInfo{},
		TypesInfo:             map[string]ITypeInfo{},
		LoadConfigurationName: map[string]bool{},
		QueueConfiguration:    []string{},
		ListDepend:            []string{},
	}
}

func (ls *loadState) setLoadState(conf *Configuration) {
	conf.polesInfo = ls.PolesInfo
	conf.typesInfo = ls.TypesInfo
	conf.loadConfigurationName = ls.LoadConfigurationName
}

type ConfigurationInfo struct {
	types []TypeInfo
	poles []PoleInfo
}

func (ci *ConfigurationInfo) load(configuration *Configuration, State *loadState) error {
	for i := 0; i < len(ci.types); i++ {
		configuration.addDocumentTypeInfo(&ci.types[i], State)
	}
	for i := 0; i < len(ci.poles); i++ {
		configuration.addPoleInfo(&ci.poles[i], State)
	}
	return nil
}
func (c *Configuration) addPoleInfo(info IPoleInfo, State *loadState) {
	var dt map[string]IPoleInfo
	var info2 IPoleInfo
	var ok bool
	if dt, ok = State.PolesInfo[info.GetDocumentType()]; !ok {
		dt = map[string]IPoleInfo{}
		State.PolesInfo[info.GetDocumentType()] = dt
	}
	if info2, ok = dt[info.GetPoleName()]; !ok {
		dt[info.GetPoleName()] = info
	} else {
		log.Println("[" + info.GetConfigurationName() + "]" + info.GetDocumentType() + ":" + info.GetPoleName() + " поропущен так как уже есть в конфигурации " + "[" + info2.GetConfigurationName() + "]" + info2.GetDocumentType() + ":" + info2.GetPoleName())
	}
}

func (c *Configuration) addDocumentTypeInfo(info ITypeInfo, State *loadState) {
	var info2 ITypeInfo
	var ok bool
	if info2, ok = State.TypesInfo[info.GetDocumentType()]; !ok {
		State.TypesInfo[info.GetDocumentType()] = info
	} else {
		log.Println("[" + info.GetConfigurationName() + "]" + info.GetDocumentType() + " поропущен так как уже есть в конфигурации " + "[" + info2.GetConfigurationName() + "]" + info2.GetDocumentType())
	}
}
func (c *Configuration) GetPoleInfo(DocumentType string, PoleName string) (IPoleInfo, error) {
	if di, ok := c.polesInfo[DocumentType]; ok {
		if pi, ok := di[PoleName]; ok {
			return pi, nil
		}
		return nil, errors.New("Pole  not found. " + DocumentType + ":" + PoleName)
	}
	return nil, errors.New("DocumentType not found. " + DocumentType)
}

/*
func (c *Configuration) GetPolesInfo() []IPoleInfo {
	Num := 0
	for _, m := range c.polesInfo {
		Num += len(m)
	}
	poles := make([]IPoleInfo, Num)
	N := 0
	for _, m := range c.polesInfo {
		for _, p := range m {
			poles[N] = p
			N++
		}
	}
	return poles
}/**/

/*
func (c *Configuration) GetPolesInfoEx(DocumentType string, poles []string) map[string]IPoleInfo {
	if m, ok := c.polesInfo[DocumentType]; ok {
		polesInfo := map[string]IPoleInfo{}
		if len(poles) != 0 {
			for _, pole := range poles {

				p := strings.Split(pole, ".")
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
}*/

func (c *Configuration) loadConfigurationDepend(ConfigurationName string, State *loadState) error {
	State.ListDepend = append(State.ListDepend, ConfigurationName)
	return nil
	/*if (c.DefaultConfiguration.ContainsKey(ConfigurationName)){
	        return nil;
				}
	      if (State.ConfigList.ContainsKey(ConfigurationName)){
	        return nil;
				}
	      State.ConfigList[ConfigurationName] = getDocuments("Configuration", new IDocumentWhere[] {
	          new DocumentWhereCompare() { PoleName = "Configuration.Name", Value = ConfigurationName, Operation= DocumentWhereCompare.CmpOperation.Equally },
	        });
	      if (State.ConfigList[ConfigurationName].Count() == 0)
	        throw new ConfigurationException("", "", "Конфигурация не нейдена " + ConfigurationName);
	      var DependList = new List<IDocument>();
	      foreach (IDocument doc in State.ConfigList[ConfigurationName])
	        if ((String)doc["Configuration.Type"] == "Depend")
	          DependList.Add(doc);
	      DependList.Sort((A, B) => (Int32)B["Configuration.Depend.Prioritet"] - (Int32)A["Configuration.Depend.Prioritet"]);
	      foreach (IDocument doc in DependList)
	        State.QueueConfiguration.Enqueue((String)doc["Configuration.Depend.ConfigurationName"]);
	      return nil;*/
}

/*
func (c *Configuration) LoadConfiguration(ConfigurationName string) error {
	var err error
	State := newLoadState()
	// Загрузка всех данный для конфигурации
	State.QueueConfiguration = append(State.QueueConfiguration, ConfigurationName)
	for len(State.QueueConfiguration) > 0 {
		confName := State.QueueConfiguration[0]
		State.QueueConfiguration = State.QueueConfiguration[1:]
		if err = c.loadConfigurationDepend(confName, State); err != nil {
			return err
		}
	}
	//Вычисление порядка загрузки
	//State.ListDepend.Add(BaseConfigurationName);
	Loading := map[string]bool{}
	OrderLoad := []string{}
	for i := len(State.ListDepend) - 1; i >= 0; i-- {
		confName := State.ListDepend[i]
		if _, ok := Loading[confName]; !ok {
			Loading[confName] = true
			OrderLoad = append(OrderLoad, confName)
		}
	}
	//OrderLoad.Reverse();
	for i := 0; i < len(OrderLoad); i++ {
		confName := OrderLoad[i]
		if conf, ok := DefaultConfiguration[confName]; ok {
			conf(c, State)
		} else {
			//LoadConfiguration(confName, State.ConfigList[confName],State);
		}
	}
	State.setLoadState(c)
	return nil
}
/**/
/*void loadPoleInfo(IDocument doc, LoadState State)
  {
    var info = new PoleInfo();
    info.ConfigurationName = (String)doc["Configuration.Name"];
    info.Title = (String)doc["Configuration.Title"];
    info.DocumentType = (String)doc["Configuration.PoleInfo.DocumentType"];
    info.PoleName = (String)doc["Configuration.PoleInfo.PoleName"];
    info.PoleType = (String)doc["Configuration.PoleInfo.PoleType"];
    info.New = Boolean.Parse((String)doc["Configuration.PoleInfo.New"]);
    info.Edit = Boolean.Parse((String)doc["Configuration.PoleInfo.Edit"]);
    info.IndexType = (eIndexType)Enum.Parse(typeof(eIndexType), (String)doc["Configuration.PoleInfo.IndexType"]);
    switch (info.PoleType)
    {
      case "StringValue": info.Checker = new PoleCheckerStringValue(); break;
      case "Int32Value": info.Checker = new PoleCheckerInt32Value(); break;
      case "DocumentLinkValue": info.Checker = null; break;
      case "DateValue": info.Checker = null; break;
      case "DateTimeValue": info.Checker = null; break;
      default: throw new NotImplementedException();
    }
    if (info.Checker != null)
      info.Checker.Load(doc);
    addPoleInfo(info, State);
  }
  void loadDocumentTypeInfo(IDocument doc, LoadState State)
  {
    var info = new DocumentTypeInfo();
    info.ConfigurationName = (String)doc["Configuration.Name"];
    info.Title = (String)doc["Configuration.Title"];
    info.DocumentType = (String)doc["Configuration.DocumentType.Name"];
    info.New =  Boolean.Parse((String)doc["Configuration.DocumentType.New"]);
    info.Read = Boolean.Parse((String)doc["Configuration.DocumentType.Read"]);
    info.Save = Boolean.Parse((String)doc["Configuration.DocumentType.Save"]);
    addDocumentTypeInfo(info, State);
  }
  delegate void dLoadConfiguration(Configuration configuration, LoadState State);
  static Dictionary<String, dLoadConfiguration> DefaultConfiguration = new Dictionary<String, dLoadConfiguration>() {
    { BaseConfigurationName, loadBaseConfiguration },
    { BaseEditConfigurationName, loadEditBaseConfiguration },
    /*{ BaseCallConfigurationName, loadCallConfiguration },
    { BaseCallEditConfigurationName, loadCallEditConfiguration },
    { BaseNodeConfigurationName, loadNodeConfiguration },
    { BaseNodeEditConfigurationName, loadNodeEditConfiguration },
    { BaseTimerConfigurationName, loadTimerConfiguration },
    { BaseTimerEditConfigurationName, loadTimerEditConfiguration },*
  };*/
