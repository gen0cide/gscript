package gscript

import (
	"bytes"
	"fmt"
	goast "go/ast"
	goparser "go/parser"
	goprinter "go/printer"
	gotoken "go/token"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/alecthomas/chroma/quick"
	"github.com/happierall/l"
)

// VMBundle defines a standalone GSE VM that will be bundled into a compiled binary
type VMBundle struct {
	ID         string         `json:"id"`
	ScriptFile string         `json:"source"`
	AssetFiles []string       `json:"imports"`
	Embeds     []EmbeddedFile `json:"-"`
}

type StringDef struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Key   rune   `json:"key"`
	Data  []rune `json:"data"`
}

// Compiler creates a skeleton structure to produce a compiled binary
type Compiler struct {
	OS           string       `json:"os"`
	Arch         string       `json:"arch"`
	OutputFile   string       `json:"output"`
	VMs          []*VMBundle  `json:"vms"`
	BuildDir     string       `json:"build_dir"`
	AssetDir     string       `json:"asset_dir"`
	OutputSource bool         `json:"output_source"`
	Logger       *l.Logger    `json:"-"`
	Source       string       `json:"-"`
	StringDefs   []*StringDef `json:"-"`
}

// NewCompiler returns a basic Compiler object
func NewCompiler(scripts []string, outfile, os, arch string, sourceOut bool) *Compiler {
	logger := l.New()
	logger.Prefix = fmt.Sprintf("%s%s%s%s%s%s ", l.Colorize("[", l.Bold+l.White), l.Colorize("GENESIS", l.Bold+l.LightRed), l.Colorize(":", "\033[0m"+l.White), l.Colorize("compiler()", l.LightYellow), l.Colorize("]", l.Bold+l.White), "\033[0m")
	logger.DisabledInfo = true
	if outfile == "-" && !sourceOut {
		logger.Crit("You need either -outfile or -source specified to build.")
	}
	vms := []*VMBundle{}
	for _, s := range scripts {
		vms = append(vms, &VMBundle{
			ID:         RandUpperAlphaString(18),
			ScriptFile: s,
			AssetFiles: []string{},
			Embeds:     []EmbeddedFile{},
		})
	}
	return &Compiler{
		VMs:          vms,
		OutputFile:   outfile,
		Logger:       logger,
		OS:           os,
		Arch:         arch,
		OutputSource: sourceOut,
		StringDefs:   []*StringDef{},
	}
}

// CreateBuildDir sets up the compiler's build directory
func (c *Compiler) CreateBuildDir() {
	dirName := RandStringRunes(16)
	bd := filepath.Join(os.TempDir(), dirName)
	err := os.MkdirAll(bd, 0744)
	if err != nil {
		c.Logger.Critf("Cannot create build directory: dir=%s error=%s", bd, err.Error())
	}
	ad := filepath.Join(bd, "assets")
	os.MkdirAll(ad, 0744)
	c.BuildDir = bd
	c.AssetDir = ad
}

// ParseAssets normalizes the import files into localized assets
func (c *Compiler) ParseAssets(filename string) []string {
	imports := []string{}
	script, err := ioutil.ReadFile(filename)
	if err != nil {
		c.Logger.Critf("Error reading genesis script: %s", err.Error())
	}
	r := regexp.MustCompile(`//import:(.*)\n`)
	matches := r.FindAllString(string(script), -1)
	for _, rawF := range matches {
		f := strings.TrimSpace(strings.Replace(rawF, "//import:", "", -1))
		if _, err := os.Stat(f); os.IsNotExist(err) {
			c.Logger.Critf("Asset file does not exist: %s", f)
			continue
		}
		imports = append(imports, f)
	}
	r = regexp.MustCompile(`//url_import:(.*)\n`)
	matches = r.FindAllString(string(script), -1)
	for _, rawF := range matches {
		f := strings.TrimSpace(strings.Replace(rawF, "//url_import:", "", -1))
		u, err := url.Parse(f)
		if err != nil {
			c.Logger.Critf("Could not parse URL: %s", err.Error())
		}
		filename := path.Base(u.Path)

		randVector := RandUpperAlphaString(12)

		dir, err := ioutil.TempDir("", randVector)
		if err != nil {
			c.Logger.Critf("Could create temp directory: %s", err.Error())
		}
		filePath := filepath.Join(dir, filename)
		out, err := os.Create(filePath)
		if err != nil {
			c.Logger.Critf("Could create temp file: %s", err.Error())
		}
		defer out.Close()

		resp, err := http.Get(u.String())
		if err != nil {
			c.Logger.Critf("Could not retreive URL: %s", err.Error())
		}
		defer resp.Body.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			c.Logger.Critf("Could not save temp file: %s", err.Error())
		}

		imports = append(imports, filePath)
	}
	return imports
}

