package main

import (
	"os"
	"urlAPI/command"
	"urlAPI/database"
	"urlAPI/handler"
)

func main() {
	command.Arg(os.Args)
	handler.Handler(command.Port)
	defer database.Disconnect()
	os.Exit(0)
}
