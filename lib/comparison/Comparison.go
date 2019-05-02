package comparison

import (
	"fmt"
	"reflect"
)

func IsZero(v reflect.Value) bool {
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func IsEqual(rVal, eVal interface{}) bool {
	result := (rVal == eVal)
	if result {
		return result
	}
	rStr := fmt.Sprintf("%v", rVal)
	eStr := fmt.Sprintf("%v", eVal)
	return rStr == eStr
}

func BelongsTo(val interface{}, list []interface{}) bool {
	for _, item := range list {
		if IsEqual(val, item) {
			return true
		}
	}
	return false
}

func VariableInfo(label string, val interface{}) {
	fmt.Printf(" - %s: [%v], type: %s\n", label, val, reflect.ValueOf(val).Type().String())
}
