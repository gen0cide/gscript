package compiler

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/gen0cide/gscript/engine"
	gparser "github.com/robertkrimen/otto/parser"
	"github.com/uudashr/gopkgs"
)

var (
	errNoVM = errors.New("compiler has no VMs to process")
)

// Compiler is the primary type for building native binaries with gscript
type Compiler struct {
	// lock to prevent race conditions during compilation
	sync.RWMutex

	// array of VMs that will be bundled into this build
	vms []*GenesisVM

	// a map that places subsets of VMs into buckets according to their priority
	sortedVMs map[int][]*GenesisVM

	// logging object to be used
	logger engine.Logger

	// source buffer used by the pre-compilation obfuscator
	sourceBuffer bytes.Buffer

	// a slice of strings enumerated by the pre-compilation obfuscator
	stringDefs []*StringDef

	// a slice of unique priorities that can be found within this VMs bundled into this build
	uniqPriorities []int

	// configuration object for the compiler
	Options
}

// NewWithDefault returns a new compiler object with default options
func NewWithDefault() *Compiler {
	return &Compiler{
		logger:  &engine.NullLogger{},
		Options: DefaultOptions(),
	}
}

// NewWithOptions returns a new compiler object with custom options
func NewWithOptions(o Options) *Compiler {
	return &Compiler{
		logger:  &engine.NullLogger{},
		Options: o,
	}
}

// SetLogger overrides the logger for the compiler (defaults to an engine.NullLogger)
func (c *Compiler) SetLogger(l engine.Logger) {
	c.logger = l
}

// AddScript attempts to create a virtual machine object based on the given parameter to be included in compilation
func (c *Compiler) AddScript(scriptPath string) error {
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script cannot be located at %s", scriptPath)
	}
	srcBytes, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		return err
	}
	fileName := filepath.Base(scriptPath)
	absPath, err := filepath.Abs(scriptPath)
	if err != nil {
		return err
	}
	prog, err := gparser.ParseFile(nil, fileName, srcBytes, 2)
	if err != nil {
		return err
	}
	newVM := NewGenesisVM(fileName, absPath, c.OS, c.Arch, srcBytes, prog)
	newVM.Compiler = c
	c.vms = append(c.vms, newVM)
	return nil
}

// Do runs all compiler functions once scripts are added to it
func (c *Compiler) Do() error {
	err := c.CreateBuildDir()
	if err != nil {
		return err
	}
	err = c.ProcessMacros()
	if err != nil {
		return err
	}
	err = c.InitializeImports()
	if err != nil {
		return err
	}
	err = c.DetectVersions()
	if err != nil {
		return err
	}
	err = c.GatherAssets()
	if err != nil {
		return err
	}
	err = c.WalkGenesisASTs()
	if err != nil {
		return err
	}
	err = c.LocateGoDependencies()
	if err != nil {
		return err
	}
	err = c.BuildGolangASTs()
	if err != nil {
		return err
	}
	err = c.SwizzleNativeCalls()
	if err != nil {
		return err
	}
	err = c.SanityCheckSwizzles()
	if err != nil {
		return err
	}
	err = c.WriteScripts()
	if err != nil {
		return err
	}
	err = c.WritePreloads()
	if err != nil {
		return err
	}
	return nil
}

// CreateBuildDir creates the compiler's build directory, with an additional asset directory as well
func (c *Compiler) CreateBuildDir() error {
	err := os.MkdirAll(c.BuildDir, 0744)
	if err != nil {
		return fmt.Errorf("cannot create build directory: %v", err)
	}
	err = os.MkdirAll(c.AssetDir(), 0744)
	if err != nil {
		return fmt.Errorf("cannot create asset directory: %v", err)
	}
	return nil
}

