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

	//通过内置方法以 "./ui/static"路径创建一个文件服务
	//文件服务器接受以"/static/"开头的请求.在具体访问前去掉"/static"前缀
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/", http.StripPrefix("/static", fileServer))
	//router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	//chain := New(app.recoverPanic, app.logRequest, secureHeaders)
	chain := New()
	chain.Append(app.recoverPanic)
	chain.Append(app.logRequest, secureHeaders)
	return chain.Then(router)
	//return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
