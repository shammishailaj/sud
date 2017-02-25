package core

type Object interface{}

var NULL = struct{}{}

func IsNull(obj interface{}) bool {
	if obj == nil {
		return true
	}
	if obj == NULL {
		return true
	}
	return false
}
