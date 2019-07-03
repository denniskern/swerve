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
	"os"
	"strings"
	"time"

	"github.com/axelspringer/swerve/src/configuration"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/axelspringer/swerve/src/db"
	"github.com/axelspringer/swerve/src/log"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

var (
	uiDomain = getOSPrefixEnv("UI_DOMAIN")
)

var secret string

const (
	envPrefix = "SWERVE_"
)

// getOSPrefixEnv get os env
func getOSPrefixEnv(s string) string {
	if e := strings.TrimSpace(os.Getenv(envPrefix + s)); len(e) > 0 {
		return e
	}

	return ""
}

func prometheusHandler() httprouter.Handle {
	h := promhttp.Handler()

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}

// NewAPIServer creates a new API server instance
func NewAPIServer(listener string, apiSecret string, dynDB *db.DynamoDB) *API {
	api := &API{
		listener: listener,
		db:       dynDB,
	}

	secret = apiSecret

	// register api router
	router := httprouter.New()
	router.GET("/health", api.health)
	router.GET("/metrics", prometheusHandler())
	router.GET("/version", api.version)
	router.POST("/login", api.login)
	router.OPTIONS("/login", api.options)

	authRouter := httprouter.New()
	authRouter.GET("/api/export", api.exportDomains)
	authRouter.POST("/api/import", api.importDomains)
	authRouter.GET("/api/domain", api.fetchAllDomains)
	authRouter.GET("/api/domain/:name", api.fetchDomain)
	authRouter.POST("/api/domain", api.registerDomain)
	authRouter.DELETE("/api/domain/:name", api.purgeDomain)
	authRouter.PUT("/api/domain/:name", api.updateDomain)
	authRouter.GET("/refresh", api.refresh)

	router.NotFound = AuthHandler(authRouter)
	// router.NotFound = static

	api.server = &http.Server{
		Addr:    listener,
		Handler: router,
	}

	return api
}

// Listen to socket
func (api *API) Listen() error {
	log.Infof("API listening to %s", api.listener)
	return api.server.ListenAndServe()
}

func (api *API) options(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", uiDomain)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept, token")
	w.WriteHeader(http.StatusOK)
}

// health handler
func (api *API) health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sendJSONMessage(w, "ok", http.StatusOK)
}

// version handler
func (api *API) version(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"version\":\"%s\"}", configuration.Version)))
	w.WriteHeader(http.StatusOK)
}

// exportDomains exports the domains
func (api *API) exportDomains(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	domains, err := api.db.FetchAll()

	if err != nil {
		sendJSONMessage(w, "Error while fetching domains", http.StatusInternalServerError)
		return
	}

	export := &db.ExportDomains{
		Domains: domains,
	}

	sendJSON(w, export, http.StatusOK)
}

// importDomains imports a domain export set
func (api *API) importDomains(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Body == nil {
		sendJSONMessage(w, "Please send a request body", http.StatusBadRequest)
		return
	}

	var export db.ExportDomains

	if err := json.NewDecoder(r.Body).Decode(&export); err != nil {
		sendJSONMessage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := api.db.DeleteAllDomains(); err != nil {
		log.Error(err)
		sendJSONMessage(w, "Database operation failed", http.StatusInternalServerError)
		return
	}

	if err := api.db.Import(&export); err != nil {
		log.Error(err)
		sendJSONMessage(w, "Database operation failed", http.StatusInternalServerError)
		return
	}

	sendJSONMessage(w, "ok", http.StatusOK)
}

// purgeDomain deletes a domain entry
func (api *API) purgeDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	domain, err := api.db.FetchByDomain(name)
	if domain == nil || err != nil {
		sendJSONMessage(w, "Not found", http.StatusNotFound)
		return
	}

	if _, err = api.db.DeleteByDomain(name); err != nil {
		log.Error(err)
		sendJSONMessage(w, "Error while deleting domain", http.StatusInternalServerError)
		return
	}

	api.db.DeleteTLSCacheEntry(name)

	sendJSONMessage(w, "ok", http.StatusNoContent)
}

