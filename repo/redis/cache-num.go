package redis

import (
	"time"
)

// CacheNum : cache num fail
func (r Repo) CacheNum(key string, numFail, ttl int) (err error) {
	_, err = r.Session.GetClient().Set(key, numFail, time.Second*time.Duration(ttl)).Result()
	return
}
