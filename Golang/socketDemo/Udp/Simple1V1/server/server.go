package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	flag.Parse()
	addr, err := net.ResolveUDPAddr("udp", ":8891")
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening: ", err)
		os.Exit(1)
	}
	defer conn.Close()
	for {
		handleClient(conn)
	}

}

func handleClient(conn *net.UDPConn) {
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("Read UDP msg err: ", err)
		return
	}
	fmt.Println("read data", string(data))
	//nowTime := time.Now().Unix()
	fmt.Println(n, remoteAddr)
	//b := make([]byte, 4)
	//binary.BigEndian.PutUint32(b, uint32(nowTime))
	sendMsg := []byte("get msg success")
	_, err = conn.WriteToUDP(sendMsg, remoteAddr)
	if err != nil {
		fmt.Println("write msg err: ", err)
		return
	}
}
