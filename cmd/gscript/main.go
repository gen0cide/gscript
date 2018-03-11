package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gen0cide/gscript"
	"github.com/gen0cide/gscript/compiler"
	"github.com/gen0cide/gscript/debugger"
	"github.com/gen0cide/gscript/engine"
	"github.com/gen0cide/gscript/logging"
	"github.com/google/go-github/github"
	update "github.com/inconshreveable/go-update"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	outputFile   string
	compilerOS   string
	compilerArch string

	outputSource   = false
	compressBinary = false
	enableLogging  = false
	enableDebug    = false
	retainBuildDir = false
)

func main() {
	cli.AppHelpTemplate = fmt.Sprintf("%s\n\n%s", logging.AsciiLogo(), cli.AppHelpTemplate)
	cli.CommandHelpTemplate = fmt.Sprintf("%s\n\n%s", logging.AsciiLogo(), cli.CommandHelpTemplate)
	app := cli.NewApp()
	app.Writer = color.Output
	app.ErrWriter = color.Output
	app.Name = "gscript"
	app.Usage = "Command Line application for interacting with the GENESIS Scripting Engine."
	app.Version = gscript.Version
	app.Authors = []cli.Author{
		{
			Name:  "Alex Levinson",
			Email: "gen0cide.threats@gmail.com",
		},
	}
	app.Copyright = "(c) 2018 Alex Levinson"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable verbose output of runtime debugging logs.",
			Destination: &enableDebug,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Check a GSE script for syntax errors.",
			Action:  TestScript,
		},
		{
			Name:    "shell",
			Aliases: []string{"s"},
			Usage:   "Run an interactive gscript REPL.",
			Action:  InteractiveShell,
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Update the gscript CLI binary to the latest version.",
			Action:  UpdateCLI,
		},
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "Writes an example gscript to either the given path or STDOUT.",
			Action:  NewScript,
		},
		{
			Name:    "compile",
			Aliases: []string{"c"},
			Usage:   "Compile genesis scripts into a stand alone binary.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "outfile",
					Value:       "-",
					Usage:       "Location of the compiled binary (STDOUT if none specified)",
					Destination: &outputFile,
				},
				cli.StringFlag{
					Name:        "os",
					Value:       runtime.GOOS,
					Usage:       "The GOOS you wish to use for your compiled binary.",
					Destination: &compilerOS,
				},
				cli.StringFlag{
					Name:        "arch",
					Value:       runtime.GOARCH,
					Usage:       "The GOARCH you wish to use for your compiled binary.",
					Destination: &compilerArch,
				},
				cli.BoolFlag{
					Name:        "source",
					Usage:       "Do not compile the generated code. Output source instead.",
					Destination: &outputSource,
				},
				cli.BoolFlag{
					Name:        "upx",
					Usage:       "Attempts to UPX the final binary to reduce file size.",
					Destination: &compressBinary,
				},
				cli.BoolFlag{
					Name:        "enable-logging",
					Usage:       "Enables debug logging in the finished binary. WARNING: Will create large binaries!",
					Destination: &enableLogging,
				},
			},
			Action: CompileScript,
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run a Genesis script locally (Careful not to infect yourself).",
			Action:  RunScript,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func TestScript(c *cli.Context) error {
	logging.PrintLogo()
	dbg := debugger.New("runner")
	dbg.SetupDebugEngine()
	filename := c.Args().Get(0)
	if len(filename) == 0 {
		dbg.Engine.Logger.Fatalf("You did not supply a filename!")
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		dbg.Engine.Logger.Fatalf("File does not exist: %s", filename)
	}
	dbg = debugger.New(filepath.Base(filename))
	dbg.SetupDebugEngine()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		dbg.Engine.Logger.Fatalf("Error reading file: %s", err.Error())
	}
	dbg.LoadScript(string(data), filepath.Base(filename))
	if err != nil {
		dbg.Engine.Logger.Errorf("Script Error: %s", err.Error())
		return err
	}
	dbg.Engine.Logger.Infof("Script is passed static analysis")
	return nil

}

func InteractiveShell(c *cli.Context) error {
	logging.PrintLogo()
	dbg := debugger.New("shell")
	dbg.SetupDebugEngine()
	dbg.InteractiveSession()
	return nil
}

func CompileScript(c *cli.Context) error {
	logger := logrus.New()
	logger.Formatter = &logging.GSEFormatter{}
	logger.Out = logging.LogWriter{Name: "compiler"}
	if enableDebug {
		logger.Level = logrus.DebugLevel
	}
	if c.NArg() == 0 {
		logger.Fatalf("You did not specify a genesis script!")
	}
	scriptFiles := []string{}
	for _, a := range c.Args() {
		f, err := filepath.Glob(a)
		if err != nil {
			logger.Fatalf("Bad file glob: %s", err.Error())
		}
		for _, n := range f {
			scriptFiles = append(scriptFiles, n)
		}
	}

	finalFile := ""

	if !outputSource && outputFile == "-" {
		finalFile = filepath.Join(os.TempDir(), fmt.Sprintf("%d_genesis.bin", time.Now().Unix()))
	} else {
		f, err := filepath.Abs(outputFile)
		finalFile = f
		if err != nil {
			logger.Fatalf("Cannot determine path to outfile: %s", err.Error())
		}
	}
	gcc := compiler.NewCompiler(scriptFiles, finalFile, compilerOS, compilerArch, outputSource, compressBinary, enableLogging)
	gcc.Logger = logger
	gcc.Do()
	if !outputSource {
		gcc.Logger.Infof("Your binary is located at: %s", finalFile)
	}
	return nil
}

