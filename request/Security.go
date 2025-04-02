package request

import "urlAPI/security"

type Security struct {
	General security.General
	TxtGen  security.TxtGen
	ImgGen  security.ImgGen
	Rand    security.Rand
	WebImg  security.WebImg
}
