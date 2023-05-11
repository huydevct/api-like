package service

import (
	"fmt"
)

// SetCache : Lấy random clone
func (r Repo) SetCache(key string, data string) (err error) {
	cache := r.Cache.Session.ConClient
	if err != nil {
		fmt.Println(err)
	}
	err = cache.Put([]byte(key), []byte(data))
	return
}

// DeleteCache : Lấy random clone
func (r Repo) DeleteCache(key string) (err error) {
	cache := r.Cache.Session.ConClient
	if err != nil {
		fmt.Println(err)
	}
	err = cache.Delete([]byte(key))
	return
}
