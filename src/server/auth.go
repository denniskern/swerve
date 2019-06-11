package server

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// AuthHandler factory
func AuthHandler(next http.Handler) http.Handler {
	var authHandler AuthMiddlewareHandler
	authHandler.next = next
	return authHandler
}

func (amh AuthMiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://swerve.tortuga.cloud")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept, token")
		w.WriteHeader(http.StatusOK)
		return
	}

	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			sendJSONMessage(w, "No token found", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if !tkn.Valid {
		sendJSONMessage(w, "Token is invalid", http.StatusUnauthorized)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			sendJSONMessage(w, "Token signature is invalid", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	amh.next.ServeHTTP(w, r)
}
