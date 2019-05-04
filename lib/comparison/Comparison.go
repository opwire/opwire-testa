package comparison

import (
	"fmt"
)

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
