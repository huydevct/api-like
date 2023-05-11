package handler

import (
	"app/constants"
	"app/model"
	"fmt"
	"net/http"

	packRepo "app/repo/mongo/pack"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PackHandler : struct init pack
type PackHandler struct{}

// NewPackHandler : Tạo mới 1 đối tượng pack handler
func NewPackHandler() *PackHandler {
	return &PackHandler{}
}

// Create : Tạo mới gói
func (PackHandler) Create(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		Money           int                   `json:"money" query:"money" validate:"required"`
		Bonus           float64               `json:"bonus" query:"bonus"`
		IsLikeSub       bool                  `json:"isLikeSub,omitempty" query:"isLikeSub"`
		IsInstagram     bool                  `json:"is_instagram,omitempty" query:"is_instagram"`
		IsYoutube       bool                  `json:"is_youtube,omitempty" query:"is_youtube"`
		IsReg           bool                  `json:"is_reg,omitempty" query:"is_reg,omitempty"`
		ShareLiveStream bool                  `json:"share_live_stream,omitempty" query:"share_live_stream,omitempty"`
		User100App      bool                  `json:"user_100_app,omitempty" query:"user_100_app,omitempty"`
		Type            constants.PackageType `json:"type" query:"type" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	{
		// Validate type
		if request.Type != constants.PackageDeposit {
			return fmt.Errorf("Loại gói không được hỗ trợ %s", request.Type)
		}
	}
	// tạo model package
	pack := model.Package{
		Money:           request.Money,
		Bonus:           request.Bonus,
		Type:            request.Type,
		IsLikeSub:       request.IsLikeSub,
		IsInstagram:     request.IsInstagram,
		IsYoutube:       request.IsYoutube,
		ShareLiveStream: request.ShareLiveStream,
		User100App:      request.User100App,
		IsReg:           request.IsReg,
		Status:          constants.Active,
	}

	result, err := packRepo.New(httpCtx).Insert(pack)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Tạo gói mới %s", err))
	}

	return c.JSON(success(result))
}

// Update : Cập nhật gói
func (PackHandler) Update(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		ID              string                  `json:"id" query:"id" validate:"required"`
		Money           *int                    `json:"money" query:"money"`
		Bonus           *float64                `json:"bonus" query:"bonus"`
		IsLikeSub       bool                    `json:"is_like_sub,omitempty" query:"is_like_sub"`
		IsInstagram     bool                    `json:"is_instagram,omitempty" query:"is_instagram"`
		IsYoutube       bool                    `json:"is_youtube,omitempty" query:"is_youtube"`
		IsReg           bool                    `json:"is_reg,omitempty" query:"is_reg"`
		ShareLiveStream bool                    `json:"share_live_stream,omitempty" query:"share_live_stream"`
		User100App      bool                    `json:"user_100_app,omitempty" query:"user_100_app"`
		Type            *constants.PackageType  `json:"type" query:"type"`
		Status          *constants.CommonStatus `json:"status" query:"status"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate packID
	packObjectID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		return fmt.Errorf("Mã gói không hợp lệ %s", err)
	}
	packRepoInstance := packRepo.New(httpCtx)
	//
	pack, err := packRepoInstance.GetOneActiveByID(packObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy mã gói ID %s", err))
	}
	if !pack.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy mã gói"))
	}
	// cập nhật gói
	{
		// Cap nhat tien
		if request.Money != nil {
			if pack.Money != *request.Money {
				pack.Money = *request.Money
			}
		}
	}
	{
		// Cap nhat bonus
		if request.Bonus != nil {
			if pack.Bonus != *request.Bonus {
				pack.Bonus = *request.Bonus
			}
		}
	}
	{
		// Cap nhat type
		if request.Type != nil {
			if *request.Type != constants.PackageDeposit {
				return fmt.Errorf("Loại gói không được hỗ trợ %s", *request.Type)
			}
			if pack.Type != *request.Type {
				pack.Type = *request.Type
			}
		}
	}
	pack.IsLikeSub = false
	pack.ShareLiveStream = false
	pack.User100App = false
	pack.IsReg = false
	pack.IsInstagram = false
	pack.IsYoutube = false
	if request.IsLikeSub == true {
		pack.IsLikeSub = true
	}

	if request.IsInstagram == true {
		pack.IsInstagram = true
	}

	if request.IsYoutube == true {
		pack.IsYoutube = true
	}

	if request.ShareLiveStream == true {
		pack.ShareLiveStream = true
	}

	if request.User100App == true {
		pack.User100App = true
	}

	if request.IsReg == true {
		pack.IsReg = true
	}

	{
		// CẬp nhật status
		if request.Status != nil {
			if *request.Status != constants.Active && *request.Status != constants.Pause {
				return fmt.Errorf("Trạng thái gói không được hỗ trợ %s", *request.Status)
			}
			if pack.Status != *request.Status {
				pack.Status = *request.Status
			}
		}
	}

	err = packRepoInstance.Update(pack)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật gói %s", err))
	}

	return c.JSON(success(pack))
}

// Delete : Xoa gói
func (PackHandler) Delete(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		ID string `json:"id" query:"id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate packID
	packObjectID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		return fmt.Errorf("Mã gói không hợp lệ %s", err)
	}
	packRepoInstance := packRepo.New(httpCtx)
	//
	pack, err := packRepoInstance.GetOneActiveByID(packObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy mã gói ID %s", err))
	}
	if !pack.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy mã gói"))
	}
	// cập nhật trạng thái của gói
	pack.Status = constants.Delete

	err = packRepoInstance.Update(pack)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Xoá gói %s", err))
	}

	return c.JSON(success(nil))
}

// Detail : chi tiết gói
func (PackHandler) Detail(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		ID string `json:"id" query:"id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate packID
	packObjectID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		return fmt.Errorf("Mã gói không hợp lệ %s", err)
	}
	packRepoInstance := packRepo.New(httpCtx)
	//
	pack, err := packRepoInstance.GetOneActiveByID(packObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy mã gói ID %s", err))
	}
	if !pack.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy mã gói"))
	}
	return c.JSON(success(pack))
}

// All : tất cả gói
// offset, limit
func (PackHandler) All(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                 `json:"total" query:"total"`
		Data       []model.Package     `json:"data" query:"data"`
	}
	type myRequest struct {
		Offset string                   `json:"offset" query:"offset"`
		Limit  int                      `json:"limit" query:"limit"`
		Status []constants.CommonStatus `json:"status" query:"status"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}

	allPackReq := model.AllPackageReq{
		Status: request.Status,
		Offset: primitive.NilObjectID,
		Limit:  request.Limit,
	}
	// validate
	if request.Offset != "" {
		allPackReq.Offset, err = primitive.ObjectIDFromHex(request.Offset)
		if err != nil {
			return fmt.Errorf("Offset không hợp lệ %s", err)
		}
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.Package, 0),
	}
	packs, err := packRepo.New(httpCtx).All(allPackReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách mã gói %s", err))
	}

	if len(packs) > 0 {
		response.Data = packs
		response.Total = len(packs)
		response.LastOffset = &packs[len(packs)-1].ID
	}

	return c.JSON(success(response))
}
