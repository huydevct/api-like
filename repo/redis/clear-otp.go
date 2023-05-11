package redis

// ClearOTP : delete otp
func (r Repo) ClearOTP(key string) {
	r.Session.GetClient().Del(key)
	return
}
