package compiler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"
	"text/template"

	"github.com/gen0cide/gscript/engine"
	"github.com/gen0cide/gscript/logging"
	"github.com/sirupsen/logrus"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
)

// VMBundle defines a standalone GSE VM that will be bundled into a compiled binary
type VMBundle struct {
	sync.RWMutex
	ID           string          `json:"id"`
	ScriptFile   string          `json:"source"`
	AssetFiles   []string        `json:"imports"`
	Embeds       []*EmbeddedFile `json:"-"`
	RequiredOS   string          `json:"required_os"`
	RequiredArch string          `json:"required_arch"`
	Priority     int             `json:"priority"`
	Timeout      int             `json:"timeout"`
}

type VMLibrary struct {
	ID         string `json:"id"`
	ScriptFile string `json:"library"`
}

type StringDef struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Key   rune   `json:"key"`
	Data  []rune `json:"data"`
}

// Compiler creates a skeleton structure to produce a compiled binary
type Compiler struct {
	sync.RWMutex
	OS             string              `json:"os"`
	Arch           string              `json:"arch"`
	OutputFile     string              `json:"output"`
	VMs            []*VMBundle         `json:"vms"`
	SortedVMs      map[int][]*VMBundle `json:"-"`
	BuildDir       string              `json:"build_dir"`
	AssetDir       string              `json:"asset_dir"`
	OutputSource   bool                `json:"output_source"`
	CompressBinary bool                `json:"compress_binary"`
	EnableLogging  bool                `json:"enable_logging"`
	Logger         *logrus.Logger      `json:"-"`
	Source         string              `json:"-"`
	SourceBuffer   bytes.Buffer        `json:"-"`
	StringDefs     []*StringDef        `json:"-"`
	UniqPriorities []int               `json:"-"`
}

// NewCompiler returns a basic Compiler object
func NewCompiler(scripts []string, outfile, os, arch string, sourceOut, compression bool, enableLogging bool) *Compiler {
	logger := logrus.New()
	logger.Formatter = &logging.GSEStrippedFormatter{}
	logger.Out = logging.LogWriter{Name: "compiler"}
	logging.PrintLogo()
	if outfile == "-" && !sourceOut {
		logger.Fatalf("You need either -outfile or -source specified to build.")
	}
	if compression {
		if sourceOut {
			logger.Fatalf("You cannot use --source and --upx in the same compile command.")
		}
		_, err := exec.LookPath("upx")
		if err != nil {
			logger.Fatalf("The upx executable could not be found in your $PATH!")
		}
	}
	vms := []*VMBundle{}
	for _, s := range scripts {
		logger.WithField("src", s).Info("** Compiler Option **")
		vms = append(vms, &VMBundle{
			ID:           RandUpperAlphaString(18),
			ScriptFile:   s,
			AssetFiles:   []string{},
			Embeds:       []*EmbeddedFile{},
			RequiredArch: "",
			RequiredOS:   "",
			Priority:     100,
			Timeout:      30,
		})
	}

	logger.WithField("os", os).Info("** Compiler Option **")
	logger.WithField("arch", arch).Info("** Compiler Option **")
	if sourceOut {
		logger.WithField("source", fmt.Sprintf("%t", sourceOut)).Info("** Compiler Option **")
	} else {
		logger.WithField("outfile", outfile).Info("** Compiler Option **")
	}
	logger.WithField("upx", fmt.Sprintf("%t", compression)).Info("** Compiler Option **")
	logger.WithField("logging", fmt.Sprintf("%t", enableLogging)).Info("** Compiler Option **")

	return &Compiler{
		VMs:            vms,
		OutputFile:     outfile,
		Logger:         logger,
		OS:             os,
		Arch:           arch,
		OutputSource:   sourceOut,
		StringDefs:     []*StringDef{},
		CompressBinary: compression,
		EnableLogging:  enableLogging,
		SortedVMs:      make(map[int][]*VMBundle),
		UniqPriorities: []int{},
	}
}

