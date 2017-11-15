package gscript

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/happierall/l"
	bindata "github.com/jteeuwen/go-bindata"
)

var CompiledTemplate = `package %s

import (
  "github.com/gen0cide/gscript"
)

func Run() {
  gse := gscript.New("")
  gse.CreateVM()
  gse.LoadScript(MustAsset("genesis_entry_point.gs"))
  gse.ExecutePlan()
}

`

const EntryPoint = `genesis_entry_point.gs`

type Compiler struct {
	ScriptFile  string
	PackageName string
	OutputFile  string
	AssetFiles  []string
	BuildDir    string
	AssetDir    string
	Logger      *l.Logger
}

func NewCompiler(script, pkg, outfile string, assetFiles []string) *Compiler {
	logger := l.New()
	logger.Prefix = fmt.Sprintf("%s%s%s%s%s%s ", l.Colorize("[", l.Bold+l.White), l.Colorize("GENESIS", l.Bold+l.LightRed), l.Colorize(":", "\033[0m"+l.White), l.Colorize("compiler()", l.LightYellow), l.Colorize("]", l.Bold+l.White), "\033[0m")
	logger.DisabledInfo = false
	return &Compiler{
		ScriptFile:  script,
		OutputFile:  outfile,
		PackageName: pkg,
		AssetFiles:  assetFiles,
		Logger:      logger,
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
	for _, f := range c.AssetFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			c.Logger.Critf("Asset file does not exist: %s", f)
			continue
		}
		err := LocalCopyFile(f, filepath.Join(c.AssetDir, filepath.Base(f)))
		if err != nil {
			c.Logger.Critf("Asset file copy error: file=%s, error=%s", f, err.Error())
		}
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
	cAssets := bindata.NewConfig()
	cAssets.Package = c.PackageName
	cAssets.Prefix = c.AssetDir
	cAssets.Output = filepath.Join(c.BuildDir, "embeds.go")
	cAssets.Input = make([]bindata.InputConfig, len(c.AssetFiles))
	for i := range cAssets.Input {
		corraspondingFile := filepath.Base(c.AssetFiles[i])
		cAssets.Input[i] = parseInput(filepath.Join(c.AssetDir, corraspondingFile))
	}
	bindata.Translate(cAssets)
}

func (c *Compiler) BuildEntryPoint() {
	entry := fmt.Sprintf(CompiledTemplate, c.PackageName)
	err := ioutil.WriteFile(filepath.Join(c.BuildDir, "entry.go"), []byte(entry), 0644)
	if err != nil {
		c.Logger.Critf("Error writing entry.go: %s", err.Error())
	}
}

func (c *Compiler) Do() {
	c.CreateBuildDir()
	c.GatherAssets()
	c.WriteScript()
	os.Chdir(c.BuildDir)
	c.CompileAssets()
	c.BuildEntryPoint()
	os.RemoveAll(c.AssetDir)
	c.Logger.Logf("Genesis Script successfully compiled to: %s", c.BuildDir)
}

func parseInput(path string) bindata.InputConfig {
	if strings.HasSuffix(path, "/...") {
		return bindata.InputConfig{
			Path:      filepath.Clean(path[:len(path)-4]),
			Recursive: true,
		}
	} else {
		return bindata.InputConfig{
			Path:      filepath.Clean(path),
			Recursive: false,
		}
	}

}
