package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gen0cide/gscript/compiler"
	"github.com/gen0cide/gscript/compiler/computil"
	"github.com/urfave/cli/v2"
)

var (
	shellOpts    = computil.DefaultOptions()
	shellCommand = &cli.Command{
		Name:      "shell",
		Usage:     "drop into an interactive REPL within the genesis runtime",
		UsageText: "gscript shell [--macro MACRO] [--macro MACRO] ...",
		Action:    interactiveShellCommand,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "build-dir",
				Usage:       "Perform the gscript compile in a specific build directory.",
				Destination: &shellOpts.BuildDir,
			},
			&cli.StringSliceFlag{
				Name:  "macro, m",
				Usage: "apply a compiler macro to the interactive shell",
			},
		},
	}
)

func interactiveShellCommand(c *cli.Context) error {
	displayBefore = false
	buf := new(bytes.Buffer)
	if len(c.StringSlice("macro")) > 0 {
		for _, m := range c.StringSlice("macro") {
			buf.WriteString("//")
			buf.WriteString(m)
			buf.WriteString("\n")
		}
	}
	buf.WriteString("\n")
	buf.WriteString(string(computil.MustAsset("debugger.gs")))
	randDirName := computil.RandLowerAlphaString(18)
	randBinName := computil.RandLowerAlphaString(18)
	tmpDir := filepath.Join(os.TempDir(), randDirName)
	scriptPath := filepath.Join(tmpDir, "debugger")
	os.MkdirAll(tmpDir, 0755)
	shellOpts.ObfuscationLevel = 3
	shellOpts.ImportAllNativeFuncs = true
	shellOpts.UseHumanReadableNames = true
	shellOpts.DebuggerEnabled = true
	shellOpts.LoggingEnabled = true
	exePath := filepath.Join(shellOpts.BuildDir, randBinName)
	if runtime.GOOS == "windows" {
		exePath = fmt.Sprintf("%s.exe", exePath)
	}
	shellOpts.OutputFile = exePath
	gc := compiler.NewWithOptions(shellOpts)
	gc.SetLogger(cliLogger)
	ioutil.WriteFile(scriptPath, buf.Bytes(), 0644)
	gc.AddScript(scriptPath)
	err := gc.Do()
	if err != nil {
		cliLogger.Errorf("Build Dir Located At: %s", gc.BuildDir)
		return err
	}
	err = runShell(exePath)
	os.RemoveAll(tmpDir)
	return err
}

func runShell(exePath string) error {
	cmd := exec.Command(exePath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}
