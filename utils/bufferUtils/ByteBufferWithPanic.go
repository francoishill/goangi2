package bufferUtils

import (
	"bytes"
)

func CreateByteBufferWithPanic(buf *bytes.Buffer) *ByteBufferWithPanic {
	return &ByteBufferWithPanic{
		buf: buf,
	}
}

type ByteBufferWithPanic struct {
	buf *bytes.Buffer
}

func (this *ByteBufferWithPanic) WriteString_PanicOnError(str string) {
	_, err := this.buf.WriteString(str)
	if err != nil {
		panic(err)
	}
}

func (this *ByteBufferWithPanic) Write_PanicOnError(data []byte) {
	_, err := this.buf.Write(data)
	if err != nil {
		panic(err)
	}
}

func (this *ByteBufferWithPanic) GetBytes() []byte {
	return this.buf.Bytes()
}

func (this *ByteBufferWithPanic) GetString() string {
	return string(this.GetBytes())
}
