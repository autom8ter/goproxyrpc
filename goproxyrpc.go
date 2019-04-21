//go:generate godocdown -o GODOC.md

package goproxyrpc

import (
	"github.com/autom8ter/authzero/grants"
	"github.com/autom8ter/goproxyrpc/pkg/config"
	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net/http"
	"os"

	_ "google.golang.org/genproto/googleapis/rpc/errdetails" // Pull in errdetails
)

//Proxy is a REST-gRPC reverse proxy server
type Proxy struct {
	http.Handler
	v    *viper.Viper
	port int
}

//RegisterFunc registers a grpc endpoint from the generated RegisterfromEndpoint function from the grpc-gateway protoc plugin
type RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

//NewRegisterFunc is a helper to create a RegisterFunc
func NewRegisterFunc(fn func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)) RegisterFunc {
	return fn
}

//Config holds the necessary non-config file configurations needed to create a Proxy
type Config struct {
	EnvPrefix    string
	DialOptions  []grpc.DialOption
	RegisterFunc RegisterFunc
}

//NewProxy creates a new REST-gRPC proxy server. If a jwt_key is found in your config file, the endpoint will be reject all requests that dont provide a valid bearer token.
func NewProxy(ctx context.Context, cfg *Config) *Proxy {
	mux := http.NewServeMux()
	v := config.SetupViper(cfg.EnvPrefix)
	c := &config.ProxyConfig{
		JWTKey:               v.GetString("jwt_key"),
		Endpoint:             v.GetString("endpoint"),
		LogLevel:             v.GetString("log_level"),
		LogHeaders:           v.GetBool("log_headers"),
		CorsAllowOrigin:      v.GetString("cors.allow-origin"),
		CorsAllowCredentials: v.GetString("cors.allow-credentials"),
		CorsAllowMethods:     v.GetString("cors.allow-methods"),
		CorsAllowHeaders:     v.GetString("cors.allow-headers"),
		ApiPrefix:            v.GetString("proxy.api-prefix"),
	}
	gw := config.SetupGateway()
	if err := cfg.RegisterFunc(ctx, gw, c.Endpoint, cfg.DialOptions); err != nil {
		logrus.Fatalf("failed to register grpc gateway from endpoint: %s", err.Error())
	}
	mux.Handle(c.ApiPrefix, handlers.CustomLoggingHandler(os.Stdout, http.StripPrefix(c.ApiPrefix[:len(c.ApiPrefix)-1], config.AllowCors(c, gw)), config.LogFormatter(c)))
	if c.JWTKey != "" {
		return &Proxy{
			Handler: grants.Middleware(c.JWTKey, true).Handler(mux),
			port:    v.GetInt("proxy.port"),
			v:       v,
		}
	}
	return &Proxy{
		Handler: mux,
		port:    v.GetInt("proxy.port"),
		v:       v,
	}
}
