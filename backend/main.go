package main

import (
	"backend/cmd/set"
	"backend/internal/router"
	"log"
	"os"
)

func main() {
	var err error
	var pwd set.SetResponse
	if os.Args[0] == "restore" {
		log.Println("Restoring...")
		pwd, err = set.Restore()
	} else if os.Args[0] == "repwd" {
		log.Println("Password Resetting...")
		pwd, err = set.RePwd()
	} else {
		router.Start()
	}
	if err != nil {
		panic(err)
	} else {
		log.Printf("The new dashboard password is %s\n", pwd)
	}
}
