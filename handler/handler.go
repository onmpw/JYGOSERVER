package handler

import (
	"server/client"
)

var Handlers []*client.Client

func RegisterHandler(c *client.Client) {
	Handlers = append(Handlers, c)
}
