package main

import (
	"math/rand"
	"net"
	"server/client"
	"time"
)

type ClientAllocator struct {
	clientPool  *client.Pool
}

func InitAllocator() *ClientAllocator {
	allocator := &ClientAllocator{
		clientPool:new(client.Pool),
	}
	return allocator
}

// registerClient 注册新的client
func (allocator *ClientAllocator) registerClient(conn net.Conn) {
	client := new(client.Client)

	client.Created = time.Unix(time.Now().Unix(),0).Format(DateFormat)
	client.ActiveTime = time.Unix(time.Now().Unix(),0).Format(DateFormat)
	client.Conn = conn

	// 生成clientId
	rand.Seed(time.Now().UnixNano())
	client.ID = rand.Uint64()

	client.ReadMessage()

	allocator.clientPool.Push(client)
}