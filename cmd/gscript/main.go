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
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/gen0cide/gscript/compiler"
	"github.com/gen0cide/gscript/debugger"
	"github.com/gen0cide/gscript/engine"
	"github.com/gen0cide/gscript/logging"
	"github.com/google/go-github/github"
	update "github.com/inconshreveable/go-update"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var outputFile string
var compilerOS string
var compilerArch string
var outputSource = false
var compressBinary = false

func main() {
	app := cli.NewApp()
	app.Name = "gscript"
	app.Usage = "Interact with the Genesis Scripting Engine (GSE)"
	app.Version = "0.0.8"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Alex Levinson",
			Email: "gen0cide.threats@gmail.com",
		},
	}
	app.Copyright = "(c) 2017 Alex Levinson"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Run gscript in debug mode.",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Suppress all logging output.",
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
			Usage:   "Run an interactive GSE console session.",
			Action:  InteractiveShell,
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Update Genesis Scripting Engine to the latest version.",
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
			},
			Action: CompileScript,
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run a Genesis script (Careful, don't infect yourself!).",
			Action:  RunScript,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func TestScript(c *cli.Context) error {
	gse := engine.New("test")
	logger := logrus.New()
	logger.Formatter = &logging.GSEFormatter{}
	logger.Out = logging.LogWriter{Name: "test"}
	gse.SetLogger(logger)
	filename := c.Args().Get(0)
	if len(filename) == 0 {
		gse.Logger.Fatalf("You did not supply a filename!")
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		gse.Logger.Fatalf("File does not exist: %s", filename)
	}
	_, err := exec.LookPath("jshint")
	if err != nil {
		gse.Logger.Fatalf("You do not have jshint in your path. Run: npm install -g jshint")
	}

	jshCmd := exec.Command("jshint", filename)
	jshOutput, err := jshCmd.CombinedOutput()
	if err != nil {
		gse.Logger.Fatalf("File Not Valid Javascript!\n -- JSHint Output:\n%s", jshOutput)
	}
	data, err := ioutil.ReadFile(filename)
	gse.SetName(filename)
	gse.CreateVM()
	err = gse.ValidateAST(data)
	if err != nil {
		gse.Logger.Errorf("Invalid Script Error: %s", err.Error())
	} else {
		gse.Logger.Infof("Script Valid: %s", filename)
	}
	return nil
}

func InteractiveShell(c *cli.Context) error {
	dbg := debugger.New("shell")
	dbg.SetupDebugEngine()
	dbg.InteractiveSession()
	return nil
}

func CompileScript(c *cli.Context) error {
	logger := logrus.New()
	logger.Formatter = new(logrus.TextFormatter)
	logger.Out = logging.LogWriter{Name: "compiler"}
	if c.NArg() == 0 {
		logger.Fatalf("You did not specify a genesis script!")
	}
	scriptFiles := c.Args()
	if !outputSource && outputFile == "-" {
		outputFile = filepath.Join(os.TempDir(), fmt.Sprintf("%d_genesis.bin", time.Now().Unix()))
	}
	gcc := compiler.NewCompiler(scriptFiles, outputFile, compilerOS, compilerArch, outputSource, compressBinary)
	gcc.Do()
	if !outputSource {
		gcc.Logger.Infof("Your binary is located at: %s", outputFile)
	}
	return nil
}

func NewScript(c *cli.Context) error {
	logger := logrus.New()
	logger.Formatter = new(logrus.TextFormatter)
	logger.Out = logging.LogWriter{Name: "compiler"}
	if c.NArg() == 0 {
		fmt.Printf("%s\n", string(compiler.RetrieveExample()))
		return nil
	}
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
	dbg := debugger.New("runner")
	dbg.SetupDebugEngine()
	filename := c.Args().Get(0)
	if len(filename) == 0 {
		dbg.Engine.Logger.Fatalf("You did not supply a filename!")
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		dbg.Engine.Logger.Fatalf("File does not exist: %s", filename)
	}
	data, err := ioutil.ReadFile(filename)
	err = dbg.Engine.LoadScript(data)
	if err != nil {
		dbg.Engine.Logger.Errorf("Script Error: %s", err.Error())
	} else {
		dbg.Engine.Logger.Infof("Script loaded successfully")
	}
	err = dbg.Engine.ExecutePlan()
	if err != nil {
		dbg.Engine.Logger.Fatalf("Hooks Failure: %s", err.Error())
	}
	dbg.Engine.Logger.Infof("Hooks executed successfully")
	return nil
}

func UpdateCLI(c *cli.Context) error {
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

	assetID := 0
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
