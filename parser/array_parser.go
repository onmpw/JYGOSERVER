package parser

import "fmt"

type ArrayParser struct {
	Parser 			*ParseVal
	identifier 		int
}

type ArrayVal struct {
	num 	int   // 元素个数
	val 	[]interface{}
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

func (parse *ArrayParser) Parse() (interface{},error){
	val := new(ParseVal)

	fmt.Println(val)
	fmt.Println("ArrayParser")

	return val,nil
}
