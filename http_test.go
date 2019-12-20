package lambda_proxy_http_adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseParams(t *testing.T) {
	testCases := map[string]struct {
		pathPattern    string
		path           string
		expectedParams map[string]string
	}{
		"multiple params": {
			pathPattern: "/abc/{param1}/def/{param2}",
			path:        "/abc/xyz/def/123",
			expectedParams: map[string]string{
				"param1": "xyz",
				"param2": "123",
			},
		},
		"one param": {
			pathPattern: "/{name}",
			path:        "/dave",
			expectedParams: map[string]string{
				"name": "dave",
			},
		},
		"one param no matches": {
			pathPattern:    "/greet/{name}",
			path:           "/",
			expectedParams: map[string]string{},
		},
		"no params": {
			pathPattern:    "/users",
			path:           "/users",
			expectedParams: map[string]string{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			params := parsePathParams(tc.pathPattern, tc.path)
			assert.Equal(t, tc.expectedParams, params)
		})
	}
}
