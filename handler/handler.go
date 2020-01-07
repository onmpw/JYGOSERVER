package handler

import (
	"log"
	"server/client"
)

type ContractHandler interface {
	Init() error
	Handle(c *client.Client) bool
}

var handler ContractHandler

var clientChan chan *client.Client

func InitHandlerPool() <-chan *client.Client {
	clientChan = make(chan *client.Client, client.ClientsMax)

	return clientChan
}

func RegisterHandler() {
	handler = new(OrderHandler)

	err := handler.Init()

	if err != nil {
		log.Panic(err.Error())
	}
}

func RegisterClient(c *client.Client) {
	if c != nil {
		clientChan <- c
	}
}

func Handle(c *client.Client) {
	handler.Handle(c)
}
