package main

import (
	"log"
	"os"
	"urlAPI/cmd/set"
	"urlAPI/internal/router"
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
		switch arg {
		case "restore":
			log.Println("Restoring...")
			pwd, err = set.Restore()
			log.Printf("The new dashboard password is %s\n", pwd.Pwd)
		case "repwd":
			log.Println("Password Resetting...")
			pwd, err = set.RePwd()
			log.Printf("Dashboard password is %s, please change it ASAP\n", pwd.Pwd)
		case "clear":
			err = set.Clear()
			log.Printf("Cleared")
		case "clear_ip_restriction":
			err = set.ClearIP()
			log.Printf("Cleared IP restriction")
		case "update":
			err = set.Update()
			log.Printf("Updated")
		case "port":
			setPort = true
		}
		if err != nil {
			panic(err)
		}
	}
	if len(os.Args) == 1 || os.Args[1] == "port" {
		log.Printf("The server will start on port %s\n", port)
		router.Start(port)
	}
}
