package handler

import (
	"app/constants"
	"app/model"
	userRepo "app/repo/mongo/user"
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"net/url"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QrContent chưa thông tin QR
type QrContent struct {
	TokenLogin string `json:"token_login"`
	DomainAPI  string `json:"domain_api"`
}

// QrCodeHandler : struct init qrcode
type QrCodeHandler struct{}

// NewQrCodeHandler : Tạo mới 1 đối tượng bình luận handler
func NewQrCodeHandler() *QrCodeHandler {
	return &QrCodeHandler{}
}

// Gen : Trả về mã qr code chứa token tạm
func (QrCodeHandler) Gen(c echo.Context) (err error) {
	httpCtx := c.Request().Context()

	type myRequest struct {
		Size      int                `json:"size" query:"size" validate:"required"`
		UserToken string             `json:"user_token" query:"user_token" validate:"required"`
		UserID    primitive.ObjectID `json:"user_id" query:"user_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if c.Get("user_token") != nil {
		request.UserToken = c.Get("user_token").(string)
		request.UserID = c.Get("user_id").(primitive.ObjectID)
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Lấy thông tin user theo userID
	userRepoInstance := userRepo.New(httpCtx)
	user, err := userRepoInstance.GetOneActiveByID(request.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng %s", err))
	}
	if !user.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}

	// Lưu thông tin lên redis
	// gen token
	tokenQR, err := uuid.NewUUID()
	if err != nil {
		err = fmt.Errorf("Gen token by UUID fail: %s", err)
		return
	}

	// Lưu thông tin lên redis với token tokenQR
	{
		redisClient := cfg.Redis["core"].GetClient()

		loginQRKey := fmt.Sprintf(constants.QRLoginKey, tokenQR)
		// Tạo data chứa thông tin login của user
		dataByte, _ := json.Marshal(user)

		// Lưu thông tin trong 3 phút
		err = redisClient.Set(loginQRKey, dataByte, time.Minute*3).Err()
		if err != nil {
			err = fmt.Errorf("Save qr code to redis %s", err)
			return
		}
	}

	// genQR code từ tokenQR
	s, err := url.QueryUnescape(tokenQR.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	code, err := qr.Encode(s, qr.L, qr.Auto)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if request.Size < 100 || request.Size > 1000 {
		request.Size = 250
	}

	// Scale the barcode to the appropriate size
	code, err = barcode.Scale(code, request.Size, request.Size)
	if err != nil {
	}

	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	response := c.Response()
	response.Header().Set("Cache-Control", "public, max-age=2592000")
	response.Header().Set("Last-Modified", time.Now().AddDate(0, 0, 1).Format(http.TimeFormat))
	response.Header().Set("Expires", time.Now().Format(http.TimeFormat))

	response.Write(buffer.Bytes())
	return
}

// ParseQRCode : Lấy nội dung qr by token temp
func (QrCodeHandler) ParseQRCode(c echo.Context) (err error) {

	type myRequest struct {
		Token string `json:"token" query:"token" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}

	loginQRKey := fmt.Sprintf(constants.QRLoginKey, request.Token)

	redisClient := cfg.Redis.Get("core").GetClient()
	value, errRedis := redisClient.Get(loginQRKey).Result()
	if errRedis != nil {
		if errRedis != redis.Nil {
			err = fmt.Errorf("Connect redis fail %s", errRedis)
			return
		} else {
			err = fmt.Errorf("QR code not valid")
			return
		}
	}

	redisClient.Del(loginQRKey)

	qrContent := model.UserInfo{}
	// Kiểm tra thông tin từ redis dùng được không ?
	err = json.Unmarshal([]byte(value), &qrContent)
	if err != nil {
		err = fmt.Errorf("Parse data redis to QrContent %s", err)
		return
	}

	return c.JSON(success(qrContent))
}

// ParseQRCodeLogin : Parse QR code lấy token user login
func (QrCodeHandler) ParseQRCodeLogin(c echo.Context) (err error) {

	type myRequest struct {
		Token string `json:"token" query:"token" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}

	loginQRKey := fmt.Sprintf(constants.QRLoginKey, request.Token)

	redisClient := cfg.Redis.Get("core").GetClient()
	value, errRedis := redisClient.Get(loginQRKey).Result()
	if errRedis != nil {
		if errRedis != redis.Nil {
			err = fmt.Errorf("Connect redis fail %s", errRedis)
			return
		} else {
			err = fmt.Errorf("QR code not valid")
			return
		}
	}

	redisClient.Del(loginQRKey)

	qrContent := model.UserInfo{}
	// Kiểm tra thông tin từ redis dùng được không ?
	err = json.Unmarshal([]byte(value), &qrContent)
	if err != nil {
		err = fmt.Errorf("Parse data redis to QrContent %s", err)
		return
	}

	// TẠo token đăng nhập trả về cho user
	// Tạo token đăng nhập mới cho user này
	tokenModel := model.Token{
		ExpiredAfterSecond: 7 * 24 * 3600,
		// token info
		UserID:    qrContent.ID,
		UserToken: qrContent.Token,
		Type:      "user",
		DeviceID:  time.Now().String(),
		UserAgent: c.Request().UserAgent(),
		RemoteIP:  c.RealIP(),
	}

	tokenLogin, err := GenTokenByUserIDDeviceID(tokenModel)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(success(map[string]interface{}{
		"token": tokenLogin,
	}))
}
