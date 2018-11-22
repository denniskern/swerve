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

	"github.com/axelspringer/swerve/src/configuration"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/axelspringer/swerve/src/db"
	"github.com/axelspringer/swerve/src/log"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

func prometheusHandler() httprouter.Handle {
	h := promhttp.Handler()

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}

// NewAPIServer creates a new API server instance
func NewAPIServer(listener string, staticDir string, dynDB *db.DynamoDB) *API {
	api := &API{
		listener: listener,
		db:       dynDB,
	}

	// register api router
	router := httprouter.New()
	router.GET("/health", api.health)
	router.GET("/metrics", prometheusHandler())
	router.GET("/version", api.version)

	router.GET("/api/export", api.exportDomains)
	router.POST("/api/import", api.importDomains)
	router.GET("/api/domain", api.fetchAllDomains)
	router.GET("/api/domain/:name", api.fetchDomain)
	router.POST("/api/domain", api.registerDomain)
	router.DELETE("/api/domain/:name", api.purgeDomain)

	static := httprouter.New()
	static.ServeFiles("/*filepath", http.Dir(staticDir))
	router.NotFound = static

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

// health handler
func (api *API) health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sendJSONMessage(w, "ok", 200)
}

// version handler
func (api *API) version(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"version\":\"%s\"}", configuration.Version)))
}

// exportDomains exports the domains
func (api *API) exportDomains(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	domains, err := api.db.FetchAll()

	if err != nil {
		sendJSONMessage(w, "Error while fetching domains", 500)
		return
	}

	export := &db.ExportDomains{
		Domains: domains,
	}

	sendJSON(w, export, 200)
}

// importDomains imports a domain export set
func (api *API) importDomains(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Body == nil {
		sendJSONMessage(w, "Please send a request body", 400)
		return
	}

	var export db.ExportDomains

	if err := json.NewDecoder(r.Body).Decode(&export); err != nil {
		sendJSONMessage(w, "Invalid request body", 400)
		return
	}

	if err := api.db.DeleteAllDomains(); err != nil {
		log.Error(err)
		sendJSONMessage(w, "Database operation failed", 500)
		return
	}

	if err := api.db.Import(&export); err != nil {
		log.Error(err)
		sendJSONMessage(w, "Database operation failed", 500)
		return
	}

	sendJSONMessage(w, "ok", 200)
}

// purgeDomain deletes a domain entry
func (api *API) purgeDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	domain, err := api.db.FetchByDomain(name)

	if domain == nil || err != nil {
		sendJSONMessage(w, "not found", 404)
		return
	}

	if _, err = api.db.DeleteByDomain(name); err != nil {
		log.Error(err)
		sendJSONMessage(w, "Error while deleting domain", 500)
		return
	}

	api.db.DeleteTLSCacheEntry(name)

	sendJSONMessage(w, "ok", 204)
}

// fetchAllDomains return a list of all domains
func (api *API) fetchAllDomains(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	domains, err := api.db.FetchAll()

	if err != nil {
		sendJSONMessage(w, "Error while fetching domains", 500)
		return
	}

	sendJSON(w, domains, 200)
}

func (api *API) fetchDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	domain, err := api.db.FetchByDomain(name)

	if err != nil {
		sendJSONMessage(w, "not found", 404)
		return
	}

	sendJSON(w, domain, 200)
}

func (api *API) registerDomain(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Body == nil {
		sendJSONMessage(w, "Please send a request body", 400)
		return
	}

	var domain db.Domain

	if err := json.NewDecoder(r.Body).Decode(&domain); err != nil {
		sendJSONMessage(w, "Invalid request body", 400)
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
		sendJSONMessage(w, errMsg, 400)
		return
	}

	// insert new domain
	if err := api.db.InsertDomain(domain); err != nil {
		log.Error(err)
		sendJSONMessage(w, "Can't store document", 500)
		return
	}

	sendJSONMessage(w, "ok", 201)
}
