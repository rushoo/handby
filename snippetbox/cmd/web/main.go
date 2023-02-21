package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	//通过内置方法以 "./ui/static"路径创建一个文件服务
	//文件服务器接受以"/static/"开头的请求.在具体访问前去掉"/static"前缀
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", Home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	go func() {
		log.Println("start sever on :4000")
		err := http.ListenAndServe(":4000", mux)
		//Note that any error returned by http.ListenAndServe() is always non-nil.
		log.Fatal(err)
	}()
	log.Println("start sever on :5000")
	err := http.ListenAndServe(":5000", mux)
	log.Fatal(err)
}
