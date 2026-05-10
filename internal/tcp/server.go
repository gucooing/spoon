package tcp

import (
	"context"
	"net"
)

type ServerOption func(*Server)

type Server struct {
	address string
	network string

	lis net.Listener
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.listen(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) listen() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	return nil
}
