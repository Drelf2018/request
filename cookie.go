package request

import (
	"database/sql/driver"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
)

var cookieExp = regexp.MustCompile("([^ =]+)=([^ ;]+);")
var tmpl, _ = template.New("cookieTemplate").Parse("{{range .}}{{.Name}}={{.Value}};{{end}}")

type Cookie struct {
	Text string
}

func (c *Cookie) Write(p []byte) (n int, err error) {
	c.Text += string(p)
	return len(p), nil
}

func (c *Cookie) SetCookies(_ *url.URL, cookies []*http.Cookie) {
	tmpl.Execute(c, cookies)
}

func (c *Cookie) Cookies(*url.URL) []*http.Cookie {
	list := cookieExp.FindAllStringSubmatch(c.Text, -1)
	cookies := make([]*http.Cookie, len(list))
	for i, match := range list {
		cookies[i] = &http.Cookie{Name: match[1], Value: match[2]}
	}
	return cookies
}

func (c *Cookie) Scan(val any) error {
	if val == nil {
		*c = Cookie{}
		return nil
	}
	switch v := val.(type) {
	case []byte:
		*c = Cookie{string(v)}
	case string:
		*c = Cookie{v}
	default:
		return ErrType
	}
	return nil
}

func (c *Cookie) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return c.Text, nil
}

func (c *Cookie) String() string {
	return c.Text
}

var _ CookieJar = new(Cookie)
