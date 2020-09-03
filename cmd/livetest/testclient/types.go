package testclient

import (
	"net/http"

	"github.com/go-resty/resty/v2"
)

type request struct {
	Method      string
	Url         string
	Headers     map[string]string
	FormData    map[string]string
	QueryParams map[string]string
	Body        interface{}
}

type response struct {
	StatusCode int
	TraceInfo  resty.TraceInfo
	Message    string
	RawData    []byte
	Header     http.Header
	Cookies    []*http.Cookie
}

type login struct {
	Username string `json:"username"`
	Pwd      string `json:"pwd"`
}

type version struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
