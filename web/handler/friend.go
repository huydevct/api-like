package handler

import (
	"fmt"
	"net/http"

	"app/constants"
	"app/model"

	friendRepo "app/repo/mongo/friend"
	friendUIDRepo "app/repo/mongo/frienduid"
	userRepo "app/repo/mongo/user"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FriendHandler : struct init friend
type FriendHandler struct{}

// NewFriendHandler : Tạo mới 1 đối tượng bạn bè handler
func NewFriendHandler() *FriendHandler {
	return &FriendHandler{}
}

// Create : Tạo mới bạn bè
// Temp ignore, need replica mongo
func (FriendHandler) Create(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
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

	friend := model.FriendInfo{
		Name:   request.Name,
		Token:  user.Token,
		Status: constants.Active,
	}

	friend, err = friendRepo.New(httpCtx).Insert(friend)
	if err != nil {
		return fmt.Errorf("Tạo bạn bè: %s", err)
	}
	existedUID := map[string]bool{}

	// Tạo uid
	for _, uid := range request.UIDS {
		if !existedUID[uid] {
			// add vào mãng để check trùng
			existedUID[uid] = true
			temp := model.FriendUIDInfo{
				FriendID: friend.ID,
				UID:      uid,
				Token:    user.Token,
				Status:   constants.Active,
			}
			err = friendUIDRepo.New(httpCtx).UpsertByUID(temp)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Upsert friend theo uid %s", err))
			}
		}
	}
	return c.JSON(success(nil))
}

// Delete : Xóa friend
func (FriendHandler) Delete(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		FriendID string `json:"friend_id" query:"friend_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	friendObjectID, err := primitive.ObjectIDFromHex(request.FriendID)
	if err != nil {
		return fmt.Errorf("Mã số bạn bè không hợp lệ %s", err)
	}
	friendRepoInstance := friendRepo.New(httpCtx)
	//
	friend, err := friendRepoInstance.GetOneByID(friendObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin bạn bè %s", err))
	}
	if !friend.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin bạn bè"))
	}
	// validate người thực thi: user id
	if c.Get("user_token") != nil {
		if friend.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Mã số bạn bè không hợp lệ")
		}
	}
	// Xóa.
	friend.Status = constants.Delete

	err = friendRepoInstance.Update(friend)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật bình luận: %s", err))
	}

	return c.JSON(success(friend))
}

// Detail : Hiển thị chi tiết bạn bè
func (FriendHandler) Detail(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		FriendID string `json:"friend_id" query:"friend_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	friendObjectID, err := primitive.ObjectIDFromHex(request.FriendID)
	if err != nil {
		return fmt.Errorf("Mã số bạn bè không hợp lệ %s", err)
	}
	friendRepoInstance := friendRepo.New(httpCtx)
	//
	friend, err := friendRepoInstance.GetOneByID(friendObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin bạn bè %s", err))
	}
	if !friend.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin bạn bè"))
	}
	// validate người thực thi: user id
	if c.Get("user_token") != nil {
		if friend.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Mã số bạn bè không hợp lệ")
		}
	}

	return c.JSON(success(friend))
}

// All : Lấy toàn bộ friend
func (FriendHandler) All(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                 `json:"total" query:"total"`
		Data       []model.FriendInfo  `json:"data" query:"data"`
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
	// Tạo all Friend request
	allFriendReq := model.AllFriendReq{
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

		allFriendReq.Offset = offsetObjectID
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.FriendInfo, 0),
	}

	friends, err := friendRepo.New(httpCtx).All(allFriendReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách bạn bè : %s", err))
	}

	if len(friends) > 0 {
		response.Data = friends
		response.Total = len(friends)
		response.LastOffset = &friends[len(friends)-1].ID
	}

	return c.JSON(success(response))
}
