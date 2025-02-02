package main

import (
	"backend/cmd/set"
	"backend/internal/router"
	"fmt"
	"os"
)

func main() {
	router.Start()
	var err error
	var pwd set.SetResponse
	//fmt.Printf("%x", sha256.Sum256([]byte("password")))
	if os.Args[0] == "restore" {
		pwd, err = set.Restore()
	} else if os.Args[0] == "repwd" {
		pwd, err = set.RePwd()
	} else {
		router.Start()
	}
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("The new dashboard password is %s\n", pwd)
	}
}
