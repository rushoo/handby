### 使用https自签证书
```
$ mkdir tls && cd tls

# 获取执行路径，这里是 /d/Go/，随后在此路径下找到generate_cert.go
$ which go
/d/Go/bin/go

$ go run  /d/Go/src/crypto/tls/generate_cert.go  --rsa-bits=2048 --host=localhost
2023/02/27 15:50:16 wrote cert.pem
2023/02/27 15:50:16 wrote key.pem

#更改服务启动
srv.ListenAndServe()    -->    srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	
#完成后注意证书的保密，这里是实验就一起丢git上了	

```