package models

import (
	"errors"
)

// 通过在数据层将错误分级，应用层可以不关心具体数据，只看数据处理结果
var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)
