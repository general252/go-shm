package go_shm

import "log"

const (
	shmName = "shm_transport"
	shmSize = int32(64 * 1024 * 1024)
)

func ExampleNewReader() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	r := NewReader()

	if err := r.Open(shmName, shmSize); err != nil {
		log.Println(err)
		return
	}
	defer r.Close()

	buffer := make([]byte, 32*1024)

	log.Println("接收数据:")
	for i := 0; i >= 0; i++ {
		n, err := r.Read(buffer)
		if err != nil {
			log.Println(err)
			return
		}
		_ = n

		if false {
			if i%25000 == 0 {
				log.Printf("recv: %v\n", string(buffer[:10]))
			}
		}
		if true {
			log.Printf("recv: %v\n", string(buffer[:n]))
		}
	}

	// output:

}
