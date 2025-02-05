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
	var port = "2233"
	setPort := false
	for _, arg := range os.Args {
		if setPort == true {
			port = arg
			setPort = false
			continue
		}
		if arg == "restore" {
			log.Println("Restoring...")
			pwd, err = set.Restore()
			if err != nil {
				panic(err)
			}
			log.Printf("The new dashboard password is %s\n", pwd.Pwd)
		}
		if arg == "repwd" {
			log.Println("Password Resetting...")
			pwd, err = set.RePwd()
			if err != nil {
				panic(err)
			}
			log.Printf("The new dashboard password is %s\n", pwd.Pwd)
		}
		if arg == "clear" {
			err := set.Clear()
			if err != nil {
				panic(err)
			}
			log.Printf("Cleared")
		}
		if arg == "clear_ip_restriction" {
			err := set.ClearIP()
			if err != nil {
				panic(err)
			}
			log.Printf("Cleared IP restriction")
		}
		if arg == "port" {
			setPort = true
		}
	}
	if len(os.Args) == 1 || os.Args[1] == "port" {
		log.Printf("The server will start on port %s\n", port)
		router.Start(port)
	}
}
