package compiler

import (
	"bytes"
	"debug/elf"
	"debug/gosym"
	"debug/macho"
	"debug/pe"
	"encoding/hex"
	"fmt"
	goast "go/ast"
	goparser "go/parser"
	goprinter "go/printer"
	gotoken "go/token"
	"html/template"
	"io/ioutil"
	"math/rand"
	"regexp"
	"sort"
	"strings"
)

var ObfuscatedBlobs = []string{
	"2f686f6d652f5b5b3a776f72643a5d5c2e5c5c2f205d2a",
	"2f55736572732f5b5b3a776f72643a5d5c2e5c5c2f205d2a",
	"2f726f6f742f5b5b3a776f72643a5d5c2e5c5c2f205d2a",
	"2f746d702f5b5b3a776f72643a5d5c2e5c5c2f205d2a",
	"2f7573722f6c6f63616c2f5b5b3a776f72643a5d5c2e5c5c2f205d2a",
	"6769746875625b5b3a776f72643a5d5c2e5c5c2f205d2a",
	"676f6f676c655b5b3a776f72643a5d5c2e5c5c2f205d2a",
	"676f6c616e675b5b3a776f72643a5d5c2e5c5c2f205d2a",
	"676f706b675b5b3a776f72643a5d5c2e5c5c2f205d2a",
	"5550585b5b3a776f72643a5d5c2e5c5c2f205d2a",
	"24496e666f5b5b3a776f72643a5d5c2e5c5c2f205d2a",
	"67656e30636964655b5b3a776f72643a5d5c2e5c5c2f205d2a",
	"677363726970745b5b3a776f72643a5d5c2e5c5c2f205d2a",
}

var (
	FunctionMap = map[string]string{}
)

func CreateReplacement(s string) []byte {
	b := make([]byte, len(s))
	return b
}

type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func (c *Compiler) ObfuscateBinary() {
	data, err := ioutil.ReadFile(c.OutputFile)
	if err != nil {
		c.Logger.Fatalf("Could not read binary file: %s", err.Error())
	}
	for _, r := range ObfuscatedBlobs {
		src := []byte(r)
		dst := make([]byte, hex.DecodedLen(len(src)))
		n, err := hex.Decode(dst, src)
		if err != nil {
			c.Logger.Fatalf("Could not decode obfuscation regex: %s", err.Error())
		}
		seeker := fmt.Sprintf("%s", dst[:n])

		re, err := regexp.Compile(seeker)
		if err != nil {
			c.Logger.Fatalf("Could not compile obfuscation regex: %s", err.Error())
		}

		data = re.ReplaceAllFunc(data, func(b []byte) []byte {
			for i := range b {
				b[i] = byte(rand.Int() % 256)
			}
			return b
		})
	}
	err = ioutil.WriteFile(c.OutputFile, data, 0755)
	if err != nil {
		c.Logger.Fatalf("Could not write obfuscated binary: %s", err.Error())
	}
}

func (c *Compiler) GenerateTangledHairs() string {
	totalBuf := ""
	for _, str := range c.StringDefs {
		tmpl := template.New("obf_str")
		tmpl.Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
		newTmpl, err := tmpl.Parse(string(MustAsset("templates/obfstring.go.tmpl")))
		if err != nil {
			c.Logger.Fatalf("Error generating obfuscated string: %s", err.Error())
		}
		var buf bytes.Buffer
		err = newTmpl.Execute(&buf, str)
		if err != nil {
			c.Logger.Fatalf("Error generating obfuscated string: %s", err.Error())
		}
		totalBuf += buf.String()
		totalBuf += "\n\n"
	}
	return totalBuf
}

func (c *Compiler) MordorifyWindows() {
	rawFuncs := []string{}
	bin, err := pe.Open(c.OutputFile)
	if err != nil {
		panic(err)
	}
	rawFile, err := ioutil.ReadFile(c.OutputFile)
	if err != nil {
		panic(err)
	}
	syms := bin.Symbols
	for _, x := range syms {
		rawFuncs = append(rawFuncs, x.Name)
	}
	sort.Sort(ByLength(rawFuncs))
	for _, x := range rawFuncs {
		rawFile = bytes.Replace(rawFile, []byte(x), CreateReplacement(x), -1)
	}
	pSec := bin.Section("__gopclntab")
	pSecRaw := rawFile[pSec.SectionHeader.Offset:(int(pSec.SectionHeader.Offset) + int(pSec.SectionHeader.Size))]
	tSec := bin.Section("__text")
	sSec := bin.Section("__gosymtab")
	sSecRaw := rawFile[sSec.SectionHeader.Offset:(int(sSec.SectionHeader.Offset) + int(sSec.SectionHeader.Size))]
	pcln := gosym.NewLineTable(pSecRaw, uint64(tSec.SectionHeader.VirtualAddress))
	tab, err := gosym.NewTable(sSecRaw, pcln)
	if err != nil {
		panic(err)
	}
	funcList := []string{}
	_ = funcList
	fileMap := tab.Files
	for x := range fileMap {
		rawFile = bytes.Replace(rawFile, []byte(x), CreateReplacement(x), -1)
	}
	err = ioutil.WriteFile(c.OutputFile, rawFile, 0755)
	if err != nil {
		panic(err)
	}
}

