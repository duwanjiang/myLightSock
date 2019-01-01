//实现以上功能的 client 端
package client

import (
"../core"
"log"
"net"
)

type LsLocal struct {
	*core.SecureSocket
}

// 新建一个本地端
func New(password *core.Password, listenAddr, remoteAddr *net.TCPAddr) *LsLocal {
	return &LsLocal{
		SecureSocket: &core.SecureSocket{
			Cipher:     core.NewCipher(password),
			ListenAddr: listenAddr,
			RemoteAddr: remoteAddr,
		},
	}
}

// 本地端启动监听，接收来自本机浏览器的连接
func (local *LsLocal) Listen(didListen func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", local.ListenAddr)
	if err != nil {
		return err
	}

	//等该函数返回后 才执行该方法
	defer listener.Close()

	if didListen != nil {
		didListen(listener.Addr())
	}

	for {
		userConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		userConn.Write([]byte("hello word"))

		// userConn被关闭时直接清除所有数据 不管没有发送的数据
		userConn.SetLinger(0)
		go local.handleConn(userConn)
	}
	return nil
}

func (local *LsLocal) handleConn(userConn *net.TCPConn) {
	defer userConn.Close()

	proxyServer, err := local.DialRemote()
	if err != nil {
		log.Println(err)
		return
	}
	defer proxyServer.Close()
	// Conn被关闭时直接清除所有数据 不管没有发送的数据
	proxyServer.SetLinger(0)

	// 进行转发
	// 从 proxyServer 读取数据发送到 localUser
	go func() {
		err := local.DecodeCopy(userConn, proxyServer)
		if err != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			userConn.Close()
			proxyServer.Close()
		}
	}()
	// 从 localUser 发送数据发送到 proxyServer，这里因为处在翻墙阶段出现网络错误的概率更大
	local.EncodeCopy(proxyServer, userConn)
}