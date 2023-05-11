package service

import (
	"context"
	"sync"

	"app/common/config"
	"app/repo"
)

// Repo type
type Repo struct {
	repo.Cache
	context.Context
}

var (
	instance *Repo
	once     sync.Once
)

// New ..
func New(ctx context.Context) *Repo {
	once.Do(func() {
		instance = &Repo{
			repo.Cache{
				Session: config.GetConfig().Cache.Get("core"),
			},
			ctx,
		}
	})
	// re-init context
	return instance
}
