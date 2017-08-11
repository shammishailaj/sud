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
