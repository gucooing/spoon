package tcp

import (
	"context"
	"net"

	"github.com/gucooing/spoon/external"
)

func (s *Server) NewConn(ctx context.Context, conn net.Conn) external.Session {
	session := &Session{
		conn:      conn,
		ctx:       ctx,
		read:      s.read,
		sessionID: 0,
	}

	return session
}
