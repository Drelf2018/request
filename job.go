package request

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

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

// 解析 cookies
func (job *Job) ParseCookies(s string) {
	for _, cookie := range strings.Split(s, ";") {
		data := strings.Split(cookie, "=")
		job.Cookies.Insert(data[0], data[1])
	}
}

// 发送请求
func (job *Job) Request() *Result {
	// 大写请求方式
	method := strings.ToUpper(job.Method)

	// 添加 POST 参数
	payload := make(url.Values)
	if method == http.MethodPost {
		job.Data.SetTo(payload)
	}

	// 新建请求
	req, err := http.NewRequest(method, job.Url, strings.NewReader(payload.Encode()))
	if err != nil {
		return &Result{err: err}
	}

	// 添加 GET 参数
	if method == http.MethodGet {
		q := req.URL.Query()
		job.Data.SetTo(q)
		req.URL.RawQuery = q.Encode()
	}

	// 添加请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", HEADERS["User-Agent"])
	job.Headers.SetTo(req.Header)

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
	if err != nil {
		return &Result{err: err}
	}

	// 读取内容
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Result{err: err}
	}
	return &Result{resp, body, nil}
}

func (job Job) Fetch(url string) *Result {
	job.Url = url
	return job.Request()
}

func (job *Job) Error() error {
	return job.Request().err
}

func (job Job) Test() *JobResult {
	result := job.Request()
	r := &JobResult{Job: job, Error: result.err}
	if result.err == nil {
		m := make(map[string]any)
		if result.Json(&m) == nil {
			r.Data = m
		} else {
			r.Content = result.Text()
		}
	}
	return r
}
