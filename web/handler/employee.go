package handler

import (
	"fmt"
	"net/http"

	"app/constants"
	"app/model"
	"app/utils"

	employeeRepo "app/repo/mongo/employee"
	tokenRepo "app/repo/mongo/token"
	userRepo "app/repo/mongo/user"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EmployeeHandler : struct init employee
type EmployeeHandler struct{}

// NewEmployeeHandler : Tạo mới 1 đối tượng employee handler
func NewEmployeeHandler() *EmployeeHandler {
	return &EmployeeHandler{}
}

// Register : Tạo mới user quản lý
// role admin trong hệ thống
func (EmployeeHandler) Register(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		Username string `json:"username" query:"username" validate:"required"`
		FullName string `json:"fullname" query:"fullname" validate:"required"`
		Password string `json:"password" query:"password" validate:"required,min=6"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// TODO: validate
	employeeRepoInstance := employeeRepo.New(httpCtx)
	{
		// TODO: Kiểm tra đã có tài khoản đăng ký với username này hay chưa ?
		user, err := employeeRepoInstance.GetOneByUsername(request.Username)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Kiểm tra nhân viên tồn tại theo username %s", err))
		}
		if user.IsExists() {
			return fmt.Errorf("Tài khoản đã tồn tại %s", request.Username)
		}
	}
	// Tạo bsrypt cho password
	// bscryptPassword, err := bscryptPassword()
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Tạo mật khẩu lỗi %s", err))
	// }
	// TODO: Tạo user model
	employeeInfo := model.EmployeeInfo{
		FullName:    request.FullName,
		Username:    request.Username,
		Password:    generateHash(request.Password),
		Status:      constants.Active,
		Permissions: make([]constants.PermissionCommand, 0),
	}
	result, err := employeeRepoInstance.Insert(employeeInfo)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Đăng ký nhân viên %s", err))
	}

	return c.JSON(success(result))
}

// Login : Đăng nhập
// Limit số lần đăng nhập thất bại
func (h EmployeeHandler) Login(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		Username string `json:"username" query:"username" validate:"required"`
		Password string `json:"password" query:"password" validate:"required,min=6"`
		DeviceID string `json:"device_id" query:"device_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// TODO: Đếm số lần login fail theo phone
	numLoginFail, errCount := countNumLoginFail(request.Username)
	if errCount != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Đếm số lần đăng nhập thất bại %s", err))
	}
	numLoginFail++

	if numLoginFail > 3 {
		return echo.NewHTTPError(http.StatusLocked, "Tài khoản của bạn đang tạm thời bị khoá trong vòng 60s do đăng nhập sai quá nhiều lần. Vui lòng chờ và thực hiện lại.")
	}

	// Lấy thông tin employee by phone
	employee, err := employeeRepo.New(httpCtx).GetOneByUsername(request.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Tài khoản hoặc mật khẩu không đúng. Vui lòng thử lại %s", err))
	}
	if !employee.IsExists() {
		return fmt.Errorf("Tài khoản hoặc mật khẩu không đúng. Vui lòng thử lại")
	}
	// TODO : Check password match

	if !checkPasswordMatch(request.Password, employee.Password) {
		// Increase num login fail
		ttl := 60
		setNumLoginFail(request.Username, numLoginFail, ttl)
		return fmt.Errorf("Tài khoản hoặc mật khẩu không đúng. Vui lòng thử lại")
	}
	// TODO: Lưu token với period 7 ngày
	tokenModel := model.Token{
		ExpiredAfterSecond: 7 * 24 * 3600,
		// token info
		UserID:    employee.ID,
		Type:      "employee",
		DeviceID:  request.DeviceID,
		UserAgent: c.Request().UserAgent(),
		RemoteIP:  c.RealIP(),
	}

	token, err := GenTokenByUserIDDeviceID(tokenModel)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(success(map[string]interface{}{
		"token":   token,
		"user_id": employee.ID,
		"phone":   employee.Username,
	}))
}

// Logout : Đăng xuất
func (EmployeeHandler) Logout(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	token := c.Get("token").(string)

	userType := "employee"
	_, err = tokenRepo.New(httpCtx).ClearToken(token, userType)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Token không hợp lệ")
	}

	return c.JSON(success(nil))
}

