package request_test

import (
	"encoding/json"
	"fmt"
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
	api := request.New("get", "https://api.bilibili.com/x/v2/reply")
	api.Data = request.M{"pn": "1", "type": "17", "oid": "643451139714449427", "sort": "0"}
	data := request.Json[ApiData](api)
	for _, r := range data.Data.Replies {
		t.Log(r)
	}
}

func TestGet(t *testing.T) {
	job := request.New(
		"GET",
		"https://postman-echo.com/get",
		request.Data("test", "123"),
		request.Cookie("buvid", "somebase64"),
		request.Header("auth", "admin"),
	)
	b, _ := json.MarshalIndent(job.Test(), "", "\t")
	fmt.Printf("%v\n", string(b))
}

// func TestPost(t *testing.T) {
// 	result := request.Post(
// 		"https://postman-echo.com/post",
// 		request.Data("test", "123"),
// 		request.Cookie("buvid", "somebase64"),
// 		request.Header("auth", "admin"),
// 	)
// 	if result.Error() != nil {
// 		t.Fatal(result.Error())
// 	}
// 	m := make(map[string]any)
// 	err := result.Json(&m)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Printf("POST: %v\n", m)
// }

// func TestWeibo(t *testing.T) {
// 	resp := request.Get(
// 		"https://wx4.sinaimg.cn/orj480/007Raq4zly8hd1vqpx3coj30u00u00uv.jpg",
// 		request.Referer("https://weibo.com/"),
// 	)
// 	if resp.Error() != nil {
// 		t.Fatal(resp.Error())
// 	}
// 	err := resp.WriteWithPath("./public", os.ModePerm)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestGorm(t *testing.T) {
// 	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
// 	db.AutoMigrate(new(request.Job))

// 	job := request.New(
// 		http.MethodPost,
// 		"https://postman-echo.com/post",
// 		request.Data("test", "123"),
// 		request.Cookie("buvid", "somebase64"),
// 		request.Header("auth", "admin"),
// 	)
// 	db.Create(job)

// 	var query request.Job
// 	db.Find(&query)

// 	result := query.Request()
// 	if result.Error() != nil {
// 		t.Fatal(result.Error())
// 	}
// 	m := make(map[string]any)
// 	err := result.Json(&m)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Printf("QUERY: %v\n", m)
// }
