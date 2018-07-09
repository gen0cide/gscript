package compiler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"sync"
	"text/template"

	"github.com/gen0cide/gscript/compiler/computil"
	gast "github.com/robertkrimen/otto/ast"
	gfile "github.com/robertkrimen/otto/file"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
	"golang.org/x/tools/imports"
)

var (
	defaultPriority = 100

	requiredBuildTemplates = []string{
		"init",
		"preload",
		"import_assets",
		"import_script",
		"import_native",
		"execute",
	}

	//CallablesByEngineVersion is a table that maps the entry points to the expected script versions
	CallablesByEngineVersion = map[int][]string{
		1: []string{
			"BeforeDeploy",
			"Deploy",
			"AfterDeploy",
		},
		2: []string{
			"Deploy",
		},
	}

	callablePointMap = map[string]string{
		"BeforeDeploy": "no",
		"Deploy":       "yes",
		"AfterDeploy":  "no",
	}
)

// GenesisVM is the object representation of a genesis script including it's imports and dynamic linking targets
type GenesisVM struct {
	// mutex for compiler manipulation
	sync.RWMutex

	// generated script ID
	ID string `json:"id"`

	// Absolute path to the script file
	SourceFile string `json:"source"`

	// map of embedded files
	Embeds map[string]*EmbeddedFile

	// name of the script (base name of file)
	Name string

	// raw script contents
	Data []byte

	// FileSet for parsing the genesis ASTs
	FileSet *gfile.FileSet

	// represents script as an AST
	AST *gast.Program

	// holds the value of parsed macros
	Macros []*Macro

	// maps the current build environment's golang imports
	// key = golang import path
	// value = reference to go package object
	GoPackageByImport map[string]*GoPackage

	// maps the namespace of genesis go imports to their
	// corrasponding go packages
	// key = genesis import namespace
	// value = reference to a go package object
	GoPackageByNamespace map[string]*GoPackage

	// required operating system for this script (GOOS)
	RequiredOS string

	// required architecture for this script (GOARCH)
	RequiredArch string

	// Object that holds the translation targets between golang and gscript
	Linker *Linker

	// maintains a map of the function names to the obfuscated references
	EntryPointMapping map[string]string

	// unique variable name to reference this scripts entry point
	PreloadAlias string

	// list of functions that need to be waterfalled for successful execution of the VM
	// usually either single element (Deploy) or the legacy compatible BeforeDeploy, Deploy, AfterDeploy
	EngineVersion int

	// reference to the parent compiler
	Compiler *Compiler

	// GenesisFile holds the intermediate representation of this VM's bundle code
	GenesisFile *bytes.Buffer

	// DecryptionKey is the key used to decrypt the embedded assets
	DecryptionKey string
}

// NewGenesisVM creates a new virtual machine object for the compiler
func NewGenesisVM(name, path, os, arch string, data []byte, prog *gast.Program) *GenesisVM {
	vm := &GenesisVM{
		ID:                   computil.RandUpperAlphaString(14),
		SourceFile:           path,
		Name:                 name,
		FileSet:              &gfile.FileSet{},
		Data:                 data,
		RequiredArch:         arch,
		RequiredOS:           os,
		AST:                  prog,
		Embeds:               map[string]*EmbeddedFile{},
		Macros:               []*Macro{},
		GoPackageByImport:    map[string]*GoPackage{},
		GoPackageByNamespace: map[string]*GoPackage{},
		EntryPointMapping:    map[string]string{},
		PreloadAlias:         computil.RandUpperAlphaString(12),
		DecryptionKey:        computil.RandMixedAlphaNumericString(32),
	}
	vm.Linker = NewLinker(vm)
	return vm
}

// scan for macros
// initialize go imports
// walk genesis AST for golang calls
// locate golang dependencies
// walk golang AST for func declarations
// link golang AST with genesis AST

