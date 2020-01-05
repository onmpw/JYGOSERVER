package handler

import "server/client"

type OrderHandler struct {
}

func (handle *OrderHandler) Handle(c *client.Client) bool {
	c.Response()
	return true
}
