package request

type Request struct {
	DB        DB
	Processor Processor
	Security  Security
}
