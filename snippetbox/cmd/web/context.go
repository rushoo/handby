package main

// 自定义类型别名，可以有效避免代码逻辑混乱
type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
