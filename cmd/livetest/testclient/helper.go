package testclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

func doRequest(r request) (response, error) {
	res := response{}
	c := resty.New()
	c.SetTimeout(time.Duration(20 * time.Second))
	c.SetRetryCount(3)
	c.SetRetryWaitTime(5 * time.Second)
	c.SetRetryMaxWaitTime(15 * time.Second)
	c.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	c.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}))
	c.SetHeaders(r.Headers)
	c.SetFormData(r.FormData)
	c.SetQueryParams(r.QueryParams)

	resp, err := c.R().EnableTrace().SetBody(r.Body).Execute(r.Method, r.Url)

	if err != nil {
		return res, err
	}

	if resp == nil {
		res.Message = fmt.Sprintf("Internal error (%v)\n", err)
		return res, err
	}
	res.StatusCode = resp.StatusCode()
	res.TraceInfo = resp.Request.TraceInfo()
	res.RawData = resp.Body()
	res.Header = resp.Header()
	res.Cookies = resp.Cookies()

	return res, err
}

func (c *Client) PrintCall(url string) {
	fmt.Printf("(%d) call %s\n", c.NextRun(), url)
}

func (c *Client) PrintOk(msg string) {
	fmt.Printf(" -> [OK] %s\n", msg)
}