// ProcessMacros runs the preprocessor to locate and extract genesis macro's
// out of the script to be used during compilation
func (g *GenesisVM) ProcessMacros() error {
	g.Macros = ScanForMacros(g.AST.Comments)
	return nil
}

// DetectTargetEngineVersion examines the genesis script's AST to determine whether required top level functions exist,
// and if so, for what version of the engine they target. This mapping can be found in CallablesByEngineVersion
func (g *GenesisVM) DetectTargetEngineVersion() error {
	cFuncs := map[string]bool{}
	for _, s := range g.AST.Body {
		fnStmt, ok := s.(*gast.FunctionStatement)
		if !ok {
			continue
		}
		fnLabel := fnStmt.Function.Name.Name
		if callablePointMap[fnLabel] != "" {
			cFuncs[fnLabel] = true
		}
	}
	cLen := len(cFuncs)
	if cLen == 3 {
		g.EngineVersion = 1
		return nil
	}
	if cLen == 1 && cFuncs["Deploy"] == true {
		g.EngineVersion = 2
		return nil
	}
	if cFuncs["Deploy"] != true {
		return fmt.Errorf("no Deploy() entry point detected in script %s", g.Name)
	}
	for _, x := range CallablesByEngineVersion[1] {
		if cFuncs[x] != true {
			return fmt.Errorf("no %s() entry point detected in script %s", x, g.Name)
		}
	}
	return fmt.Errorf("no entry point functions were found declared in the script %s", g.Name)
}

// CacheAssets indexes all //import: compiler macros and retrieves the corrasponding asset
func (g *GenesisVM) CacheAssets() error {
	importMacros := []*Macro{}
	for _, m := range g.Macros {
		if m.Key == "import" {
			importMacros = append(importMacros, m)
		}
	}
	for _, m := range importMacros {
		err := g.RetrieveAsset(m)
		if err != nil {
			return err
		}
	}
	return nil
}

// RetrieveAsset attempts to copy the asset into the build directory
func (g *GenesisVM) RetrieveAsset(m *Macro) error {
	ef, err := NewEmbeddedFile(m.Params["value"])
	if err != nil {
		return err
	}
	err = ef.CacheFile(g.Compiler.AssetDir())
	if err != nil {
		return err
	}
	g.Lock()
	g.Embeds[ef.OrigName] = ef
	g.Unlock()
	return nil
}

// DecryptionKeyArray returns the decryption key as an array
func (g *GenesisVM) DecryptionKeyArray() []byte {
	return []byte(g.DecryptionKey)
}

// EncodeBundledAssets encodes all assets within the asset pack into their compressed format
func (g *GenesisVM) EncodeBundledAssets() error {
	fns := []func() error{}
	for _, e := range g.Embeds {
		fns = append(fns, e.GenerateEmbedData)
	}
	return computil.ExecuteFuncsInParallel(fns)
}

// WriteGenesisScript writes a genesis script to the asset directory and returns a reference to an embeddedfile
// for use by the compiler
func (g *GenesisVM) WriteGenesisScript(name string, src []byte) (*EmbeddedFile, error) {
	scriptFileID := computil.RandUpperAlphaNumericString(18)
	scriptName := fmt.Sprintf("%s.gs", scriptFileID)
	scriptLocation := filepath.Join(g.Compiler.AssetDir(), scriptName)
	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)
	miniVersion := new(bytes.Buffer)
	r := bytes.NewReader(src)
	if err := m.Minify("text/javascript", miniVersion, r); err != nil {
		return nil, err
	}
	err := ioutil.WriteFile(scriptLocation, miniVersion.Bytes(), 0644)
	if err != nil {
		return nil, err
	}
	scriptEmbed := &EmbeddedFile{
		CachedPath:    scriptLocation,
		Filename:      scriptName,
		OrigName:      name,
		ID:            scriptFileID,
		EncryptionKey: []byte(g.DecryptionKey),
	}
	return scriptEmbed, nil
}

