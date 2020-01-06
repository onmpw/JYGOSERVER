package main

import (
	"server/client"
	"server/command"
	"server/handler"
	"server/parser"
)

func startParserServer() {
	for {
		c := allocator.clientPool.Pop()
		if c != nil {
			// 提取client的消息 开始进行解析
			val, suc := parser.Parser(c.Message.Content)
			c.Data = val
			if !suc {
				c.Err = val.Error
				c.Response()
				continue
			}
			SetClient(c)
		}
	}
}

func SetClient(c *client.Client) {
	if c.Data.MesType.String() == parser.CmdIdentifier {
		command.RegisterClient(c)
	} else if c.Data.MesType.String() == parser.NormalIdentifier {
		handler.RegisterClient(c)
	}
}
