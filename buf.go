package bufferPool

import (
	"io"

	lowlevelfunctions "github.com/NikoMalik/low-level-functions"
)

const bufferSize = 8192

var poolbuf = GetPool(bufferSize)

type Buffer struct {
	buf   []byte
	start int
	end   int
	size  int
}

func New() *Buffer {
	buf := poolbuf.Get().([]byte)

	if cap(buf) >= bufferSize {
		buf = buf[:bufferSize]
	} else {
		buf = make([]byte, bufferSize)
	}

	return &Buffer{
		buf: buf,
	}
}

func (b *Buffer) Release() {
	if b == nil {
		return
	}

	b.buf = nil
	b.Reset()
	if cap(b.buf) == bufferSize {
		poolbuf.Put(&b.buf)
	}
}

func (b *Buffer) Len() int {
	return b.end - b.start
}

func (b *Buffer) Cap() int {
	return cap(b.buf)
}

func (b *Buffer) Reset() {
	b.start = 0
	b.end = 0
	b.size = 0
}

func (b *Buffer) Byte(idx int) byte {
	return b.buf[b.start+idx]
}

func (b *Buffer) Bytes() []byte {
	return b.buf[b.start:b.end]
}

func (b *Buffer) String() string {
	return lowlevelfunctions.String(b.buf[b.start:b.end])
}

func (b *Buffer) WriteByte(c byte) error {
	b.buf[b.end] = c
	b.end++
	return nil
}
func (b *Buffer) WriteString(s string) (int, error) {
	nBytes := copy(b.buf[b.end:], s)
	b.end += nBytes
	return nBytes, nil
}

func (b *Buffer) grow(n int) {
	buf := lowlevelfunctions.MakeNoZero(2*cap(b.buf) + n)[:len(b.buf)]
	copy(buf, b.buf)
	b.buf = buf
}

func (b *Buffer) Grow(n int) {
	// Check if n is negative
	if n < 0 {
		// Panic with the message "fast.Buffer.Grow: negative count"
		panic("fast.Buffer.Grow: negative count")
	}
	if cap(b.buf)-len(b.buf) < n {
		b.grow(n)
	}
}

func (b *Buffer) Write(p []byte) (int, error) {

	b.buf = append(b.buf, p...)
	return len(p), nil
}

func (b *Buffer) ReadByte() (byte, error) {
	if b.start == b.end {
		return 0, io.EOF
	}

	nb := b.buf[b.start]
	b.start++
	return nb, nil
}

func (b *Buffer) Read(data []byte) (int, error) {

	if b.Len() == 0 {
		return 0, io.EOF
	}
	nBytes := copy(data, b.buf[b.start:b.end])
	if nBytes == b.Len() {
		b.Reset()
	} else {
		b.start += nBytes
	}
	return nBytes, nil
}
func (b *Buffer) ReadFrom(reader io.Reader) (int64, error) {
	if b.IsFull() {
		b.grow(1)
	}
	n, err := reader.Read(b.buf[b.end:])
	b.end += n
	return int64(n), err
}

func (b *Buffer) IsEmpty() bool {
	return b.Len() == 0
}

func (b *Buffer) ReadFullFrom(reader io.Reader, size int) (int, error) {
	if b.Cap() < size {
		b.grow(size)
	}
	n, err := reader.Read(b.buf[b.end : b.end+size])
	b.end += n
	return n, err
}

func (b *Buffer) WriteTo(writer io.Writer) (int, error) {
	return writer.Write(b.buf[b.start:b.end])
}

func (b *Buffer) IsFull() bool {
	return b.Cap() == b.Len()
}
