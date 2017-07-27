package core

import (
	"github.com/crazyprograms/sud/corebase"
)

const BaseConfigurationName = "Configuration"

func init() {

}
func initConfConfiguration(c *Core) {
	conf := NewConfiguration([]string{"ConfigurationReader"})
	conf.AddType(TypeInfo{ConfigurationName: "Configuration", RecordType: "Configuration", AccessType: "Free", AccessNew: "ConfigurationEditor", AccessRead: "ConfigurationReader", AccessSave: "ConfigurationEditor", Title: "Конфигурация"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.Name", PoleType: "StringValue", Default: corebase.NULL, IndexType: "None", Checker: &PoleCheckerStringValue{}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Имя конфигурации"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.Type", PoleType: "StringValue", Default: corebase.NULL, IndexType: "None", Checker: &PoleCheckerStringValue{List: map[string]bool{"PoleInfo": true, "RecordType": true, "Depend": true}}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Тип конфигурируемого элемента"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.Title", PoleType: "StringValue", Default: corebase.NULL, IndexType: "None", Checker: &PoleCheckerStringValue{}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание элемента"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.Depend.ConfigurationName", PoleType: "StringValue", Default: corebase.NULL, IndexType: "None", Checker: &PoleCheckerStringValue{}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Ссылка на конфигурацию"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.Depend.Prioritet", PoleType: "Int64Value", Default: 100, IndexType: "None", Checker: &PoleCheckerInt64Value{}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Приоритет зависимости чем больше значение тем выше приоритет. По умолчанию 100"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.RecordType.Name", PoleType: "StringValue", Default: corebase.NULL, IndexType: "None", Checker: &PoleCheckerStringValue{}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание типа документа. Имя типа"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.RecordType.New", PoleType: "StringValue", Default: "True", IndexType: "None", Checker: &PoleCheckerStringValue{List: map[string]bool{"True": true, "False": true}}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание типа документа. Разрешено создание таких документов"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.RecordType.Read", PoleType: "StringValue", Default: "True", IndexType: "None", Checker: &PoleCheckerStringValue{List: map[string]bool{"True": true, "False": true}}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание типа документа. Разрешено чтение таких документов"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.RecordType.Save", PoleType: "StringValue", Default: "True", IndexType: "None", Checker: &PoleCheckerStringValue{List: map[string]bool{"True": true, "False": true}}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание типа документа. Разрешено  сохранение таких документов"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.PoleInfo.RecordType", PoleType: "StringValue", Default: "True", IndexType: "None", Checker: &PoleCheckerStringValue{}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание поля. Тип документа"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.PoleInfo.PoleType", PoleType: "StringValue", Default: "True", IndexType: "None", Checker: &PoleCheckerStringValue{List: map[string]bool{"StringValue": true, "Int32Value": true, "DateTimeValue": true, "DateValue": true, "RecordLinkValue": true}}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание поля. Тип поля"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.PoleInfo.PoleName", PoleType: "StringValue", Default: "None", IndexType: "None", Checker: &PoleCheckerStringValue{}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание поля. Имя поля"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.PoleInfo.IndexType", PoleType: "StringValue", Default: "True", IndexType: "None", Checker: &PoleCheckerStringValue{List: map[string]bool{"None": true, "Index": true, "Unique": true}}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание поля. Индекс поля"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.PoleInfo.New", PoleType: "StringValue", Default: "True", IndexType: "None", Checker: &PoleCheckerStringValue{List: map[string]bool{"True": true, "False": true}}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание поля. Доступ на создание"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.PoleInfo.Edit", PoleType: "StringValue", Default: "True", IndexType: "None", Checker: &PoleCheckerStringValue{List: map[string]bool{"True": true, "False": true}}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание поля. Доступ на редактирование"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.PoleInfo.CheckerInt64ValueMin", PoleType: "Int64Value", Default: corebase.NULL, IndexType: "None", Checker: &PoleCheckerInt64Value{}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание поля. Проверка Int Min"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.PoleInfo.CheckerInt64ValueMax", PoleType: "Int64Value", Default: corebase.NULL, IndexType: "None", Checker: &PoleCheckerInt64Value{}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание поля. Проверка Int Max"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.PoleInfo.CheckerInt64ValueList", PoleType: "StringValue", Default: corebase.NULL, IndexType: "None", Checker: &PoleCheckerStringValue{}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание поля. Проверка Int список доступных значений"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.PoleInfo.CheckerStringValueList", PoleType: "StringValue", Default: corebase.NULL, IndexType: "None", Checker: &PoleCheckerStringValue{}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание поля. Проверка String список доступных значений"})
	conf.AddPole(PoleInfo{ConfigurationName: "Configuration", RecordType: "Configuration", PoleName: "Configuration.PoleInfo.CheckerStringValueAllowNull", PoleType: "StringValue", Default: corebase.NULL, IndexType: "None", Checker: &PoleCheckerStringValue{List: map[string]bool{"True": true, "False": true}}, AccessRead: "ConfigurationReader", AccessWrite: "ConfigurationEditor", Title: "Описание поля. Может равнятся NULL"})
	c.AddBaseConfiguration("Configuration", conf)
}
func initDocumentConfiguration(c *Core) {
	conf := NewConfiguration([]string{"Document"})
	conf.AddType(TypeInfo{ConfigurationName: "Document", RecordType: "Document", AccessType: "Check", AccessNew: "DocumentEditor", AccessRead: "DocumentReader", AccessSave: "DocumentEditor", Title: "Документ"})
	conf.AddPole(PoleInfo{ConfigurationName: "Document", RecordType: "Document", PoleName: "Document.DocumentType", PoleType: "StringValue", Default: corebase.NULL, IndexType: "Index", Checker: &PoleCheckerStringValue{}, AccessRead: "DocumentReader", AccessWrite: "DocumentEditor", Title: "Тип документа"})
	c.AddBaseConfiguration("Document", conf)
}

func InitBaseConfiguration(c *Core) {
	initConfConfiguration(c)
	initDocumentConfiguration(c)
}
