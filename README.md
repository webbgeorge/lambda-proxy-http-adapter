# Lambda Proxy HTTP Adapter

Wrap your API Gateway Proxy Lambda handlers in a go net/http handler. Useful to make the Lambda Proxy handlers
compatible with the huge amount of tooling that already exists for net/http.

This adapter will work for many use cases, but was built specifically with https://github.com/steinfletcher/apitest in
mind. This adapter allows you to write apitests for your AWS API Gateway Proxy Lambda applications.

## Usage

### Basic

```go
package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/gaw508/lambda-proxy-http-adapter"
)

func main() {
	httpHandler := lambda_proxy_http_adapter.GetHttpHandler(apiGatewayProxyHandler, "/", nil)
	...
}

func apiGatewayProxyHandler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	...
}
```

### With Path Params

When using API Gateway, path routing is done on behalf of the application inside the API Gateway. Because of this, if
your application uses path params, you need to pass the resource path pattern. This is so the adapter is able
to parse the URL in the request, and is able to provide path params to your proxy handler.

```go
package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/gaw508/lambda-proxy-http-adapter"
)

func main() {
	httpHandler := lambda_proxy_http_adapter.GetHttpHandler(apiGatewayProxyHandler, "/users/{userId}", nil)
	...
}

func apiGatewayProxyHandler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	...
}
```

### With Stage Variables

API Gateway Proxy handlers allows you to use stage variables. If you use these in your handler, you have to pass them
into the adapter as a `map[string]string`.

```go
package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/gaw508/lambda-proxy-http-adapter"
)

func main() {
	httpHandler := lambda_proxy_http_adapter.GetHttpHandler(apiGatewayProxyHandler, "/", map[string]string{"var1": "var1value"})
	...
}

func apiGatewayProxyHandler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	...
}
```
