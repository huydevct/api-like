package utils

import "encoding/json"

func ConvertSliceStructToSliceMap(slice interface{}) (results []map[string]interface{}) {
	j, _ := json.Marshal(slice)
	if err := json.Unmarshal(j, &results); err != nil {
		return
	}
	return
}
func ConvertStructToMap(istruct interface{}) (results map[string]interface{}) {
	j, _ := json.Marshal(istruct)
	if err := json.Unmarshal(j, &results); err != nil {
		return
	}
	return
}
