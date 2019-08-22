package main

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/gen0cide/gscript/compiler"
	"github.com/urfave/cli"
)

var (
	compileCommand = cli.Command{
		Name:      "compile",
		Usage:     "compiles the provided scripts using the genesis compiler",
		UsageText: "gscript compile [OPTIONS] SCRIPT [SCRIPT SCRIPT ...]",
		Flags: []cli.Flag{
	    	cli.StringFlag{
				Name:        "go-build-args",
				Usage:       "extra args to pass to the 'go build' command (i.e. -buildmode=c-shared)",
				Destination: &defaultCompileOptions.BuildArgs,
			},
		    cli.StringFlag{
				Name:        "os",
				Usage:       "operating system to target for native compilation",
				Value:       defaultCompileOptions.OS,
				Destination: &defaultCompileOptions.OS,
			},
			cli.StringFlag{
				Name:        "arch",
				Usage:       "architecture to target for native compilation",
				Value:       defaultCompileOptions.Arch,
				Destination: &defaultCompileOptions.Arch,
			},
			cli.StringFlag{
				Name:        "output-file, o",
				Usage:       "location to write final compiled binary",
				Value:       defaultCompileOptions.OutputFile,
				Destination: &defaultCompileOptions.OutputFile,
			},
			cli.BoolFlag{
				Name:        "windowsgui",
				Usage:       "Enable -ldflags -H=windowsgui to prevent a console window from appearing if double clicked",
				Destination: &defaultCompileOptions.WindowsGui,
			},
			cli.BoolFlag{
				Name:        "keep-build-dir",
				Usage:       "keep the build directory of the genesis intermediate representation (default: false)",
				Destination: &defaultCompileOptions.SaveBuildDir,
			},
			cli.BoolFlag{
				Name:        "enable-upx-compression",
				Usage:       "compress the final binary using UPX (default: false)",
				Destination: &defaultCompileOptions.UPXEnabled,
			},
			cli.BoolFlag{
				Name:        "enable-logging",
				Usage:       "enable logging in the final binary for debugging (default: false)",
				Destination: &defaultCompileOptions.LoggingEnabled,
			},
			cli.BoolFlag{
				Name:        "enable-debugging",
				Usage:       "enable the interactive debugger in the final binary (default: false)",
				Destination: &defaultCompileOptions.DebuggerEnabled,
			},
			cli.BoolFlag{
				Name:        "enable-human-readable-names",
				Usage:       "use human readable names in the genesis intermediate representation (default: false)",
				Destination: &defaultCompileOptions.UseHumanReadableNames,
			},
			cli.BoolFlag{
				Name:        "enable-import-all-native-funcs",
				Usage:       "link all possible native functions into the runtime, not just those called by the scripts (default: false)",
				Destination: &defaultCompileOptions.ImportAllNativeFuncs,
			},
			cli.BoolFlag{
				Name:        "disable-native-compilation",
				Usage:       "do not compile the intermediate representation to a native binary (default: false)",
				Destination: &defaultCompileOptions.SkipCompilation,
			},
			cli.BoolFlag{
				Name:        "enable-test-build",
				Usage:       "enable the test harness in the build - for testing only! (default: false)",
				Destination: &defaultCompileOptions.EnableTestBuild,
			},
			cli.BoolFlag{
				Name:        "touch-a-silmaril",
				Usage:       "Easter egg expirimental obfuscator of black magic.",
				Destination: &defaultCompileOptions.ForceUseMordorifier,
			},
			cli.IntFlag{
				Name:        "obfuscation-level",
				Usage:       "override the default obfuscation level, where argument can be 0-4 with 0 being full and 4 being none",
				Destination: &defaultCompileOptions.ObfuscationLevel,
			},
		},
		Action: compileScriptCommand,
	}

	boolText = map[bool]string{
		true:  color.New(color.FgHiCyan, color.Bold).Sprintf("%-72s", `[ENABLED]`),
		false: color.New(color.FgRed).Sprintf("%-72s", `[DISABLED]`),
	}
)

