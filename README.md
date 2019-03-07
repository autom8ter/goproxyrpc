# GoProxyRPC

# Overview

## Registering a Grpc Gaterway Service
```go
//Eazy Peazy with Golang first class functions
func RegisterFunc() goproxyrpc.RegisterFunc {
	return func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
		return echopb.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	}
}
```

## Starting a Remote Grpc Service
```go
if err := echopb.NewEchoServer().Serve(grpc.NewServer()); err != nil {
			log.Fatal(errors.New("grpc server error", err).String())
		}
```


## Starting the GoProxyRPC server
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
