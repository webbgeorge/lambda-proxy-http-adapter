package lambda_proxy_http_adapter

import (
	"context"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

type APIGatewayProxyHandler func(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func GetHttpHandler(
	lambdaHandler APIGatewayProxyHandler,
	resourcePathPattern string,
	stageVariables map[string]string,
) http.Handler {
	return getHttpHandler(func(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return lambdaHandler(r)
	}, resourcePathPattern, stageVariables)
}

type APIGatewayProxyHandlerWithContext func(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func GetHttpHandlerWithContext(
	lambdaHandler APIGatewayProxyHandlerWithContext,
	resourcePathPattern string,
	stageVariables map[string]string,
) http.Handler {
	return getHttpHandler(lambdaHandler, resourcePathPattern, stageVariables)
}

func getHttpHandler(
	lambdaHandler APIGatewayProxyHandlerWithContext,
	resourcePathPattern string,
	stageVariables map[string]string,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		proxyResponse, err := lambdaHandler(r.Context(), events.APIGatewayProxyRequest{
			Resource:                        resourcePathPattern,
			Path:                            r.URL.Path,
			HTTPMethod:                      r.Method,
			Headers:                         singleValue(r.Header),
			MultiValueHeaders:               r.Header,
			QueryStringParameters:           singleValue(r.URL.Query()),
			MultiValueQueryStringParameters: r.URL.Query(),
			PathParameters:                  parsePathParams(resourcePathPattern, r.URL.Path),
			StageVariables:                  stageVariables,
			Body:                            string(body),
		})

		if err != nil {
			// write a generic error, the same as API GW would if an error was returned by handler
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`error`))
			return
		}

		writeResponse(w, proxyResponse)
	})
}

func singleValue(multiValueMap map[string][]string) map[string]string {
	singleValueMap := make(map[string]string)
	for k, mv := range multiValueMap {
		if len(mv) > 0 {
			singleValueMap[k] = mv[0]
		}
	}
	return singleValueMap
}

func parsePathParams(pathPattern string, path string) map[string]string {
	exp := regexp.MustCompile(`{(\w+)}`)
	pathPatternExp := regexp.MustCompile(exp.ReplaceAllString(pathPattern, `(?P<$1>\w+)`))

	subMatches := pathPatternExp.FindStringSubmatch(path)
	subMatchNames := pathPatternExp.SubexpNames()

	params := make(map[string]string)
	for i, paramName := range subMatchNames {
		if paramName == "" || len(subMatches) < i {
			continue
		}
		params[paramName] = subMatches[i]
	}

	return params
}

func writeResponse(w http.ResponseWriter, proxyResponse events.APIGatewayProxyResponse) {
	for k, v := range proxyResponse.Headers {
		w.Header().Add(k, v)
	}

	for k, vs := range proxyResponse.MultiValueHeaders {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(proxyResponse.StatusCode)
	_, _ = w.Write([]byte(proxyResponse.Body))
}
