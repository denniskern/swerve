package model

import (
	"net/url"
	"strings"
)

// GetRedirect returns the calculated route
func (r *Redirect) GetRedirect(reqURL *url.URL) (string, int) {
	code := r.Code
	reURL := r.RedirectTo
	rePath := ""
	reQuery := ""
	reqPath := reqURL.EscapedPath()

	if r.Promotable == true {
		rePath = reqURL.Path

		if reqURL.RawQuery != "" {
			reQuery = "?" + reqURL.RawQuery
		}
	}

	if r.PathMaps != nil && len(r.PathMaps) > 0 {
		for _, p := range r.PathMaps {
			if p.To == "" {
				continue
			}
			if reqPath == p.From {
				rePath = p.To
				if r.Promotable && strings.Contains(rePath, "?") {
					reQuery = strings.Replace(reQuery, "?", "&", 1)
				}
				break
			}
		}
	}

	if strings.HasSuffix(reURL, "/") && strings.HasPrefix(rePath, "/") {
		rePath = strings.TrimLeft(rePath, "/")
	}
	if !strings.HasPrefix(reURL, "http://") && !strings.HasPrefix(reURL, "https://") {
		reURL = "https://" + reURL
	}
	if strings.Contains(rePath, reqURL.RawQuery) {
		reQuery = ""
	}

	return reURL + rePath + reQuery, code
}
