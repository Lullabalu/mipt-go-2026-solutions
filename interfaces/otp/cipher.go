//go:build !solution

package otp

import (
	"errors"
	"io"
)

type Reader struct {
	r    io.Reader
	prng io.Reader
}

type Writer struct {
	w    io.Writer
	prng io.Reader
}

func (reader *Reader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	rb := make([]byte, len(p))

	n, err = reader.r.Read(rb)

	if n == 0 || (err != nil && !errors.Is(err, io.EOF)) {
		return
	}

	prb := make([]byte, n)

	_, err = reader.prng.Read(prb)

	if err != nil && !errors.Is(err, io.EOF) {
		return
	}
	for i := range n {
		p[i] = rb[i] ^ prb[i]
	}
	return
}

func (writer *Writer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	prb := make([]byte, len(p))
	n, err = writer.prng.Read(prb)

	if n == 0 || (err != nil && !errors.Is(err, io.EOF)) {
		return
	}

	for i := range n {
		prb[i] = prb[i] ^ p[i]
	}
	prb = prb[:n]
	n, err = writer.w.Write(prb)
	return
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &Reader{r: r, prng: prng}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &Writer{w: w, prng: prng}
}
