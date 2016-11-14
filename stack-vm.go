package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var programName = "stack-vm"
var programVersion = "no-version"

// VMWord is a byte of the VM.
type VMWord int32

// Default stack size.
const defaultStackSize = 1024

const (
	// No parameters required for operand.
	NoParams = 0
	// One parameter required for operand.
	OneParam = 1
)

var (
	// Stops the VM, opcode 0.
	HALT = VMWord(0)

	// No operation for the VM, opcode 1.
	NOP = VMWord(1)

	// Writes the top of the stack as a char.
	PRINT = VMWord(2)
)

var DEFAULT_WRITER io.Writer = os.Stdout

var unknownOperandError error = errors.New("unknown operand")

// Jump moves the program counter to the specified location.
// Returns an error if the address is not in the progam memory bounds.
func Jump(vm *VM, addr int) (err error) {
	if addr < 0 || addr >= len(vm.program) {
		err = errors.New("jump out of memory bounds")
		return
	}
	vm.pc = addr
	return
}

// Returns the number of parameters for the specified operand.
// Returns an error if the operand it's unknown.
func GetParamNumber(op VMWord) (num int, err error) {
	switch op {
	case HALT:
		num = NoParams
	case NOP:
		num = NoParams
	case PRINT:
		num = OneParam
	default:
		err = unknownOperandError
	}
	return
}

// NewStack returns a new stack with the specified stack size.
func NewStack(stackSize int) Stack {
	return Stack{
		size:  stackSize,
		top:   0,
		items: make([]VMWord, stackSize),
	}
}

// Stack is the stack of the VM.
type Stack struct {
	size  int
	top   int
	items []VMWord
}

// Push adds an element and increments the top index, after checking
// for overflow.
func Push(s *Stack, elem VMWord) (err error) {
	if s.top == s.size {
		return errors.New("overflow")
	}
	s.items[s.top] = elem
	s.top = s.top + 1
	return
}

// Pop decrements the top index after checking for underflow
// and returns the item that was previously the top one.
func Pop(s *Stack) (elem VMWord, err error) {
	if s.top == 0 {
		err = errors.New("underflow")
		return
	}
	s.top = s.top - 1
	elem = s.items[s.top]
	return
}

// DefaultVm creates a new VM with stack size of defaultStackSize.
func DefaultVm() VM {
	return NewVM(defaultStackSize)
}

// NewVM creates a new VM with stackSize stack size of stackSize.
func NewVM(stackSize int) VM {
	return VM{
		out:     DEFAULT_WRITER,
		stack:   NewStack(stackSize),
		program: []VMWord{},
	}
}

// Runs the program loaded in the VM.
func Run(vm *VM) error {
	var err error
	var op = VMWord(-1)
	for op, err = Fetch(vm); op != HALT && err == nil; op, err = Fetch(vm) {
		nParams, err := GetParamNumber(op)
		if err != nil {
			return err
		}
		params, err := LoadParams(vm, nParams)
		if err != nil {
			return err
		}
		err = Execute(vm, op, params)
		if err != nil {
			return err
		}
	}
	return err
}

// Execute executes the operand with the specified params.
// Returns an error if something bad happens.
func Execute(vm *VM, op VMWord, params []VMWord) (err error) {
	switch op {
	case NOP:
	case PRINT:
		fmt.Fprintf(vm.out, "%c", params[0])
	default:
		err = unknownOperandError
		return
	}
	return
}

// Fetch returns n params from the VM.
func LoadParams(vm *VM, n int) ([]VMWord, error) {
	params := make([]VMWord, n)
	for i := 0; i < n; i++ {
		fetched, err := Fetch(vm)
		if err != nil {
			return params, err
		}
		params[i] = fetched
	}
	return params, nil
}

// Fetch returns the next operand.
// Returns an error if the program counter is outside program bounds.
func Fetch(vm *VM) (op VMWord, err error) {
	if vm.pc < 0 || vm.pc >= len(vm.program) {
		err = errors.New("program out of bounds")
		return
	}
	op = vm.program[vm.pc]
	vm.pc = vm.pc + 1
	return
}

// VM is a stack Virtual Machine.
type VM struct {
	out     io.Writer
	stack   Stack
	pc      int
	program []VMWord
}

// LoadProgram loads a program into the VM.
// Returns an error for an empty program.
func LoadProgram(vm *VM, program []VMWord) (err error) {

	if len(program) == 0 {
		err = errors.New("empty program")
		return
	}

	vm.program = make([]VMWord, len(program))

	for i := 0; i < len(program); i++ {
		vm.program[i] = program[i]
	}

	return
}

// entry point
func main() {
	fmt.Println(getProgramHeader())
}

func getProgramHeader() string {
	return fmt.Sprintf("%s (%s)", programName, programVersion)
}
