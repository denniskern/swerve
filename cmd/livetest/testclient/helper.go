package testclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
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
func (c *Client) PrintFail(msg string) {
	fmt.Printf(" -> [Fail] %s\n", msg)
}

func printResults(results [][]string) {
	println()
	table := tablewriter.NewWriter(os.Stdout)
	table.NumLines()
	table.SetHeader([]string{"RES", "PROM", "Call", "expected", "got", "w-c", "g-c", "description"})
	table.SetBorder(true)     // Set Border to false
	table.AppendBulk(results) // Add Bulk Data
	table.Render()
}
