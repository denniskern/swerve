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
	"time"

	"github.com/axelspringer/swerve/src/log"
)

func sendJSON(w http.ResponseWriter, obj interface{}, code int) {
	jsonBytes, _ := json.Marshal(struct {
		Data interface{} `json:"data"`
	}{
		Data: obj,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", uiDomain)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

func sendJSONMessage(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", uiDomain)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("{\"code\":%d,\"message\":\"%s\"}", code, msg)))
}

func sendPlainMessage(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", uiDomain)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("%d - %s", code, msg)))
}

func handlerWithLogging(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Infof("Request starting %s %s %s %s", r.Proto, r.Method, r.Host, r.URL.Path)
		f(w, r)
		end := time.Now()
		diff := end.Sub(start)
		log.Infof("Request finished in %d ms", diff.Nanoseconds()/int64(time.Millisecond))
	})
}
