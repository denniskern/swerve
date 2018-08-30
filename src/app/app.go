package app

import (
	"github.com/hammi85/swerve/src/log"

	"os"
	"os/signal"
	"syscall"

	"github.com/hammi85/swerve/src/api"
	"github.com/hammi85/swerve/src/certificate"
	"github.com/hammi85/swerve/src/configuration"
	"github.com/hammi85/swerve/src/db"
	"github.com/hammi85/swerve/src/http"
	"github.com/hammi85/swerve/src/https"
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
	httpsServer := https.NewServer(a.Config.HTTPSListener, a.Certificates)
	go func() {
		log.Fatal(httpsServer.Listen())
	}()

	// run the http listener
	httpServer := http.NewServer(a.Config.HTTPListener, a.Certificates)
	go func() {
		log.Fatal(httpServer.Listen())
	}()

	// run the api listener
	apiServer := api.NewServer(a.Config.APIListener, a.DynamoDB)
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
