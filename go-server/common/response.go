// Author: Steve Zhang
// Date: 2020/10/16 2:19 下午

package common

type Response struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message"`
	Data    ResponseData `json:"data"`
}

type ResponseCode int
type ResponseData map[string]interface{}

const (
	ResponseCodeOK ResponseCode = iota
	ResponseCodeRequestParamErr
	ResponseCodeInternalErr
	ResponseCodeAuthFailed
)

func NewOKResponse() *Response {
	return &Response{
		Code:    ResponseCodeOK,
		Message: "操作成功",
	}
}

func NewParamErrResponse(err error) *Response {
	return &Response{
		Code:    ResponseCodeRequestParamErr,
		Message: "请求参数错误：" + err.Error(),
	}
}
