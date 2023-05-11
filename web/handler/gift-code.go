package handler

import (
	"app/constants"
	"fmt"
	"net/http"

	"app/model"

	codeRepo "app/repo/mongo/code"
	giftRepo "app/repo/mongo/gift"
	giftUsedRepo "app/repo/mongo/giftused"
	userRepo "app/repo/mongo/user"
	walletLogRepo "app/repo/mongo/walletlog"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GiftCodeHandler : struct init gift code
type GiftCodeHandler struct{}

// NewGiftCodeHandler : Tạo mới 1 đối tượng gift code handler
func NewGiftCodeHandler() *GiftCodeHandler {
	return &GiftCodeHandler{}
}

// Gen : Tạo mới gift code
func (GiftCodeHandler) Gen(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		Quantity int `json:"quantity" query:"quantity" validate:"required"`
		Value    int `json:"value" query:"value" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// TODO: Gen ma code với quantity
	// Gen identity code
	codes, err := codeRepo.New(httpCtx).Generate("gift_code", request.Quantity)
	if err != nil || len(codes) == 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Tạo mã dịch vụ %s", err))
	}

	dataInsert := make([]model.GiftCode, 0)
	for _, code := range codes {
		dataInsert = append(dataInsert, model.GiftCode{
			Code:   code,
			Value:  request.Value,
			Status: constants.GiftActive,
		})
	}
	// TODO: insert to db
	err = giftRepo.New(httpCtx).InsertMany(dataInsert)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Tạo mã khuyến mãi %s", err))
	}

	return c.JSON(success(codes))
}

// Apply : Áp dụng mã khuyến mãi
func (GiftCodeHandler) Apply(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		UserID   string `json:"user_id" query:"user_id" validate:"required"`
		GiftCode string `json:"gift_code" query:"gift_code" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	// validate người thực thi: user id
	if c.Get("user_id_str") != nil {
		request.UserID = c.Get("user_id_str").(string)
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// validate
	userObjectID, err := primitive.ObjectIDFromHex(request.UserID)
	if err != nil {
		return fmt.Errorf("Mã số nhân viên không hợp lệ %s", err)
	}

	userRepoInstance := userRepo.New(httpCtx)
	giftRepoInstance := giftRepo.New(httpCtx)
	giftUsedRepoInstance := giftUsedRepo.New(httpCtx)
	walletRepoInstance := walletLogRepo.New(httpCtx)

	// Lấy thông tin tài khoản
	user, err := userRepoInstance.GetOneActiveByID(userObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng %s", err))
	}
	if !user.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}

	// TODO: Do many things with mongo transaction
	session, err := giftRepoInstance.Session.ConClient.StartSession()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Get mongo session: %s", err))
	}
	// Start transaction
	err = session.StartTransaction()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Start mongo transaction: %s", err))
	}
	defer session.EndSession(httpCtx)

	err = mongo.WithSession(httpCtx, session, func(sessCtx mongo.SessionContext) (err error) {
		// TODO: Apply code
		giftCode, err := giftRepoInstance.ApplyWithSessionCtx(sessCtx, request.GiftCode)
		if err != nil || !giftCode.IsExists() {
			session.AbortTransaction(sessCtx)
			return fmt.Errorf("Mã code %s không hợp lệ", request.GiftCode)
		}

		// TODO: Tạo gift code used
		giftCodeUsed := model.GiftCodeUsed{
			UserID:   &user.ID,
			Username: user.Username,
			Fullname: user.Fullname,
			Code:     giftCode.Code,
			Value:    giftCode.Value,
		}

		_, err = giftUsedRepoInstance.InsertWithSessionCtx(sessCtx, giftCodeUsed)
		if err != nil {
			session.AbortTransaction(sessCtx)
			return fmt.Errorf("Tạo mã khuyến mãi cho khách %s %s", user.Username, err)
		}

		// TODO: Cộng tiền cho khách
		balanceChanged := giftCode.Value
		user.Balance = user.Balance + balanceChanged
		err = userRepoInstance.UpdateWithSessionCtx(sessCtx, user)
		if err != nil {
			session.AbortTransaction(sessCtx)
			return fmt.Errorf("Cập nhật tài khoản cho khách %s %s", user.Username, err)
		}

		giftWallet := model.GiftWallet{
			Code:  giftCode.Code,
			Value: giftCode.Value,
		}

		// TODO: Ghi wallet log
		walletLog := model.WalletLog{
			Type:    constants.WallletGift,
			Token:   user.Token,
			Value:   balanceChanged,
			Balance: user.Balance,
			Gift:    &giftWallet,
		}
		_, err = walletRepoInstance.InsertWithSessionCtx(sessCtx, walletLog)
		if err != nil {
			session.AbortTransaction(sessCtx)
			return fmt.Errorf("Ghi log nạp tiền %s %s", user.Username, err)
		}

		session.CommitTransaction(sessCtx)
		return
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Fail in transaction: %s", err))
	}

	return c.JSON(success(nil))
}

