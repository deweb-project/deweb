package lib_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"x.x/x/deweb/lib"
)

var TestParseDEIDCases = make(map[string]lib.DEID)

func TestParseDEID(t *testing.T) {
	TestParseDEIDCases["proto:identifier[key=asdasdasdasdasdasdasd]"] = lib.DEID{
		Protocol:   "proto",
		Identifier: "identifier",
		Key:        "asdasdasdasdasdasdasd",
		Extra:      map[string]string{},
	}
	for deid, correct := range TestParseDEIDCases {
		result, err := lib.ParseDEID(deid)
		if err != nil {
			if result.Protocol != "" {
				t.Error(deid, err)
				t.Fail()
			}
		} else {
			if result.Protocol == "" {
				t.Error(deid, "no error (should be)")
				t.Fail()
			}
		}
		if result.Protocol == correct.Protocol &&
			result.Identifier == correct.Identifier &&
			result.Key == correct.Key &&
			reflect.DeepEqual(result.Extra, correct.Extra) {
			//ok
		} else {
			r, _ := json.MarshalIndent(result, "", "  ")
			c, _ := json.MarshalIndent(correct, "", "  ")
			t.Error("result != correct", deid, "\n"+string(r)+"\n"+string(c))
		}
	}
}
