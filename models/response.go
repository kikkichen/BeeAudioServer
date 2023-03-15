package model

type ResponseBody struct {
	OK      int         `json:"ok"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
