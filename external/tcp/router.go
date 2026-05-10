package tcp

import (
	"github.com/gucooing/spoon/errors"
	"github.com/gucooing/spoon/external"
)

type Router struct {
	router     map[uint32]external.Handler // 路由
	notHandler external.Handler            // 默认方法
}

func NewRouter() *Router {
	return &Router{
		router:     make(map[uint32]external.Handler),
		notHandler: defaultNotHandler,
	}
}

func (router *Router) RegisterHandler(id uint32, h external.Handler) {
	router.router[id] = h
}

func defaultNotHandler(ctx external.Context, req external.Request) (external.Response, error) {
	return nil, errors.New(errors.UnknownRouter, "NotHandler", "UnknownRouter")
}

func (router *Router) Handle(ctx external.Context, req external.Request) (external.Response, error) {
	handle, ok := router.router[req.GetMsgID()]
	if !ok {
		return router.notHandler(ctx, req)
	}
	rsp, err := handle(ctx, req)
	return rsp, err
}
