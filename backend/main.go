package main

import (
	"backend/internal/data"
)

func main() {
	data.Data()
	defer data.Disconnect()
}
