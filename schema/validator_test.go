package schema

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/stretchr/testify/assert"
)

func Test_Validator(t *testing.T) {
	testcases := []struct {
		Redirect       string
		ExpectedErrors int
		Description    string
	}{
		{
			Redirect: `{
			"redirect_from": "rd1.swervetest.de",
			"redirect_to": "www.bild.de",
			"promotable": true,
			"code": 301,
			"description": "1, +path_map, +prom",
			"path_map": [
				{ "from": "/", "to": "/t" }
			]}`,
			ExpectedErrors: 0,
			Description:    "RD1 Valid RD",
		},
		{
			Redirect: `{
			"redirect_from": "rd1.swervetest.de",
			"redirect_to": "4.de",
			"promotable": false,
			"code": 307
			}`,
			ExpectedErrors: 0,
			Description:    "RD2 Valid RD",
		},
		{
			Redirect: `{
			"redirect_from": "rd1.swervetest.de",
			"promotable": false,
			"code": 307
			}`,
			ExpectedErrors: 1,
			Description:    "RD4 Invalid RD, missing retirect_to",
		},
		{
			Redirect: `{
			"redirect_from": "rd1.swervetest.de",
			"redirect_to": "bild.de",
			"promotable": false,
			"code": 306
			}`,
			ExpectedErrors: 1,
			Description:    "RD5 Invalid RD, wrong statuscode",
		},
		{
			Redirect: `{
			"redirect_from": "rd1.swervetest.de",
			"redirect_to": "www.bild.de",
			"promotable": true,
			"code": 309,
			"description": "1, +path_map, +prom",
			"path_map": [
				{ "from": "/" }
			]}`,
			ExpectedErrors: 2,
			Description:    "RD6 Invalid RD. wrong path_map",
		},
		{
			Redirect: `{
			"foo":"bar",
			"redirect_from": "rd1.swervetest.de",
			"redirect_to": "www.bild.de",
			"promotable": true,
			"code": 301,
			"description": "1, +path_map, +prom",
			"path_map": [
				{ "from": "/" }
			]}`,
			ExpectedErrors: 2,
			Description:    "RD7 Invalid RD. wrong path_map",
		},
	}

	vali := New()
	for _, tc := range testcases {
		errors := vali.ValidateRedirect([]byte(tc.Redirect))
		if len(errors) != tc.ExpectedErrors {
			spew.Dump(errors)
		}
		assert.Len(t, vali.ValidateRedirect([]byte(tc.Redirect)), tc.ExpectedErrors, tc.Description)
	}

}
