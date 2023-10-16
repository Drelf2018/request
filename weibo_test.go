package request_test

import (
	"os"
	"testing"

	"github.com/Drelf2018/request"
)

func TestWeibo(t *testing.T) {
	resp := request.Get(
		"https://wx4.sinaimg.cn/orj480/007Raq4zly8hd1vqpx3coj30u00u00uv.jpg",
		request.Referer("https://weibo.com/"),
	)
	if resp.Error() != nil {
		t.Fatal(resp.Error())
	}
	err := resp.Write("face.jpg", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
