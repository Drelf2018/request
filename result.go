package request

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
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
func (r *Result) Json(out any) error {
	return json.Unmarshal(r.Content, out)
}

// 解析结果为 yaml
func (r *Result) Yaml(out any) error {
	return yaml.Unmarshal(r.Content, out)
}

// 重新获取 io.Reader
func (r *Result) Reader() *bytes.Reader {
	return bytes.NewReader(r.Content)
}

// 写出文件
func (r *Result) Write(name string, perm os.FileMode) error {
	path := filepath.Dir(name)
	err := os.MkdirAll(path, perm)
	if err != nil {
		return err
	}
	return os.WriteFile(name, r.Content, perm)
}

// 写出文件到指定目录
func (r *Result) WriteToDir(path string, perm os.FileMode) error {
	base := filepath.Base(r.Response.Request.URL.Path)
	return r.Write(filepath.Join(path, base), perm)
}

// 写出到指定层级下
func (r *Result) WriteWithPath(root string, perm os.FileMode) error {
	return r.Write(filepath.Join(root, r.Response.Request.URL.Path), perm)
}

type JobResult struct {
	Job     Job            `json:"job"`
	Error   error          `json:"error,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
	Content string         `json:"content,omitempty"`
}