func (c *Compiler) gatherAssets() {
	for _, vm := range c.VMs {
		assets := c.ParseAssets(vm.ScriptFile)
		for _, asset := range assets {
			tempFile := filepath.Join(c.AssetDir, filepath.Base(asset))
			err := LocalCopyFile(asset, tempFile)
			if err != nil {
				c.Logger.Critf("Asset file copy error: file=%s, error=%s", asset, err.Error())
			}
			c.Logger.Logf("Packing File: %s", filepath.Base(asset))
			vm.AssetFiles = append(vm.AssetFiles, tempFile)
		}
	}
}

func (c *Compiler) writeScript() {
	for _, vm := range c.VMs {
		if _, err := os.Stat(vm.ScriptFile); os.IsNotExist(err) {
			c.Logger.Critf("Genesis Script does not exist: %s", vm.ScriptFile)
		}
		entryFile := filepath.Join(c.AssetDir, fmt.Sprintf("%s.gs", vm.ID))
		err := LocalCopyFile(vm.ScriptFile, entryFile)
		if err != nil {
			c.Logger.Critf("Asset file copy error: file=%s, error=%s", vm.ScriptFile, err.Error())
		}
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

func (c *Compiler) buildEntryPoint() {
	tmpl := template.New("gse_builder")
	tmpl.Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
	newTmpl, err := tmpl.Parse(string(MustAsset("templates/entrypoint.go.tmpl")))
	if err != nil {
		c.Logger.Critf("Error generating source: %s", err.Error())
	}
	var buf bytes.Buffer
	err = newTmpl.Execute(&buf, &c)
	if err != nil {
		c.Logger.Critf("Error generating source: %s", err.Error())
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
			c.Logger.Critf("Error generating obfuscated string: %s", err.Error())
		}
		var buf bytes.Buffer
		err = newTmpl.Execute(&buf, str)
		if err != nil {
			c.Logger.Critf("Error generating obfuscated string: %s", err.Error())
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
		quick.Highlight(os.Stdout, newSourceB, "go", "terminal", "vim")
		return
	}
	err := ioutil.WriteFile(filepath.Join(c.BuildDir, "main.go"), []byte(newSourceB), 0644)
	if err != nil {
		c.Logger.Critf("Error writing main.go: %s", err.Error())
	}
}

func (c *Compiler) Do() {
	cwd, _ := os.Getwd()
	c.CreateBuildDir()
	c.gatherAssets()
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
		c.Logger.Logf("Obfuscating Binary...")
		c.ObfuscateBinary()
		// TODO: Fix GOTTI Binary Obfuscator #sadpanda
		// c.Logger.Logf("Mordor-ifying Binary...")
		// switch c.OS {
		// case "darwin":
		// 	c.MordorifyDarwin()
		// case "linux":
		// 	c.MordorifyLinux()
		// case "windows":
		// 	c.MordorifyWindows()
		// default:
		// 	c.Logger.Critf("GOOS specified is not supported.")
		// }
	}
	os.Chdir(cwd)
	os.RemoveAll(c.BuildDir)
}

func (c *Compiler) LollerSkateDaStringz(source []byte) []byte {
	fset := gotoken.NewFileSet()
	file, err := goparser.ParseFile(fset, "", source, 0)
	if err != nil {
		c.Logger.Critf("Could not parse Golang source: %s", err.Error())
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
