# 使用go操作数据库时一些注意事项
### 空值
go里通过row.Scan()将返回的值复制到目标结构体，在此过程中，如果遇到列值为null时会出错，   
如果假定这种空数据是正常的，可以将结构体字段类型设置为sql.NullTypeName   
可以参考这里https://gist.github.com/alexedwards/dc3145c8e2e6d2fd6cd9
```
type Book struct {
  Isbn  string
  Title  sql.NullString
  Author sql.NullString
  Price  sql.NullFloat64
}
err := rows.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price)
```

### 事务处理
执行具体的sql操作时，go从数据库连接池中取出一个连接对象，即使两个相邻的sql操作也无法保证  
这源于同一个sql连接，比如涉及到锁表操作时，就必须保证同一个连接锁表和解锁，否则会死锁。  
```
type ExampleModel struct {
	DB *sql.DB
}
func (m *ExampleModel) ExampleTransaction() error {
	// 调用 Begin() 方法开启一个事务
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	// Defer tx.Rollback() 用以保证事务未顺利执行时状态回滚，这也能确保执行结束放回连接池
	defer tx.Rollback()

	// 通过事务执行sql操作
	_, err = tx.Exec("INSERT INTO ...")
	if err != nil {
		return err
	}
	// 执行另一个事务
	_, err = tx.Exec("UPDATE ...")
	if err != nil {
		return err
	}
	// 如果上述过程执行无误，执行结果就会在此最终提交并生效
	err = tx.Commit()
	return err
}
```

### 预处理语句(Prepared statements)
一、SQL 语句的执行处理    
1、即时 SQL   
　　一条 SQL 在 DB 接收到最终执行完毕返回，大致的过程如下：   
　　1. 词法和语义解析；   
　　2. 优化 SQL 语句，制定执行计划；    
　　3. 执行并返回结果；     
　　如上，一条 SQL 直接是走流程处理，一次编译，单次运行，此类普通语句被称作 Immediate Statements （即时 SQL）。    
2、预处理 SQL    
　　但是，绝大多数情况下，某需求某一条 SQL 语句可能会被反复调用执行，或者每次执行的时候只有个别的值不同    
    （比如 select 的 where 子句值不同，update 的 set 子句值不同，insert 的 values 值不同）。    
    如果每次都需要经过上面的词法语义解析、语句优化、制定执行计划等，则效率就明显不行了。特别是多次重复执行且复杂的多join语句。       
    预编译语句就是将此类 SQL 语句中的值用占位符替代，一般称这类语句叫Prepared Statements。   
    SQL注入一般通过把SQL命令插入到Web表单提交或输入域名或页面请求的查询字符串，最终达到欺骗服务器执行恶意的SQL命令。   
    使用预编译语句可以很好地规避sql注入风险。   
3、关于prepared statement的另一方面
    从下面代码示例可以观察到，当一个prepared statement被执行时，它会先根据连接池内的某个连接创建，并会记住这个连接。   
    当该预编译语句被再次执行时，它会去寻找此前那个连接，如果连接无法使用就会再选择一个连接再次重复创建。    
    对于负载压力大的场景，这可能会导致语句它不断地选择新的连接re-prepared，这个过程会带来新的复杂性。  
    所以假使在复杂性和性能之间做一个权衡，建议常规的增删改查语句就用常规的执行方式了。    
```
// 将预编译语句和sql连接一起嵌入model中
type ExampleModel struct {
	DB *sql.DB
	InsertStmt *sql.Stmt
}

// 使用构造方法设置预编译语句
func NewExampleModel(db *sql.DB) (*ExampleModel, error) {
	// 使用 Prepare 方法基于当前连接池创建一条prepared语句
	insertStmt, err := db.Prepare("INSERT INTO ...")
	if err != nil {
		return nil, err
	}
	// 以model结构体来储存这个预编译语句
	return &ExampleModel{db, insertStmt}, nil
}
// 将model作为方法接收者便可调用预编译语句了
func (m *ExampleModel) Insert(args...) error {
	// 这里的预编译语句还支持 Query、QueryRow等方法
	_, err := m.InsertStmt.Exec(args...)
	return err
}

func main() {
	db, _ := sql.Open(...) //delebrately ignored err
	defer db.Close()
	
	// 用db实例化一个含预编译语句的model
	exampleModel, err := NewExampleModel(db)
	if err != nil {
		errorLog.Fatal(err)
	}
	// 使用 Defer 确保程序推出前关闭连接
	defer exampleModel.InsertStmt.Close()
}
```    

