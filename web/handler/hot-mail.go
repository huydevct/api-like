package handler

import (
	"app/constants"
	"app/model"
	hotmailRepo "app/repo/mongo/hotmail"
	userRepo "app/repo/mongo/user"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HotMailHandler init
type HotMailHandler struct{}

// NewHotMailHandler : Tạo đối tượng hotmail
func NewHotMailHandler() *HotMailHandler {
	return &HotMailHandler{}
}

// ImportHotMail func
func (HotMailHandler) ImportHotMail(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	type myRequest struct {
		Data []model.HotMail `json:"data" query:"data" validate:"required"`
	}

	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}

	if err = c.Validate(request); err != nil {
		return
	}

	if len(request.Data) < 1 {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Hotmail không được để trống"))
	}
	existedEmail := map[string]bool{}
	dataHotMail := make([]model.HotMail, 0)
	// Tạo email
	for _, hotmail := range request.Data {
		if !existedEmail[hotmail.Email] {
			// add vào mãng để check trùng
			existedEmail[hotmail.Email] = true
			hotmail.Status = constants.HotMailLive
			dataHotMail = append(dataHotMail, hotmail)
		}
	}
	hotmailRepo.New(httpCtx).InsertMany(dataHotMail)
	return c.JSON(success(nil))
}

// GetHotMail func
func (HotMailHandler) GetHotMail(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
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
	checkToken, err := userRepo.New(httpCtx).GetOneByToken(request.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng %s", err))
	}
	if !checkToken.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}
	hotmail, err := hotmailRepo.New(httpCtx).FindOneAndUpdate(request.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin hot mail %s", err))
	}
	if !hotmail.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy hotmail"))
	}
	return c.JSON(success(hotmail))
}

// All : Lấy toàn bộ hotmail
func (HotMailHandler) All(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                 `json:"total" query:"total"`
		Data       []model.HotMail     `json:"data" query:"data"`
	}
	type myRequest struct {
		Token  string                   `json:"token" query:"token"`
		Status []constants.CommonStatus `json:"status" query:"status"`
		Offset string                   `json:"offset" query:"offset"`
		Limit  int                      `json:"limit" query:"limit"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	if c.Get("user_token") != nil {
		request.Token = c.Get("user_token").(string)
	}
	// Tạo all HotMail request
	allHotMailReq := model.AllHotMailReq{
		Token:  request.Token,
		Offset: primitive.NilObjectID,
		Limit:  request.Limit,
		Status: request.Status,
	}
	if request.Offset != "" {
		offsetObjectID, err := primitive.ObjectIDFromHex(request.Offset)
		if err != nil {
			return fmt.Errorf("Offset không hợp lệ %s", err)
		}

		allHotMailReq.Offset = offsetObjectID
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.HotMail, 0),
	}

	hotmail, err := hotmailRepo.New(httpCtx).All(allHotMailReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách hotmail : %s", err))
	}

	if len(hotmail) > 0 {
		response.Data = hotmail
		response.Total = len(hotmail)
		response.LastOffset = &hotmail[len(hotmail)-1].ID
	}

	return c.JSON(success(response))
}
