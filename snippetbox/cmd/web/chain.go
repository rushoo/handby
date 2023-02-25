package main

import "net/http"

/*
目的是以更直观地方式来构建中间件链
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
====>
	//chain := New(app.recoverPanic, app.logRequest, secureHeaders)
	chain := New()
	chain.Append(app.recoverPanic)
	chain.Append(app.logRequest, secureHeaders)
	return chain.Then(mux)
*/

// Constructor 函数类型构造器，参数http.Handler，返回值http.Handler
type Constructor func(http.Handler) http.Handler

// Chain 存放Constructor列表
type Chain struct {
	constructors []Constructor
}

func New(constructors ...Constructor) *Chain {
	//var cs []Constructor
	//cs = append(cs, constructors...)
	return &Chain{
		//constructors: cs,
		//与任意类型的nil值一样，我们可以用[]int(nil)类型转换表达式来生成一个对应类型slice的nil值。
		append([]Constructor(nil), constructors...),
	}
}

// Append 扩充chain的底层数组
func (c *Chain) Append(constructors ...Constructor) {
	newCons := make([]Constructor, 0, len(c.constructors)+len(constructors))
	newCons = append(newCons, c.constructors...)
	newCons = append(newCons, constructors...)

	c.constructors = newCons
}
func (c *Chain) Extend(chain *Chain) {
	c.Append(chain.constructors...)
}

// Then 调用mux，以此使用注册的中间件----- New(m1, m2, m3).Then(h) ==>  m1(m2(m3(h)))
func (c *Chain) Then(h http.Handler) http.Handler {
	//如果没自定义的mux，就使用默认的DefaultServeMux
	if h == nil {
		h = http.DefaultServeMux
	}
	for i := range c.constructors {
		//从后往前依次wrap, (h) ==> m3(h) ==> m2(m3(h)) ==> m1(m2(m3(h)))
		h = c.constructors[len(c.constructors)-1-i](h)
	}

	return h
}
func (c *Chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return c.Then(nil)
	}
	return c.Then(fn)
}
