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
	"github.com/urfave/cli"
)

var (
	shellCommand = cli.Command{
		Name:      "shell",
		Usage:     "drop into an interactive REPL within the genesis runtime",
		UsageText: "gscript shell [--macro MACRO] [--macro MACRO] ...",
		Action:    interactiveShellCommand,
		Flags: []cli.Flag{
			cli.StringSliceFlag{
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
	exePath := filepath.Join(tmpDir, randBinName)
	if runtime.GOOS == "windows" {
		exePath = fmt.Sprintf("%s.exe", exePath)
	}
	scriptPath := filepath.Join(tmpDir, "debugger")
	os.MkdirAll(tmpDir, 0755)
	opts := computil.DefaultOptions()
	opts.ObfuscationLevel = 3
	opts.ImportAllNativeFuncs = true
	opts.UseHumanReadableNames = true
	opts.DebuggerEnabled = true
	opts.LoggingEnabled = true
	opts.OutputFile = exePath
	gc := compiler.NewWithOptions(opts)
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
