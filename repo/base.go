package repo

import (
	"context"

	"app/common/adapter"
)

// MongoV2 struct
type MongoV2 struct {
	Session    *adapter.Mongo
	Collection string
	ctx        context.Context
}

// InitContext set global context when init
func (r *MongoV2) InitContext(ctx context.Context) {
	r.ctx = ctx
}

// GetContext get global context when init
func (r *MongoV2) GetContext() context.Context {
	if r.ctx == nil {
		return context.Background()
	}
	return r.ctx
}

// Cache struct
type Cache struct {
	Session *adapter.Cache
}
