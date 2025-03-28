package request

import "urlAPI/security"

type Security struct {
	General security.General
	TxtGen  security.TxtGen
	TxtSum  security.TxtSum
	ImgGen  security.ImgGen
	Rand    security.Rand
	WebImg  security.WebImg
}
