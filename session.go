package request

import (
	"net/http"
	urlpkg "net/url"
	"strings"
)

// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func HasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

// removeEmptyPort strips the empty port in ":port" to ""
// as mandated by RFC 3986 Section 6.2.3.
func RemoveEmptyPort(host string) string {
	if HasPort(host) {
		return strings.TrimSuffix(host, ":")
	}
	return host
}

func GenerateURL(url string) (*urlpkg.URL, error) {
	// 解析新的URL
	parsedURL, err := urlpkg.Parse(url)
	if err != nil {
		return nil, err
	}
	parsedURL.Host = RemoveEmptyPort(parsedURL.Host)
	return parsedURL, nil
}

func ModifyURL(req *http.Request, url string) error {
	parsedURL, err := GenerateURL(url)
	if err != nil {
		return err
	}
	// 保留原有的查询参数
	parsedURL.RawQuery = req.URL.RawQuery
	// 更新请求的URL字段
	req.URL = parsedURL
	req.Host = parsedURL.Host
	return nil
}

type Session struct {
	Client  http.Client
	Request *http.Request
}

func (s *Session) Do(url ...string) *Response {
	if len(url) != 0 && url[0] != "" {
		ModifyURL(s.Request, url[0])
	}
	resp, err := s.Client.Do(s.Request)
	if err != nil {
		return NewError(err)
	}
	return NewResponse(resp)
}

func NewSession(method string, opts ...Option) (*Session, error) {
	return NewSessionFromJob(New(method, "", opts...))
}

func NewSessionFromJob(job *Job) (*Session, error) {
	req, err := job.Request()
	if err != nil {
		return nil, err
	}
	return &Session{
		Request: req,
		Client:  job.Client,
	}, nil
}

type TypeSession[T any] Session

func (t *TypeSession[T]) Session() *Session {
	return (*Session)(t)
}

func (t *TypeSession[T]) Do() (out T, err error) {
	resp := t.Session().Do("")
	err = resp.err
	if err != nil {
		return
	}
	err = resp.Json(&out)
	return
}

func (t *TypeSession[T]) Must() (out T) {
	out, _ = t.Do()
	return
}

func NewTypeSession[T any](method, url string, opts ...Option) (*TypeSession[T], error) {
	s, err := NewSession(method, opts...)
	if err != nil {
		return nil, err
	}
	ModifyURL(s.Request, url)
	return (*TypeSession[T])(s), nil
}

func NewTypeSessionFromJob[T any](job *Job) (*TypeSession[T], error) {
	s, err := NewSessionFromJob(job)
	if err != nil {
		return nil, err
	}
	return (*TypeSession[T])(s), nil
}
