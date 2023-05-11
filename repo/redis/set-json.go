package redis

import (
	"encoding/json"
	"time"
)

// SetJSON : cache otp
func (r Repo) SetJSON(key string, value interface{}, ttl int) (err error) {
	storeByte, _ := json.Marshal(value)
	_, err = r.Session.GetClient().Set(key, storeByte, time.Second*time.Duration(ttl)).Result()
	return
}
