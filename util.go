package main

import (
	"encoding/json"
	"reflect"
)

func jsonEqual(a, b string) bool {
	var ja, jb interface{}
	json.Unmarshal([]byte(a), &ja)
	json.Unmarshal([]byte(b), &jb)
	return reflect.DeepEqual(ja, jb)
}
