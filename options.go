package spoon

import (
	"context"
	"os"
)

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type Option func(*options)

type options struct {
	ctx  context.Context
	sigs []os.Signal

	id      string
	name    string
	version string

	servers []Server
}

func ID(id string) Option {
	return func(o *options) { o.id = id }
}

func Name(name string) Option {
	return func(o *options) { o.name = name }
}

func Version(version string) Option {
	return func(o *options) { o.version = version }
}

func Servers(servers ...Server) Option {
	return func(o *options) { o.servers = servers }
}
