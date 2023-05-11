package handler

import (
	"app/constants"
	"app/model"
	"fmt"
	"net/http"

	walletLogRepo "app/repo/mongo/walletlog"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WalletLogHandler : struct init service
type WalletLogHandler struct{}

// NewWalletLogHandler : Tạo mới 1 đối tượng wallet log handler
func NewWalletLogHandler() *WalletLogHandler {
	return &WalletLogHandler{}
}

// All : Lấy toàn bộ các wallet log
func (h WalletLogHandler) All(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                 `json:"total" query:"total"`
		Data       []model.WalletLog   `json:"data" query:"data"`
	}
	type myRequest struct {
		UserToken string                    `json:"user_token" query:"user_token"`
		Type      []constants.WalletLogType `json:"type" query:"type"`
		Offset    string                    `json:"offset" query:"offset"`
		Limit     int                       `json:"limit" query:"limit"`
		FromTime  int                       `json:"from_time" query:"from_time"`
		ToTime    int                       `json:"to_time" query:"to_time"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// validate
	var offsetObjectID primitive.ObjectID
	if request.Offset != "" {
		offsetObjectID, err = primitive.ObjectIDFromHex(request.Offset)
		if err != nil {
			return fmt.Errorf("Offset không hợp lệ %s", err)
		}
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.WalletLog, 0),
	}

	// Tạo search request
	allWalletLogReq := model.AllWalletLogReq{
		Types:     request.Type,
		UserToken: request.UserToken,
		Offset:    offsetObjectID,
		Limit:     request.Limit,
		FromTime:  request.FromTime,
		ToTime:    request.ToTime,
	}

	if c.Get("user_token") != nil {
		allWalletLogReq.UserToken = c.Get("user_token").(string)
	}

	results, err := walletLogRepo.New(httpCtx).All(allWalletLogReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách dịch vụ %s", err))
	}

	if len(results) > 0 {
		response.Data = results
		response.Total = len(results)
		response.LastOffset = &results[len(results)-1].ID
	}

	return c.JSON(success(response))
}
