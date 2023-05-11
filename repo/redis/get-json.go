package redis

import "github.com/go-redis/redis"

// GetJSON : Lấy thông tin json
func (r Repo) GetJSON(key string) (JSON string, err error) {
	value, err := r.Session.GetClient().Get(key).Result()
	if err != nil && err != redis.Nil {
		return
	}

	JSON = value
	return
}
