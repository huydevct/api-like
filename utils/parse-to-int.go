package utils

import (
	"strconv"
)

// ConvertToInt : chuyển data input to int
func ConvertToInt(input interface{}) (result int) {

	if value, ok := input.(int); ok {
		result = value
	} else if value, ok := input.(int32); ok {
		result = int(value)
	} else if value, ok := input.(int64); ok {
		result = int(value)
	} else if value, ok := input.(float32); ok {
		result = int(value)
	} else if value, ok := input.(float64); ok {
		result = int(value)
	} else if value, ok := input.(string); ok {
		if i, err := strconv.Atoi(value); err == nil {
			result = i
		}
	}

	return
}
