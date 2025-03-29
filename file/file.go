package file

import "embed"

//go:embed ssfonts.ttf
var Font embed.FS

//go:embed icon/*
var Icons embed.FS

//go:embed logo/*
var Logos embed.FS

//go:embed setting.json
var Settings embed.FS
