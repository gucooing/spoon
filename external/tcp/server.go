package tcp

import (
	"context"
	"errors"
	"net"

	"github.com/google/uuid"
	"github.com/gucooing/spoon/external"
)

type ServerOption func(*Server)

// SetAddress 设置监听地址
func SetAddress(address string) ServerOption {
	return func(srv *Server) {
		srv.address = address
	}
}

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc

	address string
	network string

	lis  net.Listener // tcp
	read Read         // 解析传输层数据方法
	// 打包传输层数据方法
	sessionManager external.SessionManager // 会话管理

	router   external.Router        // 路由
	handlers external.HandlersChain // 中间件
}

// NewServer 新建一个对外实例
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		address:  "",
		network:  "tcp",
		lis:      nil,
		read:     defaultRead,
		router:   NewRouter(),
		handlers: nil,
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

// Router 获取路由
func (s *Server) Router() external.Router {
	return s.router
}

// SetRouter 设置路由
func (s *Server) SetRouter(router external.Router) {
	s.router = router
}

// Use 写入中间件
func (s *Server) Use(middleware ...external.HandlerFunc) {
	s.handlers = append(s.handlers, middleware...)
}

// NewReqContext 新建一个请求上下文
func (s *Server) NewReqContext(ctx context.Context) *Context {
	reqCtx := &Context{
		Context:  ctx,
		handlers: s.handlers,
	}
	return reqCtx
}

// Start 启动实例
func (s *Server) Start(ctx context.Context) error {
	sCtx, cancel := context.WithCancel(ctx)
	s.ctx = sCtx
	s.cancel = cancel
	s.sessionManager = NewSessionManager(s.ctx)
	if err := s.listen(); err != nil {
		return err
	}
	for {
		conn, err := s.lis.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if !errors.Is(err, net.ErrClosed) {
					return err
				}
				return nil
			}
		}
		session := &Session{
			conn: conn,
			uuid: uuid.NewString(),
			read: s.read,
		}
		s.sessionManager.NewSession(session)
		go s.sessionManager.StartSession(s.ctx, session)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}
	_ = s.lis.Close()
	return nil
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
