package framework

import "context"

type Context struct {
	context.Context
}

func (ctx *Context) Render() *Render          { return &Render{} }
func (ctx *Context) Bind(container any) error { return nil }
func (ctx *Context) Run()                     {}
