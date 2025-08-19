package frontend

import (
	"embed"
	"io/fs"
)

// 前端资源
var (
	//go:embed dist/*
	static   embed.FS
	SubFS, _ = fs.Sub(static, "dist")
)
