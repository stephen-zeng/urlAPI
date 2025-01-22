package txt

import (
	"backend/internal/data"
	"fmt"
)

func Main() string {
	test := data.Test()
	fmt.Println(test)
	return test
}
