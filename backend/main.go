package main

import (
	"backend/internal/data"
	"fmt"
)

func main() {
	task, _ := data.Fetch(data.DataConfig(
		data.WithUUID("7343b8ec-90c5-4c77-85a4-f4ac24cb52b9")))
	for _, item := range task {
		fmt.Println(item)
	}
	fmt.Println("\n")
	task, _ = data.Fetch(data.DataConfig())
	for _, item := range task {
		fmt.Println(item)
	}
	defer data.Disconnect()
}
