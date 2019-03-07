package goproxyrpc_test

import (
	"context"
	"fmt"
	"github.com/autom8ter/goproxyrpc"
	"github.com/autom8ter/goproxyrpc/pkg/errors"
	"github.com/autom8ter/goproxyrpc/pkg/testing/gen/echo"
	"github.com/autom8ter/goproxyrpc/pkg/util"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

//Eazy Peazy with Golang first class functions
func RegisterFunc() goproxyrpc.RegisterFunc {
	return func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
		return echopb.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	}
}

func TestNewGoProxyRpc(t *testing.T) {
	ctx := context.Background()
	go func() {
		if err := echopb.NewEchoServer().Serve(grpc.NewServer()); err != nil {
			t.Fatal(errors.New("grpc server error", err).String())
		}
	}()
	go goproxyrpc.New(ctx, &goproxyrpc.Config{
		EnvPrefix:    "ECHO",
		DialOptions:  []grpc.DialOption{grpc.WithInsecure()},
		RegisterFunc: RegisterFunc(),
	}).ListenServe(ctx)
	ur, err := url.Parse("http://localhost:8080/v1/echo")
	if err != nil {
		t.Fatal(errors.New("failed to parse proxy url", err))
	}
	resp, err := http.Post(ur.String(), "application/json", strings.NewReader(util.ToPrettyJsonString(&echopb.EchoMessage{
		Value: "hello there",
	})))
	if err != nil {
		t.Fatal(errors.New("failed to ping proxy", err))
	}
	if resp == nil {
		t.Fatal(errors.New("empty response", errors.NewErr("empty response returned")))
	}
	fmt.Print(util.ToPrettyJsonString(resp.Body))
}
