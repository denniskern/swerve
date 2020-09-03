package testclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"

	"github.com/axelspringer/swerve/cmd/livetest/config"
	"github.com/axelspringer/swerve/log"
)

type Client struct {
	cfg   config.LiveConfig
	token string
	Runs  int
}

func New(cfg config.LiveConfig) Client {
	c := Client{}
	c.cfg = cfg
	c.token = c.Login()
	return c
}

func (c *Client) Login() string {
	var err error
	var token string

	req := request{}
	req.Url = fmt.Sprintf("%s/login", c.cfg.APIBaseURL)
	c.PrintCall(req.Url)
	req.Body, err = json.Marshal(login{
		Username: c.cfg.DynoUser,
		Pwd:      c.cfg.DynoPw,
	})
	if err != nil {
		log.Fatal(err)
	}
	req.Method = http.MethodPost
	res, err := doRequest(req)
	if err != nil {
		log.Fatal(err)
	}
	for _, cookie := range res.Cookies {
		if cookie.Name == "token" {
			token = cookie.Value
			break
		}
	}
	if token == "" {
		spew.Dump(res)
		log.Fatal("could not fetch login token")
	}
	c.PrintOk(fmt.Sprintf("code: %d token: %s", res.StatusCode, token))
	return token
}

func (c *Client) NextRun() int {
	c.Runs++
	return c.Runs
}

func (c *Client) Version() {
	r := request{}
	r.Url = fmt.Sprintf("%s/version", c.cfg.APIBaseURL)
	c.PrintCall(r.Url)
	res, err := doRequest(r)
	if err != nil {
		log.Fatal(err)
	}
	v := version{}
	err = json.Unmarshal(res.RawData, &v)
	if err != nil {
		log.Fatal(err)
	}
	c.PrintOk(fmt.Sprintf("code: %d Message: %s", v.Code, v.Message))
}

func (c *Client) InsertRedirects() {
	header := make(map[string]string)
	header["Cookie"] = fmt.Sprintf("token=%s", c.token)
	for _, v := range c.cfg.Data {
		url := fmt.Sprintf("%s/%s/redirects/", c.cfg.APIBaseURL, c.cfg.APIVersion)
		c.PrintCall(url)
		var err error
		req := request{}
		req.Body, err = json.Marshal(v.Redirect)
		if err != nil {
			log.Fatal(err)
		}
		req.Headers = header
		req.Url = url
		req.Method = http.MethodPost
		res, err := doRequest(req)
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode != 200 {
			spew.Dump(res)
			log.Fatal(fmt.Errorf("can't push redirect"))
		}
		c.PrintOk(fmt.Sprintf("redirect %s -> %s created", v.Redirect.RedirectFrom, v.Redirect.RedirectTo))
	}
}

func (c *Client) TestRedirects() {
	var errors []error
	header := make(map[string]string)
	header["Cookie"] = fmt.Sprintf("token=%s", c.token)

	for _, v := range c.cfg.Data {
		for _, test := range v.Cases {
			c.PrintCall(test.Call)
			req := request{}
			req.Url = fmt.Sprintf("%s", test.Call)
			req.Headers = header
			req.Method = http.MethodGet
			res, err := doRequest(req)
			if err != nil {
				errors = append(errors, err)
			}
			c.PrintOk(res.Header.Get("Location"))
		}
	}

	if len(errors) > 0 {
		spew.Dump(errors)
		log.Fatal("Not all calls and expected results matches")
	}
}
