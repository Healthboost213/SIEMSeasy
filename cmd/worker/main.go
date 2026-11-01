package main

import (
	"SIEMSEasy/internal"
	"fmt"
)

func main() {

	regexData := `(?m)^ERROR|SUCCESS|WARNING`
	lineData := internal.FileTail("test_log.txt", regexData)
	fmt.Println(lineData)

}
