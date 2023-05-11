package handler

import (
	"app/constants"
	"app/model"
	actionProfileRepo "app/repo/mongo/actionprofile"
	userRepo "app/repo/mongo/user"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActionProfileHandler init
type ActionProfileHandler struct{}

// NewActionProfileHandler : Tạo đối tượng action profile
func NewActionProfileHandler() *ActionProfileHandler {
	return &ActionProfileHandler{}
}

// Create func : Tạo mới action profile
func (ActionProfileHandler) Create(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		UserID        string                   `json:"user_id" query:"user_id" validate:"required"`
		Name          string                   `json:"name" query:"name" validate:"required"`
		AppName       string                   `json:"appname" query:"appname"`
		Template      constants.TemplateAction `json:"template" query:"template" validate:"required"`
		Actions       [][]model.Action         `json:"actions" query:"actions"`
		ActionDefault []model.Action           `json:"action_default" query:"action_default" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if c.Get("user_id_str") != nil {
		request.UserID = c.Get("user_id_str").(string)
	}
	if err = c.Validate(request); err != nil {
		return
	}
	userObjectID, err := primitive.ObjectIDFromHex(request.UserID)
	if err != nil {
		return fmt.Errorf("Mã số nhân viên không hợp lệ %s", err)
	}
	// Lấy thông tin tài khoản
	user, err := userRepo.New(httpCtx).GetOneActiveByID(userObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng %s", err))
	}
	if !user.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}
	if len(request.ActionDefault) <= 0 {
		return fmt.Errorf("action profile không hợp")
	}
	actionProfile := model.ActionProfile{
		Name:          request.Name,
		Token:         user.Token,
		AppName:       request.AppName,
		Template:      constants.TemplateActionUser,
		Actions:       request.Actions,
		ActionDefault: request.ActionDefault,
		Status:        constants.Active,
	}

	actionProfile, err = actionProfileRepo.New(httpCtx).Insert(actionProfile)
	if err != nil {
		return fmt.Errorf("Tạo action profile: %s", err)
	}

	return c.JSON(success(nil))
}

// Delete : Xóa action profile
func (ActionProfileHandler) Delete(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		ActionProfileID string `json:"action_profile_id" query:"action_profile_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	actionProfileObjectID, err := primitive.ObjectIDFromHex(request.ActionProfileID)
	if err != nil {
		return fmt.Errorf("aciton profile không hợp lệ %s", err)
	}
	//
	actionProfile, err := actionProfileRepo.New(httpCtx).GetOneByID(actionProfileObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin action profile %s", err))
	}
	if !actionProfile.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin action profile"))
	}
	// validate người thực thi: user id
	if c.Get("user_token") != nil {
		if actionProfile.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Mã số người dùng không hợp lệ")
		}
	}
	// Xóa.
	actionProfile.Status = constants.Delete

	err = actionProfileRepo.New(httpCtx).Update(actionProfile)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật action profile: %s", err))
	}

	return c.JSON(success(actionProfile))
}

// All : Lấy toàn bộ action
func (ActionProfileHandler) All(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID   `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                   `json:"total" query:"total"`
		Data       []model.ActionProfile `json:"data" query:"data"`
	}
	type myRequest struct {
		Token    string                   `json:"token" query:"token"`
		Template constants.TemplateAction `json:"template" query:"template"`
		AppName  string                   `json:"appname" query:"appname"`
		Status   []constants.CommonStatus `json:"status" query:"status"`
		Name     string                   `json:"name" query:"name"`
		Offset   string                   `json:"offset" query:"offset"`
		Limit    int                      `json:"limit" query:"limit"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if c.Get("user_token") != nil {
		request.Token = c.Get("user_token").(string)
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Tạo all action request
	allActionProfileReq := model.AllActionProfileReq{
		Name:     request.Name,
		Token:    request.Token,
		Template: request.Template,
		AppName:  request.AppName,
		Offset:   primitive.NilObjectID,
		Limit:    request.Limit,
		Status:   request.Status,
	}
	if request.Offset != "" {
		offsetObjectID, err := primitive.ObjectIDFromHex(request.Offset)
		if err != nil {
			return fmt.Errorf("Offset không hợp lệ %s", err)
		}

		allActionProfileReq.Offset = offsetObjectID
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.ActionProfile, 0),
	}

	actions, err := actionProfileRepo.New(httpCtx).All(allActionProfileReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách actions : %s", err))
	}

	if len(actions) > 0 {
		response.Data = actions
		response.Total = len(actions)
		response.LastOffset = &actions[len(actions)-1].ID
	}

	return c.JSON(success(response))
}

