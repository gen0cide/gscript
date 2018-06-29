package compiler

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"sync"

	gast "github.com/robertkrimen/otto/ast"
	gfile "github.com/robertkrimen/otto/file"
)

var (
	requiredBuildTemplates = []string{
		"init",
		"preload",
		"import_native_libs",
		"import_script_libs",
		"import_assets",
		"unpack_assets",
		"decrypt_assets",
		"decode_assets",
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

	// List of absolute file paths of embedded files
	AssetFiles []string `json:"imports"`

	// map of embedded files
	Embeds []*EmbeddedFile `json:"-"`

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
}

// NewGenesisVM creates a new virtual machine object for the compiler
func NewGenesisVM(name, path, os, arch string, data []byte, prog *gast.Program) *GenesisVM {
	vm := &GenesisVM{
		ID:                   RandLowerAlphaString(12),
		SourceFile:           path,
		Name:                 name,
		FileSet:              &gfile.FileSet{},
		Data:                 data,
		RequiredArch:         arch,
		RequiredOS:           os,
		AST:                  prog,
		AssetFiles:           []string{},
		Embeds:               []*EmbeddedFile{},
		Macros:               []*Macro{},
		GoPackageByImport:    map[string]*GoPackage{},
		GoPackageByNamespace: map[string]*GoPackage{},
		EntryPointMapping:    map[string]string{},
		PreloadAlias:         RandLowerAlphaString(12),
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

// WriteScript writes the VM's source to a cached location in the compiler's asset directory
func (g *GenesisVM) WriteScript() error {
	scriptLocation := filepath.Join(g.Compiler.AssetDir(), fmt.Sprintf("%s.gs", g.ID))
	return ioutil.WriteFile(scriptLocation, g.Data, 0644)
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

// GenerateFunctionKeys creates random functions for the various parts of the VM's source file
func (g *GenesisVM) GenerateFunctionKeys() {
	for _, x := range requiredBuildTemplates {
		g.EntryPointMapping[x] = RandUpperAlphaString(12)
	}
}

// FunctionKey is used by the intermediate representation generator to map specific functions
// in the virtual machine's constructors to unique identifiers in the IR
func (g *GenesisVM) FunctionKey(k string) string {
	return g.EntryPointMapping[k]
}
