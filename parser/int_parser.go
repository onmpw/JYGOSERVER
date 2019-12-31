package parser

import "fmt"

type IntParser struct {
	Parser 			*ParseVal
	identifier 		int
}

type IntVal struct {
	val 	int
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

func (parse *IntParser) Parse() (interface{},error){
	val := new(ParseVal)

	fmt.Println("IntParser")

	return val,nil
}
