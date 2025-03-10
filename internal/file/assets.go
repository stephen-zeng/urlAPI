package file

import "embed"

//go:embed assets/logo/*
var LogoFS embed.FS

//go:embed assets/icon/*
var IconFS embed.FS
