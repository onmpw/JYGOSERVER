package command

import (
	"fmt"
	"reflect"
	"server/client"
	"server/parser"
)

type Command struct {
	// name 命令名称
	name string

	// param 命令参数
	param []interface{}
	// paramLen 参数个数
	paramLen int

	client *client.Client
}

type Contract interface {
	Exec(cmd *Command)
	Succeed() bool
	Result() string
	Error() error
}

var commandChan chan *Command

var defCmd = map[string]struct {
	num    int    // 默认参数个数
	minNum int    // 至少需要参数个数
	maxNum int    // 最大参数个数
	desc   string // 命令描述
	exec   Contract
}{
	"get": {1, 1, 0, "获取指定参数的值", new(getCmd)},
}

// InitPool 初始化命令队列
func InitPool() <-chan *Command {
	commandChan = make(chan *Command, client.ClientsMax)

	return commandChan
}

// RegisterClient 解析Client的命令
func RegisterClient(c *client.Client) {
	cmd := new(Command)
	err := cmd.setCmd(c)

	if err != nil {
		fmt.Println(err.Error())
	}

	commandChan <- cmd
}

// setCmd 对命令进行设置
func (cmd *Command) setCmd(c *client.Client) error {
	cmd.client = c
	// 解析客户端传来的命令
	val, err := parser.GetArray(c.Data)

	if err != nil {
		c.Err = err
		return err
	}

	err = setCmdName(cmd, val)
	if err != nil {
		c.Err = err
		return err
	}

	setCmdParam(cmd, val)

	return nil
}

// setCmdName 设置命令名称
func setCmdName(cmd *Command, val []interface{}) error {
	name, err := fetchCmdName(val)
	cmd.name = name
	return err
}

// setCmdParam 设置命令参数
func setCmdParam(cmd *Command, val []interface{}) {
	cmd.param, cmd.paramLen = fetchCmdParam(val)
}

// fetchCmdName 获取命令名称
func fetchCmdName(val []interface{}) (cmd string, err error) {

	if len(val) == 0 {
		return cmd, fmt.Errorf("没有检测到命令！")
	}

	c := reflect.ValueOf(val[0])
	ind := reflect.Indirect(c)

	if ind.Kind() != reflect.String {
		return cmd, fmt.Errorf("命令不是字符串类型！")
	}

	cmd = reflect.Indirect(c).Interface().(string)

	return cmd, err
}

// fetchCmdParam 获取命令参数
func fetchCmdParam(val []interface{}) ([]interface{}, int) {
	p := val[1:]

	return p, len(p)
}

// Exec 执行命令
func (cmd *Command) Exec() bool {
	if cmd.client.Err != nil {
		cmd.client.ErrorResponse()
		return false
	}

	// 检测命令是否被定义，并且参数是否合法
	if !checkCmd(cmd) {
		cmd.client.ErrorResponse()
		return false
	}

	e := defCmd[cmd.name]

	e.exec.Exec(cmd)

	cmd.client.Result = e.exec.Result()

	cmd.client.Response()
	return e.exec.Succeed()
}

// checkCmd 检测命令是否被定义，并且参数是否合法
func checkCmd(cmd *Command) bool {
	if _, ok := defCmd[cmd.name]; !ok {
		cmd.client.Err = fmt.Errorf("命令%s没有定义！", cmd.name)
		return false
	}

	oriCmd := defCmd[cmd.name]
	if cmd.paramLen != oriCmd.num && !(cmd.paramLen <= oriCmd.maxNum && cmd.paramLen >= oriCmd.minNum) {
		cmd.client.Err = fmt.Errorf("该命令参数个数不正确")
		return false
	}
	return true
}
