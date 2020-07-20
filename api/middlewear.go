package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/axelspringer/swerve/log"

	"github.com/gorilla/mux"
	"github.com/rs/xid"
)

func (api *API) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		guid := xid.New()
		entry := logrus.WithFields(logrus.Fields{
			"request_id": guid.String(),
		})
		ctx := context.WithValue(r.Context(), "RequestLogger", entry)
		next.ServeHTTP(w, r.WithContext(ctx))
		start := time.Now()
		log.Infof("Request starting %s %s %s %s", r.Proto, r.Method, r.Host, r.URL.Path)
		end := time.Now()
		diff := end.Sub(start)
		log.Infof("Request finished in %d ms", diff.Nanoseconds()/int64(time.Millisecond))
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
				sendJSONMessage(w, "No token found", http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value

		if !api.Model.CheckToken(tknStr, api.Config.Secret) {
			sendJSONMessage(w, "Token is invalid", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
