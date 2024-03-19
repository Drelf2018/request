package request

import (
	"context"
	"net/http"
	"net/url"
)

// 请求任务
type Job struct {
	Method string      `form:"method" yaml:"method" json:"method"`
	Url    string      `form:"url"    yaml:"url"    json:"url"`
	Data   M           `form:"data"   yaml:"data"   json:"data"`
	Query  M           `form:"query"  yaml:"query"  json:"query"`
	Header M           `form:"header" yaml:"header" json:"header"`
	Cookie *Cookie     `form:"cookie" yaml:"cookie" json:"cookie"`
	Client http.Client `form:"-"      yaml:"-"      json:"-" gorm:"-"`
}

func (job *Job) SetURL(url string) *Job {
	job.Url = url
	return job
}

func (job *Job) SetData(data M) *Job {
	job.Data = data
	return job
}

func (job *Job) SetQuery(query M) *Job {
	job.Query = query
	return job
}

func (job *Job) SetHeader(header M) *Job {
	job.Header = header
	return job
}

func (job *Job) SetReferer(referer string) *Job {
	job.Header["Referer"] = referer
	return job
}

func (job *Job) SetCookie(cookie string) *Job {
	job.Cookie = &Cookie{cookie}
	if job.Client.Jar == nil {
		job.Client.Jar = job.Cookie
		return job
	}
	cookieURL, _ := url.Parse(job.Url)
	job.Client.Jar.SetCookies(cookieURL, job.Cookie.Cookies(cookieURL))
	return job
}

// 发送请求
func (job *Job) DoWithContext(ctx ...context.Context) *Response {
	req, err := job.Request()
	if err != nil {
		return job.error(err)
	}
	for _, c := range ctx {
		req = req.WithContext(c)
	}
	resp, err := job.Client.Do(req)
	if err != nil {
		return job.error(err)
	}
	return job.succeed(resp)
}

// 发送请求
func (job *Job) Do() *Response {
	return job.DoWithContext()
}

func (job *Job) Plain() (*http.Response, error) {
	req, err := job.Request()
	if err != nil {
		return nil, err
	}
	return job.Client.Do(req)
}

func (job *Job) Request() (req *http.Request, err error) {
	req, err = http.NewRequest(job.Method, job.Url, job.Data)
	if err != nil {
		return
	}
	if job.Query != nil {
		req.URL.RawQuery = job.Query.Encode()
	}
	if job.Header != nil {
		job.Header.WriteHeader(req.Header)
	}
	return
}

func (job *Job) succeed(resp *http.Response) *Response {
	return &Response{Response: resp, Job: job}
}

func (job *Job) error(err error) *Response {
	return &Response{err: err, Job: job}
}

func New(method, url string, opts ...Option) (job *Job) {
	job = &Job{Method: method, Url: url, Data: make(M), Query: make(M), Header: UserAgentHeader.Clone(), Cookie: new(Cookie)}
	for _, opt := range opts {
		opt(job)
	}
	return
}

func NewGet(url string, opts ...Option) *Job {
	return New(http.MethodGet, url, opts...)
}

func NewPost(url string, opts ...Option) *Job {
	return New(http.MethodPost, url, opts...)
}
