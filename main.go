package request

func Get(url string, opts ...Option) *Response {
	return NewGet(url, opts...).Do()
}

func Post(url string, opts ...Option) *Response {
	return NewPost(url, opts...).Do()
}

func Json[T any](job *Job) (out T) {
	job.Do().Json(&out)
	return
}
