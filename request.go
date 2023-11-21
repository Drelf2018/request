package request

import (
	"net/http"
)

type Option func(*Job)

// 数据
func Data(data ...string) Option {
	return func(job *Job) { job.Data.Inserts(data...) }
}

func Datas(data M) Option {
	return func(job *Job) { job.Data = data }
}

// 请求头
func Header(header ...string) Option {
	return func(job *Job) { job.Headers.Inserts(header...) }
}

func Headers(headers M) Option {
	return func(job *Job) { job.Headers = headers }
}

func Referer(referer string) Option {
	return func(job *Job) {
		HEADERS.SetTo(job.Headers)
		job.Headers["Referer"] = referer
	}
}

// Cookies
func Cookie(cookie ...string) Option {
	return func(job *Job) { job.Cookies.Inserts(cookie...) }
}

func Cookies(cookies M) Option {
	return func(job *Job) { job.Cookies = cookies }
}

func CookieString(s string) Option {
	return func(job *Job) { job.ParseCookies(s) }
}

// Client
func Client(client http.Client) Option {
	return func(job *Job) { job.Client = client }
}

// 构造函数
func New(method, url string, options ...Option) *Job {
	job := &Job{Method: method, Url: url, Data: make(M), Headers: make(M), Cookies: make(M)}
	for _, op := range options {
		op(job)
	}
	return job
}

// 简化 Get 请求
func Get(url string, options ...Option) *Result {
	return New(http.MethodGet, url, options...).Request()
}

// 简化 Post 请求
func Post(url string, options ...Option) *Result {
	return New(http.MethodPost, url, options...).Request()
}

// 简化结构体获取
func Json[T any](job *Job) (out T) {
	job.Request().Json(&out)
	return
}
