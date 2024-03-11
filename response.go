package request

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type Response struct {
	*http.Response
	once    sync.Once
	content []byte

	*Job
	err error
}

func (r *Response) Error() error {
	return r.err
}

func (r *Response) ReadCloser() io.ReadCloser {
	return r.Response.Body
}

func (r *Response) Content() []byte {
	r.once.Do(func() {
		defer r.Response.Body.Close()
		r.content, _ = io.ReadAll(r.Response.Body)
	})
	return r.content
}

func (r *Response) Text() string {
	return string(r.Content())
}

func (r *Response) Json(out any) error {
	return json.Unmarshal(r.Content(), out)
}

func (r *Response) Data() map[string]any {
	m := make(map[string]any)
	r.Json(&m)
	return m
}

func (r *Response) Yaml(out any) error {
	return yaml.Unmarshal(r.Content(), out)
}

func (r *Response) Write(name string) error {
	err := os.MkdirAll(filepath.Dir(name), os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(name, r.Content(), os.ModePerm)
}

func (r *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"job":     r.Job,
		"error":   r.err,
		"content": r.Text(),
	})
}

func NewResponse(resp *http.Response) *Response {
	return &Response{Response: resp}
}

func NewError(err error) *Response {
	return &Response{err: err}
}
