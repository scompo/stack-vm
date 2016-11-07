package main

import (
	"fmt"
)

var programName = "stack-vm"
var programVersion = "no-version"

func main() {
	fmt.Println(getProgramHeader())
}

func getProgramHeader() string {
	return fmt.Sprintf("%s (%s)", programName, programVersion)
}
