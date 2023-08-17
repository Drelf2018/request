# request

通过预定任务请求

### 使用

```go
func TestGet() {
	result := request.Get(
		"https://postman-echo.com/get",
		request.Data("test", "123"),
		request.Cookie("buvid", "somebase64"),
		request.Header("auth", "admin"),
	)
	if result.Error != nil {
		panic(result.Error)
	}
	m := make(map[string]any)
	err := result.Json(&m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("text: %v\n", result.Text())
	fmt.Printf("json: %v\n", m)
}
```
