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
