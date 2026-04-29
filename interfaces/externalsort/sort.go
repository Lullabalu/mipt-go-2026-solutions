//go:build !solution

package externalsort

import (
	"container/heap"
	"errors"
	"io"
	"os"
	"sort"
)

type item struct {
	str string
	ind int
}

type Reader struct {
	r io.Reader
}

type Writer struct {
	w io.Writer
}

type itheap []item

func (h itheap) Len() int {
	return len(h)
}

func (h itheap) Less(i, j int) bool {
	return h[i].str < h[j].str
}

func (h itheap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *itheap) Push(x any) {
	*h = append(*h, x.(item))
}

func (h *itheap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func (reader *Reader) ReadLine() (string, error) {
	bt := make([]byte, 1)
	str := make([]byte, 0)

	for {
		n, err := reader.r.Read(bt)

		if n > 0 {
			if bt[0] == '\n' {
				return string(str), nil
			}
			str = append(str, bt[0])
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(str) == 0 {
					return "", io.EOF
				}
				return string(str), nil
			}
			return "", err
		}
	}
}

func (writer *Writer) Write(str string) error {
	_, err := writer.w.Write([]byte(str))

	if err != nil {
		return err
	}
	_, err = writer.w.Write([]byte("\n"))
	return err

}

func NewReader(r io.Reader) LineReader {
	rd := &Reader{r: r}
	return rd
}

func NewWriter(w io.Writer) LineWriter {
	wr := &Writer{w: w}
	return wr
}

func Merge(w LineWriter, readers ...LineReader) error {
	rds := readers
	hp := &itheap{}

	heap.Init(hp)

	for i, reader := range rds {
		str, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			continue
		}

		if err != nil {
			return err
		}
		heap.Push(hp, item{str, i})
	}

	for hp.Len() > 0 {
		top := heap.Pop(hp).(item)
		err := w.Write(top.str)

		if err != nil {
			return err
		}

		str, err := readers[top.ind].ReadLine()
		if errors.Is(err, io.EOF) {
			continue
		}

		if err != nil {
			return err
		}
		heap.Push(hp, item{str, top.ind})

	}

	return nil
}

func sortFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	reader := NewReader(file)

	var lines []string

	for {
		line, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			file.Close()
			return err
		}
		lines = append(lines, line)
	}

	file.Close()

	sort.Strings(lines)

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := NewWriter(out)

	for _, line := range lines {
		if err := writer.Write(line); err != nil {
			return err
		}
	}

	return nil
}

func Sort(w io.Writer, in ...string) error {
	for _, path := range in {
		if err := sortFile(path); err != nil {
			return err
		}
	}

	files := make([]*os.File, 0, len(in))
	readers := make([]LineReader, 0, len(in))

	for _, path := range in {
		file, err := os.Open(path)
		if err != nil {
			for _, ffile := range files {
				_ = ffile.Close()
			}
			return err
		}
		files = append(files, file)
		readers = append(readers, NewReader(file))
	}

	for _, file := range files {
		defer file.Close()
	}

	lw := NewWriter(w)
	return Merge(lw, readers...)
}
