package main

import (
	"server/handler"
)

func startHandlerServer() {
	cc := handler.InitPool()

	handler.RegisterHandler()

	for {
		select {
		case c := <-cc:
			handler.Handle(c)
		}
	}
}
