package request

import (
	"encoding/json"
	"net/http"
)

// 结果
type Result struct {
	// 原结果
	*http.Response
	// 内容
	Content []byte
	// 错误
	Error error
}

// 是否有错
func (res *Result) hasErr(err error) bool {
	if err != nil {
		res.Error = err
		return true
	}
	return false
}

// 解析结果为文本
func (res *Result) Text() string {
	return string(res.Content)
}

// 解析结果为 json
func (res *Result) Json(data any) error {
	return json.Unmarshal(res.Content, data)
}
