package service

import (
	"app/model"
	"encoding/json"
	"fmt"
)

// DetailByCode : Lấy chi tiết serice theo service code
func DetailByCode(requestID string, serviceCode string) (services model.Service, err error) {

	type myResponse struct {
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Data    model.Service `json:"data"`
	}

	api := cfg.API.Get("autolike_detail_service_by_code")

	dataRequest := map[string]interface{}{
		"service_code": serviceCode,
	}

	api.SetParams(dataRequest)

	res := myResponse{}
	bResp, err := api.MakeRequest(requestID)
	if err != nil {
		err = fmt.Errorf("Lỗi tạo hoặc gọi API: %s - %s", api.Name, err)
		return
	}

	err = json.Unmarshal(bResp, &res)
	if err != nil || res.Code != 200 {
		err = fmt.Errorf("Lỗi tạo hoặc gọi API: %s, Message response: %s", api.Name, res.Message)
		return
	}
	// cập nhật kết qủa
	services = res.Data

	return
}
