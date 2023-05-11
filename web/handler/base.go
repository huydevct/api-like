package handler

import (
	"app/constants"
	"app/model"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"app/common/config"
	actionProfileRepo "app/repo/mongo/actionprofile"
	cloneRepo "app/repo/mongo/clone"
	deviceRepo "app/repo/mongo/device"
	settingRepo "app/repo/mongo/setting"
	settingPriceRepo "app/repo/mongo/settingprice"
	tokenRepo "app/repo/mongo/token"
	userRepo "app/repo/mongo/user"
	redisRepo "app/repo/redis"

	gHandler "app/common/gstuff/handler"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	cfg = config.GetConfig()
)

func success(data interface{}) (int, interface{}) {
	return gHandler.Success(data)
}

func successRaw(data, dataRaw interface{}) (int, interface{}) {
	return gHandler.SuccessRaw(data, dataRaw)
}

func getRequestID(c echo.Context) (requestID string) {
	return gHandler.GetRequestID(c)
}

/*
	1. check password theo hash trong hệ thống cũ truớc
	2. nếu không match thì check bscrypt của hệ thống mới
*/
func checkPasswordMatch(originalPassword, hashPassword string) (isOK bool) {

	if generateHash(originalPassword) == hashPassword {
		isOK = true
		return
	}
	// Check tiếp bằng bscrypt password cho hệ thống mới
	return checkBscryptPassword(originalPassword, hashPassword)
}

func checkBscryptPassword(originalPassword, hashPassword string) (isOK bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(originalPassword))
	if err != nil {
		return
	}
	// match password
	isOK = true
	return
}

func generateHash(input string) (result string) {
	b := md5.Sum([]byte(input))
	result = hex.EncodeToString(b[:])
	return
}

// bscryptPassword: Mã hóa bscrypt password
func bscryptPassword(originalPassword string) (result string, err error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(originalPassword), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("Generate bscrypt passwordd fail: %s", err)
		return
	}
	result = string(hashPassword)
	return
}

// countNumLoginFail : đếm số lần thất bại khi đăng nhập theo điện thoại
func countNumLoginFail(phone string) (numLogin int, err error) {
	key := fmt.Sprintf(constants.NumberLoginFailKey, phone)
	numLogin, err = redisRepo.New().CountNum(key)
	return
}

// setNumLoginFail : set số lần thất bại khi đăng nhập theo điện thoại
func setNumLoginFail(phone string, numLoginFail, ttl int) (err error) {
	key := fmt.Sprintf(constants.NumberLoginFailKey, phone)
	return redisRepo.New().CacheNum(key, numLoginFail, ttl)
}

// countNumChangePassFail : đếm số lần thay đổi pass thất bại theo điện thoại
func countNumChangePassFail(phone string) (numLogin int, err error) {
	key := fmt.Sprintf(constants.NumberChangePasswordFailKey, phone)
	numLogin, err = redisRepo.New().CountNum(key)
	return
}

// setNumChangePassFail : set số lần thất bại khi thay đổi pass theo điện thoại
func setNumChangePassFail(phone string, numLoginFail, ttl int) (err error) {
	key := fmt.Sprintf(constants.NumberChangePasswordFailKey, phone)
	return redisRepo.New().CacheNum(key, numLoginFail, ttl)
}

// countNumSendOTPResetPassword : Đếm số lần gửi mã OTP reset password thành công
func countNumSendOTPResetPassword(phone string) (numReset int, err error) {
	key := fmt.Sprintf(constants.NumberSendOTPResetPasswordKey, phone)
	numReset, err = redisRepo.New().CountNum(key)
	return
}

// setNumSendOTPResetPassword : set số lần gửi mã OTP reset password thành công
func setNumSendOTPResetPassword(phone string, numSendOTP, ttl int) (err error) {
	key := fmt.Sprintf(constants.NumberSendOTPResetPasswordKey, phone)
	return redisRepo.New().CacheNum(key, numSendOTP, ttl)
}

// genOTP : tạo mã OTP
// Mã có thời hạn 10 phút
func genOTP(typeOTP, phone string) (OTP string, ttl int, err error) {
	key := fmt.Sprintf(constants.OTPKey, typeOTP, phone)

	// Kiểm tra đã gửi mã cho khách hay chưa ?
	ttl, err = redisRepo.New().GetOTPTtl(key)
	if err != nil {
		return
	}
	if ttl > 0 { // Đã giở mã cho khách
		return
	}

	// Tạo mã mới
	var letterRunes = []rune("023456789")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	OTP = string(b)
	// OTP = "123456"

	// Lưu mã OTP trong 10 phút
	err = redisRepo.New().SetOTP(key, OTP, 600)
	if err != nil {
		return
	}

	return
}

// verifyOTP : Kiểm tra mã OTP có hợp lệ hay không ?
func verifyOTP(typeOTP, phone, otp string) (isOK bool, err error) {
	key := fmt.Sprintf(constants.OTPKey, typeOTP, phone)
	result, err := redisRepo.New().GetOTP(key)
	if err != nil {
		return
	}

	if otp == result {
		isOK = true
		return
	}

	return
}

// clearOTP : Xóa OTP
func clearOTP(typeOTP, phone string) {
	key := fmt.Sprintf(constants.OTPKey, typeOTP, phone)
	redisRepo.New().ClearOTP(key)
	return
}

