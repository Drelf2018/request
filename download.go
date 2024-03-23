package request

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Downloader struct {
	Sinaimg *Job
	Default *http.Request
}

func (d *Downloader) Get(url string) (resp *http.Response, err error) {
	switch {
	case strings.Contains(url, "sinaimg.cn"):
		return d.Sinaimg.SetURL(url).Plain()
	default:
		u, err := GenerateURL(url)
		if err != nil {
			return nil, err
		}
		d.Default.URL = u
		return http.DefaultClient.Do(d.Default)
	}
}

var DefaultDownloader = &Downloader{
	Sinaimg: New(http.MethodGet, "").SetReferer("https://weibo.com/"),
	Default: &http.Request{
		Method: http.MethodGet,
		Header: http.Header{"User-Agent": {UserAgent}},
	},
}

func SplitName(name string) (protocol, path string) {
	s := strings.SplitN(name[1:], "/", 2)
	return s[0], s[1]
}

type DownloadSystem struct {
	Root       string
	Downloader DownloaderInterface
}

func (ds DownloadSystem) Open(name string) (http.File, error) {
	if !strings.HasPrefix(name, "/http") {
		return http.Dir(ds.Root).Open(name)
	}

	protocol, path := SplitName(name)
	fullpath := filepath.Join(ds.Root, path)

	_, err := os.Stat(fullpath)
	if err == nil {
		return os.Open(fullpath)
	}

	err = os.MkdirAll(filepath.Dir(fullpath), os.ModePerm)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s://%s", protocol, path)
	resp, err := ds.Downloader.Get(url)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(fullpath, content, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return os.Open(fullpath)
}

func DefaultDownloadSystem(root string) *DownloadSystem {
	return &DownloadSystem{
		Root:       root,
		Downloader: http.DefaultClient,
	}
}

var _ DownloaderInterface = new(Downloader)
var _ http.FileSystem = new(DownloadSystem)