// Update : update action profile
func (ActionProfileHandler) Update(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		ActionProfileID string           `json:"action_profile_id" query:"action_profile_id" validate:"required"`
		Actions         [][]model.Action `json:"actions" query:"actions"`
		ActionDefault   []model.Action   `json:"action_default" query:"action_default"  validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	actionProfileObjectID, err := primitive.ObjectIDFromHex(request.ActionProfileID)
	if err != nil {
		return fmt.Errorf("aciton profile không hợp lệ %s", err)
	}
	if len(request.ActionDefault) <= 0 {
		return fmt.Errorf("aciton profile không hợp")
	}
	//
	actionProfile, err := actionProfileRepo.New(httpCtx).GetOneByID(actionProfileObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin action profile %s", err))
	}

	if !actionProfile.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin action profile"))
	}
	// validate người thực thi: user id
	if c.Get("user_token") != nil {
		if actionProfile.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Mã số người dùng không hợp lệ")
		}
	}
	// update
	actionProfile.Actions = request.Actions
	actionProfile.ActionDefault = request.ActionDefault
	err = actionProfileRepo.New(httpCtx).Update(actionProfile)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật action profile: %s", err))
	}

	return c.JSON(success(actionProfile))
}

// Detail : Thông tin chi tiết action profile
func (ActionProfileHandler) Detail(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		ActionProfileID string `json:"action_profile_id" query:"action_profile_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	actionProfileObjectID, err := primitive.ObjectIDFromHex(request.ActionProfileID)
	if err != nil {
		return fmt.Errorf("aciton profile không hợp lệ %s", err)
	}
	//
	actionProfile, err := actionProfileRepo.New(httpCtx).GetOneByID(actionProfileObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin action profile %s", err))
	}

	if !actionProfile.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin action profile"))
	}
	// validate người thực thi: user id
	if c.Get("user_token") != nil && actionProfile.Template != 1 {
		if actionProfile.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Mã số người dùng không hợp lệ")
		}
	}

	return c.JSON(success(actionProfile))
}

// AllByClone : Lấy toàn bộ action cho clone
func (ActionProfileHandler) AllByClone(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID   `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                   `json:"total" query:"total"`
		Data       []model.ActionProfile `json:"data" query:"data"`
	}
	type myRequest struct {
		Token    string                   `json:"token" query:"token"`
		AppName  string                   `json:"appname" query:"appname"`
		Template constants.TemplateAction `json:"template" query:"template"`
		Status   []constants.CommonStatus `json:"status" query:"status"`
		Name     string                   `json:"name" query:"name"`
		Offset   string                   `json:"offset" query:"offset"`
		Limit    int                      `json:"limit" query:"limit"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if c.Get("user_token") != nil {
		request.Token = c.Get("user_token").(string)
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Tạo all action request
	allActionProfileReq := model.AllActionProfileReq{
		Name:     request.Name,
		Template: request.Template,
		Token:    request.Token,
		AppName:  request.AppName,
		Offset:   primitive.NilObjectID,
		Limit:    request.Limit,
		Status:   request.Status,
	}
	if request.Offset != "" {
		offsetObjectID, err := primitive.ObjectIDFromHex(request.Offset)
		if err != nil {
			return fmt.Errorf("Offset không hợp lệ %s", err)
		}

		allActionProfileReq.Offset = offsetObjectID
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.ActionProfile, 0),
	}

	actions, err := actionProfileRepo.New(httpCtx).AllByClone(allActionProfileReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách actions : %s", err))
	}

	if len(actions) > 0 {
		response.Data = actions
		response.Total = len(actions)
		response.LastOffset = &actions[len(actions)-1].ID
	}

	return c.JSON(success(response))
}
