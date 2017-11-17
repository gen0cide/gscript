package gscript

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/happierall/l"
)

var TempStart = "package %s\n\nimport (\n\t\"bytes\"\n\t\"compress/flate\"\n\t\"io/ioutil\"\n\n\t\"github.com/gen0cide/gscript\"\n)\n"

var TempMain = "func main() {\n\tgse := gscript.New(\"\")\n\tgse.CreateVM()"

var TempEnd = "\tgse.LoadScript(gse.Imports[\"genesis_entry_point.gs\"]())\n\tgse.ExecutePlan()\n}\n"

const EntryPoint = `genesis_entry_point.gs`

type Compiler struct {
	OS          string
	Arch        string
	ScriptFile  string
	PackageName string
	OutputFile  string
	AssetFiles  []string
	Embeds      []EmbeddedFile
	BuildDir    string
	AssetDir    string
	Logger      *l.Logger
}

func NewCompiler(script, outfile, os, arch string) *Compiler {
	logger := l.New()
	logger.Prefix = fmt.Sprintf("%s%s%s%s%s%s ", l.Colorize("[", l.Bold+l.White), l.Colorize("GENESIS", l.Bold+l.LightRed), l.Colorize(":", "\033[0m"+l.White), l.Colorize("compiler()", l.LightYellow), l.Colorize("]", l.Bold+l.White), "\033[0m")
	logger.DisabledInfo = false
	return &Compiler{
		ScriptFile: script,
		OutputFile: outfile,
		Logger:     logger,
		OS:         os,
		Arch:       arch,
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

func (c *Compiler) GatherAssets() {
	script, err := ioutil.ReadFile(c.AssetFiles[0])
	if err != nil {
		c.Logger.Critf("Error reading file: %s", err.Error())
	}
	r := regexp.MustCompile(`//import:(.*)\n`)
	matches := r.FindAllString(string(script), -1)
	for _, rawF := range matches {
		f := strings.TrimSpace(strings.Replace(rawF, "//import:", "", -1))
		if _, err := os.Stat(f); os.IsNotExist(err) {
			c.Logger.Critf("Asset file does not exist: %s", f)
			continue
		}
		err := LocalCopyFile(f, filepath.Join(c.AssetDir, filepath.Base(f)))
		if err != nil {
			c.Logger.Critf("Asset file copy error: file=%s, error=%s", f, err.Error())
		}
		c.AssetFiles = append(c.AssetFiles, filepath.Join(c.AssetDir, filepath.Base(f)))
	}
}

func (c *Compiler) WriteScript() {
	if _, err := os.Stat(c.ScriptFile); os.IsNotExist(err) {
		c.Logger.Critf("Genesis Script does not exist: %s", c.ScriptFile)
	}
	err := LocalCopyFile(c.ScriptFile, filepath.Join(c.AssetDir, EntryPoint))
	if err != nil {
		c.Logger.Critf("Asset file copy error: file=%s, error=%s", c.ScriptFile, err.Error())
	}
	c.AssetFiles = append(c.AssetFiles, filepath.Join(c.AssetDir, EntryPoint))
}

func (c *Compiler) CompileAssets() {
	for _, f := range c.AssetFiles {
		e := EmbeddedFile{
			SourcePath: f,
		}
		e.Embed()
		c.Embeds = append(c.Embeds, e)
	}
}

func (c *Compiler) BuildEntryPoint() {
	sourceFile := fmt.Sprintf(TempStart, "main")
	for _, e := range c.Embeds {
		sourceFile += "\n"
		sourceFile += e.VariableDef
		sourceFile += "\n"
	}
	sourceFile += TempMain
	sourceFile += "\n"
	for _, e := range c.Embeds {
		sourceFile += fmt.Sprintf("\tgse.AddImport(\"%s\", %s)\n", e.Filename, e.NameHash)
	}
	sourceFile += TempEnd
	err := ioutil.WriteFile(filepath.Join(c.BuildDir, "main.go"), []byte(sourceFile), 0644)
	if err != nil {
		c.Logger.Critf("Error writing main.go: %s", err.Error())
	}
}

func (c *Compiler) Do() {
	cwd, _ := os.Getwd()
	c.CreateBuildDir()
	c.WriteScript()
	c.GatherAssets()
	os.Chdir(c.BuildDir)
	c.CompileAssets()
	c.BuildEntryPoint()
	os.RemoveAll(c.AssetDir)
	cmd := exec.Command("go", "build", "-o", c.OutputFile)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", c.OS))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", c.Arch))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	os.Chdir(cwd)
	c.Logger.Debugf("%s", c.BuildDir)
}
