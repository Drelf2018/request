# request

通过预定任务发送网络请求

### 使用

```go
package request_test

import (
	"testing"

	"github.com/Drelf2018/gorms"
	"github.com/Drelf2018/request"
)

var db = gorms.SetSQLite("jobs.db").AutoMigrate(&request.Job{})
var records = 0

func TestFind(t *testing.T) {
	r := gorms.MustFind[request.Job]()
	records = len(r)
	t.Log(r)
}

func TestCreate(t *testing.T) {
	if records < 2 {
		err := db.Create(
			request.NewGet(
				"https://www.baidu.com",
			).SetQuery(
				request.M{"uid": "12138"},
			).SetHeader(
				request.M{"auth": "admin"},
			).SetCookie("buvid3=BE308E31; i-wanna-go-back=-1;"),
		).Error
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestDo(t *testing.T) {
	job := request.NewPost(
		"https://postman-echo.com/post",
	).SetData(
		request.M{"data": "abc"},
	).SetQuery(
		request.M{"uid": "114514"},
	).SetHeader(
		request.M{"auth": "admin"},
	).SetCookie(
		"buvid3=ABCDE;",
	).SetReferer(
		"postman-echo.com",
	)
	resp := job.Do()
	err := resp.Error()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Data())
}
```

### 进阶用法

```go
package request_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Drelf2018/request"
)

type Replie struct {
	Member struct {
		Mid   string `json:"mid"`
		Uname string `json:"uname"`
	} `json:"member"`
	Content struct {
		Message string `json:"message"`
	} `json:"content"`
}

func (r Replie) String() string {
	return fmt.Sprintf("%v(%v): %v", r.Member.Uname, r.Member.Mid, r.Content.Message)
}

type ApiData struct {
	Code int `json:"code"`
	Data struct {
		Replies []Replie `json:"replies"`
	} `json:"data"`
}

func TestTypeSession(t *testing.T) {
	session, err := request.NewTypeSession[ApiData](http.MethodGet, "https://api.bilibili.com/x/v2/reply", request.M{
		"pn":   "1",
		"type": "11",
		"oid":  "307611044",
		"sort": "0",
	}.Query)
	if err != nil {
		t.Fatal(err)
	}
	data, err := session.Do()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}
```