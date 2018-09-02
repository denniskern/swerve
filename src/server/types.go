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
	"net/http"

	"github.com/axelspringer/swerve/src/certificate"
	"github.com/axelspringer/swerve/src/db"
)

// ListenerInterface contains the main server functions
type ListenerInterface interface {
	Listen() error
}

// API server model
type API struct {
	ListenerInterface
	db       *db.DynamoDB
	server   *http.Server
	listener string
}

// HTTP server model
type HTTP struct {
	ListenerInterface
	certManager *certificate.Manager
	server      *http.Server
	listener    string
}

// HTTPS server model
type HTTPS struct {
	ListenerInterface
	certManager *certificate.Manager
	server      *http.Server
	listener    string
}
