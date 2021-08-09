package format

import (
	"reflect"
	"strings"
)

// Print formats data in json or tabular form.
func Print(output string, data interface{}) {
	if isNil(data) {
		return
	}
	switch strings.ToLower(output) {
	case "tab":
		PrintTab(data)
	default:
		PrintJSON(data)
	}
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
