package common

import 	"io"

// NewSizedReader returns a new SizedReader from a reader and it's size.
func NewSizedReader(reader io.Reader, size int64) SizedReader {
	sr := SizedReader{reader, size}
	return sr
}

// SizedReader is a Reader but has a size.
type SizedReader struct {
	R    io.Reader
	Size int64
}