// AllUser : Lấy tất cả người dùng theo offset, limit
func (EmployeeHandler) AllUser(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myResponse struct {
		LastOffset *primitive.ObjectID `json:"last_offset,omitempty" query:"last_offset"`
		Total      int                 `json:"total" query:"total"`
		Data       []model.UserInfo    `json:"data" query:"data"`
	}
	type myRequest struct {
		Status   []constants.CommonStatus `json:"status" query:"status"`
		Username string                   `json:"username" query:"username"`
		Offset   string                   `json:"offset" query:"offset"`
		Limit    int                      `json:"limit" query:"limit"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// tạo response
	response := myResponse{
		Data: make([]model.UserInfo, 0),
	}
	allUserReq := model.AllUserReq{
		Status:   request.Status,
		Username: request.Username,
		Offset:   primitive.NilObjectID,
		Limit:    request.Limit,
	}
	if request.Offset != "" {
		allUserReq.Offset, err = primitive.ObjectIDFromHex(request.Offset)
		if err != nil {
			return fmt.Errorf("Offset không đúng format: %s", err)
		}
	}
	response.Data, err = userRepo.New(httpCtx).All(allUserReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy tất cả thông tin người dùng%s", err))
	}
	if len(response.Data) > 0 {
		// cập nhật last offset
		response.Total = len(response.Data)
		response.LastOffset = &response.Data[len(response.Data)-1].ID
	}

	return c.JSON(success(response))
}

// DetailUserByPhone : Lấy thông tin người dùng theo số điện thoại
func (EmployeeHandler) DetailUserByPhone(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		Phone string `json:"phone" query:"phone" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// TODO: validate
	{
		// validate phone
		request.Phone, err = utils.PhoneValidate(request.Phone)
		if err != nil {
			return fmt.Errorf("Số điện thoại không hợp lệ %s", err)
		}
	}
	user, err := userRepo.New(httpCtx).GetOneByPhone(request.Phone)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng theo số điện thoại %s", err))
	}
	if !user.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}

	return c.JSON(success(user))
}

// DetailUserByID : Lấy thông tin người dùng theo UserID
func (EmployeeHandler) DetailUserByID(c echo.Context) (err error) {

	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		UserID string `json:"user_id" query:"user_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	userObjectID, err := primitive.ObjectIDFromHex(request.UserID)
	if err != nil {
		return fmt.Errorf("Mã số khách hàng không hợp lệ %s", err)
	}
	user, err := userRepo.New(httpCtx).GetOneByID(userObjectID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng theo ID %s", err))
	}
	if !user.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}

	return c.JSON(success(user))
}

// UpdateUser : Chỉnh sửa thông tin  người dùng
// . Thông tin
// . Status
// . Quyền
func (EmployeeHandler) UpdateUser(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		UserID         string                        `json:"user_id" query:"user_id" validate:"required"`
		Fullname       string                        `json:"fullname" query:"fullname"`
		AccountName    string                        `json:"account_name" query:"account_name"`
		BankName       string                        `json:"bank_name" query:"bank_name"`
		BankNumber     string                        `json:"bank_number" query:"bank_number"`
		Status         *constants.CommonStatus       `json:"status,omitempty" query:"status,omitempty"`
		Permissions    []constants.PermissionCommand `json:"permissions" query:"permissions"`
		EnableAPI8     *bool                         `json:"enable_api8" query:"enable_api8"`
		IsDeCheckpoint *bool                         `json:"is_decheckpoint" query:"is_decheckpoint"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// TODO: validate
	{
		// validate permissions
		if len(request.Permissions) > 0 {
			for _, permissionCommand := range request.Permissions {
				if !constants.UserPermissions[permissionCommand] {
					return fmt.Errorf("Quyền %s không hợp lệ", permissionCommand)
				}
			}
		}
	}
	{
		// validate status
		if request.Status != nil {
			if *request.Status != constants.Active && *request.Status != constants.Pause {
				return fmt.Errorf("Trạng thái %s không hợp lệ", *request.Status)
			}
		}
	}
	userRepoInstance := userRepo.New(httpCtx)
	// Kiểm tra tài khoản có tồn tài hay không ?
	userObjectID, err := primitive.ObjectIDFromHex(request.UserID)
	if err != nil {
		return fmt.Errorf("Mã số khách hàng không hợp lệ %s", err)
	}
	oldUser, err := userRepoInstance.GetOneByID(userObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng %s", err))
	}
	if !oldUser.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}
	{
		// update name
		if request.Fullname != "" {
			oldUser.Fullname = request.Fullname
		}
	}
	{
		// update permission
		if len(request.Permissions) > 0 {
			oldUser.Permissions = request.Permissions
		}
	}
	{
		// update AccountName
		if request.AccountName != "" {
			oldUser.AccountName = request.AccountName
		}
	}
	{
		// update BankName
		if request.BankName != "" {
			oldUser.BankName = request.BankName
		}
	}
	{
		// update BankNumber
		if request.BankNumber != "" {
			oldUser.BankNumber = request.BankNumber
		}
	}
	{
		// update status
		if request.Status != nil {
			oldUser.Status = *request.Status
		}
	}
	{
		// update enable api8
		if request.EnableAPI8 != nil {
			oldUser.EnableAPI8 = *request.EnableAPI8
		}
	}
	{
		// update gỡ checkpoint
		if request.IsDeCheckpoint != nil {
			oldUser.IsDeCheckpoint = *request.IsDeCheckpoint
		}
	}
	// Cập nhật status người dùng
	err = userRepoInstance.Update(oldUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật thông tin người dùng %s", err))
	}

	return c.JSON(success(oldUser))
}

