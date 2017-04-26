package handlers

import "net/http"

type (
	RespCode int
	RespMsg  string
)

const (
	UNKNOWN RespCode = -1
	OK RespCode = iota
	DATABASE_ERROR
	DATABASE_SAVE_ERROR
	DATABASE_FIND_ERROR
	DATABASE_NOT_FOUND
	DATA_BIND_ERROR
)

var ResponseMessage = map[RespCode]RespMsg{
	UNKNOWN:             "未知错误.",
	OK:                  "OK",
	DATA_BIND_ERROR:     "数据绑定错误.",
	DATABASE_ERROR:      "数据库错误.",
	DATABASE_SAVE_ERROR: "保存数据错误.",
	DATABASE_FIND_ERROR: "查找数据错误.",
	DATABASE_NOT_FOUND:  "没找到该数据.",
}

type CommonResponse struct {
	Code    RespCode    `json:"code"`
	Message RespMsg     `json:"message"`
	Result  interface{} `json:"result"`
}

func NewResp(code RespCode, data ...interface{}) *CommonResponse {
	cr := &CommonResponse{Code: code}
	// 如果code不存在,将会是空的,那么就是完全自定义code和message
	cr.Message = ResponseMessage[code]
	l := len(data)
	if l > 0 {
		switch l {
		// 只有一个参数时,判断是否RespMsg类型,
		// 如果是,将作为Message传递,如果不是,将作为Result传递
		case 1:
			if msg, ok := data[0].(RespMsg); ok {
				// 合并默认消息
				cr.Message += msg
			} else {
				cr.Result = data[0]
			}
		// 有两个参数时,第一个参数作为Message,第二个参数作为Result
		case 2:
			if msg, ok := data[0].(RespMsg); ok {
				// 合并默认消息
				cr.Message += msg
			}
			cr.Result = data[1]
		}
	}
	return cr
}

func (c *CommonResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func GetResponseMessage(code RespCode) string {
	return string(ResponseMessage[code])
}
