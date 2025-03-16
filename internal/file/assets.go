package file

import "embed"

//go:embed img/logo/*
var LogoFS embed.FS

//go:embed img/icon/*
var IconFS embed.FS
