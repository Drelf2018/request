package request

import (
	"net/http"
)

func Get(url string, opts ...Option) *Response {
	return NewGet(url, opts...).Do()
}

func Post(url string, opts ...Option) *Response {
	return NewPost(http.MethodPost, opts...).Do()
}

func Json[T any](job *Job) (out T) {
	job.Do().Json(&out)
	return
}
