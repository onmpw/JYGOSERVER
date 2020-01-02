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

const Splitter = "\\r\\n"

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
	Data interface{}
}

type Val struct {
	ErrorFlag bool
	Error     error
	// MesType 	消息类型标识
	MesType OneByte

	// OriginMess 存储原始消息
	OriginMess string

	// ActiveMess 下一步要解析的数据
	ActiveMess string

	Value *ParseVal
}

func (mt OneByte) String() string {
	switch mt {
	case 36:
		return "$"
	case 42:
		return "*"
	default:
		return ""
	}
}

type ParseContract interface {
	myType() int
	Parse() (interface{}, error)
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

func (val *ParseVal) readBySplitter() (str string, err error) {
	index := strings.Index(val.ActiveMess, Splitter)

	if index == -1 {
		return str, fmt.Errorf("消息格式错误，缺少分割符/结束符 '%s'", Splitter)
	}

	str = val.ActiveMess[0:index]
	val.ActiveMess = val.ActiveMess[index+len(Splitter):]
	return str, nil
}

func (val *Val) setError(err error) {
	val.ErrorFlag = true
	val.Error = err
}
