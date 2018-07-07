package main

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/gen0cide/gscript/compiler"
	"github.com/urfave/cli"
)

var (
	defaultCompilerOptions = compiler.DefaultOptions()
	compileCommand         = cli.Command{
		Name:      "compile",
		Usage:     "compiles the provided scripts using the genesis compiler",
		UsageText: "gscript compile [OPTIONS] SCRIPT [SCRIPT SCRIPT ...]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "os",
				Usage:       "operating system to target for native compilation",
				Value:       defaultCompilerOptions.OS,
				Destination: &defaultCompilerOptions.OS,
			},
			cli.StringFlag{
				Name:        "arch",
				Usage:       "architecture to target for native compilation",
				Value:       defaultCompilerOptions.Arch,
				Destination: &defaultCompilerOptions.Arch,
			},
			cli.StringFlag{
				Name:        "output-file, o",
				Usage:       "location to write final compiled binary",
				Value:       defaultCompilerOptions.OutputFile,
				Destination: &defaultCompilerOptions.OutputFile,
			},
			cli.BoolFlag{
				Name:        "keep-build-dir",
				Usage:       "keep the build directory of the genesis intermediate representation (default: false)",
				Destination: &defaultCompilerOptions.SaveBuildDir,
			},
			cli.BoolFlag{
				Name:        "enable-upx-compression",
				Usage:       "compress the final binary using UPX (default: false)",
				Destination: &defaultCompilerOptions.UPXEnabled,
			},
			cli.BoolFlag{
				Name:        "enable-logging",
				Usage:       "enable logging in the final binary for debugging (default: false)",
				Destination: &defaultCompilerOptions.LoggingEnabled,
			},
			cli.BoolFlag{
				Name:        "enable-debugging",
				Usage:       "enable the interactive debugger in the final binary (default: false)",
				Destination: &defaultCompilerOptions.DebuggerEnabled,
			},
			cli.BoolFlag{
				Name:        "enable-human-readable-names",
				Usage:       "use human readable names in the genesis intermediate representation (default: false)",
				Destination: &defaultCompilerOptions.UseHumanReadableNames,
			},
			cli.BoolFlag{
				Name:        "enable-import-all-native-funcs",
				Usage:       "link all possible native functions into the runtime, not just those called by the scripts (default: false)",
				Destination: &defaultCompilerOptions.ImportAllNativeFuncs,
			},
			cli.BoolFlag{
				Name:        "disable-native-compilation",
				Usage:       "do not compile the intermediate representation to a native binary (default: false)",
				Destination: &defaultCompilerOptions.SkipCompilation,
			},
			cli.IntFlag{
				Name:        "obfuscation-level",
				Usage:       "override the default obfuscation level, where argument can be 0-4 with 0 being full and 4 being none",
				Destination: &defaultCompilerOptions.ObfuscationLevel,
			},
		},
		Action: compileScriptCommand,
	}
)

func copt(label string, val interface{}) string {
	prefix := color.HiWhiteString("%25s", label)
	middle := color.HiWhiteString(":")
	ending := color.HiGreenString("%-72v", val)
	return fmt.Sprintf("%s %s %s", prefix, middle, ending)
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
	cliLogger.Info(copt("OS", defaultCompilerOptions.OS))
	cliLogger.Info(copt("Arch", defaultCompilerOptions.Arch))
	cliLogger.Info(copt("Output File", defaultCompilerOptions.OutputFile))
	cliLogger.Info(copt("Keep Build Directory", defaultCompilerOptions.SaveBuildDir))
	cliLogger.Info(copt("UPX Compression Enabled", defaultCompilerOptions.UPXEnabled))
	cliLogger.Info(copt("Logging Enabled", defaultCompilerOptions.LoggingEnabled))
	cliLogger.Info(copt("Debugger Enabled", defaultCompilerOptions.DebuggerEnabled))
	cliLogger.Info(copt("Human Redable Names", defaultCompilerOptions.UseHumanReadableNames))
	cliLogger.Info(copt("Import All Native Funcs", defaultCompilerOptions.ImportAllNativeFuncs))
	cliLogger.Info(copt("Skip Compilation", defaultCompilerOptions.SkipCompilation))
	cliLogger.Info(copt("Obfuscation Level", defaultCompilerOptions.ObfuscationLevel))
	cliLogger.Info("")
	cliLogger.Info(color.HiRedString("***  SOURCE SCRIPTS  ***"))
	cliLogger.Info("")
	for _, a := range c.Args() {
		cliLogger.Info(cinc("Script", a))
	}
	cliLogger.Info("")
	cliLogger.Info(color.HiRedString("************************"))
	cliLogger.Info("")
	gc := compiler.NewWithOptions(defaultCompilerOptions)
	gc.SetLogger(cliLogger)
	for _, a := range c.Args() {
		gc.AddScript(a)
	}
	err := gc.Do()
	if err != nil {
		return err
	}
	cliLogger.Infof("Compiled binary located at:\n\n%s\n", gc.OutputFile)
	return nil
}
