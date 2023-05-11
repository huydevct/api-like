package handler

import (
	"app/model"
	settingRepo "app/repo/mongo/setting"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

//SettingHandler init
type SettingHandler struct{}

//NewSettingHandler : Tạo đối tượng setting
func NewSettingHandler() *SettingHandler {
	return &SettingHandler{}
}

//CreateUpdate func : Tạo mới setting price
func (SettingHandler) CreateUpdate(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		Data model.Setting `json:"data" query:"data" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}

	if err = c.Validate(request); err != nil {
		return
	}

	// kiểm tra tồn tại setting price
	settingItem, err := settingRepo.New(httpCtx).GetOneSetting()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Kiểm tra thông tin setting %s", err))
	}
	setting := request.Data
	if !settingItem.IsExists() {
		setting, err = settingRepo.New(httpCtx).Insert(setting)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Thêm thông tin setting %s", err))
		}
	} else {
		setting.ID = settingItem.ID
		err := settingRepo.New(httpCtx).Update(setting)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật thông tin setting %s", err))
		}

	}

	return c.JSON(success(setting))
}

//Detail func : Lấy thông tin setting
func (SettingHandler) Detail(c echo.Context) (err error) {
	httpCtx := c.Request().Context()

	// Lấy thông tin setting
	setting, err := settingRepo.New(httpCtx).GetOneSetting()
	fmt.Println(setting)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Kiểm tra thông tin setting price %s", err))
	}
	return c.JSON(success(setting))
}
