package handler

import (
	"app/constants"
	"app/model"
	friendRepo "app/repo/mongo/friend"
	friendUidsRepo "app/repo/mongo/frienduid"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FriendUIDHandler : struct init frienduid
type FriendUIDHandler struct{}

// NewFriendUIDHandler : Tạo mới 1 đối tượng friend uid handler
func NewFriendUIDHandler() *FriendUIDHandler {
	return &FriendUIDHandler{}
}

// All : Lấy toàn bộ friend uid
func (FriendUIDHandler) All(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		FriendID string                   `json:"friend_id" query:"friend_id" validate:"required"`
		Status   []constants.CommonStatus `json:"status" query:"status"`
		Token    string                   `json:"token" query:"token"`
		Offset   string                   `json:"offset" query:"offset"`
		Limit    int                      `json:"limit" query:"limit"`
	}

	type myResponse struct {
		LastOffset *primitive.ObjectID   `json:"last_offset,omitempty"`
		Total      int                   `json:"total" query:"total"`
		Data       []model.FriendUIDInfo `json:"data" query:"data"`
	}
	request := new(myRequest)
	response := new(myResponse)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate friend id
	firendObjectID, err := primitive.ObjectIDFromHex(request.FriendID)
	if err != nil {
		return fmt.Errorf("Mã số bạn bè không hợp lệ %s", err)
	}
	friend, err := friendRepo.New(httpCtx).GetOneByID(firendObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin bạn bè%s", err))
	}
	if !friend.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin bạn bè"))
	}

	// Lấy toàn bộ uid theo friend id
	allFriendUIDReq := model.AllFriendUIDReq{
		FriendID: friend.ID,
		Token:    request.Token,
		Status:   request.Status,
		Offset:   primitive.NilObjectID,
		Limit:    request.Limit,
	}
	if request.Offset != "" {
		offsetObjectID, err := primitive.ObjectIDFromHex(request.Offset)
		if err != nil {
			return fmt.Errorf("Offset không hợp lệ %s", err)
		}

		allFriendUIDReq.Offset = offsetObjectID
	}

	friendUIDS, err := friendUidsRepo.New(httpCtx).All(allFriendUIDReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh uid: %s", err))
	}
	if len(friendUIDS) > 0 {
		response.Data = friendUIDS
		response.Total = len(friendUIDS)
		response.LastOffset = &friendUIDS[len(friendUIDS)-1].ID
	}
	return c.JSON(success(response))
}

// Delete : Xóa friend uid
func (FriendUIDHandler) Delete(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		FriendUID string `json:"friend_uid_id" query:"friend_uid_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	friendObjectID, err := primitive.ObjectIDFromHex(request.FriendUID)
	if err != nil {
		return fmt.Errorf("Mã số uid không hợp lệ %s", err)
	}
	friendUidRepoInstance := friendUidsRepo.New(httpCtx)
	//
	friendUID, err := friendUidRepoInstance.GetOneByID(friendObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin bạn bè %s", err))
	}
	if !friendUID.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin bạn bè"))
	}
	// validate người thực thi: user id
	if c.Get("user_token") != nil {
		if friendUID.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Mã số bạn bè không hợp lệ")
		}
	}
	// Xóa.
	friendUID.Status = constants.Delete

	err = friendUidRepoInstance.Update(friendUID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật bình luận: %s", err))
	}

	return c.JSON(success(friendUID))
}
