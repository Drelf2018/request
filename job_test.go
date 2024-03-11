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
