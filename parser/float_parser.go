package parser

import "fmt"

type FloatParser struct {
	Parser 			*ParseVal
	identifier 		int
}

type FloatVal struct {
	val 	int
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

func (parse *FloatParser) Parse() (interface{},error){
	val := new(ParseVal)

	fmt.Println(val)
	fmt.Println("FloatParser")

	return val,nil
}
