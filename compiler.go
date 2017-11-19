package gscript

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/alecthomas/chroma/quick"
	"github.com/happierall/l"
)

var TempStart = "package %s\n\nimport (\n\t\"bytes\"\n\t\"compress/flate\"\n\t\"io/ioutil\"\n\n\t\"github.com/gen0cide/gscript\"\n)\n"

var TempMain = "func main() {\n\tgse := gscript.New(\"\")\n\tgse.CreateVM()"

var TempEnd = "\tgse.LoadScript(gse.Imports[\"genesis_entry_point.gs\"]())\n\tgse.ExecutePlan()\n}\n"

const EntryPoint = `genesis_entry_point.gs`

type VMBundle struct {
	ID         string         `json:"id"`
	ScriptFile string         `json:"source"`
	AssetFiles []string       `json:"imports"`
	Embeds     []EmbeddedFile `json:"-"`
}

type Compiler struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
	// ScriptFile  string
	// PackageName string
	OutputFile string `json:"output"`
	// AssetFiles  []string
	// Embeds      []EmbeddedFile
	VMs          []*VMBundle `json:"vms"`
	BuildDir     string      `json:"build_dir"`
	AssetDir     string      `json:"asset_dir"`
	OutputSource bool        `json:"output_source"`
	Logger       *l.Logger   `json:"-"`
	Source       string      `json:"-"`
}

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
	}
}

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
	return imports
}

func (c *Compiler) GatherAssets() {
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

func (c *Compiler) WriteScript() {
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

func (c *Compiler) CompileAssets() {
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

func (c *Compiler) BuildEntryPoint() {
	tmpl := template.New("gse_builder")
	tmpl.Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
	newTmpl, err := tmpl.Parse(string(MustAsset("templates/compile_template.go.tmpl")))
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

func (c *Compiler) WriteSource() {
	if c.OutputSource {
		c.Logger.Log("YEP")
		quick.Highlight(os.Stdout, string(c.Source), "go", "terminal", "vim")
		return
	}
	err := ioutil.WriteFile(filepath.Join(c.BuildDir, "main.go"), []byte(c.Source), 0644)
	if err != nil {
		c.Logger.Critf("Error writing main.go: %s", err.Error())
	}
}

func (c *Compiler) Do() {
	cwd, _ := os.Getwd()
	c.CreateBuildDir()
	c.GatherAssets()
	c.WriteScript()
	os.Chdir(c.BuildDir)
	c.CompileAssets()
	c.BuildEntryPoint()
	os.RemoveAll(c.AssetDir)
	c.WriteSource()
	if !c.OutputSource {
		cmd := exec.Command("go", "build", "-o", c.OutputFile)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", c.OS))
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", c.Arch))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
	os.Chdir(cwd)
	os.RemoveAll(c.BuildDir)
}
