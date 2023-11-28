package testhelper

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gabukuro/insta-gift-api/internal/pkg/router"
	"github.com/stretchr/testify/assert"
)

type (
	Input struct {
		Method  string
		Path    string
		Headers http.Header
		Body    io.Reader
	}

	Expected struct {
		Status    int
		Body      string
		ErrString string
	}

	Test struct {
		Description string
		Input       *Input
		Expected    *Expected
	}

	NewRequestParams struct {
		Method  string
		Path    string
		Body    io.Reader
		Headers http.Header
		Profile *string
	}
)

func TestAllRequests(tests []Test, t *testing.T, routerInstance *router.Router) {
	for _, test := range tests {
		test := test
		t.Run(test.Description, func(t *testing.T) {
			req := NewRequest(&NewRequestParams{
				Method:  test.Input.Method,
				Path:    test.Input.Path,
				Body:    test.Input.Body,
				Headers: test.Input.Headers,
			})

			res, err := Request(t, req, routerInstance)
			VerifyError(t, "", err)

			assert.Equal(t, test.Expected.Status, res.StatusCode)
			VerifyHTTPResponseBodyJSON(t, test.Expected.Body, res.Body, nil)
		})
	}
}

func NewRequest(opts *NewRequestParams) *http.Request {
	req := httptest.NewRequest(
		opts.Method,
		opts.Path,
		opts.Body,
	)

	if opts.Headers != nil {
		for key, value := range opts.Headers {
			req.Header.Add(key, value[0])
		}
	}

	if _, ok := req.Header["Content-Type"]; !ok {
		req.Header.Add("Content-Type", "application/json")
	}

	return req
}

func Request(t *testing.T, req *http.Request, routerInstance *router.Router) (*http.Response, error) {
	res, err := routerInstance.App().Test(req, -1)
	return res, err
}

func ParseResponseBody(res *http.Response, i interface{}) string {
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bodyBytes, &i)
	if err != nil {
		panic(err)
	}

	return string(bodyBytes)
}
