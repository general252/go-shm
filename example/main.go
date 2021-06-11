package main

import (
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	goShm "github.com/general252/go-shm"
	"github.com/hidez8891/shm"
	"log"
	"time"
)

const (
	shmName = "shm_transport"
	shmSize = int32(128 * 1024 * 1024)
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	mode := flag.String("mode", "debug", "run mode. (debug/write/read)")
	bufferSize := flag.Int("size", 8*1024, "buffer size [byte]")
	flag.Parse()

	if *bufferSize <= 0 {
		log.Printf("error size: %v", *bufferSize)
		return
	}

	log.Printf("mode: %v size: %v", *mode, *bufferSize)

	switch *mode {
	case "debug":
		debug()
	case "read":
		read(*bufferSize)
	case "write":
		write(*bufferSize)
	default:
		log.Printf("unknown mode %v", *mode)
	}
}

func write(bufSize int) {
	w := goShm.NewWriter()

	if err := w.Open(shmName, shmSize); err != nil {
		log.Println(err)
		return
	}
	defer w.Close()

	dataLen := bufSize
	data := make([]byte, dataLen)
	for i := 0; i < len(data); i++ {
		data[i] = byte((i % 10) + '0')
	}

	var (
		sendSize  int
		lastSize  int
		firstTime = time.Now()
		lastTime  = time.Now()
	)

	for i := 0; i >= 0; i++ {
		time.Sleep(time.Nanosecond)

		for i := 0; i < 150; i++ {
			_ = w.Write(data)

			sendSize += dataLen
		}

		if true {
			now := time.Now()
			sec := now.Sub(lastTime).Seconds()
			if sec > 1 {
				secSize := sendSize - lastSize
				msg := fmt.Sprintf("send: %10.2f MB  %10.2f MB/s [per] %10.2f MB/s [total] %20v [total bytes]",
					float64(sendSize)/1024/1024,
					float64(secSize)/1024/1024/sec,
					float64(sendSize)/1024/1024/time.Since(firstTime).Seconds(),
					sendSize)
				log.Println(msg)

				copy(data, msg)

				lastTime = now
				lastSize = sendSize
			}
		}
	}
}

func read(bufSize int) {
	r := goShm.NewReader()

	if err := r.Open(shmName, shmSize); err != nil {
		log.Println(err)
		return
	}
	defer r.Close()

	buffer := make([]byte, bufSize)

	var (
		recvSize  int
		lastSize  int
		firstTime = time.Now()
		lastTime  = time.Now()
	)

	for i := 0; i >= 0; i++ {
		n, err := r.Read(buffer)
		if err != nil {
			log.Println(err)
			return
		}
		recvSize += n

		if true {
			now := time.Now()
			sec := now.Sub(lastTime).Seconds()
			if sec > 1 {
				secSize := recvSize - lastSize

				maxLen := 0
				if len(buffer) > 20 {
					maxLen = 20
				}

				log.Printf("recv: %20v %10.2f MB  %10.2f MB/s [per] %10.2f MB/s [total] \"%v\" ",
					recvSize,
					float64(recvSize)/1024/1024,
					float64(secSize)/1024/1024/sec,
					float64(recvSize)/1024/1024/time.Since(firstTime).Seconds(),
					string(buffer[:maxLen]))

				lastTime = now
				lastSize = recvSize
			}
		}
	}
}

func debug() {
	ctrlMem, err := shm.Create(shmName+"_ctrl", 1024)
	if err != nil {
		return
	}

	var printUint64 = func(head string, data uint64) {
		log.Printf("%v %20v %20x", head, data, data)
	}

	buffer := make([]byte, 48)
	for {
		time.Sleep(time.Second)

		_, _ = ctrlMem.ReadAt(buffer, 0)

		log.Println(hex.EncodeToString(buffer[0:32]))

		printUint64("read length  : ", binary.BigEndian.Uint64(buffer[0:8]))
		printUint64("read offset  : ", binary.BigEndian.Uint64(buffer[8:16]))
		printUint64("write length : ", binary.BigEndian.Uint64(buffer[16:24]))
		printUint64("write offset : ", binary.BigEndian.Uint64(buffer[24:32]))
	}
}
