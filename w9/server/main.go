package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"studytcpzb/proto"
)

// socket_stick/server/main.go
func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		msg, err := proto.Decode(reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		fmt.Println("收到client发来的数据：", msg)
	}
}

func process2(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF && data[0] == 'V' {
			// 第一个字节为V 版本号
			if len(data) > 4 {
				length := int16(0)
				// 长度保存在 3，4字节
				binary.Read(bytes.NewReader(data[2:4]), binary.BigEndian, &length)
				if int(length)+4 <= len(data) {
					 // 说明读到的数据是完整的
					return int(length) + 4, data[:int(length)+4], nil
				}
			}
		}
		return
	})
	for scanner.Scan() {
		scannedPack := new(proto.Package)
		scannedPack.Unpack(bytes.NewReader(scanner.Bytes()))
		log.Println(scannedPack)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("无效数据包")
	}
}


func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process2(conn)
	}
}