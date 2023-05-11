package handler

import (
	"app/model"
	settingPriceRepo "app/repo/mongo/settingprice"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//SettingPriceHandler init
type SettingPriceHandler struct{}

//NewSettingPriceHandler : Tạo đối tượng setting price
func NewSettingPriceHandler() *SettingPriceHandler {
	return &SettingPriceHandler{}
}

//Create func : Tạo mới setting price
func (SettingPriceHandler) Create(c echo.Context) (err error) {

	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		Appname  string                   `json:"appname" query:"appname" validate:"required"`
		Settings []model.SettingPriceItem `json:"settings" query:"settings" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}

	if err = c.Validate(request); err != nil {
		return
	}

	// kiểm tra tồn tại setting price
	settingPrices, err := settingPriceRepo.New(httpCtx).GetOneByAppname(request.Appname)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Kiểm tra thông tin setting price %s", err))
	}
	if settingPrices.IsExists() {
		return fmt.Errorf("Thông tin này đã tồn tại với appname : " + request.Appname)
	}
	settingPrice := model.SettingPrice{
		Settings: request.Settings,
		Appname:  request.Appname,
	}
	settingPrice, err = settingPriceRepo.New(httpCtx).Insert(settingPrice)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Thêm thông tin setting price %s", err))
	}

	return c.JSON(success(settingPrice))
}

//Update func : Cập nhật setting price
func (SettingPriceHandler) Update(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		SettingPriceID string                   `json:"setting_price_id" query:"setting_price_id" validate:"required"`
		Settings       []model.SettingPriceItem `json:"settings" query:"settings" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}

	if err = c.Validate(request); err != nil {
		return
	}
	settingPriceObjectID, err := primitive.ObjectIDFromHex(request.SettingPriceID)
	if err != nil {
		return fmt.Errorf("Mã setting price hợp lệ %s", err)
	}
	// kiểm tra tồn tại setting price
	settingPrice, err := settingPriceRepo.New(httpCtx).GetOneByID(settingPriceObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Kiểm tra thông tin setting price %s", err))
	}
	if !settingPrice.IsExists() {
		return fmt.Errorf("Thông tin này không tồn tại")
	}
	settingPrice.Settings = request.Settings
	err = settingPriceRepo.New(httpCtx).Update(settingPrice)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Thêm thông tin setting price %s", err))
	}

	return c.JSON(success(settingPrice))
}

//Detail func : Lấy thông tin setting price
func (SettingPriceHandler) Detail(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		SettingPriceID string `json:"setting_price_id" query:"setting_price_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}

	if err = c.Validate(request); err != nil {
		return
	}
	settingPriceObjectID, err := primitive.ObjectIDFromHex(request.SettingPriceID)
	if err != nil {
		return fmt.Errorf("Mã setting price hợp lệ %s", err)
	}
	// kiểm tra tồn tại setting price
	settingPrice, err := settingPriceRepo.New(httpCtx).GetOneByID(settingPriceObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Kiểm tra thông tin setting price %s", err))
	}

	return c.JSON(success(settingPrice))
}

// All : Lấy toàn bộ action
func (SettingPriceHandler) All(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID  `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                  `json:"total" query:"total"`
		Data       []model.SettingPrice `json:"data" query:"data"`
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.SettingPrice, 0),
	}

	settings, err := settingPriceRepo.New(httpCtx).All()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách setting : %s", err))
	}

	if len(settings) > 0 {
		response.Data = settings
		response.Total = len(settings)
		response.LastOffset = &settings[len(settings)-1].ID
	}

	return c.JSON(success(response))
}
