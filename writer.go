package go_shm

import (
	"encoding/binary"
	"fmt"
	"github.com/hidez8891/shm"
	"io"
	"time"
)

type Writer struct {
	size    int
	m       *shm.Memory
	ctrlMem *shm.Memory
}

func NewWriter() *Writer {
	return &Writer{}
}

func (c *Writer) Open(name string, size int32) error {
	var err error
	c.m, err = shm.Create(name, size)
	if err != nil {
		return err
	}

	c.ctrlMem, err = shm.Create(name+"_ctrl", 1024)
	if err != nil {
		_ = c.m.Close()
		return err
	}

	_, _ = c.m.Seek(int64(c.getWriteOffset()), io.SeekStart)

	c.size = int(size)

	return err
}

func (c *Writer) Write(p []byte) error {
	if p == nil || len(p) > c.size {
		return fmt.Errorf("size error")
	}

	for {
		emptyCount := uint64(c.size) - (c.getWriteLength() - c.getReadLength())
		if emptyCount > uint64(len(p)) {
			break
		} else {
			time.Sleep(time.Millisecond)
		}
	}

	n, err := c.m.Write(p)
	if err == io.EOF {
		// 文件尾
		_, _ = c.m.Seek(0, io.SeekStart)
		_, _ = c.m.Write(p)
		return nil
	}

	if n != len(p) {
		// 接近文件尾
		_, _ = c.m.Seek(0, io.SeekStart)
		_, _ = c.m.Write(p[n:])
	}

	c.setWriteLength(uint64(len(p)) + c.getWriteLength())
	c.updateWriteOffset()

	return nil
}

func (c *Writer) Close() error {
	if c == nil || c.m == nil || c.ctrlMem == nil {
		return fmt.Errorf("not open")
	}

	_ = c.ctrlMem.Close()

	return c.m.Close()
}

func (c *Writer) getWriteLength() uint64 {
	var oldData [8]byte
	_, _ = c.ctrlMem.ReadAt(oldData[:], 16)

	oldLength := binary.BigEndian.Uint64(oldData[:])
	return oldLength
}

func (c *Writer) setWriteLength(dataLength uint64) {
	var newData [8]byte
	binary.BigEndian.PutUint64(newData[:], dataLength)
	_, _ = c.ctrlMem.WriteAt(newData[:], 16)
}

// updateReadOffset 更新写入位置
func (c *Writer) updateWriteOffset() {
	offset, _ := c.m.Seek(0, io.SeekCurrent)

	var readOffset [8]byte
	binary.BigEndian.PutUint64(readOffset[:], uint64(offset))
	_, _ = c.ctrlMem.WriteAt(readOffset[:], 24)
}

// getReadOffset shm中记录的写入位置
func (c *Writer) getWriteOffset() uint64 {
	var oldData [8]byte
	_, _ = c.ctrlMem.ReadAt(oldData[:], 24)

	oldLength := binary.BigEndian.Uint64(oldData[:])
	return oldLength
}

func (c *Writer) getReadLength() uint64 {
	var oldData [8]byte
	_, _ = c.ctrlMem.ReadAt(oldData[:], 0)

	oldLength := binary.BigEndian.Uint64(oldData[:])
	return oldLength
}
