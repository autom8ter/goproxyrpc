# GoProxyRPC

# Overview

## Config (./config.yaml)

### Example:
```yaml
endpoint: "localhost:3000"
cors:
  allow-origin:
  allow-credentials:
  allow-methods:
  allow-headers:
proxy:
  port: 8080
  api-prefix: "/"
```

## Registering a Grpc Gaterway Service

function signature: `func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)`

```go
//Eazy Peazy with Golang first class functions
func RegisterFunc() goproxyrpc.RegisterFunc {
	return func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
		return echopb.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
		//^^^generated from grpc gateway plugin	
    }
}
```


## Starting the GoProxyRPC server
Dials the endpoint in your config:

```go

func main() {
	ctx := context.Background()
	goproxyrpc.New(ctx, &goproxyrpc.Config{
		EnvPrefix:    "ECHO",
		DialOptions:  []grpc.DialOption{grpc.WithInsecure()},
		RegisterFunc: RegisterFunc(),
	}).ListenServe(ctx)
}

```

## Curl the Proxy Server (pkg/testing/echo)

```text
#!/usr/bin/env bash

curl --header "Content-Type: application/json"   --request POST   --data '{"say":"hello"}'   http://localhost:8080/v1/echo
```
### Response

`{"say":"echoed: hello"}`

## Usage

## Usage

#### type Config

```go
type Config struct {
	EnvPrefix    string
	DialOptions  []grpc.DialOption
	RegisterFunc RegisterFunc
}
```

Config holds the necessary non-config file configurations needed to create a
Proxy

#### type Proxy

```go
type Proxy struct {
	http.Handler
}
```

Proxy is a REST-gRPC reverse proxy server

#### func  NewProxy

```go
func NewProxy(ctx context.Context, cfg *Config) *Proxy
```
NewProxy creates a new REST-gRPC proxy server. If a jwt_key is found in your
config file, the endpoint will be reject all requests that dont provide a valid
bearer token.

#### type RegisterFunc

```go
type RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
```

RegisterFunc registers a grpc endpoint from the generated RegisterfromEndpoint
function from the grpc-gateway protoc plugin

#### func  NewRegisterFunc

```go
func NewRegisterFunc(fn func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)) RegisterFunc
```
NewRegisterFunc is a helper to create a RegisterFunc

