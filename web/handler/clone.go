package handler

import (
	"context"
	"sync"

	"app/constants"
	"app/model"

	"fmt"
	"net/http"

	actionProfileRepo "app/repo/mongo/actionprofile"
	cloneRepo "app/repo/mongo/clone"
	userRepo "app/repo/mongo/user"

	"github.com/labstack/echo/v4"
	"github.com/panjf2000/ants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CloneHandler init
type CloneHandler struct{}

// NewCloneHandler : Tạo đối tượng Clone Handle
func NewCloneHandler() *CloneHandler {
	return &CloneHandler{}
}

//
var (
	// Danh sách app nam support
	SupportedCloneAppname = map[constants.AppName]bool{
		constants.AppNameFaceBook:  true,
		constants.AppNameInstagram: true,
		constants.AppNameTiktok:    true,
		constants.AppNameYoutube:   true,
	}
)

func (CloneHandler) AllClone(c echo.Context) (err error) {
	httpCtx := c.Request().Context()

	type myResponse struct {
		// CurrentPage int               `json:"currentPage" query:"currentPage"`
		LastOffset *primitive.ObjectID `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                 `json:"totalCount" query:"totalCount"`
		Data       []model.CloneInfo   `json:"data" query:"data"`
	}
	type myRequest struct {
		DeviceSystem string `json:"device_system" query:"device_system"`
		DeviceID     string `json:"device_id" query:"device_id"`
		Page         int    `json:"page" query:"page"`
		Limit        int    `json:"limit" query:"limit"`
		OrderColumn  string `json:"orderColumn" query:"orderColumn"`
		OrderType    string `json:"orderType" query:"orderType"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	reqObj := cloneRepo.SearchClonesByDeviceReq{
		DeviceSystem: request.DeviceSystem,
		Limit:        request.Limit,
		Page:         request.Page,
		OrderColumn:  request.OrderColumn,
		OrderType:    request.OrderType,
	}

	if request.DeviceID != "" {
		reqObj.DeviceID, _ = primitive.ObjectIDFromHex(request.DeviceID)
	}

	// tạo response
	response := myResponse{
		Data: make([]model.CloneInfo, 0),
	}

	CloneRepoInstance := cloneRepo.New(httpCtx)
	response.Data, err = CloneRepoInstance.AllCloneByDevice(&reqObj)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lỗi lấy danh sách Clone %s", err))
	}

	response.Total, err = CloneRepoInstance.TotalAllCloneByDevice(&reqObj)
	// response.CurrentPage = request.Page
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Số lượng %s", err))
	}

	return c.JSON(success(response))
}

// Delete : Xóa Clone
func (CloneHandler) Delete(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		CloneID primitive.ObjectID `json:"clone_id" query:"clone_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}

	cloneRepoInstance := cloneRepo.New(httpCtx)
	// Kiểm tra clone có tồn tại không ?
	clone, err := cloneRepoInstance.GetOneActiveByID(request.CloneID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy clone theo %s %s", request.CloneID, err))
	}
	if !clone.IsExists() {
		return fmt.Errorf("Không tìm thấy thông tin clone %s", request.CloneID)
	}
	// Kiểm tra clone có thuộc về user này hay không ?
	if c.Get("user_token") != nil {
		if clone.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Clone %s không thuôc về bạn", request.CloneID)
		}
	}

	// Xóa
	clone.AliveStatus = constants.CloneDelete

	err = cloneRepoInstance.Update(clone)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Xóa clone: %s", err))
	}
	return c.JSON(success(clone))
}

// Reset : Cập nhật trạng thái clone thành live, và clear DeviceID
func (h CloneHandler) Reset(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		CloneID primitive.ObjectID `json:"clone_id" query:"clone_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}

	cloneRepoInstance := cloneRepo.New(httpCtx)
	// Kiểm tra clone có tồn tại không ?
	clone, err := cloneRepoInstance.GetOneActiveByID(request.CloneID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy clone theo %s %s", request.CloneID, err))
	}
	if !clone.IsExists() {
		return fmt.Errorf("Không tìm thấy thông tin clone %s", request.CloneID)
	}
	// Kiểm tra UID có tồn tại hay chưa ?
	getAllReq := model.AllCloneReq{
		UID: clone.UID,
	}
	checkClones, err := cloneRepoInstance.All(getAllReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy clone theo %s %s", clone.UID, err))
	}
	for _, checkClone := range checkClones {
		if checkClone.AliveStatus != constants.CloneDelete && checkClone.ID.Hex() != clone.ID.Hex() {
			return fmt.Errorf("Clone UID %s đã đuợc sử dụng", clone.UID)
		}
	}

	// Kiểm tra clone có thuộc về user này hay không ?
	if c.Get("user_token") != nil {
		if clone.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Clone %s không thuôc về bạn", request.CloneID)
		}
	}

	// CẬp nhật trạng thái và DeviceID
	clone.AliveStatus = constants.CloneLive
	clone.DeviceID = ""
	err = cloneRepoInstance.Update(clone)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Reset clone: %s", err))
	}

	return c.JSON(success(clone))
}

