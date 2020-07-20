package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/axelspringer/swerve/log"

	"github.com/gorilla/mux"
)

func (api *API) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		diff := end.Sub(start)
		log.Infof("Request finished in %d ms for %s %s %s %s ", diff.Nanoseconds()/int64(time.Millisecond), r.Proto, r.Method, r.Host, r.URL.Path)

	})
}

func (api *API) corsMiddlewear(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", api.Config.COR)
		if r.Method == http.MethodOptions {
			methods, err := mux.CurrentRoute(r).GetMethods()
			if err == nil {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (api *API) authMiddlewear(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				sendJSONMessage(r, w, "No token found", http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value

		if !api.Model.CheckToken(tknStr, api.Config.Secret) {
			sendJSONMessage(r, w, "Token is invalid", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
