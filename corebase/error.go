package corebase

type Error struct {
	ErrorType string
	Action    string
	Name      string
	RecordUID string
	Info      string
}

const ErrorTypeInfo = "Info"
const ErrorTypeAccessIsDenied = "Access is denied"
const ErrorTypeNotFound = "Not found"
const ErrorTypeAlreadyExists = "Already exists"
const ErrorTypeConvert = "Convert"
const ErrorTimeout = "Timeout"
const ErrorFormat = "Format"

func (e *Error) Error() string {
	return e.ErrorType + ";" + e.Action + ";" + e.Name + ";" + e.RecordUID + ";" + e.Info
}