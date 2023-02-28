package main

import (
	"bytes"
	"github.com/go-playground/assert/v2"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

//	func ping(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte("OK"))
//	}
func TestPing(t *testing.T) {
	//构建条件，模拟函数调用： ping(rr, r)
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	ping(rr, r)

	//分析结果
	rs := rr.Result()
	assert.Equal(t, rs.StatusCode, http.StatusOK)
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
func TestPing2(t *testing.T) {
	//模拟端到端全过程测试
	app := newTestApplication(t)
	// 使用 httptest.NewTLSServer() 创建一个测试server，随后使用 ts.Close() 将其关闭
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// 测试server的地址含于 ts.URL ，这里使用 ts.Client().Get() 模拟客户端发送一个 /ping 请求
	code, _, body := ts.get(t, "/ping")
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}
