package flowrate

import (
	"errors"
	"io"
)

var ErrLimit = errors.New("flowrate: flow rate limit exceeded")

type Limiter interface {
	Done() int64
	Status() Status
	SetTransferSize(bytes int64)
	SetLimit(new int64) (old int64)
	SetBlocking(new bool) (old bool)
}

type Reader struct {
	io.Reader
	*Monitor

	limit	int64
	block	bool
}

func NewReader(r io.Reader, limit int64) *Reader {
	return &Reader{r, New(0, 0), limit, true}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	p = p[:r.Limit(len(p), r.limit, r.block)]
	if len(p) > 0 {
		n, err = r.IO(r.Reader.Read(p))
	}
	return
}

func (r *Reader) SetLimit(new int64) (old int64) {
	old, r.limit = r.limit, new
	return
}

func (r *Reader) SetBlocking(new bool) (old bool) {
	old, r.block = r.block, new
	return
}

func (r *Reader) Close() error {
	defer r.Done()
	if c, ok := r.Reader.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

type Writer struct {
	io.Writer
	*Monitor

	limit	int64
	block	bool
}

func NewWriter(w io.Writer, limit int64) *Writer {
	return &Writer{w, New(0, 0), limit, true}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	var c int
	for len(p) > 0 && err == nil {
		s := p[:w.Limit(len(p), w.limit, w.block)]
		if len(s) > 0 {
			c, err = w.IO(w.Writer.Write(s))
		} else {
			return n, ErrLimit
		}
		p = p[c:]
		n += c
	}
	return
}

func (w *Writer) SetLimit(new int64) (old int64) {
	old, w.limit = w.limit, new
	return
}

func (w *Writer) SetBlocking(new bool) (old bool) {
	old, w.block = w.block, new
	return
}

func (w *Writer) Close() error {
	defer w.Done()
	if c, ok := w.Writer.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
