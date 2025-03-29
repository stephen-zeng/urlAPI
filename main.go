package main

import (
	"log"
	"os"
	"urlAPI/command"
	"urlAPI/database"
	"urlAPI/handler"
)

func main() {
	log.Println("The default password is 123456")
	command.Arg(os.Args)
	handler.Handler(command.Port)
	defer database.Disconnect()
	os.Exit(0)
}
