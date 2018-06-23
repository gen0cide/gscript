package compiler

import (
	"bytes"
	"sync"

	"github.com/sirupsen/logrus"
)

// Compiler is the primary type for building native binaries with gscript
type Compiler struct {
	// lock to prevent race conditions during compilation
	sync.RWMutex

	// target operating system for the final binary
	OS string `json:"os"`

	// target architecture for the final binary
	Arch string `json:"arch"`

	// location of the compiled binary to be copied to when built
	OutputFile string `json:"output"`

	// array of VMs that will be bundled into this build
	VMs []*GenesisVM `json:"vms"`

	// a map that places subsets of VMs into buckets according to their priority
	SortedVMs map[int][]*GenesisVM `json:"-"`

	// location of the temporary build directory
	BuildDir string `json:"build_dir"`

	// location of the temporary asset directory
	AssetDir string `json:"asset_dir"`

	// compiler option to output a zip file containing source code instead of
	// actually building a final binary
	OutputSource bool `json:"output_source"`

	// compiler option to enable UPX compression post compilation
	CompressBinary bool `json:"compress_binary"`

	// compiler option to enable verbose logging output within the final binary
	// note: this option will DISABLE binary obfuscation.
	EnableLogging bool `json:"enable_logging"`

	// logging object to be used
	Logger *logrus.Logger `json:"-"`

	// source buffer used by the pre-compilation obfuscator
	SourceBuffer bytes.Buffer `json:"-"`

	// a slice of strings enumerated by the pre-compilation obfuscator
	StringDefs []*StringDef `json:"-"`

	// a slice of unique priorities that can be found within this VMs bundled into this build
	UniqPriorities []int `json:"-"`
}
