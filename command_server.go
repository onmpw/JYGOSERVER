package main

import (
	"fmt"
	"server/command"
)

func startCommandServer() {
	cc := command.InitPool()

	for {
		select {
		case cmd := <-cc:
			r := cmd.Exec()
			fmt.Println(r)
		}
	}
}
