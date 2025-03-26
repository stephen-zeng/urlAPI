package request

import "urlAPI/processor"

type Processor struct {
	Dashboard processor.Dashboard
	Download  processor.Download
	TxtSum    processor.TxtSum
	TxtGen    processor.TxtGen
	ImgGen    processor.ImgGen
	WebImg    processor.WebImg
	Rand      processor.Rand
	Operation processor.Interface
}
