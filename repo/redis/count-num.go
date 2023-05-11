package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

// CountNum : count num fail by phone
func (r Repo) CountNum(key string) (num int, err error) {
	result, errRedis := r.Session.GetClient().Get(key).Result()
	if errRedis != nil && errRedis != redis.Nil {
		err = fmt.Errorf("Connect redis fail: %s", err)
		return
	}
	num, _ = strconv.Atoi(result)
	return
}