// CreateBuildDir sets up the compiler's build directory
func (c *Compiler) createBuildDir() {
	dirName := engine.RandStringRunes(16)
	bd := filepath.Join(os.TempDir(), dirName)
	err := os.MkdirAll(bd, 0744)
	if err != nil {
		c.Logger.Fatalf("Cannot create build directory: dir=%s error=%s", bd, err.Error())
	}
	ad := filepath.Join(bd, "assets")
	os.MkdirAll(ad, 0744)
	c.BuildDir = bd
	c.AssetDir = ad
}

func (c *Compiler) compileMacros() {
	for _, vm := range c.VMs {
		assets := c.ParseMacros(vm)
		for _, asset := range assets {
			tempFile := filepath.Join(c.AssetDir, filepath.Base(asset))
			err := engine.LocalCopyFile(asset, tempFile)
			if err != nil {
				c.Logger.Fatalf("Asset file copy error: file=%s, error=%s", asset, err.Error())
			}
			c.Logger.Infof("Found Asset: %s", filepath.Base(asset))
			vm.Lock()
			vm.AssetFiles = append(vm.AssetFiles, tempFile)
			vm.Unlock()
		}
	}
}

func (c *Compiler) writeScript() {
	for _, vm := range c.VMs {
		if _, err := os.Stat(vm.ScriptFile); os.IsNotExist(err) {
			c.Logger.Fatalf("Genesis Script does not exist: %s", vm.ScriptFile)
		}
		entryFile := filepath.Join(c.AssetDir, fmt.Sprintf("%s.gs", vm.ID))
		m := minify.New()
		m.AddFunc("text/javascript", js.Minify)

		miniVersion := new(bytes.Buffer)

		data, err := engine.LocalFileRead(vm.ScriptFile)
		if err != nil {
			c.Logger.Fatalf("Asset file copy error: file=%s, error=%s", vm.ScriptFile, err.Error())
		}
		r := bytes.NewReader(data)

		if err := m.Minify("text/javascript", miniVersion, r); err != nil {
			c.Logger.Fatalf("Minification error: %s", err.Error())
		}

		miniFinal := miniVersion.Bytes()
		c.Logger.Debugf("Original Size: %d bytes", len(data))
		c.Logger.Debugf("Minified Size: %d bytes", len(miniFinal))
		engine.LocalFileCreate(entryFile, miniFinal)
		vm.AssetFiles = append(vm.AssetFiles, entryFile)
	}
}

func (c *Compiler) processAsset(vm *VMBundle, f string, wg *sync.WaitGroup) {
	defer wg.Done()
	e := &EmbeddedFile{
		SourcePath: f,
	}
	c.Logger.Debugf("Embedding file: %s", f)
	e.Embed()
	vm.Lock()
	vm.Embeds = append(vm.Embeds, e)
	vm.Unlock()
}

func (c *Compiler) compileAssets() {
	var wg sync.WaitGroup
	for _, vm := range c.VMs {
		for _, f := range vm.AssetFiles {
			wg.Add(1)
			go func(f string, vm *VMBundle) {
				c.processAsset(vm, f, &wg)
			}(f, vm)
		}
	}
	wg.Wait()
}

func RetrieveExample() []byte {
	return MustAsset("templates/example.gs")
}

func (c *Compiler) buildEntryPoint() {
	for _, vm := range c.VMs {
		if c.SortedVMs[vm.Priority] == nil {
			c.SortedVMs[vm.Priority] = []*VMBundle{}
			c.UniqPriorities = append(c.UniqPriorities, vm.Priority)
		}
		c.SortedVMs[vm.Priority] = append(c.SortedVMs[vm.Priority], vm)
	}
	sort.Ints(c.UniqPriorities)
	tmpl := template.New("gse_builder")
	tmpl.Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
	newTmpl, err := tmpl.Parse(string(MustAsset("templates/entrypoint.go.tmpl")))
	if err != nil {
		c.Logger.Fatalf("Error generating source: %s", err.Error())
	}
	var buf bytes.Buffer
	err = newTmpl.Execute(&buf, &c)
	if err != nil {
		c.Logger.Fatalf("Error generating source: %s", err.Error())
	}
	c.SourceBuffer = buf
}

