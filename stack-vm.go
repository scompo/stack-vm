package main

import (
	"fmt"
)

var program_name = "stack-vm"
var program_version = "no-version"

func main() {
	fmt.Println(getProgramHeader())
}

func getProgramHeader() string {
	return fmt.Sprintf("%s (%s)", program_name, program_version)
}
