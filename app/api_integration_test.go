// +build integration

package app

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	app "github.com/axelspringer/swerve/app/testdata"

	"github.com/axelspringer/swerve/config"
	"github.com/axelspringer/swerve/log"
	"github.com/axelspringer/swerve/model"
	"github.com/stretchr/testify/assert"
)

const (
	defaultUser = "spooky"
)

type responseRedirect struct {
	Data model.Redirect `json:"data"`
}

var (
	token             string
	cfg               config.Swerve
	testdata          app.Testdata
	httpClient        http.Client
	baseUrlApi        string
	baseUrlRedirecter string
	updateRedirect    = struct {
		FromDomain  string          `json:"redirect_from"`
		ToDomain    string          `json:"redirect_to"`
		Description string          `json:"description"`
		Code        int             `json:"code"`
		Promotable  bool            `json:"promotable"`
		PathMap     []model.PathMap `json:"path_map"`
	}{
		"replaceme",
		"new.com",
		"new description",
		307,
		false,
		[]model.PathMap{{
			From: "/foo",
			To:   "/bar",
		}},
	}
)

func TestMain(m *testing.M) {
	os.Setenv("SWERVE_DYNO_ENDPOINT", "http://localhost:8000")
	os.Setenv("SWERVE_DYNO_REGION", "eu-west-1")
	os.Setenv("SWERVE_DYNO_AWS_KEY", "0")
	os.Setenv("SWERVE_DYNO_AWS_SECRET", "0")
	os.Setenv("SWERVE_DYNO_TABLE_USERS", "Swerve_Users")
	os.Setenv("SWERVE_DYNO_TABLE_REDIRECTS", "Swerve_Redirects")
	os.Setenv("SWERVE_DYNO_TABLE_DOMAINS_CERTCACHE", "Swerve_CertCache")
	os.Setenv("SWERVE_DYNO_BOOTSTRAP", "true")
	os.Setenv("SWERVE_DYNO_DEFAULT_USER", defaultUser)
	os.Setenv("SWERVE_DYNO_DEFAULT_PW", "$2y$12$0n7R3fPqu5E/UhTbxN0qNOQEsYvBzYVmC3eTw1DS5Jbt2MThoYfrG")
	os.Setenv("SWERVE_USE_PEBBLE", "true")
	os.Setenv("SWERVE_PEBBLE_CA_URL", "http://localhost:15000/roots/0")
	os.Setenv("SWERVE_LETSENCRYPT_URL", "https://localhost:14000/dir")
	os.Setenv("SWERVE_API_UI_URL", "*")
	os.Setenv("SWERVE_API_VERSION", "v1")
	os.Setenv("SWERVE_LOG_LEVEL", "error")

	fh, err := ioutil.ReadFile("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(fh, &testdata)
	if err != nil {
		log.Fatal(err)
	}
	if len(testdata.Data) == 0 {
		log.Fatal(fmt.Errorf("No testdata available"))
	}

	a := NewApplication()
	cfg = a.Config
	httpClient = http.Client{
		Timeout: time.Second * 5,
	}
	baseUrlApi = fmt.Sprintf("http://127.0.0.1:%d", cfg.API.Listener)
	if err := a.Setup(); err != nil {
		log.Fatal(err)
	}
	go a.Run()
	err = waitUntilServerIsUpAndReady(a.Config.API.Listener)
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

func Test_APILogin(t *testing.T) {
	testCases := []testcase{
		{"valid login", defaultUser, "mytestpw", http.StatusOK},
		{"invalid login", defaultUser, "noValidPW", http.StatusUnauthorized},
	}

	for _, te := range testCases {
		url := fmt.Sprintf("%s/login", baseUrlApi)
		t.Logf("run -> %s, user: %s wanted statuscode: %d, url: %s", te.name, te.user, te.expectedStatuscode, url)
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

	url := fmt.Sprintf(baseUrlApi + "/" + cfg.API.Version + "/redirects/")
	for _, data := range testdata.Data {
		payload, err := json.Marshal(data.Redirect)
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
		assert.Equal(t, resp.StatusCode, http.StatusOK)
	}

	if t.Failed() {
		t.Fatal(fmt.Errorf("POST redirect faild, skipp all other tests"))
	}
}

func Test_GetRedirects(t *testing.T) {
	if checkEmptyToken(t) {
		return
	}

	for _, data := range testdata.Data {
		url := fmt.Sprintf(baseUrlApi + "/" + cfg.API.Version + "/redirects/" + data.Redirect.RedirectFrom)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
		resp, err := httpClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		r := responseRedirect{}
		assert.Equal(t, nil, json.NewDecoder(resp.Body).Decode(&r), "unmarshal response against model.Redirect")
		assert.Equal(t, data.Redirect.RedirectFrom, r.Data.RedirectFrom, "unmarshal response against model.Redirect")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "check for equal response code")
	}
}

func Test_Redirects(t *testing.T) {
	httpClient := getHttpClientWithPebbleIntermediateCert(t, cfg.HttpListener)
	httpsClient := getHttpClientWithPebbleIntermediateCert(t, cfg.HttpsListener)

	for _, data := range testdata.Data {
		for _, testcase := range data.Cases {
			rawurl, err := url.Parse(testcase.Call)
			if err != nil {
				t.Fatal(err)
			}
			qstr := ""
			if rawurl.RawQuery != "" {
				qstr = "?"
			}
			port := cfg.HttpsListener
			client := httpsClient
			if rawurl.Scheme == "http" {
				port = cfg.HttpListener
				client = httpClient
			}
			callUrl := fmt.Sprintf("%s://%s:%d%s%s%s", rawurl.Scheme, rawurl.Host, port, rawurl.Path, qstr, rawurl.RawQuery)
			resp, err := client.Get(callUrl)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, data.Redirect.Code, resp.StatusCode, "compare statuscode")
			assert.Equal(t, testcase.Expected, resp.Header.Get("Location"), fmt.Sprintf("(Prom: %t, Call:%s", data.Redirect.Promotable, testcase.Call))
		}
	}

}

func Test_UpdateRedirects(t *testing.T) {
	if checkEmptyToken(t) {
		return
	}

	fromDomain := testdata.Data[0].Redirect.RedirectFrom
	reqUrl := fmt.Sprintf(baseUrlApi + "/" + cfg.API.Version + "/redirects/" + fromDomain)

	updateRedirect.FromDomain = fromDomain
	payload, err := json.Marshal(updateRedirect)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", reqUrl, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "")
}

func Test_CheckSuccessfulUpdate(t *testing.T) {
	if checkEmptyToken(t) {
		return
	}

	fromDomain := testdata.Data[0].Redirect.RedirectFrom
	url := fmt.Sprintf(baseUrlApi + "/" + cfg.API.Version + "/redirects/" + fromDomain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	r := responseRedirect{}
	assert.Equal(t, nil, json.NewDecoder(resp.Body).Decode(&r), "unmarshal response against model.Redirect")
	assert.Equal(t, updateRedirect.ToDomain, r.Data.RedirectTo, "check to domain are equal")
	assert.Equal(t, fromDomain, r.Data.RedirectFrom, "check previous updated redirect.RedirectFrom")
	assert.Equal(t, updateRedirect.Description, r.Data.Description, "check previous updated redirect.description")
	assert.Equal(t, updateRedirect.Promotable, r.Data.Promotable, "check previous updated redirect.Promotable")
	assert.Equal(t, updateRedirect.PathMap, r.Data.PathMaps, "check previous updated redirect.PathMap")
	assert.Equal(t, updateRedirect.Code, r.Data.Code, "check response code")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "check response code")
}

func Test_DeleteRedirect(t *testing.T) {
	if checkEmptyToken(t) {
		return
	}

	for _, data := range testdata.Data {
		url := fmt.Sprintf(baseUrlApi + "/" + cfg.API.Version + "/redirects/" + data.Redirect.RedirectFrom)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
		resp, err := httpClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, resp.StatusCode, "check response code")
	}
}

