package model

//DataReq : mảng data
type DataReq struct {
	Appname string `json:"appname"`
	Data    struct {
		Uids   []string `json:"uids"`
		Emails []string `json:"emails"`
	} `json:"data"`
}
