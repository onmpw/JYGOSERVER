package main

import "server/handler"

func startHandlerServer() {
	handler.RegisterHandler()

	for {
		c := handler.Pop()

		if c != nil {
			handler.Handle(c)
		}
	}
}
