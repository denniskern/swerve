package app

import (
	"fmt"
	nethttp "net/http"
	"time"

	"github.com/axelspringer/swerve/log"
)

// Only the HTTPHandler function activates the http01 challenge,
// so its important that on startup the
func (a *Application) ensureHttpCall() error {
	log.Debug("ensureHttpCall: make local http call to activate http01 challenge")
	for i := 0; i < 15; i++ {
		time.Sleep(time.Second * 2)
		resp, err := nethttp.Get(fmt.Sprintf("http://127.0.0.1:%d", a.Config.HttpListener))
		if err != nil {
			log.Errorf("ensureHttpCall: %v", err)
		}
		if resp != nil && resp.StatusCode < nethttp.StatusInternalServerError {
			log.Debugf("ensureHttpCall: successfully reached the http server, this is needed to enable http01 challenge")
			return nil
		}
	}
	return fmt.Errorf("ensureHttpCall: can't reach server on http://127.0.0.1:%d, http01 will not be available", a.Config.HttpListener)
}
