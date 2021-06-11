package go_shm

import (
	"encoding/binary"
	"fmt"
	"github.com/hidez8891/shm"
	"io"
	"log"
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
		var ok = false
		for i := 0; i < 1000; i++ {
			emptyCount := uint64(c.size) - (c.getWriteLength() - c.getReadLength())
			if emptyCount > uint64(len(p)) {
				ok = true
				break
			}
		}

		if ok {
			break
		}

		time.Sleep(time.Nanosecond)
	}

	n, err := c.m.Write(p)
	if err != nil {
		if err == io.EOF {
			// 文件尾
			_, _ = c.m.Seek(0, io.SeekStart)
			n, err = c.m.Write(p)
			if err != nil {
				return err
			}

			if n != len(p) {
				log.Printf("[warn] n: %v, len: %v", n, len(p))
			}

			c.setWriteLength(uint64(len(p)) + c.getWriteLength())
			c.updateWriteOffset()
			return nil
		}

		return err
	}

	if n != len(p) {
		// 接近文件尾
		_, _ = c.m.Seek(0, io.SeekStart)
		n2, err := c.m.Write(p[n:])
		if err != nil {
			return err
		}

		if n+n2 != len(p) {
			log.Printf("[warn] n: %v, n2: %v len: %v", n, n2, len(p))
		}
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
	// log.Printf("- %v", dataLength)

	var newData [8]byte
	binary.BigEndian.PutUint64(newData[:], dataLength)
	if _, err := c.ctrlMem.WriteAt(newData[:], 16); err != nil {
		log.Println(err)
	}
}

// updateReadOffset 更新写入位置
func (c *Writer) updateWriteOffset() {
	offset, _ := c.m.Seek(0, io.SeekCurrent)

	var readOffset [8]byte
	binary.BigEndian.PutUint64(readOffset[:], uint64(offset))
	if _, err := c.ctrlMem.WriteAt(readOffset[:], 24); err != nil {
		log.Println(err)
	}
}

// getReadOffset shm中记录的写入位置
func (c *Writer) getWriteOffset() uint64 {
	var oldData [8]byte
	n, err := c.ctrlMem.ReadAt(oldData[:], 24)
	if err != nil {
		log.Println(err)
	}
	if n != 8 {
		log.Printf("size error %v", n)
	}

	oldLength := binary.BigEndian.Uint64(oldData[:])
	return oldLength
}

func (c *Writer) getReadLength() uint64 {
	var oldData [8]byte
	n, err := c.ctrlMem.ReadAt(oldData[:], 0)
	if err != nil {
		log.Println(err)
	}
	if n != 8 {
		log.Printf("size error %v", n)
	}

	oldLength := binary.BigEndian.Uint64(oldData[:])
	return oldLength
}
