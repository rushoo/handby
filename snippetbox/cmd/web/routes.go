package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	//使用httprouter
	router := httprouter.New()

	//自定义空路径请求的响应
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	//创建一个文件服务，httprouter.Handler要求 *something(而非*) 以表示路径下全匹配，something不要求有实际意义
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*httprouter-filepath-rule", http.StripPrefix("/static", fileServer))

	chain := New(app.sessionManager.LoadAndSave, noSurf)
	router.Handler(http.MethodGet, "/", chain.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", chain.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", chain.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", chain.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", chain.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", chain.ThenFunc(app.userLoginPost))

	//须身份验证的路由，添加身份验证
	chain.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", chain.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", chain.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", chain.ThenFunc(app.userLogoutPost))

	standard := New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
	//return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
