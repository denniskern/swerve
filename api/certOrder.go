package api

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/axelspringer/swerve/config"

	"github.com/axelspringer/swerve/log"
)

func (api *API) CertOrder(data []byte) error {
	/* POC
	- Save Domain with Hostname into DB
	- Initiate SSL Request from redirect_from and first path mapping
	- Wait for response and
	*/
	order, err := api.Model.CreateCertOrderFromJSON(data)
	if err != nil {
		log.Debugf("CertOrder error"+": %s", err.Error())
		return err
	}

	client := getHttpClient(api.Config)
	url := getUrl(order.Domain, api.Config)

	log.Debugf("call %s for ssl cert creation", url)
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	if resp == nil {
		return fmt.Errorf("resp of http.Get %s is nil", url)
	}

	if resp != nil {
		log.Debugf("resp of https://"+order.Domain+" Statuscode: %d", resp.StatusCode)
	}

	return nil
}

func getHttpClient(cfg *config.Configuration) *http.Client {
	tlsConfig := &tls.Config{}
	if cfg.ACM.UsePebble {
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: time.Second * 25,
	}
}

func getUrl(domain string, cfg *config.Configuration) string {
	url := fmt.Sprintf("https://%s", domain)
	if cfg.ACM.UsePebble {
		url = fmt.Sprintf("https://%s:%d", domain, cfg.HTTPSListenerPort)
	}
	return url
}