func (c *Compiler) MordorifyLinux() {
	rawFuncs := []string{}
	bin, err := elf.Open(c.OutputFile)
	if err != nil {
		panic(err)
	}
	rawFile, err := ioutil.ReadFile(c.OutputFile)
	if err != nil {
		panic(err)
	}
	syms, _ := bin.Symbols()
	for _, x := range syms {
		rawFuncs = append(rawFuncs, x.Name)
	}
	sort.Sort(ByLength(rawFuncs))
	for _, x := range rawFuncs {
		rawFile = bytes.Replace(rawFile, []byte(x), CreateReplacement(x), -1)
	}
	pSec := bin.Section("__gopclntab")
	pSecRaw := rawFile[pSec.SectionHeader.Offset:(int(pSec.SectionHeader.Offset) + int(pSec.SectionHeader.Size))]
	tSec := bin.Section("__text")
	sSec := bin.Section("__gosymtab")
	sSecRaw := rawFile[sSec.SectionHeader.Offset:(int(sSec.SectionHeader.Offset) + int(sSec.SectionHeader.Size))]
	pcln := gosym.NewLineTable(pSecRaw, tSec.SectionHeader.Addr)
	tab, err := gosym.NewTable(sSecRaw, pcln)
	if err != nil {
		panic(err)
	}
	funcList := []string{}
	_ = funcList
	fileMap := tab.Files
	for x := range fileMap {
		rawFile = bytes.Replace(rawFile, []byte(x), CreateReplacement(x), -1)
	}
	err = ioutil.WriteFile(c.OutputFile, rawFile, 0755)
	if err != nil {
		panic(err)
	}
}

func (c *Compiler) MordorifyDarwin() {
	rawFuncs := []string{}
	bin, err := macho.Open(c.OutputFile)
	if err != nil {
		panic(err)
	}
	rawFile, err := ioutil.ReadFile(c.OutputFile)
	if err != nil {
		panic(err)
	}
	syms := bin.Symtab
	for _, x := range syms.Syms {
		rawFuncs = append(rawFuncs, x.Name)
	}
	sort.Sort(ByLength(rawFuncs))
	for _, x := range rawFuncs {
		rawFile = bytes.Replace(rawFile, []byte(x), CreateReplacement(x), -1)
	}
	pSec := bin.Section("__gopclntab")
	pSecRaw := rawFile[pSec.SectionHeader.Offset:(int(pSec.SectionHeader.Offset) + int(pSec.SectionHeader.Size))]
	tSec := bin.Section("__text")
	sSec := bin.Section("__gosymtab")
	sSecRaw := rawFile[sSec.SectionHeader.Offset:(int(sSec.SectionHeader.Offset) + int(sSec.SectionHeader.Size))]
	pcln := gosym.NewLineTable(pSecRaw, tSec.SectionHeader.Addr)
	tab, err := gosym.NewTable(sSecRaw, pcln)
	if err != nil {
		panic(err)
	}
	funcList := []string{}
	_ = funcList
	fileMap := tab.Files
	for x := range fileMap {
		rawFile = bytes.Replace(rawFile, []byte(x), CreateReplacement(x), -1)
	}
	err = ioutil.WriteFile(c.OutputFile, rawFile, 0755)
	if err != nil {
		panic(err)
	}
}

func (c *Compiler) LollerSkateDaStringz(source []byte) *bytes.Buffer {
	c.Logger.Debug("Initializing token parser")
	fset := gotoken.NewFileSet()
	c.Logger.Debug("Ingesting source into token parser")
	file, err := goparser.ParseFile(fset, "", source, 0)
	if err != nil {
		c.Logger.Fatalf("Could not parse Golang source: %s", err.Error())
	}
	c.Logger.Debug("Walking AST")
	goast.Walk(c, file)
	w := new(bytes.Buffer)
	c.Logger.Debug("Writing to buffer")
	goprinter.Fprint(w, fset, file)
	return w
}

func (c *Compiler) HairTangler(key rune, source string) string {
	varName := RandUpperAlphaString(14)
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

	c.StringDefs = append(c.StringDefs, &StringDef{
		ID:    varName,
		Value: source,
		Key:   key,
		Data:  varDef,
	})
	return cipher
}

func (c *Compiler) Visit(node goast.Node) goast.Visitor {
	switch n := node.(type) {
	case *goast.ImportSpec:
		return nil
	case *goast.BasicLit:
		if n.Kind == gotoken.STRING {
			k := rand.Intn(65536)
			n.Value = c.HairTangler(rune(k), n.Value[1:len(n.Value)-1])
		}
	}
	return c
}
