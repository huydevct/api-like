package handler

import (
	"app/utils"
	"fmt"
	"net/http"

	"app/constants"
	"app/model"

	codeRepo "app/repo/mongo/code"
	packRepo "app/repo/mongo/pack"
	transRepo "app/repo/mongo/transaction"
	userRepo "app/repo/mongo/user"
	walletLogRepo "app/repo/mongo/walletlog"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TransactionHandler : struct init transaction
type TransactionHandler struct{}

// NewTransactionHandler : Tạo mới 1 đối tượng transaction handler
func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{}
}

// Create : Tạo mới giao dịch
func (TransactionHandler) Create(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		Token  string `json:"token" query:"token" validate:"required"`
		PackID string `json:"pack_id" query:"pack_id" validate:"required"`
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

	// TODO: validate package
	packageObjectID, err := primitive.ObjectIDFromHex(request.PackID)
	if err != nil {
		return fmt.Errorf("Mã số gói không hợp lệ %s", err)
	}
	// Kiểm tra gói có tồn tại hay không ?
	pack, err := packRepo.New(httpCtx).GetOneActiveByID(packageObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin gói %s", err))
	}
	if !pack.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin gói"))
	}
	// Lấy thông tin tài khoản
	user, err := userRepo.New(httpCtx).GetOneActiveByToken(request.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng %s", err))
	}
	if !user.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}
	// Gen identity code
	codes, err := codeRepo.New(httpCtx).Generate("transaction_code", 1)
	if err != nil || len(codes) == 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Tạo mã giao dịch %s", err))
	}
	// Tạo transaction
	transaction := model.Transaction{
		Code:            fmt.Sprintf("AFARM%s", codes[0]),
		UserID:          &user.ID,
		Username:        user.Username,
		Fullname:        user.Fullname,
		Token:           user.Token,
		Status:          constants.TransPending,
		Value:           pack.Money,
		ValueInt:        pack.Money,
		IsLikeSub:       pack.IsLikeSub,
		IsInstagram:     pack.IsInstagram,
		IsYoutube:       pack.IsYoutube,
		IsReg:           pack.IsReg,
		ShareLiveStream: pack.ShareLiveStream,
		User100App:      pack.User100App,
		Bonus:           pack.Bonus / 100 * float64(pack.Money),
	}

	result, err := transRepo.New(httpCtx).Insert(transaction)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Tạo mới giao dịch %s", err))
	}

	return c.JSON(success(result))
}

