package main

import (
	"log"
	"net"
)

var allocator *ClientAllocator

func main() {
	// 加载配置项
	loadConfig()

	// 初始化客户端分配器
	allocator = InitAllocator()

	// 开启消息解析器
	go startParserServer()

	tcpAddr, err := net.ResolveTCPAddr(NETWORK, Host+":"+Port)

	if err != nil {
		log.Panic(err.Error())
	}

	listener, err := net.ListenTCP(NETWORK, tcpAddr)

	if err != nil {
		log.Panic(err.Error())
	}

	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			log.Panic(err.Error())
		}

		err = conn.SetKeepAlive(true)
		if err != nil {
			log.Panic(err.Error())
		}

		go allocator.registerClient(conn)
	}

}
