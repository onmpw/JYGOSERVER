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

	tcpAddr,err := net.ResolveTCPAddr(NETWORK,Host+":"+Port)

	if err != nil {
		log.Panic(err.Error())
	}

	listener ,err := net.ListenTCP(NETWORK,tcpAddr)

	if err != nil {
		log.Panic(err.Error())
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Panic(err.Error())
		}

		go allocator.registerClient(conn)

		/*if allocator.clientPool.num > 0 {
			go func() {
				for i:=0; i < allocator.clientPool.num; i++ {
					fmt.Println(allocator.clientPool.pool[i].Content)
				}
			}()
		}*/

	}

}
