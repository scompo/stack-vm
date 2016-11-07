package main

import (
	"testing"
)

func TestGetProgramHeader(t *testing.T) {
	expected := "stack-vm (no-version)"
	result := getProgramHeader()
	if result != expected {
		t.Errorf("bad program version, expected: \"%v\" but found \"%v\" \n", expected, result)
	}
}
