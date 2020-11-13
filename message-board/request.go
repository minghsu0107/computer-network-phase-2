package main

import (
	"strings"
)

type request struct {
	Method string
	Path   string
	Header map[string]string
	Body   string
}

func parseRequest(requestStr string) request {
	tmp := strings.Split(requestStr, "\r\n")
	var header = make(map[string]string)
	isBody := false
	for _, row := range tmp {
		if row == "\r\n" {
			isBody = true
			continue
		}
		if !isBody {
			headerRow := strings.Split(row, ": ")
			if len(headerRow) == 2 {
				header[headerRow[0]] = headerRow[1]
			}
		}
	}

	var method string
	var path string
	if len(tmp) > 0 {
		firstRow := strings.Split(tmp[0], " ")
		if len(firstRow) >= 2 {
			method = firstRow[0]
			path = firstRow[1]
		}
	}
	return request{
		Method: method,
		Path:   path,
		Header: header,
		Body:   tmp[len(tmp)-1],
	}
}

func (req request) getCookie(key string) string {
	rawCookieList, ok := req.Header["Cookie"]
	if !ok {
		return ""
	}
	cookieList := strings.Split(rawCookieList, "; ")
	tmp := make(map[string]string)
	for _, cookie := range cookieList {
		kv := strings.Split(cookie, "=")
		if len(kv) >= 2 {
			tmp[kv[0]] = kv[1]
		}
	}
	return tmp[key]
}
