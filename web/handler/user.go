package handler

import (
	"app/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	esmsAPI "app/apicall"
	"app/constants"
	"app/model"
	codeRepo "app/repo/mongo/code"
	tokenRepo "app/repo/mongo/token"
	userRepo "app/repo/mongo/user"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserHandler : struct init user
type UserHandler struct{}

// NewUserHandler : Tạo mới 1 đối tượng user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// Register : Đăng ký người dùng
// User trong hệ thống
func (h UserHandler) Register(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		Username       string `json:"username" query:"username" validate:"required"`
		AgencyCode     string `json:"agency_code" query:"agency_code"`
		ReferenceCode  string `json:"reference_code" query:"reference_code"`
		UserInviteCode string `json:"user_invite_code" query:"user_invite_code"`
		Fullname       string `json:"fullname" query:"fullname" validate:"required"`
		Password       string `json:"password" query:"password" validate:"required,min=6"`
	}
	request := new(myRequest)
	userInfo := new(model.UserInfo)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// if request.AgencyCode == "" && request.UserInviteCode == "" {
	// 	return fmt.Errorf("Mã giới thiệu không được để trống")
	// }
	// TODO: validate
	{
		// validate phone
		_, err := utils.PhoneValidate(request.Username)
		if err != nil {
			return fmt.Errorf("Số điện thoại không hợp lệ %s", err)
		}
	}
	userRepoInstance := userRepo.New(httpCtx)
	// userInviteToken, err := h.genUserToken()
	// if err != nil {
	// 	return
	// }

	// if request.AgencyCode != "" {
	// 	checkUserAgency, err := userRepoInstance.GetOneByAgentCode(request.AgencyCode)
	// 	if err != nil {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin user theo agent code %s", err))
	// 	}
	// 	if !checkUserAgency.IsExists() {
	// 		return fmt.Errorf("Agent code không tồn tại")
	// 	}
	// 	userInfo.UserInviteCode = userInviteToken
	// } else {
	// 	checkUserRefer, err := userRepoInstance.GetOneByInviteCode(request.UserInviteCode)
	// 	if err != nil {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin user theo invite code %s", err))
	// 	}
	// 	if !checkUserRefer.IsExists() {
	// 		return fmt.Errorf("Reference code không tồn tại")
	// 	}
	// 	userInfo.ReferenceCode = request.UserInviteCode
	// }

	{
		// TODO: Kiểm tra đã có tài khoản đăng ký với sdt này hay chưa ?
		user, err := userRepoInstance.GetOneByPhone(request.Username)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Kiểm tra user tồn tại theo số thoại %s", err))
		}
		if user.IsExists() {
			return fmt.Errorf("Tài khoản đã tồn tại với số điện thoại %s", request.Username)
		}
	}
	// TODO: hash pasword md5
	hashPash := generateHash(request.Password)

	userToken, err := h.genUserToken()
	if err != nil {
		return
	}

	userInfo.Fullname = request.Fullname
	userInfo.Username = request.Username
	userInfo.Password = hashPash
	userInfo.Token = userToken
	userInfo.Status = constants.Active
	userInfo.DebugMode = constants.DebugModeProduction
	userInfo.Balance = constants.Balance
	result, err := userRepoInstance.Insert(userInfo)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Đăng ký người dùng %s", err))
	}

	return c.JSON(success(result))
}

