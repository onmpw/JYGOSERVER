package client

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/onmpw/JYGO/config"
	"io"
	"net"
	"server/parser"
	"strings"
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
	// Err 错误信息
	// 如果在解析过程或者处理消息过程中以后错误，则存入该字段。如果没有则为nil
	Err error

	// ID 客户端ID
	ID uint64

	// Conn 客户端链接
	Conn net.Conn

	// 是否经过验证
	auth bool

	user string

	password string

	// Created 创建时间
	Created string
	// ActiveTime 活动时间
	ActiveTime string

	// Message 客户端传来的原始消息
	Message *Message

	// Data 对消息进行解析之后的值
	Data *parser.Val

	// Result 执行完之后的需要回传给客户端的结果值
	Result string
}

// Pool 客户端队列
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

// GetLen 获取队列当前长度
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
		if Len == mesLen {
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
		err := client.response(fmt.Sprintf("Success:'%s'", client.Result))

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

// ErrorResponse 错误响应
func (client *Client) ErrorResponse() {
	errMes := fmt.Sprintf("Failed:%s", client.Err.Error())

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

// SetAuth 设置客户端是否验证字段
func (client *Client) SetAuth(auth bool) {
	client.auth = auth
}

// IsAuth 是否已经验证成功
func (client *Client) IsAuth() bool {
	return client.auth
}

// SetUserInfo 设置当前客户端的用户名和密码
func (client *Client) SetUserInfo(user string, password string) {
	client.user = user
	client.password = password
}

// GetUserName 获取当前登录的用户
func (client Client) GetUserName() string {
	return client.user
}

// Auth 开始进行验证
func (client *Client) Auth() bool {
	if client.IsAuth() {
		return true
	}
	var stringArr []string
	authInfo, err := parser.GetString(client.Data)

	if err != nil {
		client.Err = err
		return false
	}

	stringArr = strings.Split(authInfo, "@")
	if len(stringArr) != 2 {
		client.Err = fmt.Errorf("在进行实际业务之前，请先验证用户！")
		return false
	}

	username := stringArr[0]
	password := stringArr[1]

	if ok := authUser(username, password); ok {
		client.SetAuth(ok)
		client.SetUserInfo(username, password)
		client.Result = "用户验证成功！"
		return true
	}

	client.Err = fmt.Errorf("用户验证失败，请检查用户名和密码")

	return false
}

// authUser 验证用户名和密码是否正确
func authUser(username string, password string) (ok bool) {
	if username == config.Conf.C("User") && password == config.Conf.C("Password") {
		ok = true
	} else {
		ok = false
	}
	return ok
}

// Close 关闭client
func (client *Client) Close(closeConn bool) (err error) {
	client.Data = nil
	client.Message = nil
	client.Err = nil

	if closeConn {
		err = client.Conn.Close()
		client.Conn = nil
	}

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
