package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// Valid 没有错误则返回true
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError 添加新的错误信息到 FieldErrors map
func (v *Validator) AddFieldError(key, message string) {
	// 若容器没初始化，先将其初始化
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	//再将不存在的错误存起来，已有的就忽略
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField 将不合要求(错误)的信息添加到map
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank 判断非空元素
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars 判断字符串是否含少于给定数量的rune
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedInt 限定合法字符元素
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
