package parser

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Str int = iota + 1
	Int
	Float
	Array
)

const (
	// 分隔符
	Splitter = "\\r\\n"

	// 命令标识符
	CmdIdentifier = "*"

	// 普通消息标识符
	NormalIdentifier = "$"
)

type OneByte byte

// ParseVal 存储解析后的数据
type ParseVal struct {
	// OriginMess 存储原始消息
	OriginMess string

	// ActiveMess 下一步要解析的数据
	ActiveMess string

	// TypeKey  数据类型标识
	TypeKey int

	// TypeVal  数据类型描述
	TypeVal string

	// Data 	消息数据
	Data ParseValContract
}

type Val struct {
	// ErrorFlag 是否有错误
	ErrorFlag bool
	// Error 错误信息
	Error error

	// MesType 	消息类型标识
	MesType OneByte

	// OriginMess 存储原始消息
	OriginMess string

	// ActiveMess 下一步要解析的数据
	ActiveMess string

	// Value 存储解析出来的消息内容
	Value *ParseVal
}

// ParseContract 接口
type ParseContract interface {
	myType() int
	Parse() (ParseValContract, error)
}

type ParseValContract interface {
	GetValue() interface{}
	GetType() int
}

var T map[int]string

// Init 初始化解析器
func Init() {
	T = make(map[int]string)
	T[Str] = "string"
	T[Int] = "int"
	T[Float] = "float"
	T[Array] = "array"
}

// String 重写string的String()方法
func (mt OneByte) String() string {
	switch mt {
	case 36:
		return NormalIdentifier
	case 42:
		return CmdIdentifier
	default:
		return ""
	}
}

// Parser 开始解析
func Parser(message string) (*Val, bool) {
	var jv = new(Val)

	jv.OriginMess = message
	jv.ActiveMess = message

	err := jv.readStart()
	if err != nil {
		jv.setError(err)
		return jv, false
	}

	val, err := startParse(jv.ActiveMess)

	if err != nil {
		jv.setError(err)
		return jv, false
	}

	jv.Value = val

	return jv, true
}

// startParse 开始进行消息内容的解析
func startParse(message string) (*ParseVal, error) {
	var val = new(ParseVal)

	val.OriginMess = message
	val.ActiveMess = message

	// 读取 消息类型
	t, err := val.readBySplitter()
	if err != nil {
		return val, err
	}

	dataType, err := strconv.Atoi(t)
	if err != nil {
		return val, err
	}

	val.TypeKey = dataType

	err = val.ParseStrategy()

	return val, err
}

// ParseStrategy 策略模式
// 根据消息类型使用相应的解析方法获取数据
func (val *ParseVal) ParseStrategy() error {
	var parsers []ParseContract
	var p ParseContract

	// 注册解析器
	parsers = append(parsers, newStringParser(val))
	parsers = append(parsers, newIntParser(val))
	parsers = append(parsers, newFloatParser(val))
	parsers = append(parsers, newArrayParser(val))

	for _, parse := range parsers {
		if parse.myType() == val.TypeKey {
			p = parse
			break
		}
	}

	if p == nil {
		return fmt.Errorf("消息类型不合法 %d", val.TypeKey)
	}

	pv, err := p.Parse()

	val.Data = pv

	return err
}

// readStart 读取消息开头
func (val *Val) readStart() (err error) {
	bv, err := readByIndex(strings.NewReader(val.ActiveMess), 0)

	if err != nil {
		return err
	}
	val.ActiveMess = val.ActiveMess[1:]
	val.MesType = OneByte(bv)

	return nil
}

// readByIndex 获取指定索引的字节
func readByIndex(reader *strings.Reader, offset int64) (rb byte, e error) {
	var b = make([]byte, 1)

	_, e = reader.ReadAt(b, offset)

	if e != nil {
		return rb, e
	}
	rb = b[0]

	return rb, e
}

// readBySplitter 通过分隔符读取数据
func (val *ParseVal) readBySplitter() (str string, err error) {
	index := strings.Index(val.ActiveMess, Splitter)

	if index == -1 {
		return str, fmt.Errorf("消息格式错误，缺少分割符/结束符 '%s'", Splitter)
	}

	str = val.ActiveMess[0:index]
	val.ActiveMess = val.ActiveMess[index+len(Splitter):]
	return str, nil
}

// setError 设置错误信息
func (val *Val) setError(err error) {
	val.ErrorFlag = true
	val.Error = err
}

// GetString 获取从消息中解析出来的string值
func GetString(val *Val) (str string, err error) {
	// 首先判断消息是否是字符串类型
	// 如果是，则返回具体值； 不是 则返回错误信息
	if val.Value.TypeKey != Str {
		return str, fmt.Errorf("您的消息不是期望的字符串类型，而是%s", T[val.Value.TypeKey])
	}

	v := convertData(val.Value.Data, Str)

	str = fmt.Sprintf("%s", v)

	return str, err
}

func convertData(data ParseValContract, t int) (r interface{}) {
	if data.GetType() == t && t != Array {
		r = data.GetValue()
	}

	return r
}
