package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	ip := net.ParseIP("172.24.14.255")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: 9981}
	conn, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	n, err := conn.WriteToUDP([]byte("hello"), dstAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make([]byte, 1024)
	n, _, err = conn.ReadFrom(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("read %s from <%s>\n", data[:n], conn.RemoteAddr())
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
