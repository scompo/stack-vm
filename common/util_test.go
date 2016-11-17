package common

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSizedReader(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader(make([]byte, 1))
	sr := NewSizedReader(r, 1)
	assert.Equal(r, sr.R)
	assert.Equal(int64(1), sr.Size)
}
