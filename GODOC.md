# goproxyrpc
--
    import "github.com/autom8ter/goproxyrpc"


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
