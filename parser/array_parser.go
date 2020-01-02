package parser

import (
	"fmt"
	"strconv"
)

type ArrayParser struct {
	Parser     *ParseVal
	identifier int
}

type ArrayVal struct {
	num int // 元素个数
	val []*ParseVal
}

func newArrayParser(p *ParseVal) *ArrayParser {
	parse := new(ArrayParser)
	parse.identifier = Array
	parse.Parser = p
	return parse
}

func (parse *ArrayParser) myType() int {
	return parse.identifier
}

// Parse 实现了 ParseContract 的 Parse方法
// Array数据解析入口方法
func (parse *ArrayParser) Parse() (interface{}, error) {
	val := new(ArrayVal)

	num, err := parse.getElementNum()

	if err != nil {
		return val, err
	}

	if num < 0 {
		return val, fmt.Errorf("数组元素个数必须为正整数，而消息中的为负数：%d", num)
	}

	val.num = num

	for i := 0; i < num; i++ {
		v, err := startParse(parse.Parser.ActiveMess)

		if err != nil {
			return val, err
		}

		val.val = append(val.val, v)

		parse.Parser.ActiveMess = v.ActiveMess
	}

	return val, nil
}

// 获取数组元素个数
func (parse *ArrayParser) getElementNum() (num int, err error) {
	v, err := parse.Parser.readBySplitter()

	if err != nil {
		return -1, err
	}

	num, err = strconv.Atoi(v)

	return num, err
}
