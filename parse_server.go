package main

import (
	"server/client"
	"server/command"
	"server/handler"
	"server/parser"
	"time"
)

func startParserServer() {
	// 初始化解析器
	parser.Init()

	// 开始解析
	for {
		parse2()

		//parse()
	}
}

// SetClient 将解析完消息的客户端放入相应的队列
func SetClient(c *client.Client) {
	if c.Data.MesType.String() == parser.CmdIdentifier {
		command.RegisterClient(c)
	} else if c.Data.MesType.String() == parser.NormalIdentifier {
		handler.RegisterClient(c)
	}
}

func parse() {
	c := allocator.clientPool.Pop()
	if c != nil {
		// 提取client的消息 开始进行解析
		val, suc := parser.Parser(c.Message.Content)
		c.Data = val
		if !suc {
			c.Err = val.Error
			c.Response()
			return
		}
		SetClient(c)
	}
	time.Sleep(time.Duration(10) * time.Microsecond)
}

func parse2() {
	select {
	case c := <-allocator.clientChan:
		// 提取client的消息 开始进行解析
		val, suc := parser.Parser(c.Message.Content)
		c.Data = val
		if !suc {
			c.Err = val.Error
			c.Response()
			return
		}
		SetClient(c)
	}
}
