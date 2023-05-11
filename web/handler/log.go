package handler

import (
	"app/common/gstuff/glog"
	"app/utils"
	"time"

	logsRepo "app/repo/mongo/logs"

	"github.com/labstack/echo/v4"
)

// LogHandler init
type LogHandler struct{}

// NewLogHandler : Tạo đối tượng device handler
func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (LogHandler) Write(c echo.Context) (err error) {
	httpCtx := c.Request().Context()
	//
	type myRequest struct {
		SearchKey    string `json:"search_key,omitempty" query:"search_key,omitempty"`
		ComputerName string `json:"computer_name,omitempty" query:"computer_name,omitempty"`
		AppName      string `json:"app_name,omitempty" query:"app_name,omitempty"`
		AppVersion   string `json:"app_version,omitempty" query:"app_version,omitempty"`
		Token        string `json:"token,omitempty" query:"token,omitempty"`
		Device       string `json:"device,omitempty" query:"device,omitempty"`
		Log          string `json:"log,omitempty" query:"log,omitempty"`
	}
	request := new(myRequest)
	if err = c.Bind(request); err != nil {
		return
	}
	if err = c.Validate(request); err != nil {
		return
	}

	logData := utils.ConvertStructToMap(request)
	logData["remote-ip"] = c.RealIP()
	logData["host"] = c.Request().Host
	logData["uri"] = c.Request().RequestURI

	logRepoInstance := logsRepo.New(httpCtx)
	logRepoInstance.Insert(map[string]interface{}{
		"action": "ClientSendLog",
		"time":   time.Now(),
		"data":   logData,
	})

	glog.Send(logData)

	return c.JSON(success(logData))
}
