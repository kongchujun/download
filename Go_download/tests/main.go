package main

import (
	"fmt"
	"strings"
)

func main() {

	filePatha := "/user/local/abc.txt.gz"

	a := strings.TrimSuffix(filePatha, ".gz")
	fmt.Println(a)
}