// Active : Cập nhật trạng thái giao dịch thành active
// . Cộng tiền cho user
func (TransactionHandler) Active(c echo.Context) (err error) {
	httpCtx := c.Request().Context()

	type myRequest struct {
		Code  string `json:"code" query:"code" validate:"required"`
		Money int    `json:"money" query:"money" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}

	transRepoInstance := transRepo.New(httpCtx)
	userRepoInstance := userRepo.New(httpCtx)
	walletRepoInstance := walletLogRepo.New(httpCtx)

	transaction, err := transRepoInstance.GetOneByCode(request.Code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin giao dịch %s", err))
	}
	if !transaction.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin giao dịch"))
	}
	// Kiểm tra số tiền truyền lên có hợp lệ hay không ?
	if utils.ConvertToInt(transaction.Value) != request.Money {
		return fmt.Errorf("Số tiền không đúng")
	}

	// Kiểm tra trạng thái transaction có là pending hay không ?
	// Cập nhật thành Active + chuyển tiền cho khách
	if transaction.Status == constants.TransPending {
		// Lấy thông tin user
		user, err := userRepoInstance.GetOneByToken(transaction.Token)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng %s", err))
		}
		if !user.IsExists() {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
		}

		// TODO: Do many things with mongo transaction
		session, err := transRepoInstance.Session.ConClient.StartSession()
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
			// TODO: Cập nhật trạng thái giao dịch
			transaction.Status = constants.TransActive
			err = transRepoInstance.UpdateWithSessionCtx(sessCtx, transaction)
			if err != nil {
				session.AbortTransaction(sessCtx)
				return fmt.Errorf("Cập nhật trạng thái giao dịch %s", err)
			}

			// TODO: Cộng tiền cho khách
			balanceChanged := transaction.ValueInt + int(transaction.Bonus)
			if !transaction.IsLikeSub {
				user.Balance = user.Balance + balanceChanged
			}
			if transaction.IsLikeSub {
				user.IsLikeSub = transaction.IsLikeSub
			}
			if transaction.IsInstagram {
				user.IsInstagram = transaction.IsInstagram
			}
			if transaction.IsYoutube {
				user.IsYoutube = transaction.IsYoutube
			}
			if transaction.ShareLiveStream {
				user.ShareLiveStream = transaction.ShareLiveStream
			}
			if transaction.IsReg {
				user.IsReg = transaction.IsReg
			}
			if transaction.User100App {
				user.User100App = transaction.User100App
			}

			err = userRepoInstance.UpdateWithSessionCtx(sessCtx, user)
			if err != nil {
				session.AbortTransaction(sessCtx)
				return fmt.Errorf("Cộng tiền tài khoản %s %s", user.Username, err)
			}

			transactionWallet := model.TransactionWallet{
				Code:  transaction.Code,
				Value: transaction.ValueInt,
				Bonus: transaction.Bonus,
			}
			// TODO: Ghi wallet log
			walletLog := model.WalletLog{
				Type:        constants.WallletRecharge,
				Token:       user.Token,
				Value:       balanceChanged,
				Balance:     user.Balance,
				Transaction: &transactionWallet,
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
	}

	return c.JSON(success(transaction))
}

// Detail : Chi tiết giao dịch
func (TransactionHandler) Detail(c echo.Context) (err error) {
	httpCtx := c.Request().Context()

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
	// TODO: validate ID
	transactionObjectID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		return fmt.Errorf("Mã số giao dịch không hợp lệ %s", err)
	}
	// Kiểm tra gói có tồn tại hay không ?
	transaction, err := transRepo.New(httpCtx).GetOneByID(transactionObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin giao dịch %s", err))
	}
	if !transaction.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin giao dịch"))
	}
	// Kiểm tra giao dịch có thuộc về người này hay không ?
	if c.Get("user_token") != nil {
		if transaction.Token != c.Get("user_token").(string) {
			return fmt.Errorf("Mã giao dịch không hợp lệ")
		}
	}

	return c.JSON(success(transaction))
}

// All : Tất cả giao dịch
func (TransactionHandler) All(c echo.Context) (err error) {
	httpCtx := c.Request().Context()

	type myResponse struct {
		LastOffset *primitive.ObjectID `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                 `json:"total" query:"total"`
		Data       []model.Transaction `json:"data" query:"data"`
	}
	type myRequest struct {
		Token    string                  `json:"token" query:"token"`
		Username string                  `json:"username" query:"username"`
		Code     string                  `json:"code" query:"code"`
		Offset   string                  `json:"offset" query:"offset"`
		Limit    int                     `json:"limit" query:"limit"`
		Status   []constants.TransStatus `json:"status" query:"status"`
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
	// validate
	var offsetObjectID primitive.ObjectID
	if request.Offset != "" {
		offsetObjectID, err = primitive.ObjectIDFromHex(request.Offset)
		if err != nil {
			return fmt.Errorf("Offset không hợp lệ %s", err)
		}
	}

	// Tạo response
	response := myResponse{
		Data: make([]model.Transaction, 0),
	}

	// Tao all request
	allTransactionReq := model.AllTransactionReq{
		Status:    request.Status,
		UserToken: request.Token,
		Username:  request.Username,
		Code:      request.Code,
		Offset:    offsetObjectID,
		Limit:     request.Limit,
	}

	transactions, err := transRepo.New(httpCtx).All(allTransactionReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy danh sách giao dịch %s", err))
	}

	if len(transactions) > 0 {
		response.Data = transactions
		response.Total = len(transactions)
		response.LastOffset = &transactions[len(transactions)-1].ID
	}

	return c.JSON(success(response))
}
