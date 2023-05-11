package redis

import (
	"time"
)

// SetOTP : cache otp
func (r Repo) SetOTP(key, otp string, ttl int) (err error) {
	_, err = r.Session.GetClient().Set(key, otp, time.Second*time.Duration(ttl)).Result()
	return
}
