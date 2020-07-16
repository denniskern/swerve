// +build integration

package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/axelspringer/swerve/config"
	"github.com/axelspringer/swerve/model"
	"github.com/davecgh/go-spew/spew"
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
	a.Config.Database.Secret = "0"
	a.Config.Database.Key = "0"
	spew.Dump(a.Config.Database)
	go a.Run()
	err := waitUntilServerIsUpAndReady(a.Config.API.Listener)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

type testcase struct {
	name               string
	user               string
	pass               string
	expectedStatuscode int
}

func Test_APILOGIN(t *testing.T) {
	testCases := []testcase{
		{"valid login", "dkern", "$2a$12$gh.TtSizoP0JFLHACOdIouPr42713m6k/8fH8jKPl0xQAUBk0OIdS", http.StatusOK},
		{"invalid login", "dkern", "noValidPW", http.StatusUnauthorized},
	}

	for _, te := range testCases {
		url := fmt.Sprintf("%s/login", baseUrlApi)
		t.Logf("run -> %s, user: %s pw: %s wanted statuscode: %d, url: %s", te.name, te.user, te.pass, te.expectedStatuscode, url)
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
		assert.Equal(t, te.expectedStatuscode, resp.StatusCode, fmt.Sprintf("Test %s", te.name))
	}
}

func Test_PostRedirects(t *testing.T) {
	if checkEmptyToken(t) {
		return
	}
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

func waitUntilServerIsUpAndReady(apiport int) error {
	for i := 0; i < 30; i++ {
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/health", apiport))
		if err != nil {
			log.Println("api server not ready yet ...")
		}
		if resp != nil && resp.StatusCode == 200 {
			log.Printf("lets start the tests, api is reachable")
			return nil
		}
		time.Sleep(time.Second * 1)
	}
	return fmt.Errorf("Can't reach api server on http://127.0.0.1:%d/health", apiport)
}

func checkEmptyToken(t *testing.T) bool {
	if token == "" {
		t.Fail()
		t.Log("test skipped because of missing token")
		return true
	}
	return false
}
