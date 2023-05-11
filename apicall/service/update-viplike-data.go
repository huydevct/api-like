package service

import (
	"encoding/json"
	"fmt"
)

// UpdateViplikeDataReq ..
type UpdateViplikeDataReq struct {
	ServiceCode string
	P1          bool
	P2          bool
	P3          bool
	P4          bool
	P5          bool
}

// UpdateViplikeData : Cập nhật thông tin data với gói viplikeService
func UpdateViplikeData(requestID string, req UpdateViplikeDataReq) (err error) {

	type myResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	api := cfg.API.Get("autolike_detail_update_viplike_data")

	dataRequest := map[string]interface{}{
		"service_code": req.ServiceCode,
		"p1":           req.P1,
		"p2":           req.P2,
		"p3":           req.P3,
		"p4":           req.P4,
		"p5":           req.P5,
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

	return
}
