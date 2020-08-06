package app

import (
	"fmt"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	phm "github.com/axelspringer/swerve/prometheus"

	"github.com/axelspringer/swerve/acm"
	"github.com/axelspringer/swerve/api"
	"github.com/axelspringer/swerve/cache"
	"github.com/axelspringer/swerve/config"
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
		Config: config.NewConfiguration(),
	}
}

// Setup sets up the application
func (a *Application) Setup() error {
	err := a.Config.FromEnv()
	if err != nil {
		return errors.WithMessage(err, ErrConfigInvalid)
	}
	a.Config.FromParameter()
	err = a.Config.Validate()
	if err != nil {
		return errors.WithMessage(err, ErrConfigInvalid)
	}

	log.SetupLogger(a.Config.LogLevel, a.Config.LogFormatter)

	db, err := database.NewDatabase(a.Config.Database)
	if err != nil {
		return errors.WithMessage(err, ErrDatabaseServiceCreate)
	}
	if a.Config.Bootstrap {
		err = db.Prepare()
		if err != nil {
			return errors.WithMessage(err, ErrTablePrepare)
		}
	}

	prom := phm.NewPHM()
	a.Cache = cache.NewCache(db)

	controlModel := model.NewModel(a.Cache)

	autocertManager := acm.NewACM(a.Cache.AllowHostPolicy,
		a.Cache,
		a.Config.ACM)

	a.HTTPServer = http.NewHTTPServer(controlModel.GetRedirectByDomain,
		autocertManager.HTTPHandler,
		a.Config.HTTPListenerPort,
		prom.WrapHandler)

	a.HTTPSServer = https.NewHTTPSServer(controlModel.GetRedirectByDomain,
		autocertManager.GetCertificate,
		a.Config.HTTPSListenerPort,
		prom.WrapHandler)
	a.APIServer = api.NewAPIServer(controlModel, a.Config.API)

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

	err := a.ensureHttpCall()
	if err != nil {
		log.Error(err)
	}

	<-sigchan
}

func (a *Application) ensureHttpCall() error {
	log.Debug("ensureHttpCall: make local http call to activate http01 challenge")
	for i := 0; i < 30; i++ {
		resp, err := nethttp.Get(fmt.Sprintf("http://127.0.0.1:%d", a.Config.HTTPListenerPort))
		if err != nil {
			log.Error(err)
		}
		if resp != nil && resp.StatusCode < nethttp.StatusInternalServerError {
			log.Debugf("successfully reached the http server, this is needed to enable http01 challenge")
			return nil
		}
		time.Sleep(time.Second * 1)
	}
	return fmt.Errorf("ensureHttpCall: can't reach server on http://127.0.0.1:%d, http01 will not be available", a.Config.HTTPListenerPort)
}
