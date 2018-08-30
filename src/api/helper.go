package api

import (
	"net/http"
	"net/url"
	"strings"
)

func sendErrBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write([]byte("{\"error\":\"bad request\",\"message\":\"" + err.Error() + "\"}"))
}

func splitDomainPath(name string) (string, string, error) {
	if !strings.HasPrefix(name, "//") {
		name = "//" + name
	}

	url, err := url.Parse(name)
	if err != nil {
		return "", "", err
	}

	host, path := url.Host, url.Path

	if path == "/" {
		path = ""
	}

	return host, path, nil
}
