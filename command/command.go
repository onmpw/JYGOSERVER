package command

import (
	"server/client"
)

type Command struct {
	Cmd    string
	Client *client.Client
}

var cmdQueue []*Command

// SetCommandClient 解析Client的命令
func RegisterCmd(c *client.Client) bool {
	cmd := new(Command)
	cmd.Cmd = "get"
	cmd.Client = c

	cmdQueue = append(cmdQueue, cmd)
	return true
}