// UpdateInfo : Chỉnh sửa thông tin  người dùng
// . Thông tin
func (UserHandler) UpdateInfo(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	userID := c.Get("user_id").(primitive.ObjectID)

	type myRequest struct {
		Fullname             string                       `json:"fullname" query:"fullname"`
		AccountName          string                       `json:"account_name" query:"account_name"`
		BankName             string                       `json:"bank_name" query:"bank_name"`
		ResonanceCode        string                       `json:"resonance_code" query:"resonance_code"`
		BankNumber           string                       `json:"bank_number" query:"bank_number"`
		FBLink               string                       `json:"fb_link" query:"fb_link"`
		ActionProfileDefault []model.ActionProfileDefault `json:"action_profile_default" query:"action_profile_default"`
		Email                string                       `json:"email" query:"email"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	userRepoInstance := userRepo.New(httpCtx)

	// Kiểm tra tài khoản có tồn tài hay không ?
	oldUser, err := userRepoInstance.GetOneActiveByID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng %s", err))
	}
	if !oldUser.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}
	{
		// update name
		if request.Fullname != "" {
			oldUser.Fullname = request.Fullname
		}
	}
	{
		// update AccountName
		if request.AccountName != "" {
			oldUser.AccountName = request.AccountName
		}
	}
	{
		// update AccountName
		if len(request.ActionProfileDefault) > 0 {
			oldUser.ActionProfileDefault = request.ActionProfileDefault
		}
	}
	{
		// update BankName
		if request.BankName != "" {
			oldUser.BankName = request.BankName
		}
	}
	{
		// update BankNumber
		if request.BankNumber != "" {
			oldUser.BankNumber = request.BankNumber
		}
	}
	{
		// update ResonanceCode
		if request.ResonanceCode != "" {
			oldUser.ResonanceCode = request.ResonanceCode
		}
	}
	{
		// update fblink
		if request.FBLink != "" {
			oldUser.FBLink = request.FBLink
		}
	}
	{
		// update Email
		if request.Email != "" {
			oldUser.Email = request.Email
		}
	}
	// Cập nhật status người dùng
	err = userRepoInstance.Update(oldUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật thông tin người dùng %s", err))
	}

	return c.JSON(success(oldUser))
}

// DetailByToken : Lấy thông tin người dùng theo token
func (UserHandler) DetailByToken(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	userID := c.Get("user_id").(primitive.ObjectID)

	user, err := userRepo.New(httpCtx).GetOneByID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng theo ID %s", err))
	}
	if !user.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}

	return c.JSON(success(user))
}

// CheckUserExists: Kiểm tra sdt và token có tồn tại hay không
func (UserHandler) CheckUserExists(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		IsExists bool `json:"is_exists"`
	}
	type myRequest struct {
		Phone string `json:"username" query:"username" validate:"required"`
	}
	request := new(myRequest)
	response := new(myResponse)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}

	user, err := userRepo.New(httpCtx).GetOneByPhone(request.Phone)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng theo phone %s", err))
	}

	if user.IsExists() {
		response.IsExists = true
	}

	return c.JSON(success(response))
}

// Login : Đăng nhập
// Limit số lần đăng nhập thất bại
func (h UserHandler) Login(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		Phone    string `json:"username" query:"username" validate:"required"`
		Password string `json:"password" query:"password" validate:"required,min=6"`
		DeviceID string `json:"device_id" query:"device_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// TODO: validate
	{
		// validate phone
		request.Phone, err = utils.PhoneValidate(request.Phone)
		if err != nil {
			return fmt.Errorf("Số điện thoại không hợp lệ %s", err)
		}
	}

	// Lấy thông tin user by phone
	user, err := userRepo.New(httpCtx).GetOneByPhone(request.Phone)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Tài khoản hoặc mật khẩu không đúng. Vui lòng thử lại %s", err))
	}
	if !user.IsExists() {
		return fmt.Errorf("Tài khoản hoặc mật khẩu không đúng. Vui lòng thử lại")
	}
	// TODO : Check password match
	if !checkPasswordMatch(request.Password, user.Password) {
		return fmt.Errorf("Tài khoản hoặc mật khẩu không đúng. Vui lòng thử lại")
	}
	// TODO: Lưu token với period 7 ngày
	tokenModel := model.Token{
		ExpiredAfterSecond: 7 * 24 * 3600,
		// token info
		UserID:    user.ID,
		UserToken: user.Token,
		Type:      "user",
		DeviceID:  request.DeviceID,
		UserAgent: c.Request().UserAgent(),
		RemoteIP:  c.RealIP(),
	}

	token, err := GenTokenByUserIDDeviceID(tokenModel)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(success(map[string]interface{}{
		"token":             token,
		"id":                user.ID,
		"user_token":        user.Token,
		"fullname":          user.Fullname,
		"username":          user.Username,
		"account_name":      user.AccountName,
		"bank_name":         user.BankName,
		"bank_number":       user.BankNumber,
		"Connectivity":      user.Connectivity,
		"balance":           user.Balance,
		"check_app":         user.CheckApp,
		"is_reg":            user.IsReg,
		"isLikeSub":         user.IsLikeSub,
		"is_like_sub":       user.IsLikeSub,
		"is_instagram":      user.IsInstagram,
		"is_youtube":        user.IsYoutube,
		"user_100_app":      user.User100App,
		"role":              user.Role,
		"status":            user.Status,
		"resonance_code":    user.ResonanceCode,
		"fblink":            user.FBLink,
		"email":             user.Email,
		"share_live_stream": user.ShareLiveStream,
	}))
}

