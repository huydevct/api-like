package handler

import (
	"app/common/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	validator "gopkg.in/go-playground/validator.v9"
)

var cfg = config.GetConfig()

// NewValidator ..
func NewValidator() *MyValidator {
	return &MyValidator{validator: validator.New()}
}

// MyValidator ..
type MyValidator struct {
	validator *validator.Validate
}

// Validate ..
func (cv *MyValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// GetRequestID ..
func GetRequestID(c echo.Context) (requestID string) {
	if c.Get("reqID") != nil {
		requestID = c.Get("reqID").(string)
	}
	return
}

// Health func
func Health(c echo.Context) (err error) {
	return c.JSON(Success(nil))
}

// ResponseContent struct
type ResponseContent struct {
	Code             int         `json:"code"`
	Message          string      `json:"message"`
	Data             interface{} `json:"data"`
	DataRaw          interface{} `json:"data_raw,omitempty"`
	CodeMessage      string      `json:"code_message,omitempty"`
	CodeMessageValue string      `json:"code_message_value,omitempty"`
}

// Error func
func Error(err error, c echo.Context) {
	code := http.StatusBadRequest
	msg := http.StatusText(code)
	codeMessage := ""
	codeMessageValue := ""

	if httpError, ok := err.(*echo.HTTPError); ok {
		code = httpError.Code
		msg = fmt.Sprintf("%v", httpError.Message)
	}

	if cfg.Debug {
		msg = err.Error()

		// get code message
		if strings.Contains(msg, "|") {
			temp := strings.Split(msg, "|")

			// errPlatform: (err | code_message_key | code_message_value)
			msg = strings.TrimSpace(temp[0])
			codeMessage = strings.TrimSpace(temp[1])

			if len(temp) > 2 {
				codeMessageValue = strings.TrimSpace(temp[2])
			}
		}
	}

	log.Error(err)

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			c.NoContent(code)
		} else {
			c.JSON(code, &ResponseContent{
				Code:             code,
				Message:          msg,
				CodeMessage:      codeMessage,
				CodeMessageValue: codeMessageValue,
			})
		}
	}
}

// Success func
func Success(data interface{}) (int, ResponseContent) {
	return http.StatusOK, ResponseContent{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	}
}

// SuccessRaW func
func SuccessRaw(data interface{}, dataRaw interface{}) (int, ResponseContent) {
	return http.StatusOK, ResponseContent{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
		DataRaw: dataRaw,
	}
}

// NotFound func
func NotFound(message string) (int, ResponseContent) {
	return http.StatusOK, ResponseContent{
		Code:    http.StatusNotFound,
		Message: message,
	}
}