func NewScript(c *cli.Context) error {
	logger := logrus.New()
	logger.Formatter = new(logging.GSEFormatter)
	logger.Out = logging.LogWriter{Name: "compiler"}
	if c.NArg() == 0 {
		fmt.Printf("%s\n", string(compiler.RetrieveExample()))
		return nil
	}
	logging.PrintLogo()
	scriptFiles := c.Args()
	for _, f := range scriptFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			ioutil.WriteFile(f, compiler.RetrieveExample(), 0644)
			logger.Infof("Wrote Example File: %s", f)
			continue
		}
		logger.Errorf("File either exists or has bad perms. Skipping %s", f)
	}
	return nil
}

func RunScript(c *cli.Context) error {
	logging.PrintLogo()
	dbg := debugger.New("runner")
	dbg.SetupDebugEngine()
	filename := c.Args().Get(0)
	if len(filename) == 0 {
		dbg.Engine.Logger.Fatalf("You did not supply a filename!")
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		dbg.Engine.Logger.Fatalf("File does not exist: %s", filename)
	}
	dbg = debugger.New(filepath.Base(filename))
	dbg.SetupDebugEngine()
	if !enableDebug {
		dbg.Logger.Level = logrus.InfoLevel
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		dbg.Engine.Logger.Fatalf("Error reading file: %s", err.Error())
	}
	dbg.LoadScript(string(data), filepath.Base(filename))
	if err != nil {
		dbg.Engine.Logger.Errorf("Script Error: %s", err.Error())
	} else {
		dbg.Engine.Logger.Infof("Script loaded successfully")
	}
	err = dbg.Engine.ExecutePlan()
	if err != nil {
		dbg.Engine.Logger.Fatalf("Hooks Failure: %s", err.Error())
	}
	dbg.Engine.Logger.Debugf("Hooks executed successfully")
	return nil
}

func UpdateCLI(c *cli.Context) error {
	logging.PrintLogo()
	gse := engine.New("updater")
	logger := logrus.New()
	logger.Formatter = &logging.GSEFormatter{}
	logger.Out = logging.LogWriter{Name: "updater"}
	gse.SetLogger(logger)

	ctx := context.Background()

	client := github.NewClient(nil)

	repRel, _, err := client.Repositories.GetLatestRelease(ctx, "gen0cide", "gscript")
	if err != nil {
		gse.Logger.Fatalf("Github Error: " + err.Error())
	}

	tagName := strings.TrimLeft(repRel.GetTagName(), "v")

	if c.App.Version == tagName {
		gse.Logger.Infof("Running Lastest Version (no update available): %s", c.App.Version)
		return nil
	}

	gse.Logger.Infof("New Version Found: %s", tagName)

	fileName := fmt.Sprintf("gscript_%s_%s_amd64.zip", tagName, runtime.GOOS)

	assetID := int64(0)
	fileSize := 0

	for _, asset := range repRel.Assets {
		if asset.GetName() == fileName {
			assetID = asset.GetID()
			fileSize = asset.GetSize()
			gse.Logger.Infof("Downloading New Binary: %s", asset.GetBrowserDownloadURL())
		}
	}

	if assetID == 0 {
		gse.Logger.Fatalf("No release found for your OS! WTF? Report this.")
	}

	_, redirectURL, err := client.Repositories.DownloadReleaseAsset(ctx, "gen0cide", "gscript", assetID)
	if err != nil {
		gse.Logger.Fatalf("Github Error: %s", err.Error())
	}

	if len(redirectURL) == 0 {
		gse.Logger.Fatalf("There was an error retriving the release from Github.")
	}

	resp, err := http.Get(redirectURL)
	if err != nil {
		gse.Logger.Fatalf("Error Retrieving Binary: %s", err.Error())
	}
	defer resp.Body.Close()
	compressedFile, err := ioutil.ReadAll(resp.Body)

	gse.Logger.Infof("Uncompressing Binary...")

	var binary bytes.Buffer
	writer := bufio.NewWriter(&binary)
	compressedReader := bytes.NewReader(compressedFile)

	r, err := zip.NewReader(compressedReader, int64(fileSize))
	if err != nil {
		gse.Logger.Fatalf("Error Buffering Zip File: %s", err.Error())
	}

	for _, zf := range r.File {
		if zf.Name != "gscript" && zf.Name != "gscript.exe" {
			continue
		}
		src, err := zf.Open()
		if err != nil {
			gse.Logger.Fatalf("Error Unzipping File: %s", err.Error())
		}
		defer src.Close()

		io.Copy(writer, src)
	}

	reader := bufio.NewReader(&binary)

	err = update.Apply(reader, update.Options{})
	if err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			gse.Logger.Fatalf("Failed to rollback from bad update: %v", rerr)
		}
		gse.Logger.Fatalf("Update Failed - original version rolled back successfully.")
	}

	gse.Logger.Infof("Successfully updated to gscript v%s", tagName)
	return nil
}