// Search : support sort
func (h CloneHandler) Search(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		Total int               `json:"total" query:"total"`
		Data  []model.CloneInfo `json:"data" query:"data"`
	}
	type myRequest struct {
		UserToken   string                  `json:"user_token" query:"user_token"`
		AliveStatus []constants.AliveStatus `json:"alive_status" query:"alive_status"`
		DeviceID    string                  `json:"android_id" query:"android_id"`
		Date        interface{}             `json:"date" query:"date"`
		IsReg       *bool                   `json:"is_reg" query:"is_reg"`
		AppName     []constants.AppName     `json:"appname" query:"appname"`
		Page        int                     `json:"page" query:"page"`
		Limit       int                     `json:"limit" query:"limit"`
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
		Data: make([]model.CloneInfo, 0),
	}

	// Tạo search request
	searchCloneReq := model.SearchClone{
		Token:       request.UserToken,
		AliveStatus: request.AliveStatus,
		DeviceID:    request.DeviceID,
		IsReg:       request.IsReg,
		Date:        request.Date,
		AppName:     request.AppName,
		Page:        request.Page,
		Limit:       request.Limit,
	}
	if c.Get("user_token") != nil {
		searchCloneReq.Token = c.Get("user_token").(string)
	}

	cloneRepoInstance := cloneRepo.New(httpCtx)
	clones, err := cloneRepoInstance.Search(searchCloneReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Search danh sách clone %s", err))
	}

	// TODO: update total service
	response.Total, err = cloneRepoInstance.Total(searchCloneReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Đếm số lượng clone %s", err))
	}

	if len(clones) > 0 {
		response.Data = clones
	}

	return c.JSON(success(response))
}

type dataResponse struct {
	UID     string           `json:"uid,omitempty"`
	Email   string           `json:"email,omitempty"`
	IsErr   bool             `json:"is_err"`
	Message string           `json:"message"`
	Clone   *model.CloneInfo `json:"clone_info,omitempty"`
}

type myRequest struct {
	UserID          string                 `json:"user_id" query:"user_id" validate:"required"`
	ActionProfileID primitive.ObjectID     `json:"action_profile_id" query:"action_profile_id"`
	Clones          []model.CloneInfo      `json:"clones" query:"clones" validate:"required"`
	Novery          bool                   `json:"novery" query:"novery"`
	Country         constants.CloneCountry `json:"country" query:"country" validate:"required,oneof=vn indo en au"`
	AppName         string                 `json:"appname" query:"appname" validate:"required"`
	IsEncodeRSA     bool                   `json:"is_encode_rsa" query:"is_encode_rsa"`
}

