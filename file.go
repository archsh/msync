package main

import (
	"github.com/cavaliercoder/grab"
	"net/http"
	"os"
	"strings"
)

func FileSize(src string) (int64, error) {
	if strings.HasPrefix(src, "https://") || strings.HasPrefix(src, "http://") {
		if resp, e := http.Head(src); nil != e {
			return -1, e
		} else {
			return resp.ContentLength, nil
		}
	} else if strings.HasPrefix(src, "ftp://") {
		//return fetchFTP(src)
		return FTP_FileSize(src)
	} else if s, e := os.Stat(src); nil != e {
		return -1, e
	} else {
		return s.Size(), nil
	}
}

func FetchFile(src string) (string, error) {
	if strings.HasPrefix(src, "https://") || strings.HasPrefix(src, "http://") {
		return fetchHTTP(src)
	} else if strings.HasPrefix(src, "ftp://") {
		return fetchFTP(src)
	}
	return src, nil
}

func fetchHTTP(src string) (string, error) {
	if resp, e := grab.Get("/tmp", src); nil != e {
		return "", e
	} else {
		return resp.Filename, nil
	}
}

func fetchFTP(src string) (string, error) {
	return FTP_GetFile("/tmp", src)
}
