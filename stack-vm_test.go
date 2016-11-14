package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJump(t *testing.T) {

	tests := []struct {
		input int
		err   string
	}{
		{
			input: 100,
			err:   "jump out of memory bounds",
		},
		{
			input: -100,
			err:   "jump out of memory bounds",
		},
		{
			input: 1,
		},
	}

	assert := assert.New(t)

	for _, test := range tests {

		vm := NewVM(0)

		LoadProgram(&vm, []VMWord{NOP, HALT})
		err := Jump(&vm, test.input)

		if test.err != "" {
			assert.EqualError(err, test.err)
		} else {
			assert.NoError(err)
			assert.Equal(test.input, vm.pc, "pc not updated")
		}
	}
}

func TestDefaultVM(t *testing.T) {

	vm := DefaultVm()

	validateVM(t, defaultStackSize, vm)
}

func TestLoadProgram(t *testing.T) {

	tests := []struct {
		program []VMWord
		err     string
	}{
		{
			program: []VMWord{},
			err:     "empty program",
		},
		{
			program: []VMWord{HALT},
		},
		{
			program: []VMWord{NOP, HALT},
		},
	}

	assert := assert.New(t)
	require := require.New(t)

	for _, test := range tests {

		vm := NewVM(0)
		err := LoadProgram(&vm, test.program)

		if test.err != "" {
			assert.EqualError(err, test.err)
		} else {
			assert.NoError(err)
			require.Equal(len(test.program), len(vm.program), "bad program lenght")
			for i := 0; i < len(test.program); i++ {
				assert.Equal(test.program[i], vm.program[i], "error in instruction[", i, "]")
			}
		}
	}
}

func TestFetch(t *testing.T) {
	assert := assert.New(t)
	vm := NewVM(0)
	LoadProgram(&vm, []VMWord{HALT})
	op, err := Fetch(&vm)
	assert.NoError(err)
	assert.Equal(HALT, op, "bad operand retrived")
	assert.Equal(1, vm.pc, "bad pc value")
	op, err = Fetch(&vm)
	assert.EqualError(err, "program out of bounds")
}

func TestRun(t *testing.T) {
	assert := assert.New(t)

	vm := NewVM(0)
	LoadProgram(&vm, []VMWord{NOP,HALT})
	err := Run(&vm)

	assert.NoError(err)
	assert.Equal(2, vm.pc)
	assert.Equal(0, vm.stack.top)
}

func TestLoadParamsNone(t *testing.T) {
	assert := assert.New(t)
	vm := DefaultVm()
	LoadProgram(&vm, []VMWord{})
	res, err := LoadParams(&vm,0)
	assert.Equal(0, len(res))
	assert.Equal(0, vm.pc)
	assert.Equal(0, vm.stack.top)
	assert.NoError(err)
}

func TestLoadParamsOne(t *testing.T) {
	assert := assert.New(t)
	vm := DefaultVm()
	LoadProgram(&vm, []VMWord{VMWord(0)})
	res, err := LoadParams(&vm,1)
	assert.Equal(1, len(res))
	assert.Equal(VMWord(0), res[0])
	assert.Equal(1, vm.pc)
	assert.Equal(0, vm.stack.top)
	assert.NoError(err)
}

func TestLoadParamsError(t *testing.T) {
	assert := assert.New(t)
	vm := DefaultVm()
	_, err := LoadParams(&vm,1)
	assert.EqualError(err, "program out of bounds")
}


func TestExecuteUnknownOperand(t *testing.T) {

	assert := assert.New(t)

	vm := DefaultVm()
	err := Execute(&vm, VMWord(-1), make([]VMWord, 0))

	assert.EqualError(err, "unknown operand")
}

func TestExecuteNop(t *testing.T) {

	assert := assert.New(t)

	vm := DefaultVm()
	err := Execute(&vm, NOP, make([]VMWord, 0))

	assert.NoError(err)
	assert.Equal(0, vm.pc)
	assert.Equal(0, vm.stack.top)
}

func TestExecutePrint(t *testing.T) {

	assert := assert.New(t)

	vm := DefaultVm()

	myOut := new(bytes.Buffer)

	vm.out = myOut

	err := Execute(&vm, PRINT, []VMWord{VMWord('a')})

	assert.Equal("a", myOut.String())
	assert.NoError(err)
	assert.Equal(0, vm.pc)
	assert.Equal(0, vm.stack.top)

	vm.out = DEFAULT_WRITER
}

func TestNewVM(t *testing.T) {

	vm := NewVM(1)

	validateVM(t, 1, vm)
}

func validateVM(t *testing.T, stackSize int, vm VM) {

	assert.Equal(t, 0, vm.pc, "pc not initialized")
	assert.Equal(t, DEFAULT_WRITER, vm.out, "default writer not set for out")

	validateStack(t, stackSize, vm.stack)
}

func TestGetParamNumber(t *testing.T) {

	assert := assert.New(t)

	tests := []struct {
		input    VMWord
		expected int
		err      string
	}{
		{
			input: VMWord(-1),
			err:   "unknown operand",
		},
		{
			input:    HALT,
			expected: NoParams,
		},
		{
			input:    NOP,
			expected: NoParams,
		},
		{
			input:    PRINT,
			expected: OneParam,
		},
	}
	for _, test := range tests {

		res, err := GetParamNumber(test.input)

		if test.err != "" {
			assert.EqualError(err, test.err)
		} else {
			assert.NoError(err)
			assert.Equal(test.expected, res)
		}
	}
}

func TestPush(t *testing.T) {

	assert := assert.New(t)

	s := NewStack(1)
	err := Push(&s, 1)

	assert.Equal(1, s.top, "top not updated")
	assert.Equal(VMWord(1), s.items[0], "item not pushed")
	assert.Nil(err)

	err = Push(&s, 0)

	assert.EqualError(err, "overflow")
}

func TestPop(t *testing.T) {

	assert := assert.New(t)

	s := NewStack(1)
	Push(&s, 1)
	res, err := Pop(&s)

	assert.Equal(0, s.top, "top not updated")
	assert.Nil(err)
	assert.Equal(VMWord(1), res, "bad element returned")

	_, err = Pop(&s)

	assert.EqualError(err, "underflow")
}

func TestNewStack(t *testing.T) {

	stackSize := 1
	res := NewStack(stackSize)

	validateStack(t, stackSize, res)
}

func validateStack(t *testing.T, stackSize int, stack Stack) {

	assert := assert.New(t)
	require := require.New(t)

	require.NotNil(stack, "stack not initialized")
	assert.Equal(stackSize, stack.size, "bad res.maxSize")
	assert.Equal(stackSize, len(stack.items), "bad items lenght")
	assert.Equal(0, stack.top, "bad top value")
	for i := 0; i < stackSize; i++ {
		assert.Equal(VMWord(0), stack.items[i], "item not initialized")
	}
}

func TestGetProgramHeader(t *testing.T) {
	expected := "stack-vm (no-version)"
	result := getProgramHeader()
	assert.Equal(t, expected, result)
}
