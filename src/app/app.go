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
	// set the table prefix
	db.DBTablePrefix = a.Config.TablePrefix
	// database connection
	var err error
	a.DynamoDB, err = db.NewDynamoDB(&a.Config.DynamoDB, a.Config.Bootstrap)
	if err != nil {
		log.Fatalf("Can't setup db connection %#v", err)
	}
	// check api static file path
	if a.Config.APIClientStaticPath == "" {
		log.Fatal("You have to specify the api client static path")
	}
	// check path folder
	if dstat, err := os.Stat(a.Config.APIClientStaticPath); err != nil || !dstat.IsDir() {
		log.Fatalf("The api client static path is invalid. err %#v", err)
	}
	// cert manager
	a.Certificates = certificate.NewManager(a.DynamoDB, a.Config.StagingCA)
	// cache preload
	a.Certificates.CertCache.UpdateDomainCache()
	// backgroud update ticker
	a.Certificates.CertCache.Observe()
}

// Run the application
func (a *Application) Run() {
	log.Info("Swerve redirector")
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
	apiServer := server.NewAPIServer(a.Config.APIListener, a.Config.APIClientStaticPath, a.DynamoDB)
	go func() {
		log.Fatal(apiServer.Listen())
	}()
	// wait for signals
	<-sigchan
}

// NewApplication creates new instance
func NewApplication() *Application {
	return &Application{
		Config: configuration.NewConfiguration(),
	}
}
