package handler

import (
	"server/client"
)

type ContractHandler interface {
	Handle(c *client.Client) bool
}

var handler ContractHandler

var clients []*client.Client

func RegisterHandler() {
	handler = new(OrderHandler)
}

func RegisterClient(c *client.Client) {
	clients = append(clients, c)
}

func Pop() *client.Client {
	if len(clients) <= 0 {
		return nil
	}

	c := clients[0]

	clients = clients[1:]

	return c
}

func Handle(c *client.Client) {
	handler.Handle(c)
}
