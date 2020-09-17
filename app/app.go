package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/axelspringer/swerve/schema"

	"github.com/axelspringer/swerve/config"
	phm "github.com/axelspringer/swerve/prometheus"

	"github.com/axelspringer/swerve/acm"
	"github.com/axelspringer/swerve/api"
	"github.com/axelspringer/swerve/cache"
	"github.com/axelspringer/swerve/database"
	"github.com/axelspringer/swerve/http"
	"github.com/axelspringer/swerve/https"
	"github.com/axelspringer/swerve/log"
	"github.com/axelspringer/swerve/model"
	"github.com/pkg/errors"
)

// NewApplication creates a new instance
func NewApplication() *Application {
	return &Application{
		Config: config.Get(),
	}
}

// Setup sets up the application
func (a *Application) Setup() error {
	log.SetupLogger(a.Config.LogLevel, a.Config.LogFormat)

	db, err := database.NewDatabase(a.Config.DynamoDB)
	if err != nil {
		return errors.WithMessage(err, ErrDatabaseServiceCreate)
	}
	if a.Config.DynamoDB.Bootstrap {
		err = db.Prepare()
		if err != nil {
			return errors.WithMessage(err, ErrTablePrepare)
		}
	}

	jsonValidator := schema.New()
	prom := phm.NewPHM()
	a.Cache = cache.NewCache(db)

	controlModel := model.NewModel(a.Cache)

	autocertManager := acm.NewACM(a.Cache.AllowHostPolicy,
		a.Cache,
		a.Config.ACM)

	a.HTTPServer = http.NewHTTPServer(controlModel.GetRedirectByDomain,
		autocertManager.HTTPHandler,
		a.Config.HttpListener,
		prom.WrapHandler)

	a.HTTPSServer = https.NewHTTPSServer(controlModel.GetRedirectByDomain,
		autocertManager.GetCertificate,
		a.Config.HttpsListener,
		prom.WrapHandler)
	a.APIServer = api.NewAPIServer(controlModel, jsonValidator, a.Config.API)

	return nil
}

// Run starts the application
func (a *Application) Run() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	if err := a.Cache.Observe(a.Config.CacheInterval); err != nil {
		log.Fatal(errors.WithMessage(err, ErrCacheObserver))
	}

	go func() {
		log.Fatal(a.APIServer.Listen())
	}()

	go func() {
		log.Fatal(a.HTTPServer.Listen())
	}()

	go func() {
		log.Fatal(a.HTTPSServer.Listen())
	}()

	if !a.Config.DisableHTTPChallenge {
		err := a.ensureHttpCall()
		if err != nil {
			log.Error(err)
		}
	}

	<-sigchan
}