func Test_RedirectExistAfterDelete(t *testing.T) {
	if checkEmptyToken(t) {
		return
	}

	for _, data := range testdata.Data {
		url := fmt.Sprintf(baseUrlApi + "/" + cfg.API.Version + "/redirects/" + data.Redirect.RedirectFrom)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Cookie", fmt.Sprintf("token=%s", token))
		resp, err := httpClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "check if redirect exists after deleting")
	}
}

// ****** Tests end here *******
func getHttpClientWithPebbleIntermediateCert(t *testing.T, port int) *http.Client {
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM([]byte(cfg.ACM.PebbleCA))

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	tlsconfig := &tls.Config{
		RootCAs:   caPool,
		ClientCAs: caPool,
	}
	tr := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if !strings.Contains(addr, ":15000") {
				addr = fmt.Sprintf("127.0.0.1:%d", port)
			}
			return dialer.DialContext(ctx, network, addr)
		},
		TLSClientConfig: tlsconfig,
	}

	pebbleClient := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := pebbleClient.Get("https://127.0.0.1:15000/roots/0")
	if err != nil {
		t.Fatal(err)
	}
	intermediateCert, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	caPool.AppendCertsFromPEM(intermediateCert)

	return pebbleClient
}

func waitUntilServerIsUpAndReady(apiport int) error {
	for i := 0; i < 30; i++ {
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/health", apiport))
		if err != nil {
		}
		if resp != nil && resp.StatusCode == http.StatusOK {
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
