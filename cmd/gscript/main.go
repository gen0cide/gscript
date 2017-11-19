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

	"github.com/gen0cide/gscript"
	"github.com/google/go-github/github"
	update "github.com/inconshreveable/go-update"
	"github.com/urfave/cli"
)

func main() {

	var outputFile string
	var compilerOS string
	var compilerArch string
	var outputSource = false

	app := cli.NewApp()
	app.Name = "gscript"
	app.Usage = "Interact with the Genesis Scripting Engine (GSE)"
	app.Version = "0.0.6"
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
			Action: func(c *cli.Context) error {
				gse := gscript.New("main")
				gse.EnableLogging()
				filename := c.Args().Get(0)
				if len(filename) == 0 {
					gse.LogCritf("You did not supply a filename!")
				}
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					gse.LogCritf("File does not exist: %s", filename)
				}
				_, err := exec.LookPath("jshint")
				if err != nil {
					gse.LogCritf("You do not have jshint in your path. Run: npm install -g jshint")
				}

				jshCmd := exec.Command("jshint", filename)
				jshOutput, err := jshCmd.CombinedOutput()
				if err != nil {
					gse.LogCritf("File Not Valid Javascript!\n -- JSHint Output:\n%s", jshOutput)
				}
				data, err := ioutil.ReadFile(filename)
				gse.SetName(filename)
				gse.CreateVM()
				err = gse.ValidateAST(data)
				if err != nil {
					gse.LogErrorf("Invalid Script Error: %s", err.Error())
				} else {
					gse.LogInfof("Script Valid: %s", filename)
				}
				return nil
			},
		},
		{
			Name:    "shell",
			Aliases: []string{"s"},
			Usage:   "Run an interactive GSE console session.",
			Action: func(c *cli.Context) error {
				gse := gscript.New("shell")
				gse.EnableLogging()
				gse.CreateVM()
				gse.InteractiveSession()
				return nil
			},
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Update Genesis Scripting Engine to the latest version.",
			Action:  UpdateCLI,
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
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					gse := gscript.NewCompiler([]string{}, "", "", "", false)
					gse.Logger.Critf("You did not specify a genesis script!")
				}
				scriptFiles := c.Args()
				if !outputSource && outputFile == "-" {
					outputFile = filepath.Join(os.TempDir(), fmt.Sprintf("%d_genesis.bin", time.Now().Unix()))
				}
				compiler := gscript.NewCompiler(scriptFiles, outputFile, compilerOS, compilerArch, outputSource)
				compiler.Do()
				if !outputSource {
					compiler.Logger.Logf("Your binary is located at: %s", outputFile)
				}
				return nil
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run a Genesis script (Careful, don't infect yourself!).",
			Action: func(c *cli.Context) error {
				gse := gscript.New("main")
				gse.EnableLogging()
				filename := c.Args().Get(0)
				if len(filename) == 0 {
					gse.LogCritf("You did not supply a filename!")
				}
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					gse.LogCritf("File does not exist: %s", filename)
				}
				data, err := ioutil.ReadFile(filename)
				gse.SetName(filename)
				gse.CreateVM()
				err = gse.LoadScript(data)
				if err != nil {
					gse.LogErrorf("Script Error: %s", err.Error())
				} else {
					gse.LogInfof("Script loaded successfully")
				}
				err = gse.ExecutePlan()
				if err != nil {
					gse.LogCritf("Hooks Failure: %s", err.Error())
				}
				gse.LogInfof("Hooks executed successfully")
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func UpdateCLI(c *cli.Context) error {
	gse := gscript.New("UPDATER")
	gse.EnableLogging()
	gse.Logger.DisabledInfo = true

	ctx := context.Background()

	client := github.NewClient(nil)

	repRel, _, err := client.Repositories.GetLatestRelease(ctx, "gen0cide", "gscript")
	if err != nil {
		gse.LogCritf("Github Error: " + err.Error())
	}

	tagName := strings.TrimLeft(repRel.GetTagName(), "v")

	if c.App.Version == tagName {
		gse.LogInfof("Running Lastest Version (no update available): %s", c.App.Version)
		return nil
	}

	gse.LogInfof("New Version Found: %s", tagName)

	fileName := fmt.Sprintf("gscript_%s_%s_amd64.zip", tagName, runtime.GOOS)

	assetID := 0
	fileSize := 0

	for _, asset := range repRel.Assets {
		if asset.GetName() == fileName {
			assetID = asset.GetID()
			fileSize = asset.GetSize()
			gse.LogInfof("Downloading New Binary: %s", asset.GetBrowserDownloadURL())
		}
	}

	if assetID == 0 {
		gse.LogCritf("No release found for your OS! WTF? Report this.")
	}

	_, redirectURL, err := client.Repositories.DownloadReleaseAsset(ctx, "gen0cide", "gscript", assetID)
	if err != nil {
		gse.LogCritf("Github Error: %s", err.Error())
	}

	if len(redirectURL) == 0 {
		gse.LogCritf("There was an error retriving the release from Github.")
	}

	resp, err := http.Get(redirectURL)
	if err != nil {
		gse.LogCritf("Error Retrieving Binary: %s", err.Error())
	}
	defer resp.Body.Close()
	compressedFile, err := ioutil.ReadAll(resp.Body)

	gse.LogInfof("Uncompressing Binary...")

	var binary bytes.Buffer
	writer := bufio.NewWriter(&binary)
	compressedReader := bytes.NewReader(compressedFile)

	r, err := zip.NewReader(compressedReader, int64(fileSize))
	if err != nil {
		gse.LogCritf("Error Buffering Zip File: %s", err.Error())
	}

	for _, zf := range r.File {
		if zf.Name != "gscript" && zf.Name != "gscript.exe" {
			continue
		}
		src, err := zf.Open()
		if err != nil {
			gse.LogCritf("Error Unzipping File: %s", err.Error())
		}
		defer src.Close()

		io.Copy(writer, src)
	}

	reader := bufio.NewReader(&binary)

	err = update.Apply(reader, update.Options{})
	if err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			gse.LogCritf("Failed to rollback from bad update: %v", rerr)
		}
		gse.LogCritf("Update Failed - original version rolled back successfully.")
	}

	gse.LogInfof("Successfully updated to gscript v%s", tagName)
	return nil
}
