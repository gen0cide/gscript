package compiler

import "github.com/robertkrimen/otto/parser"

func ValidateAST(source []byte) error {
	_, err := parser.ParseFile(nil, "", source, 0)
	return err
}
