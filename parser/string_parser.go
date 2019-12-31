package parser

import (
	"strconv"
)

type StringParser struct {
	Parser 			*ParseVal
	identifier 		int
}

type StringVal struct {
	len 	int
	val 	string
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

// Parse 解析字符串类型的消息
func (parse *StringParser) Parse() (interface{},error) {
	val := new(StringVal)
	var err error

	val.len, err = parse.getLength()

	if err != nil {
		return val, err
	}

	val.val, err = parse.getContent()

	if err != nil {
		return val,err
	}

	return val,nil
}

func (parse *StringParser) getLength() (int,error) {
	v,err := parse.Parser.readBySplitter()

	if err != nil {
		return 0,err
	}
	strLen,err := strconv.Atoi(v)

	if err != nil {
		return 0,err
	}

	return strLen,nil
}

func (parse *StringParser) getContent() (string,error) {
	v,err := parse.Parser.readBySplitter()
	if err != nil {
		return "",err
	}
	return v,nil
}
