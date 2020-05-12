package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

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
	shellOpts.ObfuscationLevel = 3
	shellOpts.ImportAllNativeFuncs = true
	shellOpts.UseHumanReadableNames = true
	shellOpts.DebuggerEnabled = true
	shellOpts.LoggingEnabled = true
	gc := compiler.New(&shellOpts)
	scriptPath := filepath.Join(gc.BuildDir, "debugger")
	gc.SetLogger(cliLogger)
	err := ioutil.WriteFile(scriptPath, buf.Bytes(), 0644)
	if err != nil {
		cliLogger.Errorf("Error writing script to file path: %s", scriptPath)
		return err
	}
	err = gc.AddScript(scriptPath)
	if err != nil {
		cliLogger.Errorf("Error adding to runtime: %s", scriptPath)
		return err
	}
	err = gc.Do()
	if err != nil {
		cliLogger.Errorf("Build Dir Located At: %s", gc.BuildDir)
		return err
	}
	err = runShell(gc.OutputFile)
	os.RemoveAll(gc.BuildDir)
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
