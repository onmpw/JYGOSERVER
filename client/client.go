package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"server/parser"
	"sync"
)

const (
	BufferSize = 100
	ClientsMax = 1024
)

type Message struct {
	Content string
	Len     int
}

type Client struct {
	Err        error
	ID         uint64
	Conn       net.Conn
	Created    string
	ActiveTime string
	Message    *Message
	Data       *parser.Val
}

type Pool struct {
	sync.RWMutex
	num  int
	pool []*Client
}

// push 将client加入队列
func (clientPool *Pool) Push(client *Client) bool {
	clientPool.Lock()
	clientPool.pool = append(clientPool.pool, client)
	clientPool.num++
	clientPool.Unlock()

	return true
}

func (clientPool *Pool) GetLen() int {
	return clientPool.num
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

// ReadMessage 读取客户端传来的消息
func (client *Client) ReadMessage() (num int, err error) {
	num, mess, err := client.readMessageFromConn(client.Conn)

	message := new(Message)
	message.Len = num
	message.Content = mess
	client.Message = message

	return num, err
}

// readMessageFromConn 从连接中读取数据
func (client *Client) readMessageFromConn(conn net.Conn) (Len int, str string, err error) {
	var buffer = bytes.NewBuffer(make([]byte, 0))
	var rb = make([]byte, BufferSize)
	var length = make([]byte, 8)
	var num int

	reader := bufio.NewReader(conn)

	num, err = reader.Read(length)

	mesLen := byteToInt(length)

	for {
		num, err = reader.Read(rb)
		Len += num

		if err != nil {
			if err == io.EOF {
				break
			}
			client.Err = err
			return -1, "", err
		}

		buffer.Write(rb[0:num])
		if num < BufferSize || Len == mesLen {
			break
		}
	}

	str = buffer.String()

	return Len, str, err
}

// Response 向客户端返回信息
func (client *Client) Response() {
	if client.Err != nil {
		client.ErrorResponse()
	} else {
		err := client.response(fmt.Sprintf("成功,您发送的消息为%s", client.Data.Value.Data))

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

// ErrorResponse 错误响应
func (client *Client) ErrorResponse() {
	errMes := fmt.Sprintf("错误:%s", client.Err.Error())

	err := client.response(errMes)

	if err != nil {
		fmt.Println(err.Error())
	}
}

// response 响应客户端
func (client *Client) response(mes string) error {
	num, err := client.Conn.Write([]byte(mes))

	if err != nil {
		fmt.Println(client, err)
		return err
	}

	if num < len(mes) {
		return fmt.Errorf("数据写入不完整，需发送数据长度：%d,已发送数据长度：%d", len(mes), num)
	}

	err = client.Close(false)

	return err
}

// Close 关闭client
func (client *Client) Close(closeConn bool) (err error) {
	client.Data = nil
	client.Message = nil

	if closeConn {
		err = client.Conn.Close()
	}
	client.Conn = nil

	return err
}

// byteToInt 用于解析客户端传来的消息长度
func byteToInt(bytes []byte) (r int) {
	var j = 7
	var d = 48
	var a = 39
	for i := 0; i < len(bytes); i++ {
		b := int(bytes[i])
		if b > d+a {
			b = b - d - a
		} else {
			b = b - d
		}
		r |= b << uint((j-i)*4) & 0xffffffff
	}
	return r
}
