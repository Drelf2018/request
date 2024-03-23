package request

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.37"

var UserAgentHeader = M{"User-Agent": UserAgent}
var Headers = M{
	"Accept-Language": "zh-CN,zh;q=0.9",
	"Accept-Encoding": "gzip, deflate, br",
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
	"User-Agent":      UserAgent,
}

var (
	ErrType    = errors.New("request: failed to unmarshal JSONB value")
	ErrOddArgs = errors.New("request: odd number of parameters passed in")
)

type M map[string]string

func (m M) Set(k, v string) {
	m[k] = v
}

func (m M) SetTrimmed(k, v string) {
	m[strings.TrimSpace(k)] = strings.TrimSpace(v)
}

func (m M) SetAny(k string, v any) error {
	b, err := json.Marshal(v)
	m[k] = string(b)
	return err
}

func (m M) SetAll(s ...string) error {
	l := len(s)
	if l&1 == 1 {
		return ErrOddArgs
	}
	for i := 0; i < l; i += 2 {
		m.SetTrimmed(s[i], s[i+1])
	}
	return nil
}

func (m M) SetMap(p M) {
	for k, v := range p {
		m[k] = v
	}
}

func (m M) Add(k, v string) {
	m[k] += v
}

func (m M) Del(k string) {
	delete(m, k)
}

func (m M) Clone() (p M) {
	p = make(M)
	p.SetMap(m)
	return
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
func (m M) GormDataType() string {
	return "jsonmap"
}

func (m M) Read(p []byte) (int, error) {
	b := new(bytes.Buffer)
	for k, v := range m {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(v)
		b.WriteByte('&')
	}
	if b.Len() != 0 {
		b.Truncate(b.Len() - 1)
	}
	n, _ := b.Read(p)
	return n, io.EOF
}

func (m M) Encode() string {
	u := make(url.Values)
	for k, v := range m {
		u[k] = []string{v}
	}
	return u.Encode()
}

func (m M) WriteHeader(header http.Header) {
	for k, v := range m {
		header.Set(k, v)
	}
}

func (m M) Data(job *Job) {
	job.SetData(m)
}

func (m M) Query(job *Job) {
	job.SetQuery(m)
}

func (m M) Header(job *Job) {
	job.SetHeader(m)
}

var _ ScanValuer = new(M)
