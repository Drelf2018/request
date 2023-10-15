package request

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var ErrType = errors.New("数据类型错误")

type M map[string]string

// gorm 读取
//
// 参考: https://github.com/go-gorm/datatypes/blob/master/json_map.go
func (m *M) Scan(val any) error {
	if val == nil {
		*m = make(M)
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return ErrType
	}
	return json.Unmarshal(ba, m)
}

func (m M) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// GormDataType gorm common data type
func (m M) GormDataType() string {
	return "jsonmap"
}

// GormDBDataType gorm db data type
func (M) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	case "sqlserver":
		return "NVARCHAR(MAX)"
	}
	return ""
}

type Values interface {
	Set(string, string)
	Add(string, string)
}

// url.Values 赋值
func (m M) CopyTo(vs Values) {
	for k, v := range m {
		vs.Set(k, v)
	}
}

// url.Values 添加
func (m M) AddTo(vs Values) {
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

var HEADERS = M{
	"Accept-Language":           "zh-CN,zh;q=0.9",
	"Accept-Encoding":           "gzip, deflate, br",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.37",
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
	// Client
	Client http.Client `form:"-" yaml:"-" json:"-" gorm:"-"`
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
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.37")
	job.Headers.CopyTo(req.Header)

	// 添加 Cookies
	if len(job.Cookies) != 0 {
		if job.Client.Jar == nil {
			job.Client.Jar = job.Cookies
		} else {
			cookieURL, _ := url.Parse(job.Url)
			job.Client.Jar.SetCookies(cookieURL, job.Cookies.Cookies(nil))
		}
	}

	// 正式请求
	resp, err := job.Client.Do(req)
	if r.hasErr(err) {
		return
	}

	// 读取内容
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if r.hasErr(err) {
		return
	}
	return &Result{resp, body, nil}
}
