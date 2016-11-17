package common

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

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
	require.Equal(0, stack.top, "bad top value")
	for i := 0; i < stackSize; i++ {
		assert.Equal(VMWord(0), stack.items[i], "item not initialized")
	}
}
