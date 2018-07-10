package compiler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
	"text/template"

	"golang.org/x/tools/imports"
)

// StringDef represents an instance of an obfuscated string within
// the gscript compilers intermediate representation
type StringDef struct {
	// unique ID of the string def in relation to the target source tree
	ID string `json:"id"`

	// original value of the string as defined in source
	Value string `json:"value"`

	// key used to encrypt string with
	Key rune `json:"key"`

	// the encrypted data to represent this string
	Data []rune `json:"data"`
}

// Stylist creates a new pre-obfuscation tangle manager
type Stylist struct {
	sync.RWMutex

	// Defs is the StringDef slice for all the source code
	Combs map[*ast.File]*tangleWalker

	// Compiler references the owning compiler object
	Compiler *Compiler
}

type tangleWalker struct {
	sync.RWMutex
	Stylist  *Stylist
	AST      *ast.File
	Filename string
	Defs     []*StringDef
	// SourceBuf is what gets all the rendered string defs appended to it
	SourceBuf *bytes.Buffer
}

// NewStylist creates a new Stylist
func NewStylist(c *Compiler) *Stylist {
	return &Stylist{
		Combs:    map[*ast.File]*tangleWalker{},
		Compiler: c,
	}
}

// LollerSkateDaStringz enumerates the build directory code
func (s *Stylist) LollerSkateDaStringz() error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, s.Compiler.BuildDir, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	for filename, file := range pkgs["main"].Files {
		wg.Add(1)
		go s.ShoutoutToThaHomies(filename, fset, file, errChan, &wg)
	}
	go func() {
		wg.Wait()
		close(finChan)
	}()
	select {
	case <-finChan:
	case err := <-errChan:
		if err != nil {
			return err
		}
	}
	return nil
}

// ShoutoutToThaHomies does stuff
func (s *Stylist) ShoutoutToThaHomies(filename string, fs *token.FileSet, file *ast.File, errChan chan error, wg *sync.WaitGroup) {
	comb := &tangleWalker{
		Filename:  filename,
		AST:       file,
		Stylist:   s,
		Defs:      []*StringDef{},
		SourceBuf: new(bytes.Buffer),
	}
	s.Combs[file] = comb
	ast.Walk(comb, file)
	fileWriter, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		errChan <- err
		wg.Done()
		return
	}
	err = printer.Fprint(fileWriter, fs, file)
	if err != nil {
		errChan <- err
	}
	wg.Done()
	return
}

// BurnTheShitOuttaThisWeave does more stuff
func (comb *tangleWalker) BurnTheShitOuttaThisWeave(key rune, source string) string {
	varName := RandUpperAlphaNumericString(32)
	cipher := fmt.Sprintf("g(%d, %s)", key, varName)
	reader := strings.NewReader(source)
	varDef := []rune{}
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		varDef = append(varDef, ch^key)
		key ^= ch
	}
	comb.Lock()
	comb.Defs = append(comb.Defs, &StringDef{
		ID:    varName,
		Value: source,
		Key:   key,
		Data:  varDef,
	})
	comb.Unlock()
	return cipher
}

// Visit walks the AST
func (comb *tangleWalker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ImportSpec:
		return nil
	case *ast.GenDecl:
		if n.Tok == token.CONST {
			return nil
		}
	case *ast.BasicLit:
		if n.Kind == token.STRING {
			k := rand.Intn(65536)
			n.Value = comb.BurnTheShitOuttaThisWeave(rune(k), n.Value[1:len(n.Value)-1])
		}
	}
	return comb
}

// EasyBakeOven actually renders the tangled hairs
func (comb *tangleWalker) EasyBakeOven(sd *StringDef, errChan chan error, wg *sync.WaitGroup) {
	tmpl := template.New(fmt.Sprintf("obfstring_%s", sd.ID))
	tmpl.Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
	tmplData, err := Asset("obfstring.go.tmpl")
	if err != nil {
		errChan <- err
		wg.Done()
		return
	}
	tmpl, err = tmpl.Parse(string(tmplData))
	if err != nil {
		errChan <- err
		wg.Done()
		return
	}
	rawRender := new(bytes.Buffer)
	rawRender.WriteString("\n")
	err = tmpl.Execute(rawRender, sd)
	if err != nil {
		errChan <- err
		wg.Done()
		return
	}
	rawRender.WriteString("\n")
	comb.Lock()
	_, err = comb.SourceBuf.Write(rawRender.Bytes())
	if err != nil {
		comb.Unlock()
		errChan <- err
		wg.Done()
		return
	}
	comb.Unlock()
	wg.Done()
	return
}

// AddPurpleHairDyeToRoots creates the tangled hair arrays
func (s *Stylist) AddPurpleHairDyeToRoots() error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	for _, tw := range s.Combs {
		for _, str := range tw.Defs {
			wg.Add(1)
			go tw.EasyBakeOven(str, errChan, &wg)
		}
	}
	go func() {
		wg.Wait()
		close(finChan)
	}()
	select {
	case <-finChan:
	case err := <-errChan:
		if err != nil {
			return err
		}
	}
	return nil
}

func (comb *tangleWalker) Rinse() error {
	fr, err := os.Open(comb.Filename)
	if err != nil {
		return err
	}
	tempBuf := new(bytes.Buffer)
	_, err = io.Copy(tempBuf, fr)
	if err != nil {
		return err
	}
	err = fr.Close()
	if err != nil {
		return err
	}
	tempBuf.WriteString("\n")
	tempBuf.Write(comb.SourceBuf.Bytes())
	retOpts := imports.Options{
		Comments:  true,
		AllErrors: true,
		TabIndent: false,
		TabWidth:  2,
	}
	newData, err := imports.Process(comb.Filename, tempBuf.Bytes(), &retOpts)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(comb.Filename, newData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetTheQueenToHerThrown finishes it off
func (s *Stylist) GetTheQueenToHerThrown() error {
	fns := []func() error{}
	for _, tw := range s.Combs {
		fns = append(fns, tw.Rinse)
	}
	return ExecuteFuncsInParallel(fns)
}
