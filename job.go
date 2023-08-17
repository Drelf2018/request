package request

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var ErrType = errors.New("数据类型错误")

type M map[string]string

// gorm 读取
func (m *M) Scan(val any) error {
	if val, ok := val.([]byte); ok {
		return json.Unmarshal(val, m)
	}
	return ErrType
}

func (m M) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// url.Values 赋值
func (m M) CopyTo(vs interface{ Add(string, string) }) {
	for k, v := range m {
		vs.Add(k, v)
	}
}

// CookieJar 实现
func (m M) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for _, c := range cookies {
		m[c.Name] = c.Value
	}
}

func (m M) Cookies(u *url.URL) (r []*http.Cookie) {
	r = make([]*http.Cookie, len(m))
	i := 0
	for k, v := range m {
		r[i] = &http.Cookie{Name: k, Value: v}
		i++
	}
	return
}

// 请求任务
type Job struct {
	// GET 或 POST
	Method string `form:"method" yaml:"method" json:"method"`
	// 请求网址
	Url string `form:"url" yaml:"url" json:"url"`
	// 附带数据
	Data M `form:"data" yaml:"data" json:"data"`
	// 请求头
	Headers M `form:"headers" yaml:"headers" json:"headers"`
	// Cookies
	Cookies M `form:"cookies" yaml:"cookies" json:"cookies"`
}

// 发送请求
func (job *Job) Request() (r *Result) {
	// 添加 POST 参数
	ploady := make(url.Values)
	if job.Method == http.MethodPost {
		job.Data.CopyTo(ploady)
	}

	// 新建请求
	req, err := http.NewRequest(job.Method, job.Url, strings.NewReader(ploady.Encode()))
	if r.hasErr(err) {
		return
	}

	// 添加 GET 参数
	if job.Method == http.MethodGet {
		q := req.URL.Query()
		job.Data.CopyTo(q)
		req.URL.RawQuery = q.Encode()
	}

	// 添加请求头
	job.Headers.CopyTo(req.Header)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 新建客户端
	client := &http.Client{Jar: job.Cookies}
	resp, err := client.Do(req)
	if r.hasErr(err) {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if r.hasErr(err) {
		return
	}
	return &Result{resp.StatusCode, body, nil}
}
