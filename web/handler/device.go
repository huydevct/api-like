package handler

import (
	"app/constants"
	"app/model"
	cloneRepo "app/repo/mongo/clone"
	deviceRepo "app/repo/mongo/device"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DeviceHandler init
type DeviceHandler struct{}

// NewDeviceHandler : Tạo đối tượng device handler
func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{}
}

// All : Lấy toàn bộ Device
func (DeviceHandler) All(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                 `json:"total" query:"total"`
		Data       []model.DeviceInfo  `json:"data" query:"data"`
	}
	type myRequest struct {
		Token  string                   `json:"token" query:"token"`
		Status []constants.CommonStatus `json:"status" query:"status"`
		Name   string                   `json:"name" query:"name"`
		Offset string                   `json:"offset" query:"offset"`
		Limit  int                      `json:"limit" query:"limit"`
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
	// Tạo all Device request
	allDeviceeReq := model.AllDeviceReq{
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

		allDeviceeReq.Offset = offsetObjectID
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.DeviceInfo, 0),
	}

	devices, err := deviceRepo.New(httpCtx).All(allDeviceeReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách Device : %s", err))
	}
	fmt.Println(allDeviceeReq)
	if len(devices) > 0 {
		response.Data = devices
		response.Total = len(devices)
		response.LastOffset = &devices[len(devices)-1].ID
	}

	return c.JSON(success(response))
}

// Delete : Xóa Device
func (DeviceHandler) Delete(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		DeviceID string `json:"device_id" query:"device_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	deviceObjectID, err := primitive.ObjectIDFromHex(request.DeviceID)
	if err != nil {
		return fmt.Errorf("Mã device không hợp lệ %s", err)
	}
	deviceRepoInstance := deviceRepo.New(httpCtx)
	//

	// TODO: Do many things with mongo transaction
	session, err := deviceRepoInstance.Session.ConClient.StartSession()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Get mongo session: %s", err))
	}
	// Start transaction
	err = session.StartTransaction()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Start mongo transaction: %s", err))
	}
	defer session.EndSession(httpCtx)
	device, err := deviceRepoInstance.GetOneByID(deviceObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin device %s", err))
	}
	if !device.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin device"))
	}
	// validate người thực thi: user id
	if c.Get("user_token") != nil {
		if device.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Token không hợp lệ")
		}
	}
	err = mongo.WithSession(httpCtx, session, func(sessCtx mongo.SessionContext) (err error) {
		// Xóa.
		device.Status = constants.Delete

		err = deviceRepoInstance.UpdateWithSessionCtx(sessCtx, device)
		if err != nil {
			session.AbortTransaction(sessCtx)
			return fmt.Errorf("Cập nhật device: %s", err)
		}
		err = cloneRepo.New(httpCtx).UpdateManyLiveWithSessionCtx(sessCtx, device.Token, device.PCName)

		if err != nil {
			session.AbortTransaction(sessCtx)
			return fmt.Errorf("Cập nhật clone: %s", err)
		}

		session.CommitTransaction(sessCtx)
		return
	})

	return c.JSON(success(device))
}

// Detail : LẤy detail by ID
// Approve : duyệt chuyển trạng thái từ Active -> approved
func (DeviceHandler) Detail(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	type myRequest struct {
		DeviceID string `json:"device_id" query:"device_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// validate
	deviceObjectID, err := primitive.ObjectIDFromHex(request.DeviceID)
	if err != nil {
		return fmt.Errorf("Mã device không hợp lệ %s", err)
	}
	deviceRepoInstance := deviceRepo.New(httpCtx)

	// Kiểm tra DeviceID có tồn tại hay không ?
	device, err := deviceRepoInstance.GetOneByID(deviceObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin device %s", err))
	}
	if !device.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin device"))
	}
	// Kiểm tra device này có thuộc về token hay hông?
	// validate người thực thi: user token
	if c.Get("user_token") != nil {
		if device.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Device không thuộc về bạn. Vui lòng thử lại")
		}
	}

	return c.JSON(success(device))
}

// Search : support sort
func (h DeviceHandler) Search(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		Total int                `json:"total" query:"total"`
		Data  []model.DeviceInfo `json:"data" query:"data"`
	}
	type myRequest struct {
		Status []constants.CommonStatus `json:"status" query:"status"`
		Page   int                      `json:"page" query:"page"`
		Limit  int                      `json:"limit" query:"limit"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.DeviceInfo, 0),
	}

	// Tạo search request
	searchDeviceReq := model.SearchDevice{
		Page:   request.Page,
		Limit:  request.Limit,
		Status: request.Status,
	}

	if c.Get("user_token") != nil {
		searchDeviceReq.Token = c.Get("user_token").(string)
	}

	deviceRepoInstance := deviceRepo.New(httpCtx)

	devices, err := deviceRepoInstance.Search(searchDeviceReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Search danh sách device %s", err))
	}

	// TODO: update total device
	response.Total, err = deviceRepoInstance.Total(searchDeviceReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Đếm số lượng device %s", err))
	}

	if len(devices) > 0 {
		response.Data = devices
	}

	return c.JSON(success(response))
}

// TotalAllClone : tổng số clone của mỗi device
func (h DeviceHandler) TotalAllClone(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		Total int `json:"total" query:"total"`
	}

	type myRequest struct {
		DeviceID string `json:"device_id" query:"device_id" validate:"required"`
	}

	request := new(myRequest)
	response := myResponse{}

	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}

	allCloneReq := model.AllDeviceReq{
		DeviceID: request.DeviceID,
	}

	total, err := cloneRepo.New(httpCtx).AllByDevice(allCloneReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy tổng clone của device %s", err))
	}

	response.Total = total

	return c.JSON(success(response))
}

// All : Lấy toàn bộ Device cho admin
func (DeviceHandler) AllDevice(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                 `json:"total" query:"total"`
		Data       []model.DeviceInfo  `json:"data" query:"data"`
	}
	type myRequest struct {
		Token    string                   `json:"token" query:"token"`
		DeviceID string                   `json:"device_id" query:"device_id"`
		Status   []constants.CommonStatus `json:"status" query:"status"`
		Page     int                      `json:"page" query:"page"`
		Limit    int                      `json:"limit" query:"limit"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Tạo all Device request
	allDeviceeReq := model.AllDeviceReq{
		Token:    request.Token,
		DeviceID: request.DeviceID,
		Page:     request.Page,
		Limit:    request.Limit,
		Status:   request.Status,
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.DeviceInfo, 0),
	}

	devices, err := deviceRepo.New(httpCtx).All(allDeviceeReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách Device : %s", err))
	}
	response.Total, err = deviceRepo.New(httpCtx).TotalForAdmin(allDeviceeReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách Device : %s", err))
	}
	if len(devices) > 0 {
		response.Data = devices
	}

	return c.JSON(success(response))
}
