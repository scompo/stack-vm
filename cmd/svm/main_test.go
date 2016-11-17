package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetProgramHeader(t *testing.T) {
	expected := "stack-vm (no-version)"
	result := getProgramHeader()
	assert.Equal(t, expected, result)
}
