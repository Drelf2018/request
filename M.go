package request

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var HEADERS = M{
	"Accept-Language": "zh-CN,zh;q=0.9",
	"Accept-Encoding": "gzip, deflate, br",
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.37",
}

var (
	ErrType    = errors.New("failed to unmarshal JSONB value")
	ErrOddArgs = errors.New("odd number of parameters passed in")
)

type M map[string]string

func (m M) Any(k string, v any) error {
	b, err := json.Marshal(v)
	m[k] = string(b)
	return err
}

func (m M) Set(k, v string) {
	m[k] = v
}

// url.Values 赋值
func (m M) SetTo(p interface{ Set(string, string) }) {
	for k, v := range m {
		p.Set(k, v)
	}
}

func (m M) Add(k, v string) {
	m[k] = m[k] + v
}

// url.Values 添加
func (m M) AddTo(p interface{ Add(string, string) }) {
	for k, v := range m {
		p.Add(k, v)
	}
}

func (m M) Del(k string) {
	delete(m, k)
}

func (m M) Copy(p M) {
	for k, v := range m {
		p[k] = v
	}
}

func (m M) New() (p M) {
	p = make(M)
	m.Copy(p)
	return
}

func (m M) Insert(k, v string) {
	m[strings.TrimSpace(k)] = strings.TrimSpace(v)
}

func (m M) Inserts(s ...string) {
	l := len(s)
	if l&1 == 1 {
		panic(ErrOddArgs)
	}
	for i := 0; i < l; i += 2 {
		m.Insert(s[i], s[i+1])
	}
}

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
func (M) GormDataType() string {
	return "json"
}

// CookieJar 实现
func (m M) SetCookies(_ *url.URL, cookies []*http.Cookie) {
	for _, c := range cookies {
		m[c.Name] = c.Value
	}
}

func (m M) Cookies(_ *url.URL) (cookies []*http.Cookie) {
	for k, v := range m {
		cookies = append(cookies, &http.Cookie{Name: k, Value: v})
	}
	return
}
