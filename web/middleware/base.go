package middleware

import (
	"app/common/config"
	ghandler "app/common/gstuff/handler"
	"app/constants"
	"context"
	"fmt"

	tokenRepo "app/repo/mongo/token"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var cfg = config.GetConfig()

func getRequestID(c echo.Context) (requestID string) {
	return ghandler.GetRequestID(c)
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
