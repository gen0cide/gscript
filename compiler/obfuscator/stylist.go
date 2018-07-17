package obfuscator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"math/rand"
	"strings"
	"sync"
	"text/template"

	"github.com/gen0cide/gscript/compiler/computil"
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
	Combs map[*ast.File]*comb

	// BuildDir is the target build directory
	BuildDir string

	// Contains global fileset information for the package this comb belongs to
	FileSet *token.FileSet
}

type comb struct {
	sync.RWMutex

	// Reference back to the parent stylist
	Stylist *Stylist

	// AST of the file this comb is mangling
	AST *ast.File

	// File location of the file this comb is mangling
	Filename string

	// List of managled strings
	Defs []*StringDef

	// SourceBuf is what gets all the rendered string defs appended to it
	SourceBuf *bytes.Buffer

	// Contains global fileset information for the package this comb belongs to
	FileSet *token.FileSet
}

// NewStylist creates a new Stylist
func NewStylist(buildDir string) *Stylist {
	return &Stylist{
		Combs:    map[*ast.File]*comb{},
		BuildDir: buildDir,
	}
}

// LollerSkateDaStringz walks the build directory package AST
func (s *Stylist) LollerSkateDaStringz() error {
	s.FileSet = token.NewFileSet()
	pkgs, err := parser.ParseDir(s.FileSet, s.BuildDir, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	for filename, file := range pkgs["main"].Files {
		wg.Add(1)
		c := &comb{
			Filename:  filename,
			FileSet:   s.FileSet,
			AST:       file,
			Stylist:   s,
			Defs:      []*StringDef{},
			SourceBuf: new(bytes.Buffer),
		}
		s.Combs[file] = c
		go c.ShoutoutToThaHomies(errChan, &wg)
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

// ShoutoutToThaHomies walks each build directory source file's AST and creates a new comb for each
func (c *comb) ShoutoutToThaHomies(errChan chan error, wg *sync.WaitGroup) {
	ast.Walk(c, c.AST)
	err := printer.Fprint(c.SourceBuf, c.Stylist.FileSet, c.AST)
	c.SourceBuf.WriteString("\n\n")
	if err != nil {
		errChan <- err
		wg.Done()
		return
	}
	wg.Done()
	return
}

// BurnTheShitOuttaThisWeave creates a tangled hair between the AST and the stylist
func (c *comb) BurnTheShitOuttaThisWeave(key rune, source string) string {
	varName := computil.RandUpperAlphaNumericString(32)
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
	c.Lock()
	c.Defs = append(c.Defs, &StringDef{
		ID:    varName,
		Value: source,
		Key:   key,
		Data:  varDef,
	})
	c.Unlock()
	return cipher
}

// Visit walks the AST
func (c *comb) Visit(node ast.Node) ast.Visitor {
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
			n.Value = c.BurnTheShitOuttaThisWeave(rune(k), n.Value[1:len(n.Value)-1])
		}
	}
	return c
}

// EasyBakeOven actually renders the tangled hairs
func (c *comb) EasyBakeOven(sd *StringDef, errChan chan error, wg *sync.WaitGroup) {
	tmpl := template.New(fmt.Sprintf("obfstring_%s", sd.ID))
	tmpl.Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
	tmplData, err := computil.Asset("obfstring.go.tmpl")
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
	c.Lock()
	_, err = c.SourceBuf.Write(rawRender.Bytes())
	if err != nil {
		c.Unlock()
		errChan <- err
		wg.Done()
		return
	}
	c.Unlock()
	wg.Done()
	return
}

// AddPurpleHairDyeToRoots creates the tangled hair arrays
func (s *Stylist) AddPurpleHairDyeToRoots() error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	for _, c := range s.Combs {
		for _, str := range c.Defs {
			wg.Add(1)
			go c.EasyBakeOven(str, errChan, &wg)
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

func (c *comb) Rinse() error {
	retOpts := imports.Options{
		Comments:  true,
		AllErrors: true,
		TabIndent: false,
		TabWidth:  2,
	}
	newData, err := imports.Process(c.Filename, c.SourceBuf.Bytes(), &retOpts)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.Filename, newData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetTheQueenToHerThrown re-writes the intermediate representation with the embedded tangled hairs
func (s *Stylist) GetTheQueenToHerThrown() error {
	fns := []func() error{}
	for _, c := range s.Combs {
		fns = append(fns, c.Rinse)
	}
	err := computil.ExecuteFuncsInParallel(fns)
	if err != nil {
		return err
	}
	return nil
}

// GetIDLiterals returns all IDs used by all combs in this Stylist
func (s *Stylist) GetIDLiterals() []string {
	lits := []string{}
	for _, c := range s.Combs {
		for _, d := range c.Defs {
			_ = d
			lits = append(lits, d.ID)
		}
	}
	return lits
}
