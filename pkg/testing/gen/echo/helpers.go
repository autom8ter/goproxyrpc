package echopb

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type EchoServer struct{}

func NewEchoServer() *EchoServer {
	return &EchoServer{}
}

func (e *EchoServer) Echo(ctx context.Context, r *EchoMessage) (*EchoMessage, error) {
	fmt.Printf("rpc request Echo(%q)\n", r.Say)
	r.Say = "echoed: " + r.Say
	return r, nil
}

func (e *EchoServer) Serve(s *grpc.Server) error {
	RegisterEchoServiceServer(s, e)
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		return err
	}
	return s.Serve(lis)
}