// Lấy thông tin actionProfile theo lazyload
func getActionProfileByID(actionProfileID *primitive.ObjectID) (actionProfile model.ActionProfile, err error) {

	actionProfile, err = actionProfileRepo.New(context.Background()).GetOneByID(*actionProfileID)
	if err != nil {
		err = fmt.Errorf("Get action profile %s", err)
		return
	}

	return
}

// getUserInfoByToken : Lấy thông tin người dùng theo lazyload
func getUserInfoByToken(token string) (userInfo model.UserInfo, err error) {

	userInfo, err = userRepo.New(context.Background()).GetOneActiveByToken(token)
	if err != nil {
		err = fmt.Errorf("Get user info %s", err)
		return
	}
	return
}

// Lấy thông tin clone theo cloneID theo lazyload
func getCloneByID(cloneID primitive.ObjectID) (clone model.CloneInfo, err error) {

	clone, err = cloneRepo.New(context.Background()).GetOneActiveByID(cloneID)
	if err != nil {
		err = fmt.Errorf("Get clone %s", err)
		return
	}
	return
}

// getSetting : Lấy thông tin cài đặt giá theo lazyload
func getSetting() (setting model.Setting, err error) {
	setting, err = settingRepo.New(context.Background()).GetOneSetting()
	if err != nil {
		err = fmt.Errorf("Get setting %s", err)
		return
	}
	return
}

// getDeviceByPCName: lấy thông tin device theo pc name
func getDeviceByPCName(pcName string) (device model.DeviceInfo, err error) {
	time := time.Now()
	deviceInfo := model.DeviceInfo{
		PCName:      pcName,
		Status:      constants.Active,
		CreatedDate: &time,
		UpdatedDate: &time,
	}
	err = deviceRepo.New(context.Background()).Upsert(deviceInfo)
	if err != nil {
		err = fmt.Errorf("Upsert device %s", err)
		return
	}
	return
}

func getSettingPrices(appname string) (settingPrices model.SettingPrice, err error) {
	// re-init
	settingPrices, err = settingPriceRepo.New(context.Background()).GetOneByAppname(appname)
	if err != nil {
		err = fmt.Errorf("Get setting prices %s", err)
		return
	}
	return
}

func delGetDeviceByDeviceID(DeviceID string) {
	redisClient := cfg.Redis["core"].GetClient()
	deviceKey := fmt.Sprintf(constants.DeviceKey, DeviceID)
	redisClient.Del(deviceKey)
}

// validateDeviceIDBelongToken:
// 1. Kiểm tra DeviceID này đã duyệt hay chưa ?
// 2. Kiểm tra andoirdID có thuộc về token này hay không ?
func validateDeviceIDBelongToken(pcName string, token string) (err error) {

	// lấy thông tin device theo lazyload
	_, err = getDeviceByPCName(pcName)
	if err != nil {
		return
	}

	return
}

// GenTokenByUserIDDeviceID : Tạo mới token theo userID, deviceID
// . Kiểm tra có token nào còn hợp lệ với userID, deviceID hay không ?
// . Nếu có trả luôn token này, không tạo mới token
func GenTokenByUserIDDeviceID(data model.Token) (token string, err error) {
	tokenRepoInstance := tokenRepo.New(context.Background())

	// TODO: Kiểm tra có token nào hợp lệ với userID, deviceID này không ?
	condition := bson.M{
		"user_id":   data.UserID,
		"type":      data.Type,
		"device_id": data.DeviceID,
		"status":    constants.Active,
		"expired_date": bson.M{
			"$gt": time.Now(),
		},
	}
	result, err := tokenRepoInstance.Detail(condition)
	if err != nil {
		err = fmt.Errorf("Get detail token fail: %s", err)
		return
	}
	if result.IsExists() { // Đã tồn tại token
		token = result.Token
		return
	}

	// gen token
	temp, err := uuid.NewUUID()
	if err != nil {
		err = fmt.Errorf("Gen token by UUID fail: %s", err)
		return
	}

	tokenModel := model.Token{
		Token:              temp.String(),
		ExpiredAfterSecond: data.ExpiredAfterSecond,
		UserID:             data.UserID,
		UserToken:          data.UserToken,
		Type:               data.Type,
		Status:             constants.Active,
		DeviceID:           data.DeviceID,
		UserAgent:          data.UserAgent,
		RemoteIP:           data.RemoteIP,
		Source:             data.Source,
	}

	err = tokenRepoInstance.Create(tokenModel)
	if err != nil {
		err = fmt.Errorf("Tạo mới token: %s", err)
		return
	}

	token = tokenModel.Token
	return
}

// ValidateToken : Kiểm tra token hợp lệ hay không ?
// . Token tồn tại
// . Còn hạn sử dụng
func ValidateToken(token, userType string) (userID primitive.ObjectID, userToken string, err error) {

	condition := bson.M{
		"token":  token,
		"type":   userType,
		"status": constants.Active,
	}
	result, err := tokenRepo.New(context.Background()).Detail(condition)
	if err != nil {
		err = fmt.Errorf("Get detail token fail: %s", err)
		return
	}
	if !result.IsExists() {
		err = fmt.Errorf("Token không hợp lệ")
		return
	}
	// // Kiểm tra xem token còn hạn sử dụng hay không ?
	// if value := time.Now().Sub(*result.ExpiredDate); value > 0 {
	// 	err = fmt.Errorf("Phiên đăng nhập đã kết thúc")
	// 	return
	// }

	userID = result.UserID
	userToken = result.UserToken
	return
}
