package goproxyrpc_test

import (
	"context"
	"github.com/autom8ter/goproxyrpc"
	"github.com/autom8ter/goproxyrpc/pkg/errors"
	"github.com/autom8ter/goproxyrpc/pkg/testing/gen/echo"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"testing"
)

//Eazy Peazy with Golang first class functions
func RegisterFunc() goproxyrpc.RegisterFunc {
	return func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
		return echopb.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	}
}

func TestGoProxyRpc(t *testing.T) {
	ctx := context.Background()
	go func() {
		if err := echopb.NewEchoServer().Serve(grpc.NewServer()); err != nil {
			t.Fatal(errors.New("grpc server error", err).String())
		}
	}()
	goproxyrpc.New(ctx, &goproxyrpc.Config{
		EnvPrefix:    "ECHO",
		DialOptions:  []grpc.DialOption{grpc.WithInsecure()},
		RegisterFunc: RegisterFunc(),
	}).ListenServe(ctx)
}
