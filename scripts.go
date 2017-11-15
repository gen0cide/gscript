package gscript

import "github.com/robertkrimen/otto/parser"

func (e *Engine) LoadScript(source []byte) error {
	_, err := e.VM.Run(string(source))
	return err
}

func (e *Engine) ValidateAST(source []byte) error {
	_, err := parser.ParseFile(nil, "", source, 0)
	return err
}
