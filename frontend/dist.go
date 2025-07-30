package frontend

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var static embed.FS

var SubFS, _ = fs.Sub(static, "dist")
