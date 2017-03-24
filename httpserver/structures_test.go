package httpserver

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func testPack(t *testing.T) {
	var Templates = []interface{}{int64(123), "ABC", time.Now(), []byte("bytes")}
	var err error
	for _, value := range Templates {
		var data []byte
		var p *jsonParam
		var p2 *jsonParam
		p, err = jsonPack(value)
		if data, err = json.Marshal(&p); err == nil {
			if err = json.Unmarshal(data, &p2); err == nil {
				v, _ := jsonUnPack(p2)
				if reflect.TypeOf(v) != reflect.TypeOf(value) {
					t.Error("packe un pack error", reflect.TypeOf(v), reflect.TypeOf(value))
				}
			}
		}
		if err != nil {
			t.Error(err)
		}
	}
}
func TestPackB(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		testPack(t)
	}
}
