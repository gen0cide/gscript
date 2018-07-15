package computil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	// FullObfuscation tells the compiler to perform all available obfuscation techniques
	FullObfuscation = iota

	// NoPostObfuscation tells the compiler to only perform post compilation obfuscation techniques
	NoPostObfuscation

	// NoPreObfuscation tells the compiler to only perform pre compilation obfuscation techniques
	NoPreObfuscation

	// NoObfuscation tells the compiler not to perform any available obfuscation techniques
	NoObfuscation
)

var (
	validTargetOperatingSystems = []string{
		"windows",
		"linux",
		"darwin",
	}

	validTargetArchitectures = []string{
		"amd64",
		"386",
	}
)

// Options defines a type to allow customization of a compiler at runtime
type Options struct {
	// Target operating system
	// DEFAULT: current operating system
	OS string `json:"os,omitempty"`

	// Target architecture
	// DEFAULT: current operating system's architecture
	Arch string `json:"arch,omitempty"`

	// Location of final binary output
	// DEFAULT: a file located in your OS's temp directory
	OutputFile string `json:"output_file,omitempty"`

	// Location of build directory
	// DEFAULT: a folder within your OS's temp directory
	BuildDir string `json:"build_dir,omitempty"`

	// Do not delete your build directory after successful compilation
	// DEFAULT: false
	SaveBuildDir bool `json:"save_build_dir,omitempty"`

	// Compress the final binary with UPX to reduce file size
	// DEFAULT: false
	UPXEnabled bool `json:"upx_enabled,omitempty"`

	// Inject a genesis logger into the engine to show console output (good for testing/bad for production)
	// DEFAULT: false
	LoggingEnabled bool `json:"logging_enabled,omitempty"`

	// Inject an interactive debugger into the executable so DebugConsole() can be called.
	// DEFAULT: false
	DebuggerEnabled bool `json:"debugger_enabled,omitempty"`

	// Do not actually compile into a native binary for the target OS and Arch - stop after generating intermediate representation.
	// DEFAULT: false
	SkipCompilation bool `json:"skip_compilation,omitempty"`

	// Do not obfuscate the IDs of the various packages and gscript VM bundle IDs
	// DEFAULT: false
	UseHumanReadableNames bool `json:"use_human_readable_names,omitempty"`

	// Import all native functions into the virtual machine from native go packages, not just ones used in the script
	// DEFAULT: false
	ImportAllNativeFuncs bool `json:"import_all_native_funcs,omitempty"`

	// Determines the compilers level of obfuscation performed on the final binary
	// DEFAULT: 0 (look at compiler const for available options)
	ObfuscationLevel int `json:"obfuscation_level,omitempty"`

	// Used to attach the test harness into the genesis VM
	// DEFAULT: false
	EnableTestBuild bool `json:"enable_test_build,omitempty"`

	// Used to describe the genesis dir on this machine
	// DEFAULT: Discovered through GOPATH
	GenesisDir string `json:"genesis_dir"`
}

// ValidOSList returns the list of valid target operating systems
func ValidOSList() []string {
	return validTargetOperatingSystems
}

// ValidArchList returns the list of valid architectures
func ValidArchList() []string {
	return validTargetArchitectures
}

// IsValidOS checks to see if the supplied string is a valid operating system target
func IsValidOS(s string) bool {
	for _, o := range validTargetOperatingSystems {
		if s == o {
			return true
		}
	}
	return false
}

// IsValidArch checks to see if the supplied string is a valid operating system target
func IsValidArch(s string) bool {
	for _, o := range validTargetArchitectures {
		if s == o {
			return true
		}
	}
	return false
}

// CheckForConfigErrors examines the options to determine if there is any conflicting settings
func (o Options) CheckForConfigErrors() error {
	if !IsValidOS(o.OS) {
		return fmt.Errorf("%s is not a valid operating system", o.OS)
	}
	if !IsValidArch(o.Arch) {
		return fmt.Errorf("%s is not a valid architecture", o.Arch)
	}
	// if o.ObfuscationLevel < 3 && (o.LoggingEnabled || o.DebuggerEnabled) {
	// 	return fmt.Errorf("cannot enable obfuscation at the same time as enabling logging or a debugger")
	// }
	if o.SkipCompilation && !o.SaveBuildDir {
		return fmt.Errorf("cannot skip compilation without saving the build directory")
	}
	return nil
}

// DefaultOptions returns an Options object with all default options pre-filled
func DefaultOptions() Options {
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH
	rootDir, err := ResolveGenesisPackageDir()
	if err != nil {
		panic(errors.New("could not find the genesis package in your go path"))
	}
	ext := "bin"
	if currentOS == "windows" {
		ext = "exe"
	}
	finalFile := filepath.Join(os.TempDir(), fmt.Sprintf("%d_gscript.%s", time.Now().Unix(), ext))
	dirName := RandMixedAlphaNumericString(16)
	buildDir := filepath.Join(os.TempDir(), dirName)
	return Options{
		OS:                    currentOS,
		Arch:                  currentArch,
		OutputFile:            finalFile,
		BuildDir:              buildDir,
		SaveBuildDir:          false,
		UPXEnabled:            false,
		LoggingEnabled:        false,
		DebuggerEnabled:       false,
		SkipCompilation:       false,
		UseHumanReadableNames: false,
		ImportAllNativeFuncs:  false,
		ObfuscationLevel:      FullObfuscation,
		GenesisDir:            rootDir,
	}
}

// AssetDir returns the file path to the asset build directory of the compiler
func (o Options) AssetDir() string {
	return filepath.Join(o.BuildDir, "assets")
}
