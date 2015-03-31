package throttle

import "io"

type Reader struct {
	r io.Reader
	b *Bucket
}

func NewReader(r io.Reader, rate, cap int64) *Reader {
	return &Reader{
		r: r,
		b: NewBucket(rate, cap),
	}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	r.b.TakeExactly(int64(n))
	return
}

type Writer struct {
	w io.Writer
	b *Bucket
}

func NewWriter(w io.Writer, rate, cap int64) *Writer {
	return &Writer{
		w: w,
		b: NewBucket(rate, cap),
	}
}

func (r *Writer) Write(p []byte) (n int, err error) {
	n, err = r.w.Write(p)
	r.b.TakeExactly(int64(n))
	return
}
