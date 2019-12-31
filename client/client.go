package client

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"server/command"
	"server/parser"
	"server/storage"
	"sync"
)

const (
	BufferSize = 3
)

type Client struct {
	Err 	error
	ID 		uint64
	Conn 	net.Conn
	Created string
	ActiveTime	string
	Message 	*storage.Message
	Data 	*parser.Val
	Cmd 	*command.Command
}

type Pool struct {
	sync.RWMutex
	num 	int
	pool 	[]*Client
}

// push 将client加入队列
func (clientPool *Pool) Push(client *Client) bool{
	clientPool.Lock()
	clientPool.pool = append(clientPool.pool,client)
	clientPool.num++
	clientPool.Unlock()

	return true
}

// pop 从队列中获取client
func (clientPool *Pool) Pop() (client *Client) {
	clientPool.Lock()
	if clientPool.num > 0 {
		client = clientPool.pool[0]
		clientPool.pool = clientPool.pool[1:]
		clientPool.num--
	}
	clientPool.Unlock()
	return client
}

// readMessage 读取客户端传来的消息
func (client *Client) ReadMessage () {
	num,mess := client.readMessageFromConn(client.Conn)

	message := new(storage.Message)
	message.Len = num
	message.Content = mess
	client.Message = message
}

// readMessageFromConn 从连接中读取数据
func (client *Client)readMessageFromConn(conn net.Conn) (int,string) {
	var buffer = bytes.NewBuffer(make([]byte,0, BufferSize))
	var rb = make([]byte, BufferSize)
	var len = 0

	for {
		num, err := conn.Read(rb)
		len+=num

		if err != nil {
			if err == io.EOF {
				break
			}
			client.Err = err
			return 0,""
		}

		buffer.Write(rb[0:num])

		if num < BufferSize {
			break
		}
	}

	return len,buffer.String()
}

func (client *Client) Response () {
	if client.Err != nil {
		client.ErrorResponse()
	}
}

func (client *Client) ErrorResponse() {
	errMes := fmt.Sprintf("错误:%s",client.Err.Error())

	err := client.response(errMes)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func (client *Client) response (mes string) error {
	num , err := client.Conn.Write([]byte(mes))

	if err != nil {
		return err
	}
	
	if num < len(mes) {
		return fmt.Errorf("数据写入不完整，需发送数据长度：%d,已发送数据长度：%d",len(mes),num)
	}

	err = client.Close()

	return err
}

func (client *Client) Close() error{
	client.Data = nil
	err := client.Conn.Close()

	return err
}
