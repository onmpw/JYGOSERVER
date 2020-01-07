package parser

import (
	"strconv"
)

type IntParser struct {
	Parser     *ParseVal
	identifier int
}

type IntVal struct {
	val int
}

func newIntParser(p *ParseVal) *IntParser {
	parse := new(IntParser)
	parse.identifier = Int
	parse.Parser = p
	return parse
}

func (parse *IntParser) myType() int {
	return parse.identifier
}

// Parse 实现了 ParseContract 的 Parse方法
// Int数据解析入口方法
func (parse *IntParser) Parse() (ParseValContract, error) {
	val := new(IntVal)

	con, err := parse.getContent()

	val.val = con

	if err != nil {
		return val, err
	}

	return val, nil
}

// 获取Int数据
func (parse *IntParser) getContent() (content int, err error) {
	v, err := parse.Parser.readBySplitter()

	if err != nil {
		return 0, err
	}
	content, err = strconv.Atoi(v)

	if err != nil {
		return 0, err
	}

	return content, nil
}

// GetValue 实现 ParseValContract接口GetValue
// 获取解析后的结果中的整数值
func (val *IntVal) GetValue() interface{} {
	return val.val
}

// GetType 实现 ParseValContract 接口 GetType
// 获取值的类型
func (val *IntVal) GetType() int {
	return Int
}
