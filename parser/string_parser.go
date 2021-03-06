package parser

import (
	"strconv"
)

type StringParser struct {
	Parser     *ParseVal
	identifier int
}

type StringVal struct {
	len int
	val string
}

func newStringParser(p *ParseVal) *StringParser {
	parse := new(StringParser)
	parse.identifier = Str
	parse.Parser = p
	return parse
}

func (parse *StringParser) myType() int {
	return parse.identifier
}

// Parse 实现了 ParseContract 的 Parse方法
// 解析字符串类型的消息
func (parse *StringParser) Parse() (ParseValContract, error) {
	val := new(StringVal)
	var err error

	val.len, err = parse.getLength()

	if err != nil {
		return val, err
	}

	val.val, err = parse.getContent()

	if err != nil {
		return val, err
	}

	return val, nil
}

// 获取消息的长度
func (parse *StringParser) getLength() (int, error) {
	v, err := parse.Parser.readBySplitter()

	if err != nil {
		return 0, err
	}
	strLen, err := strconv.Atoi(v)

	if err != nil {
		return 0, err
	}

	return strLen, nil
}

// 获取消息内容
func (parse *StringParser) getContent() (string, error) {
	v, err := parse.Parser.readBySplitter()
	if err != nil {
		return "", err
	}
	return v, nil
}

// GetValue 实现 ParseValContract接口GetValue
// 获取解析后的结果中的字符串值
func (val *StringVal) GetValue() interface{} {
	return val.val
}

// GetType 实现 ParseValContract 接口 GetType
// 获取值的类型
func (val *StringVal) GetType() int {
	return Str
}