// updateDomain updates a domain entry
func (api *API) updateDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body == nil {
		sendJSONMessage(w, "Please send a request body", http.StatusBadRequest)
		return
	}

	name := ps.ByName("name")
	oldDomain, err := api.db.FetchByDomain(name)

	if oldDomain == nil || err != nil {
		sendJSONMessage(w, "Not found", http.StatusNotFound)
		return
	}

	var domain db.Domain

	if err := json.NewDecoder(r.Body).Decode(&domain); err != nil {
		sendJSONMessage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	domain.ID = oldDomain.ID
	domain.Created = oldDomain.Created
	domain.Modified = time.Now().Format(time.RFC3339)

	// validate
	if errList := domain.Validate(); len(errList) > 0 {
		errMsg := ""
		for _, err := range errList {
			errMsg = errMsg + err.Error() + ". "
		}
		sendJSONMessage(w, errMsg, http.StatusBadRequest)
		return
	}

	// insert new domain
	if err := api.db.InsertDomain(domain); err != nil {
		log.Error(err)
		sendJSONMessage(w, "Can't store document", http.StatusInternalServerError)
		return
	}

	// api.db.DeleteTLSCacheEntry(id)

	sendJSONMessage(w, "ok", http.StatusOK)
}

// fetchAllDomains return a list of all domains
func (api *API) fetchAllDomains(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var cursor *string
	queryparam, ok := r.URL.Query()["cursor"]
	if !ok || len(queryparam[0]) < 1 {
		cursor = nil
	} else {
		cursor = &queryparam[0]
	}

	domains, cursor, err := api.db.FetchAllPaginated(cursor)
	if err != nil {
		log.Error(err)
		sendJSONMessage(w, "Error while fetching domains", http.StatusInternalServerError)
		return
	}

	sendJSON(w, struct {
		Domains []db.Domain `json:"domains"`
		Cursor  *string     `json:"cursor"`
	}{
		domains,
		cursor,
	}, http.StatusOK)
}

func (api *API) fetchDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	domain, err := api.db.FetchByDomain(name)
	if err != nil {
		sendJSONMessage(w, "Not found", http.StatusNotFound)
		return
	}

	sendJSON(w, domain, http.StatusOK)
}

func (api *API) registerDomain(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Body == nil {
		sendJSONMessage(w, "Please send a request body", http.StatusBadRequest)
		return
	}

	var domain db.Domain

	if err := json.NewDecoder(r.Body).Decode(&domain); err != nil {
		sendJSONMessage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	alreadyExisting, err := api.db.FetchByDomain(domain.Name)
	if err != nil {
		sendJSONMessage(w, "Could not check database for already existing entry", http.StatusInternalServerError)
		return
	}
	if alreadyExisting.ID != "" {
		sendJSONMessage(w, "Already exists", http.StatusBadRequest)
		return
	}

	domain.ID = uuid.Must(uuid.NewV4()).String()
	domain.Created = time.Now().Format(time.RFC3339)
	domain.Modified = domain.Created

	// validate
	if errList := domain.Validate(); len(errList) > 0 {
		errMsg := ""
		for _, err := range errList {
			errMsg = errMsg + err.Error() + ". "
		}
		sendJSONMessage(w, errMsg, http.StatusBadRequest)
		return
	}

	// insert new domain
	if err := api.db.InsertDomain(domain); err != nil {
		log.Error(err)
		sendJSONMessage(w, "Can't store document", http.StatusInternalServerError)
		return
	}

	sendJSONMessage(w, "ok", http.StatusCreated)
}

func (api *API) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		sendJSONMessage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = api.db.CheckPassword(creds.Username, creds.Password)
	if err != nil {
		sendJSONMessage(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Error(err)
		sendJSONMessage(w, "JWT token could not be signed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", uiDomain)
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	w.WriteHeader(http.StatusOK)
}

func (api *API) refresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			sendJSONMessage(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		sendJSONMessage(w, "Invalid token", http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if !tkn.Valid {
		sendJSONMessage(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			sendJSONMessage(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		sendJSONMessage(w, "Invalid token", http.StatusBadRequest)
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		sendJSONMessage(w, "Not yet", http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Error(err)
		sendJSONMessage(w, "Could not sign new token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
