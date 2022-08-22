package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", ":8889")
	if err != nil {
		fmt.Println("listen err ", err)
		return
	}
	defer conn.Close()

	//等待客户端连接
	for {
		fmt.Println("请输入信息, 以回车结束")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("readString err ", err)
		}
		line = strings.Trim(line, "\r\n")
		if line == "exit" {
			fmt.Println("客户端退出")
			break
		}
		line = strings.TrimSpace(line)
		n, err := conn.Write([]byte(line))
		if err != nil {
			fmt.Println("conn.Write err ", err)
		}
		fmt.Printf("发送的内容 %d 文字\n", n)
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
