package bootstrap

import (
	"log"
	"urlAPI/internal/database"
	"urlAPI/internal/op"
)

func Init() error {
	if err := database.Init(); err != nil {
		return err
	}
	if err := op.Init(); err != nil {
		Release()
		return err
	}
	return nil
}

func Release() {
	if err := database.Disconnect(); err != nil {
		log.Println(err)
	}
}
