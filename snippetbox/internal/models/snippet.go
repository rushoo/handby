package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	//数据插入，返回插入的结果的id值
	//关于placeholder，MySQL,SQL-Server和SQLite使用 ?, PostgreSQL使用 $N 这种方式
	stmt := `INSERT INTO snippet (title, content, created, expires)
			VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	// 获取最新插入记录的id
	//id, err := result.RowsAffected()返回影响的行数
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// 这里有个类型转换，数据样本小不必担心精度损失
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippet WHERE expires > UTC_TIMESTAMP() AND id = ?`

	s := &Snippet{}
	// QueryRow最多返回一行结果，通过row.Scan()将返回的值复制到snippet结构体
	// 在scan过程中如果遇到列值为null时会出错，如果假定这种空数据是正常的，可以将结构体字段类型设置为sql.NullTypeName
	// 可以参考这里https://gist.github.com/alexedwards/dc3145c8e2e6d2fd6cd9
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//查找的数据不存在
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippet WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Query() 可以返回多行数据
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	var snippets []*Snippet
	// 在此defer前先要上面的错误检查，对于空对象close会导致panic
	// 这里的close很有必要，因为结果集在打开时会保活潜在的sql连接
	defer rows.Close()
	for rows.Next() {
		s := &Snippet{}
		// 注意参数的地址引用
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	// Err()记录了在上述iteration过程中产生的错误
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