// createBuildDir
// compileMacros
// writeScript
// compileAssets
// buildEntryPoint
// tumbleAST
// writeSource
// compileSource
// obfuscateBinary
// compressBinary
//
// CreateBuildDir
// ProcessMacros
// DetectVersions
// GatherAssets
// WalkGenesisASTs
// LocateGoDependencies
// BuildGolangASTs
// SwizzleNativeCalls
// SanityCheckSwizzles
// WriteScripts
// EncodeAssets
// GenerateIR
// WriteVMBundles
// CreateEntryPoint
// BuildNativeBinary

// ProcessMacros enumerates the compilers virtual machines with the pre-processor to extract
// compiler macros for each virtual machine
func (c *Compiler) ProcessMacros() error {
	if len(c.vms) == 0 {
		return errNoVM
	}
	var wg sync.WaitGroup
	for _, vm := range c.vms {
		wg.Add(1)
		go func(vm *GenesisVM) {
			vm.ProcessMacros()
			wg.Done()
		}(vm)
	}
	wg.Wait()
	return nil
}

// DetectVersions enumerates all VMs to determine the engine version based on the entrypoint.
// For more information on this, look at GenesisVM.DetectTargetEngineVersion()
func (c *Compiler) DetectVersions() error {
	fns := []func() error{}
	for _, vm := range c.vms {
		fns = append(fns, vm.DetectTargetEngineVersion)
	}
	return c.ExecuteVMActionInParallel(fns)
}

// createBuildDir
// compileMacros
// writeScript
// compileAssets
// buildEntryPoint
// tumbleAST
// writeSource
// compileSource
// obfuscateBinary
// compressBinary

// GatherAssets enumerates all bundled virtual machines for any embedded assets and copies them
// into the build directory's asset cache
func (c *Compiler) GatherAssets() error {
	fns := []func() error{}
	for _, vm := range c.vms {
		fns = append(fns, vm.CacheAssets)
	}
	return c.ExecuteVMActionInParallel(fns)
}

// WriteScripts enumerates the compiler's genesis VMs and writes a cached version of the
// genesis source to the asset directory to prevent race condiditons with script filesystem locations
func (c *Compiler) WriteScripts() error {
	fns := []func() error{}
	for _, vm := range c.vms {
		fns = append(fns, vm.WriteScript)
	}
	return c.ExecuteVMActionInParallel(fns)
}

// InitializeImports enumerates the compiler's genesis VMs and writes a cached version of the
// genesis source to the asset directory to prevent race condiditons with script filesystem locations
func (c *Compiler) InitializeImports() error {
	fns := []func() error{}
	for _, vm := range c.vms {
		fns = append(fns, vm.InitializeGoImports)
	}
	return c.ExecuteVMActionInParallel(fns)
}

// WalkGenesisASTs scans all genesis VMs scripts to identify Golang packages that have been
// called from inside the script using the namespace identifier from the compiler macros
func (c *Compiler) WalkGenesisASTs() error {
	fns := []func() error{}
	for _, vm := range c.vms {
		fns = append(fns, vm.WalkGenesisAST)
	}
	return c.ExecuteVMActionInParallel(fns)
}

// LocateGoDependencies gathers a list of all installed golang packages, hands a copy to each VM,
// then has every VM resolve it's own golang dependencies from that package list
func (c *Compiler) LocateGoDependencies() error {
	// grab a list of currently installed golang packages
	gopks, err := gopkgs.Packages(gopkgs.Options{NoVendor: true})
	if err != nil {
		return err
	}

	// enumerate the packages, identifying all VMs that use them
	for _, gopkg := range gopks {
		for _, vm := range c.vms {
			if gop, ok := vm.GoPackageByImport[gopkg.ImportPath]; ok {
				gop.Dir = gopkg.Dir
				gop.ImportPath = gopkg.ImportPath
				gop.Name = gopkg.Name
			}
		}
	}

	// now to check for any unmet dependencies
	packages := map[string]bool{}
	for _, vm := range c.vms {
		for _, p := range vm.UnresolvedGoPackages() {
			packages[p] = true
		}
	}

	// handle the error
	if len(packages) > 0 {
		c.logger.Errorf("a number of golang dependencies could not be resolved:")
		for k := range packages {
			c.logger.Errorf("\t%s", k)
		}
		return fmt.Errorf("unresolved golang packages discovered")
	}
	return nil
}

