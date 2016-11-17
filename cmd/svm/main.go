package main

import (
	"fmt"
	"github.com/scompo/stack-vm/common"
	"os"
)

var programName = "stack-vm"
var programVersion = "no-version"

// entry point
func main() {
	fmt.Println(getProgramHeader())
	vm := common.DefaultVM()
	programFile, err := os.Open(os.Args[1])
	defer programFile.Close()
	if err != nil {
		handleErrorAndReturn("error loading file", err)
	}
	programFileInfo, err := programFile.Stat()
	programReader := common.NewSizedReader(programFile, programFileInfo.Size())
	program, err := common.ReadProgram(programReader)
	if err != nil {
		handleErrorAndReturn("error reading program", err)
	}
	err = common.LoadProgram(&vm, program)
	if err != nil {
		handleErrorAndReturn("error loading program", err)
	}
	err = common.Run(&vm)
	if err != nil {
		handleErrorAndReturn("error running program", err)
	}
	os.Exit(0)
}

func handleErrorAndReturn(message string, err error) {
	fmt.Printf("%v: %v\n", message, err)
	os.Exit(1)
}

// getProgramHeader returns the program header "name (version)"
func getProgramHeader() string {
	return fmt.Sprintf("%s (%s)", programName, programVersion)
}
