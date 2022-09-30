package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	addr, err := net.ResolveUDPAddr("udp", ":8891")
	if err != nil {
		fmt.Println("resolve addr err: ", err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("dial err: ", err)
		os.Exit(1)
	}
	defer conn.Close()
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

		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Println("failed: ", err)
			os.Exit(1)
		}
		data := make([]byte, 1024)
		_, err = conn.Read(data)
		if err != nil {
			fmt.Println("read err: ", err)
			os.Exit(1)
		}
		fmt.Println("get msg: ", string(data))
	}
	os.Exit(0)

}
