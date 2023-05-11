package service

// GetCacheBy : Lấy random clone
func (r Repo) GetCacheBy(key string) (data string, err error) {
	cache := r.Cache.Session.ConClient
	value, err := cache.Get([]byte(key))
	data = string(value)
	return
}
