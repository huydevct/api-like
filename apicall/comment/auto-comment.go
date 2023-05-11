package comment

import (
	"encoding/json"
	"fmt"
)

// AutoComment : Lấy comment tự hệ thống auto farmer
type AutoComment struct {
	Content  string `json:"textContent"`
	Language string `json:"language"`
}

// AutoContent : Lấy comment từ hệ thống autoComment Autofarmer
func AutoContent(requestID string, language string, quantity int) (commentContents []AutoComment, err error) {

	type myResponse struct {
		Status  int           `json:"stt"`
		Message string        `json:"msg"`
		Data    []AutoComment `json:"data"`
	}

	api := cfg.API.Get("autocomment_random_comment_content")

	dataRequest := map[string]interface{}{
		"language": language,
		"limit":    quantity,
	}

	api.SetParams(dataRequest)

	res := myResponse{}
	bResp, err := api.MakeRequest(requestID)
	if err != nil {
		err = fmt.Errorf("Lỗi tạo hoặc gọi API: %s - %s", api.Name, err)
		return
	}

	err = json.Unmarshal(bResp, &res)
	if err != nil || res.Status != 200 {
		err = fmt.Errorf("Lỗi tạo hoặc gọi API: %s, Message response: %s, %s", api.Name, res.Message, err)
		return
	}
	// cập nhật kết qủa
	commentContents = res.Data

	return
}
