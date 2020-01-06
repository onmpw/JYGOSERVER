package main

import (
	"server/handler"
)

func startHandlerServer() {
	cc := handler.InitHandlerPool()

	handler.RegisterHandler()

	for {
		select {
		case c := <-cc:
			handler.Handle(c)
		}
	}
}
