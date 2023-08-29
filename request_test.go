package request_test

import (
	"fmt"
	"testing"

	"github.com/Drelf2018/request"
)

func TestGet(t *testing.T) {
	result := request.Get(
		"https://postman-echo.com/get",
		request.Data("test", "123"),
		request.Cookie("buvid", "somebase64"),
		request.Header("auth", "admin"),
	)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	m := make(map[string]any)
	err := result.Json(&m)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("GET: %v\n", m)
}

func TestPost(t *testing.T) {
	result := request.Post(
		"https://postman-echo.com/post",
		request.Data("test", "123"),
		request.Cookie("buvid", "somebase64"),
		request.Header("auth", "admin"),
	)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	m := make(map[string]any)
	err := result.Json(&m)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("POST: %v\n", m)
}

// func TestGorm(t *testing.T) {
// 	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
// 	db.AutoMigrate(new(request.Job))

// 	job := request.New(http.MethodGet, "https://postman-echo.com/get", request.Data("test", "123"))
// 	db.Create(job)

// 	var query request.Job
// 	db.Find(&query)

// 	result := query.Request()
// 	if result.Error != nil {
// 		t.Fatal(result.Error)
// 	}
// 	m := make(map[string]any)
// 	err := result.Json(&m)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Printf("QUERY: %v\n", m)
// }