func (c *Compiler) tumbleAST() {
	c.Logger.Info("Obfuscating strings")
	newSource := c.LollerSkateDaStringz(c.SourceBuffer.Bytes())
	c.Logger.Info("Generating runtime decryption keys")
	newSource.WriteString("\n\n")
	newSource.WriteString(c.GenerateTangledHairs())
	c.Logger.Info("Injecting embedded assets into source")
	tmpl := template.New("embeds")
	newTmpl, err := tmpl.Parse(string(MustAsset("templates/embed.go.tmpl")))
	if err != nil {
		c.Logger.Fatalf("Failed to parse embed template: %s", err.Error())
	}
	var buf bytes.Buffer
	err = newTmpl.Execute(&buf, &c)
	if err != nil {
		c.Logger.Fatalf("Failed to render embed template: %s", err.Error())
	}
	_, err = newSource.Write(buf.Bytes())
	if err != nil {
		c.Logger.Fatalf("Failed to append embeds: %s", err.Error())
	}
	c.Logger.Info("Writing final source")
	c.SourceBuffer.Reset()
	c.SourceBuffer.Write(newSource.Bytes())
}

func (c *Compiler) writeSource() {
	err := ioutil.WriteFile(filepath.Join(c.BuildDir, "main.go"), c.SourceBuffer.Bytes(), 0644)
	if err != nil {
		c.Logger.Fatalf("Error writing main.go: %s", err.Error())
	}
}

func (c *Compiler) printSource() {
	PrettyPrintSource(c.SourceBuffer.String())
}

func (c *Compiler) compileSource() {
	os.Chdir(c.BuildDir)
	cmd := exec.Command("go", "build", `-ldflags`, `-s -w`, "-o", c.OutputFile)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", c.OS))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", c.Arch))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		c.Logger.Fatalf("Compilation error for %s", filepath.Join(c.BuildDir, "main.go"))
	}
}

func (c *Compiler) obfuscateBinary() {
	if c.EnableLogging == true {
		c.Logger.Warnf("Not obfuscating binary because logging is enabled.")
		return
	}
	c.Logger.Infof("Obfuscating binary")
	c.ObfuscateBinary()
}

func (c *Compiler) compressBinary() {
	if c.CompressBinary == false {
		c.Logger.Warnf("Binary compression NOT enabled (default). Enable with --upx.")
		return
	}
	c.Logger.Info("Compressing binary with UPX")
	cmd := exec.Command("upx", `-9`, `-f`, `-q`, c.OutputFile)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", c.OS))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", c.Arch))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func (c *Compiler) Do() {
	cwd, _ := os.Getwd()
	c.Logger.Info("Creating build directory")
	c.createBuildDir()
	c.Logger.Info("Processing compiler macros")
	c.compileMacros()
	c.Logger.Info("Configuring build directory")
	c.writeScript()
	os.Chdir(c.BuildDir)
	c.Logger.Info("Normalizing assets")
	c.compileAssets()
	c.Logger.Info("Building entry point")
	c.buildEntryPoint()
	os.RemoveAll(c.AssetDir)
	c.Logger.Info("Randomizing AST nodes")
	c.tumbleAST()
	if c.OutputSource {
		c.printSource()
	} else {
		c.Logger.Debug("Writing final source")
		c.writeSource()
		c.Logger.Info("Compiling final binary")
		c.compileSource()
		c.obfuscateBinary()
		c.compressBinary()
	}
	os.Chdir(cwd)
	os.RemoveAll(c.BuildDir)
}
