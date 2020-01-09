package main

import (
	"io"
	"math/rand"
	"net"
	"server/client"
	"time"
)

type ClientAllocator struct {
	clientPool *client.Pool
	clientChan chan *client.Client
}

func InitAllocator() *ClientAllocator {
	allocator := &ClientAllocator{
		clientPool: new(client.Pool),
		clientChan: make(chan *client.Client, client.ClientsMax),
	}
	return allocator
}

// registerClient 注册新的client
func (allocator *ClientAllocator) registerClient(conn net.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(30 * time.Minute))
	for {
		c := makeClient(conn)

		mesLen, err := c.ReadMessage()

		if mesLen == -1 || (mesLen == 0 && err == io.EOF) {
			_ = c.Close(true)
			break
		}

		//allocator.clientPool.Push(c)
		allocator.clientChan <- c
	}
}

// makeClient 创建Client
func makeClient(conn net.Conn) (c *client.Client) {
	c = new(client.Client)
	// 生成clientId
	rand.Seed(time.Now().UnixNano())
	c.ID = rand.Uint64()
	c.Created = time.Unix(time.Now().Unix(), 0).Format(DateFormat)
	c.ActiveTime = time.Unix(time.Now().Unix(), 0).Format(DateFormat)

	c.Conn = conn

	return c
}
