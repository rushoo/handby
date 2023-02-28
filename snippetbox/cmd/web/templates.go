package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"snippetbox/internal/models"
	"snippetbox/ui"
	"time"
)

// 分设两个字段解析单一或多个snippet,
// 增加一层嵌套,这样在模板中解析时字段名称就是{{.Snippet.xxx}},而非直接{{.xxx}},方便处理多种数据来源
type templateData struct {
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// 将模板文件先解析到内存中,这样就不至于对于每个来自客户端的请求都从磁盘解析渲染一次
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	//拿到给定路径下的所有文件、、使用embed files
	//pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		// 获取文件名
		name := filepath.Base(page)
		//// 根据目的将待解析的文件合在一起,一次解析,但这样就写死了,所以采用分步解析
		files := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, files...)
		if err != nil {
			return nil, err
		}
		// 建立文件名映射
		cache[name] = ts
	}
	return cache, nil
}

/*
自定义模板方法,例:
1、先定义一个方法humanDate(),用于时间格式转换
2、将自定义方法加入template.FuncMap，得到一个含自定义方法的funcMap
3、自定义template对象，将此funcMap通过template.Funcs()在解析模板前注册
*/
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	//format前先转为utc格式
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// 稍后将此自定义funMap注册
var functions = template.FuncMap{
	"humanDate": humanDate,
}
