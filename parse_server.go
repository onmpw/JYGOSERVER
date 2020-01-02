package main

import (
	"server/parser"
)

func startParserServer() {
	for {
		client := allocator.clientPool.Pop()
		if client != nil {
			// 提取client的消息 开始进行解析
			val, suc := parser.Parser(client.Message.Content)
			client.Data = val
			if !suc {
				client.Err = val.Error
				client.Response()
				continue
			}
			client.Response()
		}
	}
}
