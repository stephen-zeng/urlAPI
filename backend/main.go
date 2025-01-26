package main

import (
	"backend/internal/data"
	"backend/internal/file"
	"fmt"
)

func main() {
	id := "c7a2a2b6-dbf6-11ef-91d2-42594833521e"
	err := file.Del(id)
	fmt.Println(err)
	defer data.Disconnect()
}
