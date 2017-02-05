package core

const BaseConfigurationName = "Configuration"

func init() {
	initConfConfiguration()
	initConfEditConfiguration()
	initDocumentConfiguration()
}
func initConfConfiguration() {
	conf := NewConfiguration()
	conf.AddType("Configuration", "Configuration", false, true, true, "Конфигурация")
	conf.AddPole("Configuration", "Configuration", "Configuration.Name", "StringValue", NewObject(nil), "None", &PoleCheckerStringValue{}, true, true, "Имя конфигурации")
	conf.AddPole("Configuration", "Configuration", "Configuration.Type", "StringValue", NewObject(nil), "None", &PoleCheckerStringValue{List: map[string]bool{"PoleInfo": true, "DocumentType": true, "Depend": true}}, true, true, "Тип конфигурируемого элемента")
	conf.AddPole("Configuration", "Configuration", "Configuration.Title", "StringValue", NewObject(nil), "None", &PoleCheckerStringValue{}, true, true, "Описание элемента")
	conf.AddPole("Configuration", "Configuration", "Configuration.Depend.ConfigurationName", "StringValue", NewObject(nil), "None", &PoleCheckerStringValue{}, true, true, "Ссылка на конфигурацию")
	conf.AddPole("Configuration", "Configuration", "Configuration.Depend.Prioritet", "Int64Value", NewObject(100), "None", &PoleCheckerInt64Value{}, true, true, "Приоритет зависимости чем больше значение тем выше приоритет. По умолчанию 100")
	conf.AddPole("Configuration", "Configuration", "Configuration.DocumentType.Name", "StringValue", NewObject(nil), "None", &PoleCheckerStringValue{}, true, true, "Описание типа документа. Имя типа")
	conf.AddPole("Configuration", "Configuration", "Configuration.DocumentType.New", "StringValue", NewObject("True"), "None", &PoleCheckerStringValue{List: map[string]bool{"True": true, "False": true}}, true, true, "Описание типа документа. Разрешено создание таких документов")
	conf.AddPole("Configuration", "Configuration", "Configuration.DocumentType.Read", "StringValue", NewObject("True"), "None", &PoleCheckerStringValue{List: map[string]bool{"True": true, "False": true}}, true, true, "Описание типа документа. Разрешено чтение таких документов")
	conf.AddPole("Configuration", "Configuration", "Configuration.DocumentType.Save", "StringValue", NewObject("True"), "None", &PoleCheckerStringValue{List: map[string]bool{"True": true, "False": true}}, true, true, "Описание типа документа. Разрешено  сохранение таких документов")
	conf.AddPole("Configuration", "Configuration", "Configuration.PoleInfo.DocumentType", "StringValue", NewObject("True"), "None", &PoleCheckerStringValue{}, true, true, "Описание поля. Тип документа")
	conf.AddPole("Configuration", "Configuration", "Configuration.PoleInfo.PoleType", "StringValue", NewObject("True"), "None", &PoleCheckerStringValue{List: map[string]bool{"StringValue": true, "Int32Value": true, "DateTimeValue": true, "DateValue": true, "DocumentLinkValue": true}}, true, true, "Описание поля. Тип поля")
	conf.AddPole("Configuration", "Configuration", "Configuration.PoleInfo.PoleName", "StringValue", NewObject("None"), "None", &PoleCheckerStringValue{}, true, true, "Описание поля. Имя поля")
	conf.AddPole("Configuration", "Configuration", "Configuration.PoleInfo.IndexType", "StringValue", NewObject("True"), "None", &PoleCheckerStringValue{List: map[string]bool{"None": true, "Index": true, "Unique": true}}, true, true, "Описание поля. Индекс поля")
	conf.AddPole("Configuration", "Configuration", "Configuration.PoleInfo.New", "StringValue", NewObject("True"), "None", &PoleCheckerStringValue{List: map[string]bool{"True": true, "False": true}}, true, true, "Описание поля. Доступ на создание")
	conf.AddPole("Configuration", "Configuration", "Configuration.PoleInfo.Edit", "StringValue", NewObject("True"), "None", &PoleCheckerStringValue{List: map[string]bool{"True": true, "False": true}}, true, true, "Описание поля. Доступ на редактирование")
	conf.AddPole("Configuration", "Configuration", "Configuration.PoleInfo.CheckerInt64ValueMin", "Int64Value", NewObject(nil), "None", &PoleCheckerInt64Value{}, true, true, "Описание поля. Проверка Int Min")
	conf.AddPole("Configuration", "Configuration", "Configuration.PoleInfo.CheckerInt64ValueMax", "Int64Value", NewObject(nil), "None", &PoleCheckerInt64Value{}, true, true, "Описание поля. Проверка Int Max")
	conf.AddPole("Configuration", "Configuration", "Configuration.PoleInfo.CheckerInt64ValueList", "StringValue", NewObject(nil), "None", &PoleCheckerStringValue{}, true, true, "Описание поля. Проверка Int список доступных значений")
	conf.AddPole("Configuration", "Configuration", "Configuration.PoleInfo.CheckerStringValueList", "StringValue", NewObject(nil), "None", &PoleCheckerStringValue{}, true, true, "Описание поля. Проверка String список доступных значений")
	conf.AddPole("Configuration", "Configuration", "Configuration.PoleInfo.CheckerStringValueAllowNull", "StringValue", NewObject(nil), "None", &PoleCheckerStringValue{List: map[string]bool{"True": true, "False": true}}, true, true, "Описание поля. Может равнятся NULL")
	InitAddBaseConfiguration("Configuration", conf)
}
func initConfEditConfiguration() {
	conf := NewConfiguration()
	conf.AddType("ConfigurationEdit", "Configuration", true, true, true, "Конфигурация")
	conf.AddDependConfiguration("Configuration")
	InitAddBaseConfiguration("ConfigurationEdit", conf)
}
func initDocumentConfiguration() {
	conf := NewConfiguration()
	conf.AddType("Document", "Document", true, true, true, "Документ")
	conf.AddPole("Document", "Document", "Document.DocumentType", "StringValue", NewObject(nil), "Index", &PoleCheckerStringValue{}, true, true, "Тип документа")
	InitAddBaseConfiguration("Document", conf)
}
