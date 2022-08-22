package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	conn, err := net.Listen("tcp", ":8889")
	if err != nil {
		fmt.Println("conn err ", err)
		return
	}
	defer conn.Close()

	//等待客户端连接
	for {
		conn, err := conn.Accept()
		if err != nil {
			fmt.Println("Accept err ", err)
		} else {
			fmt.Printf("Accept succ conn = %v client ip = %v\n", conn, conn.RemoteAddr().String())
		}
		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		// 1. 等待客户端通过conn发送信息
		// 2. 如果客户端没有write[发送], 那么协程就阻塞在这里
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("connect is closed")
				conn.Close()
			} else {
				fmt.Printf("Read Err: %s \n", err)
			}
			return
		}
		fmt.Printf("客户端%s 发送信息%s \n", conn.RemoteAddr().String(), string(buf[:n]))
	}
}
