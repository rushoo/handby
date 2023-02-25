package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"strconv"
)

//通过app的依赖注入，这里的handler就可以获取到数据库连接并使用model{}操作数据库

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//避免根路径全匹配
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	snippets, err := app.snippet.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := newTemplateData()
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	//从请求中获取id
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	//根据id去数据库select
	snippet, err := app.snippet.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := newTemplateData()
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	//定义接受的请求类型
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7
	id, err := app.snippet.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// 303重定向到页面"/snippet/view?id=id"显示刚刚插入的结果
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
