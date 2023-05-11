package utils

import (
	"app/model"
	"encoding/json"

	b64 "encoding/base64"
)

// EncodeBase64 : encode input to format base64
func EncodeBase64(input interface{}) (result string) {

	dataByte, _ := json.Marshal(input)
	// encode base64 result and response
	result = b64.StdEncoding.EncodeToString([]byte(dataByte))
	return
}

// DecodeBase64 : ..
func DecodeBase64(input string) (result string) {
	if input != "" {
		temp, err := b64.StdEncoding.DecodeString(input)
		if err == nil {
			result = string(temp)
		}
	}
	return
}

// DecodeBase64ToClone : decode input to struct cloneInfo
func DecodeBase64ToClone(input string) (result model.CloneInfo, err error) {

	sDec, _ := b64.StdEncoding.DecodeString(input)
	err = json.Unmarshal(sDec, &result)

	return
}
