package hotmail

import (
	"app/model"
	"encoding/json"
	"fmt"
)

// GetHotMail : Lấy random hotmail
func GetHotMail(requestID string, Token string) (hotmail model.HotMail, err error) {

	type myResponse struct {
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Data    model.HotMail `json:"data"`
	}

	api := cfg.API.Get("autofarmer_get_random_hotmail")

	dataRequest := map[string]interface{}{
		"token": Token,
	}

	api.SetParams(dataRequest)
	res := myResponse{}
	bResp, err := api.MakeRequest(requestID)

	if err != nil {
		err = fmt.Errorf("Lỗi tạo hoặc gọi API: %s - %s", api.Name, err)
		return
	}
	err = json.Unmarshal(bResp, &res)
	if err != nil {
		err = fmt.Errorf("Lỗi tạo hoặc gọi API: %s, Message response: %s", api.Name, res.Message)
		return
	}
	fmt.Println(res)
	// cập nhật kết qủa
	hotmail = res.Data

	return
}