// Logout : Đăng xuất
func (h UserHandler) Logout(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	token := c.Get("token").(string)

	userType := "user"
	_, err = tokenRepo.New(httpCtx).ClearToken(token, userType)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Token không hợp lệ")
	}

	return c.JSON(success(nil))
}

// ChangePassword : Thay đổi password
// . Required old password
func (UserHandler) ChangePassword(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	userID := c.Get("user_id").(primitive.ObjectID)

	type myRequest struct {
		OldPassword     string `json:"old_password" query:"old_password" validate:"required"`
		NewPassword     string `json:"new_password" query:"new_password" validate:"required"`
		ComfirmPassword string `json:"confirm_password" query:"confirm_password" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// TODO: validate
	{
		if len(request.NewPassword) < 6 {
			return fmt.Errorf("Chiều dài mật khẩu không hợp lệ")
		}
		// validate NewPassword, ComfirmPassword
		if request.NewPassword != request.ComfirmPassword {
			return fmt.Errorf("Mật khẩu xác nhận không đúng")
		}
	}
	// TODO: Đếm số lần change password fail theo phone
	numChangePassFail, errCount := countNumChangePassFail(userID.Hex())
	if errCount != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Đếm số lần thay đổi mật khẩu thất bại %s", err))
	}
	numChangePassFail++

	if numChangePassFail > 3 {
		return echo.NewHTTPError(http.StatusLocked, "Tài khoản của bạn đang tạm thời bị khoá trong vòng 60s do thay đổi mật khẩu sai quá nhiều lần. Vui lòng chờ và thực hiện lại.")
	}

	// Lấy thông tin user by phone
	userRepoInstance := userRepo.New(httpCtx)

	oldUser, err := userRepoInstance.GetOneByID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Tài khoản hoặc mật khẩu không đúng. Vui lòng thử lại %s", err))
	}
	if !oldUser.IsExists() {
		return fmt.Errorf("Tài khoản hoặc mật khẩu không đúng. Vui lòng thử lại")
	}
	// TODO : Check password match
	if !checkPasswordMatch(request.OldPassword, oldUser.Password) {
		// Increase num login fail
		ttl := 60
		setNumChangePassFail(userID.Hex(), numChangePassFail, ttl)
		return fmt.Errorf("Tài khoản hoặc mật khẩu không đúng. Vui lòng thử lại")
	}
	// Cập nhật password
	// TODO: hash pasword md5
	hashPash := generateHash(request.NewPassword)
	//
	now := time.Now()
	oldUser.Password = hashPash
	oldUser.LastChangePasswordDate = &now
	err = userRepoInstance.Update(oldUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật thông tin người dùng %s", err))
	}

	return c.JSON(success(nil))
}

// ResetPassword : Thay đổi password
// . Required OTP
func (UserHandler) ResetPassword(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		Phone           string `json:"phone" query:"phone" validate:"required"`
		OTP             string `json:"otp" query:"otp" validate:"required"`
		NewPassword     string `json:"new_password" query:"new_password" validate:"required"`
		ComfirmPassword string `json:"confirm_password" query:"confirm_password" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// TODO: validate
	{
		// validate phone
		request.Phone, err = utils.PhoneValidate(request.Phone)
		if err != nil {
			return fmt.Errorf("Số điện thoại không hợp lệ %s", err)
		}
	}
	{
		// validate NewPassword, ComfirmPassword
		if request.NewPassword != request.ComfirmPassword {
			return fmt.Errorf("Mật khẩu xác nhận không đúng")
		}
	}
	// TODO: Verify OTP
	typeOTP := "reset_password"
	isOK, err := verifyOTP(typeOTP, request.Phone, request.OTP)
	if err != nil || !isOK {
		return fmt.Errorf("OTP đã hết hạn hoặc không hợp lệ")

	}

	// Lấy thông tin user by phone
	userRepoInstance := userRepo.New(httpCtx)
	tokenRepoInstance := tokenRepo.New(httpCtx)
	oldUser, err := userRepoInstance.GetOneByPhone(request.Phone)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Tài khoản hoặc mật khẩu không đúng. Vui lòng thử lại %s", err))
	}
	if !oldUser.IsExists() {
		return fmt.Errorf("Tài khoản hoặc mật khẩu không đúng. Vui lòng thử lại")
	}
	// Cập nhật password
	bscryptNewPassword, err := bscryptPassword(request.NewPassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Tạo mật khẩu lỗi %s", err))
	}
	//
	now := time.Now()
	oldUser.Password = bscryptNewPassword
	oldUser.LastChangePasswordDate = &now
	err = userRepoInstance.Update(oldUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật thông tin người dùng %s", err))
	}
	// clear otp
	clearOTP(typeOTP, request.Phone)

	// clear all token with user
	_, err = tokenRepoInstance.ClearAll(oldUser.ID, "user")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Xóa token người dùng %s", err))
	}

	return c.JSON(success(nil))
}

