package comparison

import (
	"fmt"
)

func IsEqualTo(rVal, eVal interface{}) (bool, error) {
	result := (rVal == eVal)
	if result {
		return result, nil
	}
	rStr := fmt.Sprintf("%v", rVal)
	eStr := fmt.Sprintf("%v", eVal)
	return rStr == eStr, nil
}

func BelongsTo(val interface{}, list []interface{}) bool {
	for _, item := range list {
		if val, _ := IsEqualTo(val, item); val {
			return true
		}
	}
	return false
}
