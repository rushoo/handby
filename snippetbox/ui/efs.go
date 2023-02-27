package ui

import "embed"

// 特殊的注释指令，编译时引导将ui/html、ui/static嵌入Files变量

//go:embed "html" "static"
var Files embed.FS
