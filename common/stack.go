package common

import (
	"errors"
)

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
