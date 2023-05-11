package redis

import "github.com/go-redis/redis"

// GetOTPTtl : kiểm tra mã OTP đã được gửi đến khách hay chưa ?
func (r Repo) GetOTPTtl(key string) (ttl int, err error) {
	exist, errRedis := r.Session.GetClient().Exists(key).Result()
	if errRedis != nil {
		err = errRedis
		return
	}
	if exist > 0 { //đã tồn tại trước đó
		//get TTL
		duration, _ := r.Session.GetClient().TTL(key).Result()
		ttl = int(duration.Seconds())
	}
	return
}

// GetOTP : Lấy mã otp
func (r Repo) GetOTP(key string) (otp string, err error) {
	value, err := r.Session.GetClient().Get(key).Result()
	if err != nil && err != redis.Nil {
		return
	}

	otp = value
	return
}
