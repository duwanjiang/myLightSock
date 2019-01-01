package main

import (
	"./core"
	"./server"
	"fmt"
	"net"
)
func main()  {
	password := core.SetPassword("123456")
	startServer(password)
}

//启动远程服务
func startServer(password *core.Password)  {
	localTcpAddr, err := net.ResolveTCPAddr("tcp4", ":9000")
	ser := server.New(password,localTcpAddr)
	ser.Listen(func(listenAddr net.Addr) {
		fmt.Print("start server successfully, listenAddr:",listenAddr)
	})
	fmt.Println(err)
}