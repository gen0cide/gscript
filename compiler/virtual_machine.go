package compiler

import (
	"sync"

	gast "github.com/robertkrimen/otto/ast"
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
}

// scan for macros
// initialize go imports
// walk genesis AST for golang calls
// locate golang dependencies
// walk golang AST for func declarations
// link golang AST with genesis AST

// ProcessMacros runs the preprocessor to locate and extract genesis macro's
// out of the script to be used during compilation
func (g *GenesisVM) ProcessMacros() {
	g.Macros = ScanForMacros(g.AST.Comments)
}

// InitializeGoImports enumerates the go_import macros to initialize mappings
// for dynamic linking
func (g *GenesisVM) InitializeGoImports() error {
	return nil
}

// WalkGenesisAST walks the genesis script in order to inspect function calls
// that should be targeted for both legacy dynamic linking as well as native
// golang dynamic linking
func (g *GenesisVM) WalkGenesisAST() error {
	return nil
}

// LocateGoPackages enumerates the local golang packages to map golang packages
// to referenced golang native packages imported into the genesis script
func (g *GenesisVM) LocateGoPackages() error {
	return nil
}

// BuildGolangAST walks the golang packages imported into the script to build a mapping
// of functions, the files they're in, imports to each file (for aliases), and locations
// in the genesis script where these are referenced
func (g *GenesisVM) BuildGolangAST() error {
	return nil
}

func (g *GenesisVM) GenerateFunctionKeys() {
	for _, x := range requiredBuildTemplates {
		g.EntryPointMapping[x] = RandUpperAlphaString(12)
	}
}

func (g *GenesisVM) FunctionKey(k string) string {
	return g.EntryPointMapping[k]
}
