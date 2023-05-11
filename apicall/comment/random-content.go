package comment

import (
	"app/model"
	"encoding/json"
	"fmt"
)

// RandomContent : Lấy random comment content theo commentID
func RandomContent(requestID string, commentID string, quantity int) (commentContents []model.CommentContent, err error) {

	type myResponse struct {
		Code    int                    `json:"code"`
		Message string                 `json:"message"`
		Data    []model.CommentContent `json:"data"`
	}

	api := cfg.API.Get("autolike_random_comment_content")

	dataRequest := map[string]interface{}{
		"comment_id": commentID,
		"quantity":   quantity,
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
	commentContents = res.Data

	return
}
