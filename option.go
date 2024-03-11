package request

type Option func(*Job)

func SetURL(url string) Option {
	return func(job *Job) {
		job.Url = url
	}
}

func SetReferer(referer string) Option {
	return func(job *Job) {
		job.SetReferer(referer)
	}
}

func SetData(data M) Option {
	return func(job *Job) {
		job.SetData(data)
	}
}

func SetQuery(query M) Option {
	return func(job *Job) {
		job.SetQuery(query)
	}
}

func SetHeader(header M) Option {
	return func(job *Job) {
		job.SetHeader(header)
	}
}

func SetCookie(cookie string) Option {
	return func(job *Job) {
		job.SetCookie(cookie)
	}
}
