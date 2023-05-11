package api4

import (
	"app/model"
	"encoding/json"
	"fmt"
)

// GetCloneByIDReq : ..
type GetCloneByIDReq struct {
	Token   string
	CloneID string
}

// GetCloneByID : Get clone from api4
func GetCloneByID(requestID string, data GetCloneByIDReq) (result model.CloneOldInfo, err error) {

	type myResponse struct {
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    model.CloneOldInfo `json:"data"`
	}

	api := cfg.API.Get("api4_autofarmer_clone")
	{
		// set body
		dataRequest := map[string]interface{}{
			"action":  "GetDetail",
			"appname": "facebook",
			"token":   data.Token,
			"id":      data.CloneID,
			"info": map[string]string{
				"IMEI":     "xxx",
				"DeviceID": "xxx",
			},
		}

		api.SetParams(dataRequest)
	}

	res := myResponse{}
	bResp, err := api.MakeRequest(requestID)
	if err != nil {
		err = fmt.Errorf("Call api: %s - %s", api.Name, err)
		return
	}

	err = json.Unmarshal(bResp, &res)
	if err != nil {
		err = fmt.Errorf("Call api: %s, Message response: %s", api.Name, res.Message)
		return
	}
	if !res.Data.IsExists() {
		err = fmt.Errorf("Call api: %s, Message response: %s", api.Name, "Not found clone")
		return
	}

	// cập nhật kết qủa
	result = res.Data

	return
}
