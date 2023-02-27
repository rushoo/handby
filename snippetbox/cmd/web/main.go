package main

import (
	"flag"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox/internal/models"
)

// 为了让自定义logger其它函数也可以使用，一种方式是创建全局变量，但是依赖项多了就会很混乱
// 这里使用依赖注入，会更整洁一些,application这里可以理解为是一个依赖组
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippet       *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	//在命令行参数中添加端口配置,另外还可以通过配置文件来做这件事
	addr := flag.String("addr", ":4001", "HTTP network address")
	dsn := flag.String("dsn",
		"xm:123456@tcp(10.0.1.17:3306)/snippetbox?parseTime=True",
		"MySQL data source name")
	flag.Parse()

	//自定义日志，分别是INFO、ERROR前缀输出，Lshortfile用于显示错误日志位置，这里的日志是并发安全的
	//通过将日志记录在os.Stdout、os.Stderr，还方便以2>>error.log这种方式重定向收集，也可直接使用一个file表示日志直接写到文件里
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	//将要解析的html模板文件加载到内存
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	formDecoder := form.NewDecoder()
	//实例化依赖组，包含两条日志依赖项,和数据库连接对象、formDecoder
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippet:       &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	//通过自定义httpServer，将HTTP server产生的日志用自定义的日志收集
	infoLog.Printf("start sever on %s", *addr)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	err = srv.ListenAndServe()
	//应该避免在main函数之外使用Panic()、Fatal()，把错误传回来main，然后在此exit
	//另外在前期开发过程中一些可能如空指针引用等人为错误可以在函数内fatal或panic以方便调试
	errorLog.Fatal(err)
}