// DeleteUser : Xóa người dùng theo user id
func (EmployeeHandler) DeleteUser(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		UserID string `json:"user_id" query:"user_id" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	userRepoInstance := userRepo.New(httpCtx)
	// Kiểm tra tài khoản có tồn tài hay không ?
	userObjectID, err := primitive.ObjectIDFromHex(request.UserID)
	if err != nil {
		return fmt.Errorf("Mã số khách hàng không hợp lệ %s", err)
	}
	oldUser, err := userRepoInstance.GetOneByID(userObjectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng %s", err))
	}
	if !oldUser.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}
	// Cập nhật status người dùng
	oldUser.Status = constants.Delete
	err = userRepoInstance.Update(oldUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật thông tin người dùng %s", err))
	}

	return c.JSON(success(nil))
}

// SetUserToken : Cài đặt token số đẹp
// . 11 ký tự
// . Kiểm tra trùng token
func (EmployeeHandler) SetUserToken(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		UserID    primitive.ObjectID `json:"user_id" query:"user_id" validate:"required"`
		UserToken string             `json:"user_token" query:"user_token" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	// Validate token
	if len(request.UserToken) > 11 {
		return fmt.Errorf("Token tối đa 11 ký tự. Vui lòng thử lại sau")
	}

	// Kiểm tra token này đã tồn tại hay chưa
	userRepoInstance := userRepo.New(httpCtx)

	oldUser, err := userRepoInstance.GetOneByToken(request.UserToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng theo token %s", err))
	}
	if oldUser.IsExists() {
		return fmt.Errorf("Token %s đã đuợc đăng ký. Vui lòng thử lại sau", request.UserToken)
	}
	// Kiểm tra tài khoản có tồn tài hay không ?
	oldUser, err = userRepoInstance.GetOneByID(request.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng %s", err))
	}
	if !oldUser.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}
	// Cập nhật token người dùng
	oldUser.Token = request.UserToken
	err = userRepoInstance.Update(oldUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Cập nhật thông tin người dùng %s", err))
	}

	return c.JSON(success(nil))
}

// DetailByToken : Lấy thông tin người dùng theo token
func (h EmployeeHandler) GetTokenNewByOld(c echo.Context) (err error) {
	httpCtx := c.Request().Context()

	type myRequest struct {
		TokenOld string `json:"token_old" query:"token_old" validate:"required"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	//
	user, err := userRepo.New(httpCtx).GetTokenNewByOld(request.TokenOld)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Lấy thông tin người dùng theo token %s", err))
	}
	if !user.IsExists() {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Không tìm thấy thông tin tài khoản"))
	}

	return c.JSON(success(user))
}
