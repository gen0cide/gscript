package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gen0cide/gscript/compiler/computil"
)

const (
	// FullObfuscation tells the compiler to perform all available obfuscation techniques
	FullObfuscation = iota

	// PostCompilationOnly tells the compiler to only perform post compilation obfuscation techniques
	PostCompilationOnly

	// PreCompilationOnly tells the compiler to only perform pre compilation obfuscation techniques
	PreCompilationOnly

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
	OS string

	// Target architecture
	// DEFAULT: current operating system's architecture
	Arch string

	// Location of final binary output
	// DEFAULT: a file located in your OS's temp directory
	OutputFile string

	// Location of build directory
	// DEFAULT: a folder within your OS's temp directory
	BuildDir string

	// Do not delete your build directory after successful compilation
	// DEFAULT: false
	SaveBuildDir bool

	// Compress the final binary with UPX to reduce file size
	// DEFAULT: false
	UPXEnabled bool

	// Inject a genesis logger into the engine to show console output (good for testing/bad for production)
	// DEFAULT: false
	LoggingEnabled bool

	// Inject an interactive debugger into the executable so DebugConsole() can be called.
	// DEFAULT: false
	DebuggerEnabled bool

	// Do not actually compile into a native binary for the target OS and Arch - stop after generating intermediate representation.
	// DEFAULT: false
	SkipCompilation bool

	// Do not obfuscate the IDs of the various packages and gscript VM bundle IDs
	// DEFAULT: false
	UseHumanReadableNames bool

	// Import all native functions into the virtual machine from native go packages, not just ones used in the script
	// DEFAULT: false
	ImportAllNativeFuncs bool

	// Determines the compilers level of obfuscation performed on the final binary
	// DEFAULT: 0 (look at compiler const for available options)
	ObfuscationLevel int
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
	if o.ObfuscationLevel < 3 && (o.LoggingEnabled || o.DebuggerEnabled) {
		return fmt.Errorf("cannot enable obfuscation at the same time as enabling logging or a debugger")
	}
	if o.SkipCompilation && !o.SaveBuildDir {
		return fmt.Errorf("cannot skip compilation without saving the build directory")
	}
	return nil
}

// DefaultOptions returns an Options object with all default options pre-filled
func DefaultOptions() Options {
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH
	ext := "bin"
	if currentOS == "windows" {
		ext = "exe"
	}
	finalFile := filepath.Join(os.TempDir(), fmt.Sprintf("%d_gscript.%s", time.Now().Unix(), ext))
	dirName := computil.RandMixedAlphaNumericString(16)
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
	}
}

// AssetDir returns the file path to the asset build directory of the compiler
func (o Options) AssetDir() string {
	return filepath.Join(o.BuildDir, "assets")
}
