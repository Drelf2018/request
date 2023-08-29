package request

import (
	"net/http"
)

type Option func(*Job)

// 数据
func Data(data ...string) Option {
	return func(job *Job) {
		for i := 0; i < len(data); i += 2 {
			job.Data[data[i]] = data[i+1]
		}
	}
}

// 请求头
func Header(header ...string) Option {
	return func(job *Job) {
		for i := 0; i < len(header); i += 2 {
			job.Headers[header[i]] = header[i+1]
		}
	}
}

// Cookies
func Cookie(cookie ...string) Option {
	return func(job *Job) {
		for i := 0; i < len(cookie); i += 2 {
			job.Cookies[cookie[i]] = cookie[i+1]
		}
	}
}

// 构造函数
func New(method, url string, options ...Option) *Job {
	job := Job{method, url, make(M), make(M), make(M), nil}
	for _, op := range options {
		op(&job)
	}
	return &job
}

// 简化 Get 请求
func Get(url string, options ...Option) *Result {
	return New(http.MethodGet, url, options...).Request()
}

// 简化 Post 请求
func Post(url string, options ...Option) *Result {
	return New(http.MethodPost, url, options...).Request()
}
