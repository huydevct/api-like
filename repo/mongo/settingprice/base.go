package settingprice

import (
	"context"
	"sync"

	"app/common/config"
	"app/repo"
)

// Repo type
type Repo struct {
	repo.MongoV2
}

var (
	instance *Repo
	once     sync.Once
)

// New ..
func New(ctx context.Context) *Repo {
	once.Do(func() {
		instance = &Repo{
			repo.MongoV2{
				Session:    config.GetConfig().Mongo.Get("core"),
				Collection: "setting_prices",
			},
		}
	})
	// re-init context
	instance.InitContext(ctx)
	return instance
}
