package compiler

import (
	"bytes"
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"text/template"

	"github.com/fatih/color"
	"github.com/gen0cide/gscript/compiler/computil"
	"github.com/gen0cide/gscript/logger"
	gast "github.com/robertkrimen/otto/ast"
	gfile "github.com/robertkrimen/otto/file"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
	"golang.org/x/tools/imports"
)

var (
	defaultPriority = 100
	defaultTimeout  = 30

	requiredBuildTemplates = []string{
		"init",
		"preload",
		"import_assets",
		"import_standard_library",
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

	// represents script as an GenesisAST
	GenesisAST *gast.Program

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

	// StandardLibraries holds all references to used genesis standard library
	// packages that will be included in the build
	EnabledStandardLibs map[string]*GoPackage

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

	// reference to compiler options
	computil.Options

	// GenesisFile holds the intermediate representation of this VM's bundle code
	GenesisFile *bytes.Buffer

	// Logger to publish output from
	Logger logger.Logger

	// EnginePackage is the go package this engine should use
	EnginePackage *GoPackage

	// StandardLibs holds references to all possible standard libs
	StandardLibs map[string]*GoPackage
}

// NewGenesisVM creates a new virtual machine object for the compiler
func NewGenesisVM(name, path string, data []byte, prog *gast.Program, opts computil.Options, logger logger.Logger) *GenesisVM {
	vm := &GenesisVM{
		ID:                   computil.RandUpperAlphaString(14),
		SourceFile:           path,
		Name:                 name,
		FileSet:              &gfile.FileSet{},
		Data:                 data,
		GenesisAST:           prog,
		Options:              opts,
		Logger:               logger,
		Embeds:               map[string]*EmbeddedFile{},
		Macros:               []*Macro{},
		GoPackageByImport:    map[string]*GoPackage{},
		GoPackageByNamespace: map[string]*GoPackage{},
		EntryPointMapping:    map[string]string{},
		PreloadAlias:         computil.RandUpperAlphaString(12),
		EnabledStandardLibs:  map[string]*GoPackage{},
		StandardLibs:         map[string]*GoPackage{},
	}
	vm.Linker = NewLinker(vm)
	return vm
}

// scan for macros
// initialize go imports
// walk genesis GenesisAST for golang calls
// locate golang dependencies
// walk golang GenesisAST for func declarations
// link golang GenesisAST with genesis GenesisAST

// ProcessMacros runs the preprocessor to locate and extract genesis macro's
// out of the script to be used during compilation
func (g *GenesisVM) ProcessMacros() error {
	g.Macros = ScanForMacros(g.GenesisAST.Comments)
	return nil
}

// DetectTargetEngineVersion examines the genesis script's GenesisAST to determine whether required top level functions exist,
// and if so, for what version of the engine they target. This mapping can be found in CallablesByEngineVersion
func (g *GenesisVM) DetectTargetEngineVersion() error {
	cFuncs := map[string]bool{}
	for _, s := range g.GenesisAST.Body {
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

// GetTimeout attempts to get the timeout value set in the macro (if it was set)
func (g *GenesisVM) GetTimeout() int {
	ret := defaultTimeout
	for _, m := range g.Macros {
		if m.Key == "timeout" {
			ret, err := strconv.Atoi(m.Params["value"])
			if err != nil {
				panic(fmt.Errorf("script %s has an invalid timeout set: %s", g.Name, m.Params["value"]))
			}
			return ret
		}
	}
	return ret
}

// GetNewDecryptionKey creates a new decryption key
func (g *GenesisVM) GetNewDecryptionKey() string {
	return computil.RandMixedAlphaNumericString(32)
}

// RetrieveAsset attempts to copy the asset into the build directory
func (g *GenesisVM) RetrieveAsset(m *Macro) error {
	ef, err := NewEmbeddedFile(m.Params["value"], []byte(g.GetNewDecryptionKey()))
	if err != nil {
		return err
	}
	err = ef.CacheFile(g.Options.AssetDir())
	if err != nil {
		return err
	}
	g.Lock()
	g.Embeds[ef.OrigName] = ef
	g.Unlock()
	g.Logger.Debugf("  %s -> %s", color.HiWhiteString(g.Name), color.YellowString(ef.OrigName))
	return nil
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
	scriptLocation := filepath.Join(g.AssetDir(), scriptName)
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
		EncryptionKey: []byte(g.GetNewDecryptionKey()),
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
	scriptEmbed, err := g.WriteGenesisScript("preload.gs", computil.MustAsset("preload.gs"))
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
		gop := NewGoPackage(g, m.Params["namespace"], m.Params["gopkg"], false)
		g.GoPackageByImport[m.Params["gopkg"]] = gop
		g.GoPackageByNamespace[m.Params["namespace"]] = gop
		if !IsDefaultImport(m.Params["gopkg"]) {
			g.Linker.MaskedImports[m.Params["gopkg"]] = &MaskedImport{
				ImportPath: m.Params["gopkg"],
				OldAlias:   gop.Name,
				NewAlias:   gop.MaskedName,
			}
		}
	}
	for l := range computil.GenesisLibs {
		pkg, err := computil.ResolveStandardLibraryDir(l)
		if err != nil {
			return err
		}
		gop := NewGoPackage(g, pkg.Name, pkg.ImportPath, true)
		gop.Dir = pkg.Dir
		gop.Name = pkg.Name
		g.StandardLibs[l] = gop
		if g.DebuggerEnabled {
			g.EnableStandardLibrary(l)
		}
	}
	return nil
}

// LocateGoPackages enumerates the installed go packages on the current system and appends
// directory and namespace information to golang packages being declared by this script
func (g *GenesisVM) LocateGoPackages() error {
	for _, gpkg := range computil.InstalledGoPackages {
		if gop, ok := g.GoPackageByImport[gpkg.ImportPath]; ok {
			gop.Dir = gpkg.Dir
			gop.ImportPath = gpkg.ImportPath
			gop.Name = gpkg.Name
		}
	}
	return nil
}

// WalkGenesisAST walks the genesis script in order to inspect function calls
// that should be targeted for both legacy dynamic linking as well as native
// golang dynamic linking. Reference type genesisWalker and it's associated functions
// inside genesis_ast.go
func (g *GenesisVM) WalkGenesisAST() error {
	walker := &genesisWalker{
		vm:     g,
		source: string(g.Data),
	}
	gast.Walk(walker, g.GenesisAST)
	return walker.err
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

// WalkGoPackageAST parses the GoPackage directory for all AST files and concurrently walks each child file's AST
// looking for functions that should be included by the linker
func (g *GenesisVM) WalkGoPackageAST(gop *GoPackage, wg *sync.WaitGroup, errChan chan error) {
	ctxt := build.Default
	ctxt.GOOS = g.OS
	ctxt.GOARCH = g.Arch
	pkg, err := ctxt.Import(gop.ImportKey, gop.Dir, build.ImportComment)
	if err != nil {
		errChan <- err
		wg.Done()
		return
	}
	// NOTE: maybe we want to include pkg.CgoFiles?
	validSrcFiles := map[string]bool{}

	for _, f := range pkg.GoFiles {
		validSrcFiles[f] = true
	}

	pkgFilter := func(fi os.FileInfo) bool {
		return validSrcFiles[fi.Name()]
	}

	gop.FileSet = token.NewFileSet()
	pkgs, err := parser.ParseDir(gop.FileSet, gop.Dir, pkgFilter, parser.ParseComments)
	if err != nil {
		errChan <- err
		return
	}
	if _, ok := pkgs[gop.Name]; !ok {
		errChan <- fmt.Errorf("should have found golang package %s but didnt", gop.ImportPath)
		return
	}
	var filewg sync.WaitGroup
	fileErrChan := make(chan error, 1)
	fileFinChan := make(chan bool, 1)

	for filename, file := range pkgs[gop.Name].Files {
		if computil.SourceFileIsTest(filename) {
			continue
		}
		filewg.Add(1)
		go gop.WalkGoFileAST(file, &filewg, fileErrChan)
	}

	go func() {
		filewg.Wait()
		close(fileFinChan)
	}()

	select {
	case <-fileFinChan:
		wg.Done()
		return
	case err := <-fileErrChan:
		errChan <- err
		wg.Done()
		return
	}
}

// SanityCheckNativeFunctionCalls enumerates the VMs go packages ensuring that there are no script
// callers who do not exist within the GoPackage symbol table
func (g *GenesisVM) SanityCheckNativeFunctionCalls() error {
	for _, gop := range g.EnabledStandardLibs {
		err := gop.SanityCheckScriptCallers()
		if err != nil {
			return err
		}
	}
	return nil
}

// BuildGolangAST walks the golang packages imported into the script to build a mapping
// of functions, the files they're in, imports to each file (for aliases), and locations
// in the genesis script where these are referenced
func (g *GenesisVM) BuildGolangAST() error {
	var wg sync.WaitGroup
	numOfPackages := len(g.GoPackageByImport) + len(g.EnabledStandardLibs)
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)

	wg.Add(numOfPackages)

	for _, gop := range g.GoPackageByImport {
		go g.WalkGoPackageAST(gop, &wg, errChan)
	}
	for _, gop := range g.EnabledStandardLibs {
		go g.WalkGoPackageAST(gop, &wg, errChan)
	}

	go func() {
		wg.Wait()
		close(finChan)
	}()

	select {
	case <-finChan:
		return nil
	case err := <-errChan:
		return err
	}
}

// SwizzleNativeFunctionCalls enumerates all LinkedFunctions held by the linker and generates
// structured mappings of both arguments (left swizzle) and returns (right swizzle) so the compiler
// can map the function's shim in the intermediate representation
func (g *GenesisVM) SwizzleNativeFunctionCalls() error {
	for id, lf := range g.Linker.Funcs {
		if lf.GoDecl.Recv != nil {
			return fmt.Errorf("golang function %s in package %s declares a method receiver which is unsupported by genesis at this time", lf.Function, lf.GoPackage.ImportPath)
		}
		err := lf.SwizzleToTheLeft()
		if err != nil {
			lf.SwizzleError = err
			lf.SwizzleSuccessful = false
			g.Logger.Debugf("Could not swizzle native function %s: %v", id, err)
			if lf.Caller != nil {
				return fmt.Errorf("script %s calls %s which is not linkable", g.Name, id)
			}
			continue
		}
		err = lf.SwizzleToTheRight()
		if err != nil {
			lf.SwizzleError = err
			lf.SwizzleSuccessful = false
			g.Logger.Debugf("Could not swizzle native function %s: %v", id, err)
			if lf.Caller != nil {
				return fmt.Errorf("script %s calls %s which is not linkable", g.Name, id)
			}
			continue
		}
		lf.SwizzleSuccessful = true
	}
	return nil
}

// SanityCheckLinkedSymbols checks to make sure all linked functions do not violate
// caller conventions between the javascript and golang method signatures.
func (g *GenesisVM) SanityCheckLinkedSymbols() error {
	for _, lf := range g.Linker.Funcs {
		if lf.Caller == nil {
			continue
		}
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
	fileLocation := filepath.Join(g.BuildDir, filename)
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
		ioutil.WriteFile(fileLocation, g.GenesisFile.Bytes(), 0644)
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
		g.Logger.Errorf("problem parsing priority value: %v", err)
		return defaultPriority
	}
	return num
}

// HasDebuggingEnabled is an convienience method for checking to see if the debugger should be included
func (g *GenesisVM) HasDebuggingEnabled() bool {
	return g.DebuggerEnabled
}

// HasLoggingEnabled is an convienience method for checking to see if logging should be included
func (g *GenesisVM) HasLoggingEnabled() bool {
	return g.LoggingEnabled
}

// GetIDLiterals returns all interesting IDs used by this GenesisVM
func (g *GenesisVM) GetIDLiterals() []string {
	lits := []string{g.Name, g.ID, g.PreloadAlias}
	for _, v := range g.EntryPointMapping {
		lits = append(lits, []string{v}...)
	}
	for k, e := range g.Embeds {
		lits = append(lits, []string{k, e.ID, e.OrigName}...)
	}
	for k, gop := range g.GoPackageByImport {
		lits = append(lits, []string{k, gop.ImportPath, gop.Dir}...)
	}
	for _, lf := range g.Linker.Funcs {
		lits = append(lits, []string{lf.ID}...)
	}

	return lits
}

// EnableStandardLibrary attempts to resolve a discovered standard library and returns either the package or an error
func (g *GenesisVM) EnableStandardLibrary(name string) (*GoPackage, error) {
	if pkg, ok := g.StandardLibs[name]; ok {
		g.EnabledStandardLibs[name] = pkg
		if !IsDefaultImport(pkg.ImportPath) {
			g.Linker.MaskedImports[pkg.ImportPath] = &MaskedImport{
				ImportPath: pkg.ImportPath,
				OldAlias:   pkg.Name,
				NewAlias:   pkg.MaskedName,
			}
		}
		return pkg, nil
	}
	return nil, fmt.Errorf("invalid standard library detected: %s", name)
}

// ShouldIncludeAssetPackage is a helper function for VM bundle rendering to asset whether it needs to create
// asset functions in the intermediate representation
func (g *GenesisVM) ShouldIncludeAssetPackage() bool {
	if g.EnabledStandardLibs["asset"] == nil {
		return false
	}
	return true
}

// GetMaskedImports is a helper function to gather the masked imports during rendering of the vm bundle
func (g *GenesisVM) GetMaskedImports() []*MaskedImport {
	mi := []*MaskedImport{}
	for _, mi2 := range g.Linker.MaskedImports {
		_ = mi2
		mi = append(mi, mi2)
	}
	return mi
}
