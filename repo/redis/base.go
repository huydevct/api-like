package redis

import (
	"sync"

	"app/common/adapter"
	"app/common/config"
)

// Redis struct
type Redis struct {
	Session *adapter.Redis
}

// Repo type
type Repo Redis

var (
	cfg      = config.GetConfig()
	instance *Repo
	once     sync.Once
)

// New ..
func New() *Repo {
	once.Do(func() {
		instance = &Repo{
			Session: config.GetConfig().Redis.Get("core"),
		}
	})

	return instance
}
