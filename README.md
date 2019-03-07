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
