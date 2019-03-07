package goproxyrpc

import (
	"fmt"
	"github.com/autom8ter/goproxyrpc/pkg/config"
	"github.com/autom8ter/goproxyrpc/pkg/errors"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	_ "google.golang.org/genproto/googleapis/rpc/errdetails" // Pull in errdetails
)

type GoProxyRpc struct {
	http.Handler
	v    *viper.Viper
	port int
}

type RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

type Config struct {
	EnvPrefix    string
	DialOptions  []grpc.DialOption
	RegisterFunc RegisterFunc
}

func New(ctx context.Context, cfg *Config) *GoProxyRpc {
	mux := http.NewServeMux()
	v := config.SetupViper(cfg.EnvPrefix)
	c := &config.ProxyConfig{
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

	return &GoProxyRpc{
		Handler: mux,
		port:    v.GetInt("proxy.port"),
		v:       v,
	}
}

func (g *GoProxyRpc) ListenServe(ctx context.Context) {
	if g.Handler == nil {
		logrus.Fatalf(`nil handler: use "goproxyrpc.NewGoProxyRpc(ctx context.Context, cfg *Config) *GoProxyRpc" to initialize the proxy`)
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", g.port),
		Handler: g,
	}
	signalRunner(
		func() {
			logrus.Infof("launching http server on %v", server.Addr)
			if err := server.ListenAndServe(); err != nil {
				errors.New("failed to launch proxy server", err).FailIfErr()
			}
		},
		func() {
			shutdown, _ := context.WithTimeout(ctx, 10*time.Second)
			server.Shutdown(shutdown)
		})
}

// SignalRunner runs a runner function until an interrupt signal is received, at which point it
// will call stopper.
func signalRunner(runner, stopper func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	go func() {
		runner()
	}()

	logrus.Info("hit Ctrl-C to shutdown proxy")
	select {
	case <-signals:
		stopper()
	}
}
