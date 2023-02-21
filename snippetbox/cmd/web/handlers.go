package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func Home(w http.ResponseWriter, r *http.Request) {
	//避免根路径全匹配
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	//解析模板文件，这是从磁盘中直接读的
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}
	//遇到传入参数不定的情况,三点省略号写法...接受可变数量的参数
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err)
		http.Error(w, "解析模板文件失败", http.StatusInternalServerError)
		return
	}
	//将模板文件"base"的解析结果作为response返回
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "解析模板文件失败", http.StatusInternalServerError)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	log.Println("snippetView-RequestURI", r.RequestURI)
	//从请求中获取id
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("请求的id：%s有误", idStr), http.StatusBadRequest)
		return
	}
	//w.Write([]byte(fmt.Sprintf("Display a snippet%d", id)))
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Create a new snippet..."))
}
