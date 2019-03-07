package health

import (
	"github.com/autom8ter/goproxyrpc/pkg/errors"
	"net"
	"time"
)

type Pinger struct {
	Endpoint string
	Do       func() error
}

func New(endpoint string) *Pinger {
	return &Pinger{Endpoint: endpoint}
}

func (p *Pinger) Once() *Pinger {
	p.Do = func() error {
		_, err := net.DialTimeout("tcp", p.Endpoint, 250*time.Second)
		if err != nil {
			errors.New("failed to ping grpc endpoint", err).FailIfErr()
		}
		return nil
	}
	return p
}
