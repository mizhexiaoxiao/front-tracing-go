package common

type Response struct {
	Code int64       `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}
