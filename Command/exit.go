package Command

import (
	"fmt"
	"os"
)

func Exit() {
	fmt.Println("Exiting...")
	os.Exit(0)
}
