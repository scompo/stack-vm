package common

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

// VMWord is a byte of the VM.
type VMWord int32

// VMWordSize is the size of a VMWord in bytes (4 bytes).
const VMWordSize = 4

const (
	// NoParams it's used when an operand uses no params.
	NoParams = 0
	// OneParam it's used when an operand requires a parameter.
	OneParam = 1
)

var (
	// HALT stops the VM, opcode 0.
	HALT = VMWord(0)

	// NOP does no operation in the VM, opcode 1.
	NOP = VMWord(1)

	// PRINT writes the top of the stack as a char to DefaultWriter, opcode 2.
	PRINT = VMWord(2)

	// PUSH pushes a value to the top of the stack, opcode 3.
	PUSH = VMWord(3)

	// POP pops a value from the top of the stack, opcode 4.
	POP = VMWord(4)

	// ADD adds 2 the 2 elements from the top of the stack and pushes the result,
	// opcode 5.
	ADD = VMWord(5)

	// JMP jumps to the location pointed by the top of the stack, opcode 6.
	JMP = VMWord(6)

	// JZ jumps to a location if the top of the stack is zero, opcode 7.
	JZ = VMWord(7)

	// JNZ jumps to a location if the top of the stack is not zero, opcode 8.
	JNZ = VMWord(8)

	// CALL jumps to a function address, saving the current pc, opcode 9.
	CALL = VMWord(9)

	// RET returns from a function call retstoring the last pc, opcode 10.
	RET = VMWord(10)
)

// Default stack size.
const defaultStackSize = 1024

// DefaultWriter is the default writer for the VM.
var DefaultWriter io.Writer = os.Stdout

var errUnknownOperand = errors.New("unknown operand")

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

// GetParamsNumber returns the number of parameters for the specified operand.
// Returns an error if the operand it's unknown.
func GetParamsNumber(op VMWord) (num int, err error) {
	switch op {
	case HALT:
		num = NoParams
	case NOP:
		num = NoParams
	case PRINT:
		num = NoParams
	case PUSH:
		num = OneParam
	case POP:
		num = NoParams
	case ADD:
		num = NoParams
	case JMP:
		num = NoParams
	case JZ:
		num = OneParam
	case JNZ:
		num = OneParam
	case CALL:
		num = OneParam
	case RET:
		num = NoParams
	default:
		err = errUnknownOperand
	}
	return
}

// DefaultVM creates a new VM with stack and return stack size of
// defaultStackSize.
func DefaultVM() VM {
	return NewVM(defaultStackSize, defaultStackSize)
}

// NewVM creates a new VM with stack size of stackSize, and return stack
// size of returnStackSize.
func NewVM(stackSize int, returnStackSize int) VM {
	return VM{
		out:         DefaultWriter,
		stack:       NewStack(stackSize),
		returnStack: NewStack(returnStackSize),
		program:     []VMWord{},
	}
}

// Run runs the program loaded in the VM.
func Run(vm *VM) error {
	var err error
	var op = VMWord(-1)
	for op, err = Fetch(vm); op != HALT && err == nil; op, err = Fetch(vm) {
		nParams, err := GetParamsNumber(op)
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
		data, err := Pop(&vm.stack)
		if err != nil {
			return err
		}
		return Print(vm.out, data)
	case PUSH:
		return Push(&vm.stack, params[0])
	case POP:
		_, err = Pop(&vm.stack)
	case ADD:
		first, err := Pop(&vm.stack)
		if err != nil {
			return err
		}
		second, err := Pop(&vm.stack)
		if err != nil {
			return err
		}
		return Push(&vm.stack, first+second)
	case JMP:
		addr, err := Pop(&vm.stack)
		if err != nil {
			return err
		}
		return Jump(vm, int(addr))
	case JZ:
		value, err := Pop(&vm.stack)
		if err != nil {
			return err
		}
		if value == VMWord(0) {
			return Jump(vm, int(params[0]))
		}
	case JNZ:
		value, err := Pop(&vm.stack)
		if err != nil {
			return err
		}
		if value != VMWord(0) {
			return Jump(vm, int(params[0]))
		}
	case CALL:
		err = Push(&vm.returnStack, VMWord(vm.pc))
		if err != nil {
			return err
		}
		return Jump(vm, int(params[0]))
	case RET:
		addr, err := Pop(&vm.returnStack)
		if err != nil {
			return err
		}
		return Jump(vm, int(addr))
	default:
		err = errUnknownOperand
		return
	}
	return
}

// Print prints to out the param as char.
func Print(out io.Writer, param VMWord) (err error) {
	_, err = fmt.Fprintf(out, "%c", param)
	return
}

// LoadParams loads n parameters from the proram memory of the machine
// using Fetch.
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
	out         io.Writer
	stack       Stack
	returnStack Stack
	pc          int
	program     []VMWord
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

// ReadProgram returns a slice of VMword from the programReader.
// Returns an error if something bad happens.
func ReadProgram(programReader SizedReader) ([]VMWord, error) {
	prgSize := programReader.Size / VMWordSize
	remSize := programReader.Size % VMWordSize
	if prgSize <= 0 || remSize != 0 {
		return make([]VMWord, 0), errors.New("bad program lenght")
	}
	result := make([]VMWord, prgSize)
	err := binary.Read(programReader.R, binary.BigEndian, &result)
	return result, err
}