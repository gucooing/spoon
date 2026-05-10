package tcp

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"

	"github.com/gucooing/spoon/external"
)

type SessionManager struct {
	ctx context.Context

	sync          sync.RWMutex
	bindCnt       atomic.Uint64               // 已登录的玩家数量
	loginCnt      atomic.Uint64               // 登录中的会话数量
	bindSessions  map[uint64]external.Session // 已绑定玩家的会话表
	loginSessions map[string]external.Session // 未绑定玩家的会话表
}

func NewSessionManager(ctx context.Context) *SessionManager {
	return &SessionManager{
		ctx:           ctx,
		sync:          sync.RWMutex{},
		bindCnt:       atomic.Uint64{},
		loginCnt:      atomic.Uint64{},
		bindSessions:  make(map[uint64]external.Session),
		loginSessions: make(map[string]external.Session),
	}
}

// SessionCount 获取当前绑定的玩家会话数量
func (sm *SessionManager) SessionCount() uint64 {
	return sm.bindCnt.Load()
}

// LoginCount 获取登录中的玩家会话数量
func (sm *SessionManager) LoginCount() uint64 {
	return sm.loginCnt.Load()
}

// NewSession 添加一个新会话
func (sm *SessionManager) NewSession(session external.Session) {
	sm.sync.Lock()
	defer sm.sync.Unlock()
	sm.loginSessions[session.UUID()] = session
}

// CheckSessionByUUID 通过uuid判断会话是否存在
func (sm *SessionManager) CheckSessionByUUID(uuid string) bool {
	sm.sync.RLock()
	defer sm.sync.RUnlock()
	_, ok := sm.loginSessions[uuid]
	return ok
}

func (sm *SessionManager) StartSession(ctx context.Context, session external.Session) {
	if !sm.CheckSessionByUUID(session.UUID()) {
		return
	}
	err := session.Start(ctx)
	if err != nil {

	}
}

type Session struct {
	conn   net.Conn
	uuid   string
	ctx    context.Context
	cancel context.CancelFunc

	read   Read            // 接收传输层方法
	router external.Router // 路由

	sessionID uint64 // 玩家id
}

// UUID 获取当前会话唯一id
func (s *Session) UUID() string {
	return s.uuid
}

// GetSessionID 获取会话id
func (s *Session) GetSessionID() uint64 {
	return s.sessionID
}

func (s *Session) Start(ctx context.Context) error {
	sCtx, cancel := context.WithCancel(ctx)
	s.ctx = sCtx
	s.cancel = cancel

	defer func() {
		err := recover()
		if err != nil {

		}
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			req, err := s.read(s.ctx, s.conn)
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					return nil
				}
				return err
			}
			reqCtx := &Context{
				req:      req,
				session:  s,
				index:    0,
				handlers: nil,
			}
			rsp, err := s.router.Handle(reqCtx, req)
			if err != nil {
				return err
			}
		}
	}
}

func (s *Session) Stop(ctx context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}
	_ = s.conn.Close()
	return nil
}
