package handler

import (
	"app/constants"
	"app/model"
	pageRepo "app/repo/mongo/page"
	userRepo "app/repo/mongo/user"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PageHandler init
type PageHandler struct{}

// NewPageHandler : Tạo đối tượng Page Handle
func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

// Create func
func (PageHandler) Create(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	type myRequest struct {
		UserID string   `json:"user_id" query:"user_id" validate:"required"`
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
	existedUID := map[string]bool{}
	for _, uid := range request.UIDS {
		// Kiểm tra trùng uid
		if !existedUID[uid] {
			// add vào mảng để check trùng
			existedUID[uid] = true
			temp := model.PageInfo{
				UID:    uid,
				Token:  user.Token,
				Status: constants.Active,
			}
			err = pageRepo.New(httpCtx).UpsertByUID(temp)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Upsert page theo uid %s", err))
			}
		}
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Fail in transaction: %s", err))
	}

	return c.JSON(success(nil))
}

// All : Lấy toàn bộ Page
func (PageHandler) All(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                 `json:"total" query:"total"`
		Data       []model.PageInfo    `json:"data" query:"data"`
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
	// Tạo all Page request
	allPageReq := model.AllPageReq{
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

		allPageReq.Offset = offsetObjectID
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.PageInfo, 0),
	}

	Pages, err := pageRepo.New(httpCtx).All(allPageReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách page : %s", err))
	}

	if len(Pages) > 0 {
		response.Data = Pages
		response.Total = len(Pages)
		response.LastOffset = &Pages[len(Pages)-1].ID
	}

	return c.JSON(success(response))
}

// Delete : Xóa Page
func (PageHandler) Delete(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		PageID string `json:"page_id" query:"page_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	PagebjectID, err := primitive.ObjectIDFromHex(request.PageID)
	if err != nil {
		return fmt.Errorf("Mã số page không hợp lệ %s", err)
	}
	PageRepoInstance := pageRepo.New(httpCtx)
	//
	Page, err := PageRepoInstance.GetOneByID(PagebjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin page %s", err))
	}
	if !Page.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin page"))
	}
	// validate người thực thi: user id
	if c.Get("user_token") != nil {
		if Page.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Token không hợp lệ")
		}
	}

	// Xóa.
	Page.Status = constants.Delete

	err = PageRepoInstance.Update(Page)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật page: %s", err))
	}

	// Del redis key
	redisClient := cfg.Redis["core"].GetClient()
	pageKey := fmt.Sprintf(constants.PageKey, Page.ID)

	redisClient.Del(pageKey)

	return c.JSON(success(Page))
}

// Detail : Thông tin page
func (PageHandler) Detail(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		PageID string `json:"page_id" query:"page_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate
	PagebjectID, err := primitive.ObjectIDFromHex(request.PageID)

	if err != nil {
		return fmt.Errorf("Mã số page không hợp lệ %s", err)
	}
	PageRepoInstance := pageRepo.New(httpCtx)
	//
	Page, err := PageRepoInstance.GetOneByID(PagebjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin page %s", err))
	}
	if !Page.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin page"))
	}
	// validate người thực thi: user id
	if c.Get("user_token") != nil {
		if Page.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Token không hợp lệ")
		}
	}

	return c.JSON(success(Page))
}
