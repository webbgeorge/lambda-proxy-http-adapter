package lambda_proxy_http_adapter

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestGetHttpHandler(t *testing.T) {
	handler := func(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		assert.Equal(t, "POST", r.HTTPMethod)
		assert.Equal(t, "/users/123", r.Path)
		assert.Equal(t, "123", r.PathParameters["userId"])
		assert.Equal(t, "123", r.QueryStringParameters["abc"])
		assert.Equal(t, "123", r.MultiValueQueryStringParameters["abc"][0])
		assert.Equal(t, "application/json", r.Headers["Content-Type"])
		assert.Equal(t, "application/json", r.MultiValueHeaders["Content-Type"][0])
		assert.Equal(t, "req_body", r.Body)
		assert.Equal(t, "varValue1", r.StageVariables["var1"])

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Single-Value-Key": "single_value",
				"Mixed-Value-Key":  "single_value",
			},
			MultiValueHeaders: map[string][]string{
				"Multi-Value-Key": {"multi_value_1", "multi_value_2"},
				"Mixed-Value-Key": {"multi_value_1", "multi_value_2"},
			},
			Body: "response_body",
		}, nil
	}

	httpHandler := GetHttpHandler(handler, "/users/{userId}", map[string]string{"var1": "varValue1"})
	testServer := httptest.NewServer(httpHandler)
	defer testServer.Close()

	res, err := http.Post(testServer.URL+"/users/123?abc=123", "application/json", strings.NewReader("req_body"))

	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, []string{"single_value"}, res.Header["Single-Value-Key"])
	assert.Equal(t, []string{"multi_value_1", "multi_value_2"}, res.Header["Multi-Value-Key"])
	assert.Equal(t, []string{"single_value", "multi_value_1", "multi_value_2"}, res.Header["Mixed-Value-Key"])

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "response_body", string(body))
}

func TestGetHttpHandlerWithContext(t *testing.T) {
	handler := func(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		assert.Equal(t, "POST", r.HTTPMethod)
		assert.Equal(t, "/users/123", r.Path)
		assert.Equal(t, "123", r.PathParameters["userId"])
		assert.Equal(t, "123", r.QueryStringParameters["abc"])
		assert.Equal(t, "123", r.MultiValueQueryStringParameters["abc"][0])
		assert.Equal(t, "application/json", r.Headers["Content-Type"])
		assert.Equal(t, "application/json", r.MultiValueHeaders["Content-Type"][0])
		assert.Equal(t, "req_body", r.Body)
		assert.Equal(t, "varValue1", r.StageVariables["var1"])

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Single-Value-Key": "single_value",
				"Mixed-Value-Key":  "single_value",
			},
			MultiValueHeaders: map[string][]string{
				"Multi-Value-Key": {"multi_value_1", "multi_value_2"},
				"Mixed-Value-Key": {"multi_value_1", "multi_value_2"},
			},
			Body: "response_body",
		}, nil
	}

	httpHandler := GetHttpHandlerWithContext(handler, "/users/{userId}", map[string]string{"var1": "varValue1"})
	testServer := httptest.NewServer(httpHandler)
	defer testServer.Close()

	res, err := http.Post(testServer.URL+"/users/123?abc=123", "application/json", strings.NewReader("req_body"))

	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, []string{"single_value"}, res.Header["Single-Value-Key"])
	assert.Equal(t, []string{"multi_value_1", "multi_value_2"}, res.Header["Multi-Value-Key"])
	assert.Equal(t, []string{"single_value", "multi_value_1", "multi_value_2"}, res.Header["Mixed-Value-Key"])

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "response_body", string(body))
}

func TestGetHttpHandler_Error(t *testing.T) {
	handler := func(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{}, assert.AnError
	}

	httpHandler := GetHttpHandler(handler, "/", nil)
	testServer := httptest.NewServer(httpHandler)
	defer testServer.Close()

	res, err := http.Post(testServer.URL, "application/json", strings.NewReader("req_body"))

	assert.Nil(t, err)
	assert.Equal(t, 500, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "error", string(body))
}

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
