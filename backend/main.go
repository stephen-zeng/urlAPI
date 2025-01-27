package main

import (
	"backend/internal/data"
	"fmt"
)

func main() {
	pwd, _ := data.InitSetting()
	fmt.Println(pwd)
	defer data.Disconnect()
}
