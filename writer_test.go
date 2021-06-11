package go_shm

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"github.com/hidez8891/shm"
	"log"
	"os"
	"time"
)

func ExampleNewWriter() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	w := NewWriter()

	if err := w.Open(shmName, shmSize); err != nil {
		log.Println(err)
		return
	}
	defer w.Close()

	log.Println("输入数据, 回车发送:")

	rd := bufio.NewReader(os.Stdin)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			log.Println(err)
			return
		}

		_ = w.Write(line)
	}

	// output:

}

func ExampleShowInfo() {
	ctrlMem, err := shm.Create(shmName+"_ctrl", 1024)
	if err != nil {
		return
	}

	var printUint64 = func(data uint64) {
		log.Printf("%v %x", data, data)
	}

	buffer := make([]byte, 48)
	for {
		time.Sleep(time.Second)

		_, _ = ctrlMem.ReadAt(buffer, 0)

		log.Println(hex.EncodeToString(buffer[0:32]))

		printUint64(binary.BigEndian.Uint64(buffer[0:8]))
		printUint64(binary.BigEndian.Uint64(buffer[8:16]))
		printUint64(binary.BigEndian.Uint64(buffer[16:24]))
		printUint64(binary.BigEndian.Uint64(buffer[24:32]))
	}

	// output:

}