// Create : Tạo clone
func (h CloneHandler) Create(c echo.Context) (err error) {

	httpCtx := c.Request().Context()

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

	// validate
	if len(request.Clones) < 1 {
		return fmt.Errorf("Cần tối thiểu 1 clone")
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
	result := make([]dataResponse, 0)
	wg := new(sync.WaitGroup)
	pool, _ := ants.NewPoolWithFunc(50, func(data interface{}) {
		res := h.SaveClone(user, data, request, c.RealIP())
		result = append(result, res)
		wg.Done()
	})

	defer pool.Release()
	// TODO: Với mỗi service code chạy 1 luồng riêng để move

	for _, cloneReq := range request.Clones {
		wg.Add(1)

		pool.Invoke(cloneReq)
	}
	wg.Wait()

	return c.JSON(success(result))
}

// Create : Tạo clone
func (h CloneHandler) SaveClone(user model.UserInfo, data interface{}, request *myRequest, CreatedIP string) (result dataResponse) {
	httpCtx := context.Background()
	cloneReq := data.(model.CloneInfo)
	cloneRepoInstance := cloneRepo.New(httpCtx)
	temp := dataResponse{
		UID:     cloneReq.UID,
		Email:   cloneReq.Email,
		Message: "Success",
	}

	if cloneReq.UID == "" {
		temp.IsErr = true
		temp.Message = fmt.Sprintf("Not allow empty uid")
		result = temp
		return
	}

	if cloneReq.Password == "" {
		temp.IsErr = true
		temp.Message = fmt.Sprintf("Not allow empty password")
		result = temp
		return
	}

	// Kiểm tra uid đã tồn tại hay chưa ?
	clone, err := cloneRepoInstance.GetOneActiveByUID(cloneReq.UID)
	if err != nil {
		temp.IsErr = true
		temp.Message = fmt.Sprintf("Lấy thông tin clone theo uid %s %s", cloneReq.UID, err)
		result = temp
		return
	}
	if clone.IsExists() {
		temp.IsErr = true
		temp.Message = fmt.Sprintf("clone uid %s is existed", cloneReq.UID)
		result = temp
		return
	}
	// Tạo model clone insert
	cloneInsert := model.CloneInfo{
		UID:             cloneReq.UID,
		Email:           cloneReq.Email,
		IP:              cloneReq.IP,
		Password:        cloneReq.Password,
		Secretkey:       cloneReq.Secretkey,
		Token:           user.Token,
		Cookie:          cloneReq.Cookie,
		AliveStatus:     constants.CloneLive,
		Country:         request.Country,
		ActionProfileID: &request.ActionProfileID,
		AppName:         request.AppName,
		//
	}
	clone, err = cloneRepoInstance.Insert(cloneInsert)
	if err != nil {
		temp.IsErr = true
		temp.Message = fmt.Sprintf("Insert clone %s %s err %s", cloneReq.Email, cloneReq.UID, cloneReq.Email)
		result = temp
		return
	}
	// Cập nhật thông tin reponse
	temp.Clone = &clone

	result = temp
	return
}

// SetAction : Cài đặt action cho clone
func (CloneHandler) SetAction(c echo.Context) (err error) {

	httpCtx := c.Request().Context()

	type myRequest struct {
		Token           string `json:"token" query:"token" validate:"required"`
		ActionProfileID string `json:"action_id" query:"action_id"`

		CloneID primitive.ObjectID `json:"clone_id" query:"clone_id" validate:"required"`
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

	cloneRepoInstance := cloneRepo.New(httpCtx)

	// Kiểm tra clone có tồn tại không ?
	clone, err := cloneRepoInstance.GetOneActiveByID(request.CloneID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy clone theo %s %s", request.CloneID, err))
	}
	if !clone.IsExists() {
		return fmt.Errorf("Không tìm thấy thông tin clone %s", request.CloneID)
	}
	// Kiểm tra clone có thuộc về user này hay không ?
	if clone.Token != request.Token {
		return fmt.Errorf("Clone %s không thuôc về bạn", request.CloneID)
	}
	ActionProfileID := primitive.NilObjectID
	if request.ActionProfileID != "" {
		ActionProfileID, err = primitive.ObjectIDFromHex(request.ActionProfileID)
		if err != nil {
			return fmt.Errorf("Mã device không hợp lệ %s", err)
		}
	}
	actionProfile, err := actionProfileRepo.New(httpCtx).GetOneByID(ActionProfileID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin action profile %s", err))
	}
	// Cập nhật action cho clone
	clone.ActionProfileID = &ActionProfileID
	clone.AppName = actionProfile.AppName
	err = cloneRepoInstance.Update(clone)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật action profile cho clone %s %s", request.CloneID, err))
	}

	return c.JSON(success(clone))
}

// Detail : Detail by cloneID
func (CloneHandler) Detail(c echo.Context) (err error) {
	httpCtx := c.Request().Context()

	type myRequest struct {
		CloneID primitive.ObjectID `json:"clone_id" query:"clone_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}

	cloneRepoInstance := cloneRepo.New(httpCtx)
	// Kiểm tra clone có tồn tại không ?
	clone, err := cloneRepoInstance.GetOneActiveByID(request.CloneID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy clone theo %s %s", request.CloneID, err))
	}
	if !clone.IsExists() {
		return fmt.Errorf("Không tìm thấy thông tin clone %s", request.CloneID)
	}
	// Kiểm tra clone có thuộc về user này hay không ?
	if c.Get("user_token") != nil {
		if clone.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Clone %s không thuôc về bạn", request.CloneID)
		}
	}

	return c.JSON(success(clone))
}

// UpdateServiceInfo : Cập nhật thông tin gói
// func (h CloneHandler) UpdateCloneInfo(c echo.Context) (err error) {

// 	httpCtx := c.Request().Context()

// 	type myRequest struct {
// 		CloneID      primitive.ObjectID `json:"clone_id" query:"clone_id" validate:"required"`
// 		Name         string             `json:"name" query:"name"`
// 		Birthday     string             `json:"birthday" query:"birthday"`
// 		LinkAvatar   string             `json:"link_avatar" query:"link_avatar"`
// 		Sex          string             `json:"sex" query:"sex"`
// 		NumberFriend int                `json:"number_friend" query:"number_friend"`
// 		IsAvatar     *bool              `json:"is_avatar" query:"is_avatar"`
// 		IsCheckpoint string             `json:"is_checkpoint" query:"is_checkpoint"`
// 		Friends      []string           `json:"friends" query:"friends"`
// 	}
// 	request := new(myRequest)
// 	if err = c.Bind(request); err != nil {
// 		return
// 	}
// 	if err = c.Validate(request); err != nil {
// 		return
// 	}
// 	birthday, err := time.Parse("2006-01-02T15:04:05.000Z", request.Birthday)

// 	if err != nil {
// 		fmt.Println("birthday can not covert %s", request.Birthday)
// 	}

// 	clonesRepoInstance := cloneRepo.New(httpCtx)

// 	// tạo update Request
// 	updateInfoReq := cloneRepo.UpdateCloneInfoReq{
// 		CloneID:      request.CloneID,
// 		Name:         request.Name,
// 		Birthday:     birthday,
// 		LinkAvatar:   request.LinkAvatar,
// 		Sex:          request.Sex,
// 		NumberFriend: request.NumberFriend,
// 		IsAvatar:     request.IsAvatar,
// 		IsCheckpoint: request.IsCheckpoint,
// 		Friends:      request.Friends,
// 	}

// 	clone, err := clonesRepoInstance.UpdateInfo(updateInfoReq)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Update info clone by uid %s %s", request.CloneID, err))
// 	}

// 	return c.JSON(success(clone))
// }
