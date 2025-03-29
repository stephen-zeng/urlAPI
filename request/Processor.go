package request

import "urlAPI/processor"

type Processor struct {
	Session  processor.Session
	Download processor.Download
	TxtGen   processor.TxtGen
	ImgGen   processor.ImgGen
	WebImg   processor.WebImg
	Rand     processor.Rand
}
