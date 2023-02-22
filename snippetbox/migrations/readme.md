该目录下包含项目所用数据库的sql migration文件
```
# CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# mysql分配数据库读写用户--常规操作
#  create user xiaohong identified by '123456';
#  grant select,insert,update,delete on snippetbox.* to 'xiaohong'@'%';

# mysql分配数据库全表用户--有建表删表权限
#  create user xm identified by '123456';
#  grant all privileges on snippetbox.* to 'xm'@'%';

# migrate create -seq -ext=.sql -dir=./migrations create_snippet_table
# DSN="mysql://xm:123456@tcp(10.0.1.17:3306)/snippetbox?charset=utf8&parseTime=True&loc=Local"
# migrate -path=./migrations -database=$DSN up
```