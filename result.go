package request

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

// 结果
type Result struct {
	// 原结果
	*http.Response
	// 内容
	Content []byte
	// 错误
	err error
}

func (r *Result) Error() error {
	return r.err
}

// 解析结果为文本
func (r *Result) Text() string {
	return string(r.Content)
}

// 解析结果为 json
func (res *Result) Json(data any) error {
	return json.Unmarshal(res.Content, data)
}

// 重新获取 io.Reader
func (r *Result) Reader() *bytes.Reader {
	return bytes.NewReader(r.Content)
}

// 写出文件
func (r *Result) Write(name string, perm os.FileMode) error {
	return os.WriteFile(name, r.Content, perm)
}
