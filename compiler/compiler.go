package compiler

import (
	"bytes"
	"fmt"
	goast "go/ast"
	goparser "go/parser"
	goprinter "go/printer"
	gotoken "go/token"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/gen0cide/gscript/engine"
	"github.com/gen0cide/gscript/logging"
	"github.com/sirupsen/logrus"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
)

// VMBundle defines a standalone GSE VM that will be bundled into a compiled binary
type VMBundle struct {
	ID           string         `json:"id"`
	ScriptFile   string         `json:"source"`
	AssetFiles   []string       `json:"imports"`
	Embeds       []EmbeddedFile `json:"-"`
	RequiredOS   string         `json:"required_os"`
	RequiredArch string         `json:"required_arch"`
	Priority     int            `json:"priority"`
	Timeout      int            `json:"timeout"`
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
	StringDefs     []*StringDef        `json:"-"`
	UniqPriorities []int               `json:"-"`
}

// NewCompiler returns a basic Compiler object
func NewCompiler(scripts []string, outfile, os, arch string, sourceOut, compression bool, enableLogging bool) *Compiler {
	logger := logrus.New()
	logger.Formatter = &logging.GSEFormatter{}
	logger.Out = logging.LogWriter{Name: "compiler"}
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
		vms = append(vms, &VMBundle{
			ID:           RandUpperAlphaString(18),
			ScriptFile:   s,
			AssetFiles:   []string{},
			Embeds:       []EmbeddedFile{},
			RequiredArch: "",
			RequiredOS:   "",
			Priority:     100,
			Timeout:      30,
		})
	}
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
func (c *Compiler) CreateBuildDir() {
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

// ParseMacros normalizes the import files into localized assets
func (c *Compiler) ParseMacros(vm *VMBundle) []string {
	imports := []string{}
	script, err := ioutil.ReadFile(vm.ScriptFile)
	if err != nil {
		c.Logger.WithField("file", filepath.Base(vm.ScriptFile)).Fatalf("Error reading genesis script: %s", err.Error())
	}

	macroList := ParseMacros(string(script), c.Logger.WithField("file", filepath.Base(vm.ScriptFile)))
	if macroList == nil {
		c.Logger.WithField("file", filepath.Base(vm.ScriptFile)).Fatalf("Could not parse macros for script!")
	}

	vm.Timeout = macroList.Timeout
	vm.Priority = macroList.Priority

	for _, i := range macroList.LocalFiles {
		imports = append(imports, i)
	}

	for _, i := range macroList.RemoteFiles {
		imports = append(imports, i)
	}

	return imports
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
			c.Logger.Infof("Packing File: %s", filepath.Base(asset))
			vm.AssetFiles = append(vm.AssetFiles, tempFile)
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
		c.Logger.Infof("Original Size: %d bytes", len(data))
		c.Logger.Infof("Minified Size: %d bytes", len(miniFinal))
		engine.LocalFileCreate(entryFile, miniFinal)
		vm.AssetFiles = append(vm.AssetFiles, entryFile)
	}
}

func (c *Compiler) compileAssets() {
	for _, vm := range c.VMs {
		for _, f := range vm.AssetFiles {
			e := EmbeddedFile{
				SourcePath: f,
			}
			e.Embed()
			vm.Embeds = append(vm.Embeds, e)
		}
	}
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
	c.Source = buf.String()
}

func (c *Compiler) GenerateTangledHairs() string {
	totalBuf := ""
	for _, str := range c.StringDefs {
		tmpl := template.New("obf_str")
		tmpl.Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
		newTmpl, err := tmpl.Parse(string(MustAsset("templates/obfstring.go.tmpl")))
		if err != nil {
			c.Logger.Fatalf("Error generating obfuscated string: %s", err.Error())
		}
		var buf bytes.Buffer
		err = newTmpl.Execute(&buf, str)
		if err != nil {
			c.Logger.Fatalf("Error generating obfuscated string: %s", err.Error())
		}
		totalBuf += buf.String()
		totalBuf += "\n\n"
	}
	return totalBuf
}

func (c *Compiler) writeSource() {
	newSource := c.LollerSkateDaStringz([]byte(c.Source))
	newSourceB := fmt.Sprintf("%s\n\n%s\n", string(newSource), c.GenerateTangledHairs())
	if c.OutputSource {
		PrettyPrintSource(newSourceB)
		return
	}
	err := ioutil.WriteFile(filepath.Join(c.BuildDir, "main.go"), []byte(newSourceB), 0644)
	if err != nil {
		c.Logger.Fatalf("Error writing main.go: %s", err.Error())
	}
}

func (c *Compiler) Do() {
	cwd, _ := os.Getwd()
	c.CreateBuildDir()
	c.compileMacros()
	c.writeScript()
	os.Chdir(c.BuildDir)
	c.compileAssets()
	c.buildEntryPoint()
	os.RemoveAll(c.AssetDir)
	c.writeSource()
	if !c.OutputSource {
		cmd := exec.Command("go", "build", `-ldflags`, `-s -w`, "-o", c.OutputFile)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", c.OS))
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", c.Arch))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		if c.EnableLogging == true {
			c.Logger.Warnf("Not obfuscating binary because logging is enabled.")
		} else {
			c.Logger.Infof("Obfuscating Binary...")
			c.ObfuscateBinary()
		}
		if c.CompressBinary {
			c.Logger.Infof("Compressing binary with UPX")
			cmd = exec.Command("upx", `-9`, `-f`, `-q`, c.OutputFile)
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", c.OS))
			cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", c.Arch))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
		// TODO: Fix GOTTI Binary Obfuscator #sadpanda
		// c.Logger.Infof("Mordor-ifying Binary...")
		// switch c.OS {
		// case "darwin":
		// 	c.MordorifyDarwin()
		// case "linux":
		// 	c.MordorifyLinux()
		// case "windows":
		// 	c.MordorifyWindows()
		// default:
		// 	c.Logger.Fatalf("GOOS specified is not supported.")
		// }
	}
	os.Chdir(cwd)
	os.RemoveAll(c.BuildDir)
}

func (c *Compiler) LollerSkateDaStringz(source []byte) []byte {
	fset := gotoken.NewFileSet()
	file, err := goparser.ParseFile(fset, "", source, 0)
	if err != nil {
		c.Logger.Fatalf("Could not parse Golang source: %s", err.Error())
	}
	goast.Walk(c, file)
	w := new(bytes.Buffer)
	goprinter.Fprint(w, fset, file)
	return w.Bytes()
}

func (c *Compiler) HairTangler(key rune, source string) string {
	varName := RandUpperAlphaString(14)
	cipher := fmt.Sprintf("g(%d, %s)", key, varName)
	reader := strings.NewReader(source)
	varDef := []rune{}
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		varDef = append(varDef, ch^key)
		key ^= ch
	}

	c.StringDefs = append(c.StringDefs, &StringDef{
		ID:    varName,
		Value: source,
		Key:   key,
		Data:  varDef,
	})
	return cipher
}

func (c *Compiler) Visit(node goast.Node) goast.Visitor {
	switch n := node.(type) {
	case *goast.ImportSpec:
		return nil
	case *goast.BasicLit:
		if n.Kind == gotoken.STRING {
			k := rand.Intn(65536)
			n.Value = c.HairTangler(rune(k), n.Value[1:len(n.Value)-1])
		}
	}
	return c
}
