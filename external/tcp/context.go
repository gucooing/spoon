package tcp

import (
	"context"
	"math"

	"github.com/gucooing/spoon/external"
)

// Context 请求上下文
type Context struct {
	context.Context
	req     external.Request // 解析后的传输层
	session *Session         // 网络层实例

	index    int8
	handlers external.HandlersChain // 中间件
}

func (c *Context) GetSession() external.Session {
	return c.session
}

func (c *Context) Next() {
	c.index++
	for c.index < safeInt8(len(c.handlers)) {
		if c.handlers[c.index] != nil {
			c.handlers[c.index](c)
		}
		c.index++
	}
}

func safeInt8(n int) int8 {
	if n > math.MaxInt8 {
		return math.MaxInt8
	}
	return int8(n)
}
