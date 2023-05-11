package service

import (
	"encoding/json"
	"fmt"
)

// ESMS : create service req
type ESMS struct {
	Content string
	Phone   string
}

// SendOTPEsms : Gửi mã ESMS
func SendOTPEsms(requestID string, req ESMS) (err error) {

	type myResponse struct {
		CodeResult      string `json:"CodeResult"`
		CountRegenerate int    `json:"CountRegenerate"`
		SMSID           string `json:"SMSID"`
	}

	api := cfg.API.Get("esms_send_otp")

	dataRequest := map[string]interface{}{
		"ApiKey":    cfg.Other.Get("esms-api-key"),
		"SecretKey": cfg.Other.Get("esms-secret-key"),
		"Content":   req.Content,
		"Phone":     req.Phone,
		"SmsType":   "2",
		"Brandname": "Verify",
	}

	api.SetParams(dataRequest)

	res := myResponse{}
	bResp, err := api.MakeRequest(requestID)
	if err != nil {
		err = fmt.Errorf("Lỗi tạo hoặc gọi API: %s - %s", api.Name, err)
		return
	}

	err = json.Unmarshal(bResp, &res)
	if err != nil || res.CodeResult != "100" {
		err = fmt.Errorf("Lỗi tạo hoặc gọi API: %s", api.Name)
		return
	}

	return
}
