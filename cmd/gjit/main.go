package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type FuncVisitor struct{}

func (v *FuncVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.FuncDecl:
		if t.Name.IsExported() == false {
			break
		}
		var buf bytes.Buffer
		buf.WriteString("func ")
		if t.Recv != nil {
			buf.WriteString("(")
			buf.WriteString(t.Recv.List[0].Names[0].Name)
			switch rn := t.Recv.List[0].Type.(type) {
			case *ast.StarExpr:
				buf.WriteString(fmt.Sprintf(" *%s) ", rn.X.(*ast.Ident).Name))
			case *ast.Ident:
				buf.WriteString(fmt.Sprintf(" %s) ", rn.Name))
			}
		}
		buf.WriteString(t.Name.Name)
		buf.WriteString("(")
		paramList := []string{}
		for _, x := range t.Type.Params.List {
			var newBuf bytes.Buffer
			if len(x.Names) != 1 {
				n := []string{}
				for _, z := range x.Names {
					n = append(n, z.Name)
				}
				newBuf.WriteString(strings.Join(n, ", "))
			} else {
				newBuf.WriteString(x.Names[0].Name)
			}
			newBuf.WriteString(" ")
			switch rn := x.Type.(type) {
			case *ast.ArrayType:
				newBuf.WriteString("[]")
				newBuf.WriteString(rn.Elt.(*ast.Ident).Name)
			case *ast.MapType:
				newBuf.WriteString(fmt.Sprintf("map[%s]%s", rn.Key.(*ast.Ident).Name, rn.Value.(*ast.Ident).Name))
			case *ast.Ident:
				newBuf.WriteString(rn.Name)
			}
			paramList = append(paramList, newBuf.String())
		}
		buf.WriteString(strings.Join(paramList, ", "))
		buf.WriteString(")")
		if t.Type.Results.NumFields() > 0 {
			buf.WriteString(" ")
			if t.Type.Results.NumFields() > 1 {
				buf.WriteString("(")
			}
			fieldList := []string{}
			for _, x := range t.Type.Results.List {
				var newBuf bytes.Buffer
				if x.Names != nil {
					newBuf.WriteString(x.Names[0].Name)
					newBuf.WriteString(" ")
				}
				switch rn := x.Type.(type) {
				case *ast.StarExpr:
					newBuf.WriteString("*")
					switch sep := rn.X.(type) {
					case *ast.SelectorExpr:
						newBuf.WriteString(fmt.Sprintf("%s.%s", sep.X.(*ast.Ident).Name, sep.Sel.Name))
					case *ast.Ident:
						newBuf.WriteString(sep.Name)
					}
				case *ast.Ident:
					newBuf.WriteString(fmt.Sprintf("%s", rn.Name))
				}
				fieldList = append(fieldList, newBuf.String())
			}
			buf.WriteString(strings.Join(fieldList, ", "))
			if t.Type.Results.NumFields() > 1 {
				buf.WriteString(")")
			}
		}

		fmt.Printf("%s\n", buf.String())
	}

	return v
}

func main() {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, os.Args[1], nil, 0)
	if err != nil {
		panic(err)
	}

	ast.Walk(new(FuncVisitor), file)
}
