package go_shm

import (
	"encoding/binary"
	"fmt"
	"github.com/hidez8891/shm"
	"io"
	"time"
)

type Reader struct {
	size    int
	m       *shm.Memory
	ctrlMem *shm.Memory
}

func NewReader() *Reader {
	return &Reader{}
}

func (c *Reader) Open(name string, size int32) error {
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

	_, _ = c.m.Seek(int64(c.getReadOffset()), io.SeekStart)

	c.size = int(size)

	return err
}

func (c *Reader) Read(p []byte) (int, error) {
	if c.m == nil {
		return 0, fmt.Errorf("not open")
	}
	if p == nil || len(p) > c.size {
		return 0, fmt.Errorf("size error")
	}

	for {
		dataCount := c.getWriteLength() - c.getReadLength()
		if dataCount > 0 {
			if max := dataCount; uint64(len(p)) > max {
				p = p[:max]
			}

			break
		} else {
			time.Sleep(time.Millisecond)
		}
	}

	n, err := c.m.Read(p)
	if err == io.EOF {
		// 文件尾
		_, _ = c.m.Seek(0, io.SeekStart)
		n2, err := c.m.Read(p)

		c.setReadLength(c.getReadLength() + uint64(n2+n))
		c.updateReadOffset()
		return n + n2, err
	}

	if n != len(p) {
		// 接近文件尾
		_, _ = c.m.Seek(0, io.SeekStart)
		n2, err := c.m.Read(p[n:])

		c.setReadLength(c.getReadLength() + uint64(n2+n))
		c.updateReadOffset()
		return n + n2, err
	}

	c.setReadLength(c.getReadLength() + uint64(n))
	c.updateReadOffset()
	return n, err
}

func (c *Reader) Close() error {
	if c == nil || c.m == nil || c.ctrlMem == nil {
		return fmt.Errorf("not open")
	}

	_ = c.ctrlMem.Close()

	return c.m.Close()
}

func (c *Reader) getReadLength() uint64 {
	var oldData [8]byte
	_, _ = c.ctrlMem.ReadAt(oldData[:], 0)

	oldLength := binary.BigEndian.Uint64(oldData[:])
	return oldLength
}

func (c *Reader) setReadLength(dataLength uint64) {
	var newData [8]byte
	binary.BigEndian.PutUint64(newData[:], dataLength)
	_, _ = c.ctrlMem.WriteAt(newData[:], 0)
}

// updateReadOffset 更新读取位置
func (c *Reader) updateReadOffset() {
	offset, _ := c.m.Seek(0, io.SeekCurrent)

	var readOffset [8]byte
	binary.BigEndian.PutUint64(readOffset[:], uint64(offset))
	_, _ = c.ctrlMem.WriteAt(readOffset[:], 8)
}

// getReadOffset shm中记录的读取位置
func (c *Reader) getReadOffset() uint64 {
	var oldData [8]byte
	_, _ = c.ctrlMem.ReadAt(oldData[:], 8)

	oldLength := binary.BigEndian.Uint64(oldData[:])
	return oldLength
}

// getWriteLength 数据写的总长度
func (c *Reader) getWriteLength() uint64 {
	var oldData [8]byte
	_, _ = c.ctrlMem.ReadAt(oldData[:], 16)

	oldLength := binary.BigEndian.Uint64(oldData[:])
	return oldLength
}
