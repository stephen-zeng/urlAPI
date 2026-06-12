package static

import "embed"

/** @brief 内嵌前端静态构建产物文件系统。 */
//go:embed dist/*
var StaticFS embed.FS
