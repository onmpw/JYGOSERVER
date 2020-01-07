package parser

import (
	"strconv"
)

type FloatParser struct {
	Parser     *ParseVal
	identifier int
}

type FloatVal struct {
	val float64
}

func newFloatParser(p *ParseVal) *FloatParser {
	parse := new(FloatParser)
	parse.identifier = Float
	parse.Parser = p
	return parse
}

func (parse *FloatParser) myType() int {
	return parse.identifier
}

// Parse 实现了 ParseContract 的 Parse方法
// Float数据解析入口方法
func (parse *FloatParser) Parse() (ParseValContract, error) {
	val := new(FloatVal)

	con, err := parse.getContent()

	val.val = con

	if err != nil {
		return val, err
	}

	return val, nil
}

// 获取Float数据
func (parse *FloatParser) getContent() (content float64, err error) {
	v, err := parse.Parser.readBySplitter()

	if err != nil {
		return 0, err
	}
	content, err = strconv.ParseFloat(v, 64)

	if err != nil {
		return 0, err
	}

	return content, nil
}

// GetValue 实现 ParseValContract接口GetValue
// 获取解析后的结果中的浮点数值
func (val *FloatVal) GetValue() interface{} {
	return val.val
}

// GetType 实现 ParseValContract 接口 GetType
// 获取值的类型
func (val *FloatVal) GetType() int {
	return Float
}
