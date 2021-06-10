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

	if true {
		rd := bufio.NewReader(os.Stdin)
		for {
			line, _, err := rd.ReadLine()
			if err != nil {
				log.Println(err)
				return
			}

			_ = w.Write(line)
		}
	}

	if false {
		var sendSize int
		a := time.Now()
		_ = a

		data := make([]byte, 256*1024)
		for i := 0; i < 256*1024; i++ {
			data[i] = byte((i % 16) + 'a')
		}
		for {
			//time.Sleep(time.Millisecond)

			//data := fmt.Sprintf("---- %v", time.Now().Format(time.RFC3339Nano))
			for i := 0; i < 10000; i++ {
				_ = w.Write([]byte(data))
			}

			sendSize += 1000 * len(data)

			log.Printf("发送: %v  %v MB", sendSize, float64(sendSize)/1024/1024/time.Since(a).Seconds()) //
		}
	}

	// output:

}

func ExampleShowInfo() {
	ctrlMem, err := shm.Create(shmName+"_ctrl", 1024)
	if err != nil {
		return
	}

	buffer := make([]byte, 48)
	for {
		time.Sleep(time.Second)

		_, _ = ctrlMem.ReadAt(buffer, 0)

		log.Println(hex.EncodeToString(buffer[0:32]))

		log.Println(binary.BigEndian.Uint64(buffer[0:8]))
		log.Println(binary.BigEndian.Uint64(buffer[8:16]))
		log.Println(binary.BigEndian.Uint64(buffer[16:24]))
		log.Println(binary.BigEndian.Uint64(buffer[24:32]))
	}

	// output:

}
