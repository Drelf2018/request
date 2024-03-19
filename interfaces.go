package request

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
)

type ScanValuer interface {
	sql.Scanner
	driver.Valuer
}

type CookieJar interface {
	ScanValuer
	http.CookieJar
}

type DownloaderInterface interface {
	Get(url string) (*http.Response, error)
}
