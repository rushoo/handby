package models

import (
	"errors"
)

// ErrNoRecord 通过在数据层将错误分级，应用层可以不关心具体数据，只看数据处理结果
var ErrNoRecord = errors.New("models: no matching record found")