// BuildGolangASTs enumerates each genesis vm's golang native packages and matches exported
// function declarations to their genesis script caller. This creates a reference in the VM's
// linker object which will be used to generate native interfaces between the genesis VM and
// the underlying golang packages.
func (c *Compiler) BuildGolangASTs() error {
	fns := []func() error{}
	for _, vm := range c.vms {
		fns = append(fns, vm.BuildGolangAST)
	}
	return c.ExecuteVMActionInParallel(fns)
}

// SwizzleNativeCalls enumerates all native golang function calls mapped to genesis script
// function calls and generates the type declarations for both arguments and return values.
func (c *Compiler) SwizzleNativeCalls() error {
	fns := []func() error{}
	for _, vm := range c.vms {
		fns = append(fns, vm.SwizzleNativeFunctionCalls)
	}
	return c.ExecuteVMActionInParallel(fns)
}

// SanityCheckSwizzles enumerates all VMs to make sure their linked native functions
// are being called correctly by the corrasponding javascript callers
func (c *Compiler) SanityCheckSwizzles() error {
	fns := []func() error{}
	for _, vm := range c.vms {
		fns = append(fns, vm.SanityCheckLinkedSymbols)
	}
	return c.ExecuteVMActionInParallel(fns)
}

// createBuildDir
// compileMacros
// writeScript
// compileAssets
// buildEntryPoint
// tumbleAST
// writeSource
// compileSource
// obfuscateBinary
// compressBinary
//
// CreateBuildDir
// ProcessMacros
// DetectVersions
// GatherAssets
// WalkGenesisASTs
// LocateGoDependencies
// BuildGolangASTs
// SwizzleNativeCalls
// SanityCheckSwizzles
// WriteScripts
// WritePreloads
// EncodeAssets
// GenerateDylibs
// WriteVMBundles
// CreateEntryPoint
// BuildNativeBinary

// WritePreloads renders preload libraries for every virtual machine in the compilers asset directory
func (c *Compiler) WritePreloads() error {
	fns := []func() error{}
	for _, vm := range c.vms {
		fns = append(fns, vm.WritePreload)
	}
	return c.ExecuteVMActionInParallel(fns)
}

// EncodeAssets renders all embedded assets into intermediate representation
func (c *Compiler) EncodeAssets() error {
	return nil
}

// GenerateDylibs generates the dynamic links for each virtual machine's symbol table
func (c *Compiler) GenerateDylibs() error {
	return nil
}

// WriteVMBundles writes the intermediate representation for each virtual machine to it's
// vm bundle file within the build directory
func (c *Compiler) WriteVMBundles() error {
	return nil
}

// CreateEntryPoint renders the final main() entry point for the final binary in the build directory
func (c *Compiler) CreateEntryPoint() error {
	return nil
}

// BuildNativeBinary uses the golang compiler to attempt to build a native binary for
// the target platform specified in the compiler options
func (c *Compiler) BuildNativeBinary() error {
	os.Chdir(c.BuildDir)
	cmd := exec.Command("go", "build", `-v`, `-ldflags`, `-s -w`, "-o", c.OutputFile)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", c.OS))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", c.Arch))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// ExecuteVMActionInParallel is a meta function that takes an array of function pointers (hopefully for each VM)
// and executes them in parallel to decrease compile times. This is setup to handle errors within
// each VM gracefully and not allow a goroutine to fail silently.
func (c *Compiler) ExecuteVMActionInParallel(fns []func() error) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	for _, fn := range fns {
		wg.Add(1)
		go func(f func() error) {
			err := f()
			if err != nil {
				errChan <- err
			}
			wg.Done()
		}(fn)
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

// GetVMs is a test function
func (c *Compiler) GetVMs() []*GenesisVM {
	return c.vms
}
