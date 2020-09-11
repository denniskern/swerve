package model

import (
	"net/url"
	"strings"
)

// GetRedirect returns the calculated route
func (r *Redirect) GetRedirect(reqURL *url.URL, scheme string) (string, int) {
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
				if strings.HasPrefix(p.To, "http://") || strings.HasPrefix(p.To, "https://") {
					reURL = p.To
				} else {
					rePath = p.To
				}
				break
			}
		}
	}

	if strings.HasSuffix(reURL, "/") && strings.HasPrefix(rePath, "/") {
		rePath = strings.TrimLeft(rePath, "/")
	}
	if !strings.HasPrefix(reURL, "http://") && !strings.HasPrefix(reURL, "https://") {
		reURL = scheme + reURL
	}

	return reURL + rePath + reQuery, code
}
