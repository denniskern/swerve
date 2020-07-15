// +build integration

package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/axelspringer/swerve/config"
	"github.com/axelspringer/swerve/model"
	"github.com/stretchr/testify/assert"
)

var (
	token             string
	cfg               *config.Configuration
	httpClient        http.Client
	baseUrlApi        string
	baseUrlRedirecter string
)

func TestMain(m *testing.M) {
	os.Setenv("SWERVE_STAGING", "true")
	os.Setenv("SWERVE_DB_ENDPOINT", "http://localhost:8000")
	os.Setenv("SWERVE_DB_REGION", "eu-west-1")
	os.Setenv("SWERVE_USERS", "Users")
	os.Setenv("SWERVE_DOMAINS", "Domains")
	os.Setenv("SWERVE_DOMAINS_TLSSWERVE_DOMAINS_TLS_CACHE", "DomainsTLSCache")
	os.Setenv("SWERVE_API_UI_URL", "*")
	os.Setenv("SWERVE_API_VERSION", "v1")
	a := NewApplication()
	cfg = a.Config
	httpClient = http.Client{}
	baseUrlApi = fmt.Sprintf("http://127.0.0.1:%d", cfg.API.Listener)
	if err := a.Setup(); err != nil {
		log.Fatal(err)
	}
	go a.Run()
	waitUntilServerIsUpAndReady([]int{a.Config.API.Listener, a.Config.HTTPListenerPort, a.Config.HTTPSListenerPort})
	m.Run()
}

func Test_APILOGIN(t *testing.T) {
	testCases := []struct {
		name               string
		user               string
		pass               string
		expectedStatuscode int
	}{
		{"valid login", "dkern", "$2a$12$gh.TtSizoP0JFLHACOdIouPr42713m6k/8fH8jKPl0xQAUBk0OIdS", http.StatusOK},
		{"invalid login", "dkern", "noValidPW", http.StatusUnauthorized},
	}

	for _, te := range testCases {
		url := fmt.Sprintf("%s/login", baseUrlApi)
		payload := []byte(fmt.Sprintf(`{"username":"%s", "pwd":"%s"}`, te.user, te.pass))
		resp, err := httpClient.Post(url, "content-type: application/json", bytes.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}
		if te.name == "valid login" && resp.StatusCode == http.StatusOK {
			for _, cookie := range resp.Cookies() {
				if cookie.Name == "token" {
					token = cookie.Value
					break
				}
			}
		}
		assert.Equal(t, resp.StatusCode, te.expectedStatuscode, fmt.Sprintf("Test %s", te.name))
	}
}

func Test_PostRedirects(t *testing.T) {
	testCases := []struct {
		name               string
		payload            model.Redirect
		expectedStatuscode int
	}{
		{
			"domain ilk.org",
			model.Redirect{
				RedirectFrom: "ilk.org",
				Description:  "Testdomain1",
				RedirectTo:   "bild.de",
				Promotable:   false,
				Code:         301,
				PathMaps: []model.PathMap{
					{
						From: "/",
						To:   "/digital",
					},
				},
			},
			200,
		},
	}

	url := fmt.Sprintf(baseUrlApi + "/" + cfg.API.Version + "/redirects/")
	for _, te := range testCases {
		payload, err := json.Marshal(te.payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
		resp, err := httpClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, resp.StatusCode, te.expectedStatuscode, te.name)
	}

}

func waitUntilServerIsUpAndReady(ports []int) error {
	counter := 0
	for {
		log.Printf("counter %d", counter)
		portreachable := true
		if counter >= 2 {
			return fmt.Errorf("not all ports are available")
		}
		for _, p := range ports {
			conn, err := net.Listen("tcp4", fmt.Sprintf(":%d", p))
			if err != nil {
				portreachable = false
				continue
			}
			conn.Close()
		}
		if portreachable {
			return nil
		}
		counter++
		time.Sleep(time.Second * 1)
	}
	return fmt.Errorf("reach end of wait function ... should not happen!")
}
