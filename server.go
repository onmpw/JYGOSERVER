package main

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

var allocator *ClientAllocator

func main() {

	go startHttpProf()

	// 加载配置项
	loadConfig()

	// 初始化客户端分配器
	allocator = InitAllocator()

	// 开启消息解析器
	go startParserServer()

	// 开启handler处理器
	go startHandlerServer()

	// 开启命令处理器
	go startCommandServer()

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

// startHttpProf 开启http服务用于进行pprof 分析
func startHttpProf() {
	if err := http.ListenAndServe(":6060", nil); err != nil {
		log.Fatalf("pprof failed: %v", err)
	}
}
