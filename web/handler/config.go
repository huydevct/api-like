package handler

import (
	"fmt"
	"net/http"

	"app/constants"
	"app/model"

	"github.com/labstack/echo/v4"

	cloneRepo "app/repo/mongo/clone"
)

// Config init
type Config struct{}

// NewConfig : tạo đối tường Config
func NewConfig() *Config {
	return &Config{}
}

func (h Config) GetConfig(c echo.Context) (err error) {
	httpCtx := c.Request().Context()

	type myRequest struct {
		Token   string `json:"token" query:"token" validate:"required"`
		Action  string `json:"action" query:"action" validate:"required"`
		AppName string `json:"app_name" query:"app_name" validate:"required"`
		// GetConfig
		DeviceInfo model.DeviceInfo `json:"device_info" query:"device_info"` // report device info
		// UpdateClone
		CloneInfo model.CloneInfo `json:"clone_info" query:"clone_info"` // data clone đã mã hoá base64
	}

	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}
	userInfo, err := getUserInfoByToken(request.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Get user info by token %s %s", request.Token, err))
	}
	if !userInfo.IsExists() {
		return fmt.Errorf("Token %s is invalid", request.Token)
	}
	// Lấy thông tin device by DeviceID
	if request.DeviceInfo.PCName == "" {
		return fmt.Errorf("PC name is not allow empty")
	}

	err = validateDeviceIDBelongToken(request.DeviceInfo.PCName, request.Token)
	if err != nil {
		return fmt.Errorf("DeviceID is not valid: %s", err)
	}

	// Dựa theo action, trả về response tuơng ứng
	switch request.Action {

	case "GetClone":

		cloneRepoInstance := cloneRepo.New(httpCtx)
		cloneRequest := cloneRepo.FindCloneLiveByDeviceReq{
			Token:   request.Token,
			PCName:  request.DeviceInfo.PCName,
			AppName: request.AppName,
		}
		// Get one live clone by token and device & appname
		clone, err := cloneRepoInstance.FindOneLiveAndUpdateByRequest(cloneRequest)
		if err != nil {
			return err
		}
		if clone.IsExists() {
			// Encode base64
			return c.JSON(success(clone))
		}

		return c.JSON(success(nil))

	case "UpdateAliveStatus":
		// validate
		cloneReq := request.CloneInfo
		if err != nil {
			return err
		}
		if cloneReq.ID.IsZero() {
			return fmt.Errorf("CloneID is not allow empty")
		}

		cloneRepoInstance := cloneRepo.New(httpCtx)

		// Tùy theo trạng thái mong muốn của mobile, cập nhật clone thành trạng thái tương ứng
		if cloneReq.AliveStatus == constants.CloneStored {
			_, err := cloneRepoInstance.UpdateStored(cloneReq.ID, request.Token, request.DeviceInfo.PCName)
			if err != nil {
				return fmt.Errorf("Update clone %s to stored: %s", cloneReq.ID, err)
			}

		} else if cloneReq.AliveStatus == constants.CloneCheckpoint {
			_, err = cloneRepoInstance.UpdateCheckpoint(cloneReq.ID, request.Token, request.DeviceInfo.PCName)
			if err != nil {
				return fmt.Errorf("Update clone %s checkpoint: %s", cloneReq.ID, err)
			}
		}

		return c.JSON(success(nil))
	case "UpdateCloneInfo":
		// validate
		cloneReq := request.CloneInfo
		if err != nil {
			return fmt.Errorf("Clone info is invalid: %s", err)
		}
		if cloneReq.ID.IsZero() {
			return fmt.Errorf("CloneID is not allow empty")
		}

		cloneRepoInstance := cloneRepo.New(httpCtx)
		_, err := cloneRepoInstance.UpdateInfo(cloneReq)
		if err != nil {
			return fmt.Errorf("Update clone %s to stored: %s", cloneReq.ID, err)
		}

		return c.JSON(success(nil))

	default:
		fmt.Println("Not map action: ", request.Action)
		return c.JSON(success(nil))
	}
}
