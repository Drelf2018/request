# request

通过预定任务发送网络请求

### 使用

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

func TestBili(t *testing.T) {
	api := request.New(
		http.MethodGet,
		"https://api.bilibili.com/x/v2/reply",
		request.Datas(request.M{"pn": "1", "type": "17", "oid": "643451139714449427", "sort": "0"}),
	)
	resp := api.Request()
	if resp.Error() != nil {
		t.Fatal(resp.Error())
	}
	var data ApiData
	resp.Json(&data)
	for _, r := range data.Data.Replies {
		t.Log(r)
	}
}
```