package api

import (
	"encoding/json"
	"net/http"

	"github.com/axelspringer/swerve/src/db"
	"github.com/axelspringer/swerve/src/log"
	"github.com/julienschmidt/httprouter"
)

// NewServer creates a new instance
func NewServer(listener string, dynDB *db.DynamoDB) *Server {
	server := &Server{
		Listener: listener,
		db:       dynDB,
	}

	// register api router
	router := httprouter.New()
	router.GET("/health", server.health)
	router.GET("/domain", server.fetchAllDomains)
	router.GET("/domain/:name", server.fetchDomain)
	router.POST("/domain", server.registerDomain)
	router.DELETE("/domain/:name", server.purgeDomain)

	server.Server = &http.Server{
		Addr:    listener,
		Handler: router,
	}

	return server
}

// Listen to socket
func (s *Server) Listen() error {
	log.Infof("API listening to %s", s.Listener)
	return s.Server.ListenAndServe()
}

// health handler
func (s *Server) health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"status\":\"ok\"}"))
}

// purgeDomain deletes a domain entry
func (s *Server) purgeDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	domain, err := s.db.FetchByDomain(ps.ByName("name"))

	if domain == nil || err != nil {
		http.NotFound(w, r)
		return
	}

	if _, err = s.db.DeleteByDomain(ps.ByName("name")); err != nil {
		log.Infof("Can't delete entity %#v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"status\":\"ok\"}"))
}

// fetchAllDomains return a list of all domains
func (s *Server) fetchAllDomains(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	domains, err := s.db.FetchAll()
	if err != nil {
		http.Error(w, "Error while fetching domains", 500)
		return
	}

	jsonBytes, err := json.Marshal(domains)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func (s *Server) fetchDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	domain, err := s.db.FetchByDomain(ps.ByName("name"))

	if err != nil {
		log.Infof("Error while fetch domain '%s' %#v", ps.ByName("name"), err)
		http.NotFound(w, r)
		return
	}

	jsonBytes, err := json.Marshal(domain)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func (s *Server) registerDomain(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	var domain db.Domain

	if err := json.NewDecoder(r.Body).Decode(&domain); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	host, path, err := splitDomainPath(domain.Name)
	if err != nil {
		sendErrBadRequest(w, r, err)
	}

	domain.Name = host
	domain.PathPattern = path

	if err := s.db.InsertDomain(domain); err != nil {
		sendErrBadRequest(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte("{\"status\":\"ok\"}"))
}