// All : Lấy danh sách các mã khuyến mãi
func (GiftCodeHandler) All(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                 `json:"total" query:"total"`
		Data       []model.GiftCode    `json:"data" query:"data"`
	}
	type myRequest struct {
		GiftCode string                 `json:"gift_code" query:"gift_code"`
		Status   []constants.GiftStatus `json:"status" query:"status"`
		Offset   string                 `json:"offset" query:"offset"`
		Limit    int                    `json:"limit" query:"limit"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}

	// Tạo all Gift request
	allGiftReq := model.AllGiftReq{
		Code:   request.GiftCode,
		Status: request.Status,
		Offset: primitive.NilObjectID,
		Limit:  request.Limit,
	}
	if request.Offset != "" {
		offsetObjectID, err := primitive.ObjectIDFromHex(request.Offset)
		if err != nil {
			return fmt.Errorf("Offset không hợp lệ %s", err)
		}
		allGiftReq.Offset = offsetObjectID
	}

	response := myResponse{
		Data: make([]model.GiftCode, 0),
	}
	results, err := giftRepo.New(httpCtx).All(allGiftReq)
	if err != nil {
		return fmt.Errorf("Lấy danh sách mã khuyến mãi thất bại %s", err)
	}

	if len(results) > 0 {
		response.Data = results
		response.Total = len(results)
		response.LastOffset = &results[len(results)-1].ID
	}

	return c.JSON(success(response))
}

// AllHistory : Lấy danh sách các mã khuyến mãi đã sử dụng
func (GiftCodeHandler) AllHistory(c echo.Context) (err error) {

	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID  `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                  `json:"total" query:"total"`
		Data       []model.GiftCodeUsed `json:"data" query:"data"`
	}
	type myRequest struct {
		GiftCode string `json:"gift_code" query:"gift_code"`
		UserID   string `json:"user_id" query:"user_id"`
		Offset   string `json:"offset" query:"offset"`
		Limit    int    `json:"limit" query:"limit"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Tạo all Gift used request
	allGiftUsedReq := model.AllGiftUsedReq{
		Code:   request.GiftCode,
		UserID: primitive.NilObjectID,
		Offset: primitive.NilObjectID,
		Limit:  request.Limit,
	}
	if request.Offset != "" {
		offsetObjectID, err := primitive.ObjectIDFromHex(request.Offset)
		if err != nil {
			return fmt.Errorf("Offset không hợp lệ %s", err)
		}
		allGiftUsedReq.Offset = offsetObjectID
	}
	if request.UserID != "" {
		userObjectID, err := primitive.ObjectIDFromHex(request.UserID)
		if err != nil {
			return fmt.Errorf("Mã số người dùng không hợp lệ %s", err)
		}
		allGiftUsedReq.UserID = userObjectID
	}

	response := myResponse{
		Data: make([]model.GiftCodeUsed, 0),
	}
	results, err := giftUsedRepo.New(httpCtx).All(allGiftUsedReq)
	if err != nil {
		return fmt.Errorf("Lấy danh sách mã khuyến mãi đã sử dụng thất bại %s", err)
	}
	if len(results) > 0 {
		response.Data = results
		response.Total = len(results)
		response.LastOffset = &results[len(results)-1].ID
	}

	return c.JSON(success(response))
}

// UpdateStatus : Cập nhật trang thái
// . Active -> Waiting
// . Waiting -> Active
func (GiftCodeHandler) UpdateStatus(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		GiftCode string               `json:"gift_code" query:"gift_code" validate:"required"`
		Status   constants.GiftStatus `json:"status" query:"status" validate:"required,oneof=Active Waiting"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// TODO: Lấy code
	gift, err := giftRepo.New(httpCtx).GetOneByCode(request.GiftCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy mã khuyến mãi %s", err))
	}
	if !gift.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy mã khuyến mãi"))
	}

	gift.Status = request.Status
	err = giftRepo.New(httpCtx).Update(gift)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật trạng thái mã khuyến mãi %s", err))
	}

	return c.JSON(success(gift))
}

// Delete : Cập nhật trang thái
func (GiftCodeHandler) Delete(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		GiftCode string `json:"gift_code" query:"gift_code" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// TODO: Lấy code
	gift, err := giftRepo.New(httpCtx).GetOneByCode(request.GiftCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy mã khuyến mãi %s", err))
	}
	if !gift.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy mã khuyến mãi"))
	}

	gift.Status = constants.GiftDelete
	err = giftRepo.New(httpCtx).Update(gift)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Xóa mã khuyến mãi %s", err))
	}

	return c.JSON(success(gift))
}
