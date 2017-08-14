package structures

func MapGetStringList(data map[string]interface{}, name string) ([]string, bool) {
	var ok bool
	var mapItem interface{}
	if mapItem, ok = data[name]; ok {
		var listInterf []interface{}
		if listInterf, ok = mapItem.([]interface{}); ok {
			result := make([]string, len(listInterf))
			for i, item := range listInterf {
				var str string
				if str, ok = item.(string); ok {
					result[i] = str
				} else {
					return nil, false
				}
			}
			return result, true
		}
	}
	return nil, false
}
func MapSetStringList(data map[string]interface{}, name string, list []string) {
	ilist := make([]interface{}, len(list))
	for i, item := range list {
		ilist[i] = item
	}
	data[name] = ilist
}
func MapGetInt64List(data map[string]interface{}, name string) ([]int64, bool) {
	var ok bool
	var mapItem interface{}
	if mapItem, ok = data[name]; ok {
		var listInterf []interface{}
		if listInterf, ok = mapItem.([]interface{}); ok {
			result := make([]int64, len(listInterf))
			for i, item := range listInterf {
				var v int64
				if v, ok = item.(int64); ok {
					result[i] = v
				} else {
					return nil, false
				}
			}
			return result, true
		}
	}
	return nil, false
}
func MapSetInt64List(data map[string]interface{}, name string, list []int64) {
	ilist := make([]interface{}, len(list))
	for i, item := range list {
		ilist[i] = item
	}
	data[name] = ilist
}
func MapGetString(data map[string]interface{}, name string) (string, bool) {
	var ok bool
	var mapItem interface{}
	if mapItem, ok = data[name]; ok {
		var str string
		if str, ok = mapItem.(string); ok {
			return str, true
		}
	}
	return "", false
}
func MapGetValue(data map[string]interface{}, name string) (interface{}, bool) {
	var ok bool
	var mapItem interface{}
	if mapItem, ok = data[name]; ok {
		return mapItem, true
	}
	return nil, false
}