func obfText(lvl int) string {
	switch lvl {
	case 0:
		return color.New(color.FgRed).Sprintf("%-72s", `ALL OBFUSCATION ENABLED`)
	case 1:
		return color.New(color.FgHiMagenta, color.Bold).Sprintf("%-72s", `POST COMPILATION DISABLED`)
	case 2:
		return color.New(color.FgHiYellow, color.Bold).Sprintf("%-72s", `PRE & POST COMPILATION DISABLED`)
	default:
		return color.New(color.FgHiCyan, color.Bold, color.BlinkRapid).Sprintf("%-72s", `ALL OBFUSCATION DISABLED`)
	}
}

func copt(label string, val interface{}) string {
	prefix := color.HiWhiteString("%25s", label)
	middle := color.YellowString(":")
	ending := bopt(val)
	return fmt.Sprintf("%s%s %s", prefix, middle, ending)
}

func bopt(val interface{}) string {
	switch v := val.(type) {
	case string:
		return color.HiGreenString("%-72s", v)
	case int:
		return obfText(v)
	case bool:
		return boolText[v]
	default:
		return fmt.Sprintf("%-72v", val)
	}
}

func cinc(label string, val interface{}) string {
	prefix := fmt.Sprintf("%10s", label)
	middle := color.HiWhiteString(":")
	ending := color.HiYellowString("%-72v", val)
	return fmt.Sprintf("%s %s %s", prefix, middle, ending)
}

func compileScriptCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return errors.New("must supply at least one gscript to this command")
	}
	cliLogger.Info(color.HiRedString("*** COMPILER OPTIONS ***"))
	cliLogger.Info("")
	cliLogger.Infof(copt("OS", defaultCompileOptions.OS))
	cliLogger.Infof(copt("Arch", defaultCompileOptions.Arch))
	cliLogger.Infof(copt("Output File", defaultCompileOptions.OutputFile))
	cliLogger.Infof(copt("Keep Build Directory", defaultCompileOptions.SaveBuildDir))
	cliLogger.Infof(copt("UPX Compression", defaultCompileOptions.UPXEnabled))
	cliLogger.Infof(copt("Logging Support", defaultCompileOptions.LoggingEnabled))
	cliLogger.Infof(copt("Debugger Support", defaultCompileOptions.DebuggerEnabled))
	cliLogger.Infof(copt("Human Redable Names", defaultCompileOptions.UseHumanReadableNames))
	cliLogger.Infof(copt("Import All Native Funcs", defaultCompileOptions.ImportAllNativeFuncs))
	cliLogger.Infof(copt("Skip Compilation", defaultCompileOptions.SkipCompilation))
	cliLogger.Infof(copt("Obfuscation Level", defaultCompileOptions.ObfuscationLevel))
	cliLogger.Info("")
	cliLogger.Info(color.HiRedString("***  SOURCE SCRIPTS  ***"))
	cliLogger.Info("")
	for _, a := range c.Args() {
		cliLogger.Info(cinc("Script", a))
	}
	cliLogger.Info("")
	cliLogger.Info(color.HiRedString("************************"))
	cliLogger.Info("")
	gc := compiler.NewWithOptions(defaultCompileOptions)
	gc.SetLogger(cliLogger)
	for _, a := range c.Args() {
		addErr := gc.AddScript(a)
		if addErr != nil {
			cliLogger.Errorf("Error adding script %s: %v", a, addErr)
		}
	}
	err := gc.Do()
	if err != nil {
		cliLogger.Errorf("Build Dir Located At: %s", gc.BuildDir)
		return err
	}
	cliLogger.Infof("Compiled binary located at:\n\n%s\n", gc.OutputFile)
	if defaultCompileOptions.SaveBuildDir {
		cliLogger.Infof("Build Dir Located At: %s", gc.BuildDir)
	}
	return nil
}
