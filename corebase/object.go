package corebase

//type Object interface{}
type TNULL struct{}

var NULL = TNULL{}

func IsNull(obj interface{}) bool {
	if obj == nil {
		return true
	}
	switch obj.(type) {
	case TNULL:
		return true

	}
	return false
}