// WriteScript writes the initial user supplied script to the asset directory and tags
// it in the embed table as the entry point for the user defined functions
func (g *GenesisVM) WriteScript() error {
	scriptEmbed, err := g.WriteGenesisScript(g.Name, g.Data)
	if err != nil {
		return err
	}
	g.Lock()
	g.Embeds["__ENTRYPOINT"] = scriptEmbed
	g.Unlock()
	return nil
}

// WritePreload writes the preload library to the asset directory and tags
// it in the embed table as the preload library for the virtual machine
func (g *GenesisVM) WritePreload() error {
	scriptEmbed, err := g.WriteGenesisScript("preload.gs", []byte(Preload))
	if err != nil {
		return err
	}
	g.Lock()
	g.Embeds["__PRELOAD"] = scriptEmbed
	g.Unlock()
	return nil
}

// InitializeGoImports enumerates the go_import macros to initialize mappings
// for dynamic linking
func (g *GenesisVM) InitializeGoImports() error {
	for _, m := range g.Macros {
		if m.Key != "go_import" {
			continue
		}
		gop := &GoPackage{
			Script:         g,
			Namespace:      m.Params["namespace"],
			ImportKey:      m.Params["gopkg"],
			ScriptCallers:  map[string]*FunctionCall{},
			ImportsByFile:  map[string][]*ast.ImportSpec{},
			ImportsByAlias: map[string]*ast.ImportSpec{},
			FuncToFileMap:  map[string]string{},
			FuncTable:      map[string]*ast.FuncDecl{},
			LinkedFuncs:    []*LinkedFunction{},
		}
		g.GoPackageByImport[m.Params["gopkg"]] = gop
		g.GoPackageByNamespace[m.Params["namespace"]] = gop
	}
	return nil
}

// WalkGenesisAST walks the genesis script in order to inspect function calls
// that should be targeted for both legacy dynamic linking as well as native
// golang dynamic linking. Reference type genesisWalker and it's associated functions
// inside genesis_ast.go
func (g *GenesisVM) WalkGenesisAST() error {
	walker := &genesisWalker{
		script: g,
		source: string(g.Data),
	}
	gast.Walk(walker, g.AST)
	return nil
}

// UnresolvedGoPackages enumerates the import table to determine if any packages
// have not been resolved to local dependencies yet.
// Returns a string slice of go import paths that have yet to be resolved.
func (g *GenesisVM) UnresolvedGoPackages() []string {
	unresolved := []string{}
	for name, gpkg := range g.GoPackageByImport {
		if gpkg.Dir != "" {
			continue
		}
		unresolved = append(unresolved, name)
	}
	return unresolved
}

// BuildGolangAST walks the golang packages imported into the script to build a mapping
// of functions, the files they're in, imports to each file (for aliases), and locations
// in the genesis script where these are referenced
func (g *GenesisVM) BuildGolangAST() error {
	for _, gop := range g.GoPackageByImport {
		fs := token.NewFileSet()
		pkgs, err := parser.ParseDir(fs, gop.Dir, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		if _, ok := pkgs[gop.Name]; !ok {
			return fmt.Errorf("should have found golang package %s but didnt", gop.ImportPath)
		}
		for _, file := range pkgs[gop.Name].Files {
			var walkError error
			ast.Inspect(file, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if ok {
					funcName := fn.Name.Name
					// TODO: swizzle all exported functions so go functions can be
					// resolved at runtime (aka in the debugger)
					// if fn.Name.IsExported() {
					// 	gop.FuncTable[fn.Name.Name] = fn
					// }
					if fn.Name.IsExported() && gop.ScriptCallers[funcName] != nil {
						if len(gop.ImportsByFile[file.Name.Name]) == 0 {
							gop.ImportsByFile[file.Name.Name] = file.Imports
						}
						gop.FuncTable[funcName] = fn
						gop.FuncToFileMap[funcName] = file.Name.Name
						_, err := g.Linker.NewLinkedFunction(
							gop.ScriptCallers[funcName],
							file,
							fn,
							file.Imports,
							gop,
						)
						if err != nil {
							walkError = err
						}
					}
				}
				return true
			})
			if walkError != nil {
				return walkError
			}
		}
	}
	return nil
}

