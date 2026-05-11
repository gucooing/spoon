package external

import (
	"context"
	"github.com/gucooing/spoon/logger"
)

type Request interface {
	GetMsgID() uint32
	GetBody() []byte
}

type Response interface {
	GetMsgID() uint32
	GetBody() []byte
}

type Session interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	UUID() string
	GetSessionID() uint64
}

type Router interface {
	RegisterHandler(id uint32, h Handler)
	Handle(ctx Context, req Request) (Response, error)
}

type Context interface {
	context.Context
	GetSession() Session
}

type Handler func(ctx Context, req Request) (Response, error)

type HandlerFunc func(Context)
type HandlersChain []HandlerFunc

func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

type SessionManager interface {
	NewSession(session Session)
	StartSession(ctx context.Context, session Session)
	Logger(log logger.Logger)
}
