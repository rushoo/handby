package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	//通过内置方法以 "./ui/static"路径创建一个文件服务
	//文件服务器接受以"/static/"开头的请求.在具体访问前去掉"/static"前缀
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	//chain := New(app.recoverPanic, app.logRequest, secureHeaders)
	chain := New()
	chain.Append(app.recoverPanic)
	chain.Append(app.logRequest, secureHeaders)
	return chain.Then(mux)
	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
