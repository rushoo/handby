package main

import (
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

/*
//http.Handler的定义，有唯一方法ServeHTTP(ResponseWriter, *Request)

	type Handler interface {
		ServeHTTP(ResponseWriter, *Request)
	}

//*http.ServeMux实现了方法，所以它也是http.Handler
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {...}
*/
func myMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: Execute our middleware logic here...
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Security-Policy",
				"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
			w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "deny")
			w.Header().Set("X-XSS-Protection", "0")
			next.ServeHTTP(w, r)
		})
}
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
			next.ServeHTTP(w, r)
		})
}
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					//设置响应头Connection: Close，1、触发go关闭此连接2、通知浏览器连接已关闭
					w.Header().Set("Connection", "close")
					app.serverError(w, fmt.Errorf("%s", err))
				}
			}()
			app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
			next.ServeHTTP(w, r)
		})
}
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//判断登录状态，若未登录则redirect to login
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		//设置Cache-Control为no-store，可以防止页面被缓存
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
