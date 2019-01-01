package main

import (
	"./client"
	"./core"
	"fmt"
	"net"
)
func main()  {
	password := core.SetPassword("123456")
	localTcpAddr, err := net.ResolveTCPAddr("tcp4", ":9001")
	remoteAddr, err := net.ResolveTCPAddr("tcp4", ":9000")
	loc := client.New(password,localTcpAddr,remoteAddr)
	loc.Listen(func(listenAddr net.Addr) {
		fmt.Print("start local successfully, listenAddr:",listenAddr)

	})
	fmt.Println(err)
}

