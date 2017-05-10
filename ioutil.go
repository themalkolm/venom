package venom

import (
	"io"
)

type nopReadCloser struct {
	io.Reader
}

func (nopReadCloser) Close() error { return nil }

// NopCloser returns a ReadCloser with a no-op Close method wrapping
// the provided Reader r.
func NopReadCloser(r io.Reader) io.ReadCloser {
	return nopReadCloser{r}
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

// NopCloser returns a WriteCloser with a no-op Close method wrapping
// the provided Writer w.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}
