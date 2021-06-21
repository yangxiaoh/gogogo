package main

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"studytcpzb/proto"
	"time"
)

// socket_stick/client/main.go

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()



	for i := 0; i < 20; i++ {
		//msg := `Hello, Hello. How are you?` + strconv.Itoa(i)
		//b, err := proto.Encode(msg)
		//if err != nil {
		//	continue
		//}
		//conn.Write(b)

		pack := &proto.Package{
			Version:        [2]byte{'V', '1'},
			Timestamp:      time.Now().Unix(),
			TagLength:      4,
			Tag:            []byte("demo"),
			Msg:            []byte("message:" + strconv.Itoa(i)),
		}

		// 数据长度不包括 version和length
		pack.Length = 8 + 2 + pack.TagLength + int16(len(pack.Msg))

		buf := new(bytes.Buffer)
		pack.Pack(buf)
		conn.Write(buf.Bytes())
	}
}