// SendOTPResetPassword : Gửi mã OTP đến sdt cần reset password
// Mổi sdt giới hạn gửi OTP 5 lần trong 1 ngày
func (UserHandler) SendOTPResetPassword(c echo.Context) (err error) {
	httpCtx := c.Request().Context()

	type myResponse struct {
		TTL int `json:"TTL"`
	}
	type myRequest struct {
		Phone string `json:"phone" query:"phone" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// TODO: validate
	{
		// validate phone
		request.Phone, err = utils.PhoneValidate(request.Phone)
		if err != nil {
			return fmt.Errorf("Số điện thoại không hợp lệ %s", err)
		}
	}

	// check số điện thoại có tồn tại trong db
	{
		userRepoInstance := userRepo.New(httpCtx)
		user, err := userRepoInstance.GetActiveByPhone(request.Phone)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("get user %s", err))
		}
		if !user.IsExists() {
			return echo.NewHTTPError(http.StatusLocked, "Số điện thoại không tồn tại. Vui lòng nhập lại số điện thoại")
		}
	}

	// TODO: Đếm số lần send otp thành công đến sdt
	numSendOTP, errCount := countNumSendOTPResetPassword(request.Phone)
	if errCount != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Đếm số lần gửi mã OTP đến số %s %s", request.Phone, err))
	}
	numSendOTP++

	if numSendOTP > 3 {
		return echo.NewHTTPError(http.StatusLocked, "Bạn đã sử dụng hết lươt reset password hôm nay. Vui lòng thử lại sau 24h.")
	}

	typeOTP := "reset_password"
	OTP, ttl, err := genOTP(typeOTP, request.Phone)
	if ttl > 0 {
		// Đã gữi mã cho khách
		return c.JSON(success(myResponse{ttl}))
	} else if OTP != "" {
		// Gửi mã mới cho khách. Gọi api qua eSMS để gửi OTP đến sdt khách
		fmt.Println("OTP: ", OTP)

		reqESMS := esmsAPI.ESMS{
			Phone:   request.Phone,
			Content: OTP,
		}

		err := esmsAPI.SendOTPEsms(getRequestID(c), reqESMS)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Gửi mã OTP đến số %s %s", request.Phone, err))
		}
		// Cập nhật số lần gửi mã thành công đến khách hàng, lưu 1 ngày
		setNumSendOTPResetPassword(request.Phone, numSendOTP, 24*60*60)

		return c.JSON(success(myResponse{600}))
	}

	return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Tạo mã OTP reset password %s", err))
}

func (UserHandler) genUserToken() (result string, err error) {

	codes, err := codeRepo.New(context.Background()).GenerateUnitNum("autofarmer_user_token", 1)
	if err != nil || len(codes) == 0 {
		err = fmt.Errorf("Tạo token user %s", err)
	}
	result = codes[0]
	return
}
