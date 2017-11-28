package gscript

import (
	"bytes"
	"debug/elf"
	"debug/gosym"
	"debug/macho"
	"debug/pe"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"sort"

	"github.com/davecgh/go-spew/spew"
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
		c.Logger.Critf("Could not read binary file: %s", err.Error())
	}
	for _, r := range ObfuscatedBlobs {
		src := []byte(r)
		dst := make([]byte, hex.DecodedLen(len(src)))
		n, err := hex.Decode(dst, src)
		if err != nil {
			c.Logger.Critf("Could not decode obfuscation regex: %s", err.Error())
		}
		seeker := fmt.Sprintf("%s", dst[:n])

		re, err := regexp.Compile(seeker)
		if err != nil {
			c.Logger.Critf("Could not compile obfuscation regex: %s", err.Error())
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
		c.Logger.Critf("Could not write obfuscated binary: %s", err.Error())
	}
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
	spew.Dump(funcList)
	for x := range fileMap {
		rawFile = bytes.Replace(rawFile, []byte(x), CreateReplacement(x), -1)
	}
	err = ioutil.WriteFile(c.OutputFile, rawFile, 0755)
	if err != nil {
		panic(err)
	}
}
