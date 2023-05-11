package handler

import (
	"app/constants"
	"app/model"
	ActionProfileTemplateRepo "app/repo/mongo/actionprofile"
	employeeRepo "app/repo/mongo/employee"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//ActionProfileTemplateHandler init
type ActionProfileTemplateHandler struct{}

//NewActionProfileTemplateHandler : Tạo đối tượng action profile
func NewActionProfileTemplateHandler() *ActionProfileTemplateHandler {
	return &ActionProfileTemplateHandler{}
}

//Create func : Tạo mới action profile
func (ActionProfileTemplateHandler) Create(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		EmployeeID    string                   `json:"employee_id" query:"employee_id" validate:"required"`
		Name          string                   `json:"name" query:"name" validate:"required"`
		AppName       string                   `json:"appname" query:"appname" validate:"required"`
		Template      constants.TemplateAction `json:"template" query:"template" validate:"required"`
		Actions       [][]model.Action         `json:"actions" query:"actions"`
		ActionDefault []model.Action           `json:"action_default" query:"action_default" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if c.Get("employee_id_str") != nil {
		request.EmployeeID = c.Get("employee_id_str").(string)
	}
	if err = c.Validate(request); err != nil {
		return
	}
	employeeObjectID, err := primitive.ObjectIDFromHex(request.EmployeeID)
	if err != nil {
		return fmt.Errorf("Mã số nhân viên không hợp lệ %s", err)
	}
	// Lấy thông tin tài khoản
	employee, err := employeeRepo.New(httpCtx).GetOneActiveByID(employeeObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng %s", err))
	}
	if !employee.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}
	if len(request.ActionDefault) <= 0 {
		return fmt.Errorf("aciton profile không hợp")
	}
	ActionProfileTemplate := model.ActionProfile{
		Name:          request.Name,
		Token:         employee.Token,
		AppName:       request.AppName,
		Template:      constants.TemplateActionAdmin,
		Actions:       request.Actions,
		ActionDefault: request.ActionDefault,
		Status:        constants.Active,
	}

	ActionProfileTemplate, err = ActionProfileTemplateRepo.New(httpCtx).Insert(ActionProfileTemplate)
	if err != nil {
		return fmt.Errorf("Tạo action profile: %s", err)
	}

	return c.JSON(success(nil))
}

// Delete : Xóa action profile
func (ActionProfileTemplateHandler) Delete(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		ActionProfileTemplateID string `json:"action_profile_id" query:"action_profile_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	ActionProfileTemplateObjectID, err := primitive.ObjectIDFromHex(request.ActionProfileTemplateID)
	if err != nil {
		return fmt.Errorf("aciton profile không hợp lệ %s", err)
	}
	//
	ActionProfileTemplate, err := ActionProfileTemplateRepo.New(httpCtx).GetOneByID(ActionProfileTemplateObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin action profile %s", err))
	}
	if !ActionProfileTemplate.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin action profile"))
	}
	// validate người thực thi: employee id
	if c.Get("employee_token") != nil {
		if ActionProfileTemplate.Token != c.Get("employee_token").(string) {
			return fmt.Errorf("Mã số người dùng không hợp lệ")
		}
	}
	// Xóa.
	ActionProfileTemplate.Status = constants.Delete

	err = ActionProfileTemplateRepo.New(httpCtx).Update(ActionProfileTemplate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật action profile: %s", err))
	}

	return c.JSON(success(ActionProfileTemplate))
}

// All : Lấy toàn bộ action
func (ActionProfileTemplateHandler) All(c echo.Context) (err error) {
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
	if c.Get("employee_token") != nil {
		request.Token = c.Get("employee_token").(string)
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Tạo all action request
	allActionProfileTemplateReq := model.AllActionProfileReq{
		Name:     request.Name,
		AppName:  request.AppName,
		Template: request.Template,
		Offset:   primitive.NilObjectID,
		Limit:    request.Limit,
		Status:   request.Status,
	}
	if request.Offset != "" {
		offsetObjectID, err := primitive.ObjectIDFromHex(request.Offset)
		if err != nil {
			return fmt.Errorf("Offset không hợp lệ %s", err)
		}

		allActionProfileTemplateReq.Offset = offsetObjectID
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.ActionProfile, 0),
	}

	actions, err := ActionProfileTemplateRepo.New(httpCtx).All(allActionProfileTemplateReq)
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
func (ActionProfileTemplateHandler) Update(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		ActionProfileTemplateID string           `json:"action_profile_id" query:"action_profile_id" validate:"required"`
		Actions                 [][]model.Action `json:"actions" query:"actions"`
		ActionDefault           []model.Action   `json:"action_default" query:"action_default"  validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	ActionProfileTemplateObjectID, err := primitive.ObjectIDFromHex(request.ActionProfileTemplateID)
	if err != nil {
		return fmt.Errorf("aciton profile không hợp lệ %s", err)
	}
	if len(request.ActionDefault) <= 0 {
		return fmt.Errorf("aciton profile không hợp")
	}
	//
	ActionProfileTemplate, err := ActionProfileTemplateRepo.New(httpCtx).GetOneByID(ActionProfileTemplateObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin action profile %s", err))
	}

	if !ActionProfileTemplate.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin action profile"))
	}
	// validate người thực thi: user id
	if c.Get("employee_token") != nil {
		if ActionProfileTemplate.Token != c.Get("employee_token").(string) {
			return fmt.Errorf("Mã số người dùng không hợp lệ")
		}
	}
	// update
	ActionProfileTemplate.Actions = request.Actions
	ActionProfileTemplate.ActionDefault = request.ActionDefault
	err = ActionProfileTemplateRepo.New(httpCtx).Update(ActionProfileTemplate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật action profile: %s", err))
	}

	return c.JSON(success(ActionProfileTemplate))
}

//Detail : Thông tin chi tiết action profile
func (ActionProfileTemplateHandler) Detail(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		ActionProfileTemplateID string `json:"action_profile_id" query:"action_profile_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	ActionProfileTemplateObjectID, err := primitive.ObjectIDFromHex(request.ActionProfileTemplateID)
	if err != nil {
		return fmt.Errorf("aciton profile không hợp lệ %s", err)
	}
	//
	ActionProfileTemplate, err := ActionProfileTemplateRepo.New(httpCtx).GetOneByID(ActionProfileTemplateObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin action profile %s", err))
	}

	if !ActionProfileTemplate.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin action profile"))
	}
	// validate người thực thi: user id
	if c.Get("employee_token") != nil {
		if ActionProfileTemplate.Token != c.Get("employee_token").(string) {
			return fmt.Errorf("Mã số người dùng không hợp lệ")
		}
	}

	err = ActionProfileTemplateRepo.New(httpCtx).Update(ActionProfileTemplate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật action profile: %s", err))
	}

	return c.JSON(success(ActionProfileTemplate))
}
