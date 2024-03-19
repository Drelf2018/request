package request_test

import (
	"net/http"
	"testing"

	"github.com/Drelf2018/request"
)

func TestFileSystem(t *testing.T) {
	http.ListenAndServe("localhost:8080", http.FileServer(request.DefaultDownloadSystem("./downloads", nil)))
}
