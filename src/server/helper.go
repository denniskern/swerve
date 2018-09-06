// Copyright 2018 Axel Springer SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/axelspringer/swerve/src/db"
)

func sendJSON(w http.ResponseWriter, obj interface{}, code int) {
	jsonBytes, _ := json.Marshal(obj)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func sendJSONMessage(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"code\":%d,\"message\":\"%s\"}", code, msg)))
}

func sendPlainMessage(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("%d - %s", code, msg)))
}

// promoteRedirect append the path + querystring to the redirect host
func promoteRedirect(redirect string, reqURL *url.URL) string {
	newRedirect := path.Join(redirect, reqURL.Path)
	if len(reqURL.RawQuery) > 0 {
		newRedirect = newRedirect + "?" + reqURL.RawQuery
	}

	return newRedirect
}

// pathMappingRedirect tries to match the path of the request to the mapping list
func pathMappingRedirect(pathList *db.PathList, redirect string, reqURL *url.URL) string {
	if pathList == nil {
		return redirect
	}
	// look for matching paths
	for _, p := range *pathList {
		if p.To == "" {
			continue
		}
		// we match the path prefix
		if strings.HasPrefix(reqURL.Path, p.From) {
			// path redirect
			if strings.HasPrefix(p.To, "/") {
				return path.Join(redirect, p.To)
			}
			// domain redirect
			return p.To
		}
	}

	return redirect
}
