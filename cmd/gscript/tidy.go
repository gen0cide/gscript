package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ditashi/jsbeautifier-go/jsbeautifier"

	"github.com/urfave/cli/v2"
)

var (
	tidyCommand = &cli.Command{
		Name:      "tidy",
		Usage:     "tidy the syntax of a provided gscript",
		UsageText: "gscript tidy GSCRIPT",
		Action:    tidyScriptCommand,
	}

	defaultTidyOptions = map[string]interface{}{
		"indent_size":               2,
		"indent_char":               " ",
		"indent_with_tabs":          false,
		"preserve_newlines":         true,
		"max_preserve_newlines":     10,
		"space_in_paren":            false,
		"space_in_empty_paren":      false,
		"e4x":                       false,
		"jslint_happy":              false,
		"space_after_anon_function": false,
		"brace_style":               "collapse",
		"keep_array_indentation":    false,
		"keep_function_indentation": false,
		"eval_code":                 false,
		"unescape_strings":          false,
		"wrap_line_length":          0,
		"break_chained_methods":     false,
		"end_with_newline":          true,
	}
)

func tidyScriptCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return errors.New("must supply a gscript to this command")
	}
	target := c.Args().First()
	if _, err := os.Stat(target); err != nil {
		return err
	}
	data, err := ioutil.ReadFile(target)
	if err != nil {
		return err
	}
	dataString := string(data)

	beautified, err := jsbeautifier.Beautify(&dataString, defaultTidyOptions)
	if err != nil {
		return err
	}

	fmt.Println(beautified)
	return nil
}