// SwizzleNativeFunctionCalls enumerates all LinkedFunctions held by the linker and generates
// structured mappings of both arguments (left swizzle) and returns (right swizzle) so the compiler
// can map the function's shim in the intermediate representation
func (g *GenesisVM) SwizzleNativeFunctionCalls() error {
	for fnName, lf := range g.Linker.Funcs {
		if lf.GoDecl.Recv != nil {
			return fmt.Errorf("golang function %s in package %s declares a method receiver which is unsupported by genesis at this time", fnName, lf.GoPackage.ImportPath)
		}
		err := lf.SwizzleToTheLeft()
		if err != nil {
			return err
		}
		err = lf.SwizzleToTheRight()
		if err != nil {
			return err
		}
	}
	return nil
}

// SanityCheckLinkedSymbols checks to make sure all linked functions do not violate
// caller conventions between the javascript and golang method signatures.
func (g *GenesisVM) SanityCheckLinkedSymbols() error {
	for _, lf := range g.Linker.Funcs {
		if len(lf.GoArgs) != len(lf.Caller.ArgumentList) {
			return fmt.Errorf("function call %s.%s in script %s does not match golang method signature (argument mismatch)", lf.Caller.Namespace, lf.Caller.FuncName, g.Name)
		}
	}
	return nil
}

// GenerateFunctionKeys creates random functions for the various parts of the VM's source file
func (g *GenesisVM) GenerateFunctionKeys() {
	for _, x := range requiredBuildTemplates {
		g.EntryPointMapping[x] = computil.RandUpperAlphaString(12)
	}
}

// FunctionKey is used by the intermediate representation generator to map specific functions
// in the virtual machine's constructors to unique identifiers in the IR
func (g *GenesisVM) FunctionKey(k string) string {
	return g.EntryPointMapping[k]
}

// RenderVMBundle generates the virtual machine's bundled intermediate representation file
func (g *GenesisVM) RenderVMBundle(templateFile string) error {
	g.GenerateFunctionKeys()
	tmpl := template.New(g.ID)
	tmpl.Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
	tmpl2, err := tmpl.Parse(templateFile)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = tmpl2.Execute(buf, g)
	if err != nil {
		return err
	}
	g.GenesisFile = buf
	return nil
}

// WriteVMBundle generates the VM bundle's intermediate representation using RenderVMBundle and then writes it
// to the compilers build directory
func (g *GenesisVM) WriteVMBundle() error {
	t, err := computil.Asset("vm_file.go.tmpl")
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s.go", g.ID)
	fileLocation := filepath.Join(g.Compiler.BuildDir, filename)
	err = g.RenderVMBundle(string(t))
	if err != nil {
		return err
	}
	retOpts := imports.Options{
		Comments:  true,
		AllErrors: true,
		TabIndent: false,
		TabWidth:  2,
	}
	newData, err := imports.Process(filename, g.GenesisFile.Bytes(), &retOpts)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileLocation, newData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetSimpleMacroValue returns a string of the VM's macro defined by key argument
func (g *GenesisVM) GetSimpleMacroValue(key string) string {
	for _, m := range g.Macros {
		if m.Key == key {
			return m.Params["value"]
		}
	}
	return ""
}

// Priority returns the priority value if defined in the macros, else returns default
func (g *GenesisVM) Priority() int {
	val := g.GetSimpleMacroValue("priority")
	if val == "" {
		return defaultPriority
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		g.Compiler.Logger.Errorf("problem parsing priority value: %v", err)
		return defaultPriority
	}
	return num
}

// HasDebuggingEnabled is an convienience method for checking to see if the debugger should be included
func (g *GenesisVM) HasDebuggingEnabled() bool {
	return g.Compiler.DebuggerEnabled
}
