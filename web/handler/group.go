package handler

import (
	"app/constants"
	"app/model"
	groupRepo "app/repo/mongo/group"
	groupUIDRepo "app/repo/mongo/groupuid"
	userRepo "app/repo/mongo/user"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GroupHandler init
type GroupHandler struct{}

// NewGroupHandler : Tạo đối tượng Group Handle
func NewGroupHandler() *GroupHandler {
	return &GroupHandler{}
}

// Create func
func (GroupHandler) Create(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	type myRequest struct {
		UserID string   `json:"user_id" query:"user_id" validate:"required"`
		Name   string   `json:"name" query:"name" validate:"required"`
		UIDS   []string `json:"uids" query:"uids" validate:"required"`
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

	// validate
	if len(request.UIDS) < 1 {
		return fmt.Errorf("Cần tối thiểu 1 uid")
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
	group := model.GroupInfo{
		Name:   request.Name,
		Token:  user.Token,
		Status: constants.Active,
	}

	group, err = groupRepo.New(httpCtx).Insert(group)
	existedUID := map[string]bool{}

	// Tạo uid
	for _, uid := range request.UIDS {
		if !existedUID[uid] {
			temp := model.GroupUIDInfo{
				GroupID: group.ID,
				UID:     uid,
				Token:   user.Token,
				Status:  constants.Active,
			}
			err = groupUIDRepo.New(httpCtx).UpsertByUID(temp)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Upsert group theo uid %s", err))
			}
		}
	}

	return c.JSON(success(nil))
}

// All : Lấy toàn bộ group
func (GroupHandler) All(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                 `json:"total" query:"total"`
		Data       []model.GroupInfo   `json:"data" query:"data"`
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
	if err = c.Validate(request); err != nil {
		return
	}
	if c.Get("user_token") != nil {
		request.Token = c.Get("user_token").(string)
	}
	// Tạo all Group request
	allGroupReq := model.AllGroupReq{
		Name:   request.Name,
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

		allGroupReq.Offset = offsetObjectID
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.GroupInfo, 0),
	}

	groups, err := groupRepo.New(httpCtx).All(allGroupReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách bạn bè : %s", err))
	}

	if len(groups) > 0 {
		response.Data = groups
		response.Total = len(groups)
		response.LastOffset = &groups[len(groups)-1].ID
	}

	return c.JSON(success(response))
}

// Delete : Xóa group
func (GroupHandler) Delete(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		GroupID string `json:"group_id" query:"group_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	groupbjectID, err := primitive.ObjectIDFromHex(request.GroupID)
	if err != nil {
		return fmt.Errorf("Mã số nhóm không hợp lệ %s", err)
	}
	groupRepoInstance := groupRepo.New(httpCtx)
	//
	group, err := groupRepoInstance.GetOneByID(groupbjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin nhóm %s", err))
	}
	if !group.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin nhóm"))
	}
	// validate người thực thi: user id

	if c.Get("user_token") != nil {
		if group.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Mã số bạn bè không hợp lệ")
		}
	}

	// Xóa.
	group.Status = constants.Delete

	err = groupRepoInstance.Update(group)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật nhóm: %s", err))
	}

	return c.JSON(success(group))
}

// Detail : Thông tin chi tiết group
func (GroupHandler) Detail(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		GroupID string `json:"group_id" query:"group_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	groupbjectID, err := primitive.ObjectIDFromHex(request.GroupID)
	if err != nil {
		return fmt.Errorf("Mã số nhóm không hợp lệ %s", err)
	}
	groupRepoInstance := groupRepo.New(httpCtx)
	//
	group, err := groupRepoInstance.GetOneByID(groupbjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin nhóm %s", err))
	}
	if !group.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin nhóm"))
	}
	// validate người thực thi: user id

	if c.Get("user_token") != nil {
		if group.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Mã số bạn bè không hợp lệ")
		}
	}

	return c.JSON(success(group))
}
