package models

type RestResult struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Result  interface{} `json:"data"`
}
