package handler

import (
	"app/constants"
	"app/model"
	groupRepo "app/repo/mongo/group"
	groupUIDRepo "app/repo/mongo/groupuid"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GroupUIDHandler : struct init Groupuid
type GroupUIDHandler struct{}

// NewGroupUIDHandler : Tạo mới 1 đối tượng group uid handler
func NewGroupUIDHandler() *GroupUIDHandler {
	return &GroupUIDHandler{}
}

// All : Lấy toàn bộ group uid
func (GroupUIDHandler) All(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		GroupID string                   `json:"group_id" query:"group_id" validate:"required"`
		Status  []constants.CommonStatus `json:"status" query:"status"`
		Token   string                   `json:"token" query:"token"`
		Offset  string                   `json:"offset" query:"offset"`
		Limit   int                      `json:"limit" query:"limit"`
	}

	type myResponse struct {
		LastOffset *primitive.ObjectID  `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                  `json:"total" query:"total"`
		Data       []model.GroupUIDInfo `json:"data" query:"data"`
	}
	request := new(myRequest)
	response := new(myResponse)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	if c.Get("user_token") != nil {
		request.Token = c.Get("user_token").(string)
	}
	// Validate group id
	groupObjectID, err := primitive.ObjectIDFromHex(request.GroupID)
	if err != nil {
		return fmt.Errorf("Mã số nhóm không hợp lệ %s", err)
	}
	group, err := groupRepo.New(httpCtx).GetOneByID(groupObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin nhóm %s", err))
	}
	if !group.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin nhóm"))
	}

	// Lấy toàn bộ uid theo group id
	allGroupUIDReq := model.AllGroupUIDReq{
		GroupID: group.ID,
		Token:   request.Token,
		Status:  request.Status,
		Limit:   request.Limit,
	}
	if request.Offset != "" {
		offsetObjectID, err := primitive.ObjectIDFromHex(request.Offset)
		if err != nil {
			return fmt.Errorf("Offset không hợp lệ %s", err)
		}

		allGroupUIDReq.Offset = offsetObjectID
	}

	groupUids, err := groupUIDRepo.New(httpCtx).All(allGroupUIDReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách uid: %s", err))
	}
	if len(groupUids) > 0 {
		response.Data = groupUids
		response.Total = len(groupUids)
		response.LastOffset = &groupUids[len(groupUids)-1].ID
	}
	return c.JSON(success(response))
}

// Delete : Xóa group uid
func (GroupUIDHandler) Delete(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		GroupUID string `json:"group_uid_id" query:"group_uid_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	groupObjectID, err := primitive.ObjectIDFromHex(request.GroupUID)
	if err != nil {
		return fmt.Errorf("Mã số uid không hợp lệ %s", err)
	}
	groupUidRepoInstance := groupUIDRepo.New(httpCtx)
	//
	groupUID, err := groupUidRepoInstance.GetOneByID(groupObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thôngtin group %s", err))
	}
	if !groupUID.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin group"))
	}
	// validate người thực thi: user id
	if c.Get("user_token") != nil {
		if groupUID.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Mã số bạn bè không hợp lệ")
		}
	}
	// Xóa.
	groupUID.Status = constants.Delete

	err = groupUidRepoInstance.Update(groupUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật bình luận: %s", err))
	}

	return c.JSON(success(groupUID))
}
