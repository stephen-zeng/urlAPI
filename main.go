package urlAPI

import (
	"log"
	"os"
	"urlAPI/command"
	"urlAPI/handler"
)

func main() {
	log.Println("The default password is 123456")
	command.Arg(os.Args)
	handler.Handler(command.Port)
	os.Exit(0)
}
