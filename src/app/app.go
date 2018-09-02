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

package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/axelspringer/swerve/src/certificate"
	"github.com/axelspringer/swerve/src/configuration"
	"github.com/axelspringer/swerve/src/db"
	"github.com/axelspringer/swerve/src/log"
	"github.com/axelspringer/swerve/src/server"
)

// Setup the application configuration
func (a *Application) Setup() {
	// read config
	a.Config.FromEnv()
	a.Config.FromParameter()

	// setup logger
	log.SetupLogger(a.Config.LogLevel, a.Config.LogFormatter)

	// database connection
	var err error
	a.DynamoDB, err = db.NewDynamoDB(&a.Config.DynamoDB, a.Config.Bootstrap)
	if err != nil {
		log.Fatalf("Can't setup db connection %#v", err)
	}

	// certificate pool
	a.Certificates = certificate.NewManager(a.DynamoDB)
}

// Run the application
func (a *Application) Run() {
	// signal channel
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	// run the https listener
	httpsServer := server.NewHTTPSServer(a.Config.HTTPSListener, a.Certificates)
	go func() {
		log.Fatal(httpsServer.Listen())
	}()

	// run the http listener
	httpServer := server.NewHTTPServer(a.Config.HTTPListener, a.Certificates)
	go func() {
		log.Fatal(httpServer.Listen())
	}()

	// run the api listener
	apiServer := server.NewAPIServer(a.Config.APIListener, a.DynamoDB)
	go func() {
		log.Fatal(apiServer.Listen())
	}()

	log.Info("Swerve redirector")

	// wait for signals
	<-sigchan

	log.Info("Exit application")
}

// NewApplication creates new instance
func NewApplication() *Application {
	return &Application{
		Config: configuration.NewConfiguration(),
	}
}
