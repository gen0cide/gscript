package compiler

import (
	"fmt"
	"regexp"

	gast "github.com/robertkrimen/otto/ast"
)

var (
	basicMacro    = `(?P<value>\S*)\z`
	goImportMacro = `(?P<gopkg>\S*?) as (?P<namespace>\w*)\z`

	validMacros = map[string]*regexp.Regexp{
		"priority":  generateMacroRegexp("priority", basicMacro),
		"os":        generateMacroRegexp("os", basicMacro),
		"timeout":   generateMacroRegexp("timeout", basicMacro),
		"arch":      generateMacroRegexp("arch", basicMacro),
		"version":   generateMacroRegexp("version", basicMacro),
		"go_import": generateMacroRegexp("go_import", goImportMacro),
	}

	// DefaultPriority is the default value for a script's priority during execution
	DefaultPriority = 100

	// DefaultTimeout is the default time in seconds a script will be allowed to run
	DefaultTimeout = 30
)

// Macro defines the object created by each macro parsed by the preprocessor
type Macro struct {
	// the genesis compiler option's name for this macro
	Key string

	// map of the values associated with this macro
	Params map[string]string
}

// ScanForMacros takes the genesis comments from the AST and parses known
// macros out for the compiler
func ScanForMacros(commentMap gast.CommentMap) []*Macro {
	macros := []*Macro{}
	for _, tmp := range commentMap {
		for _, comment := range tmp {
			for name, res := range validMacros {
				if res.MatchString(comment.Text) {
					keys := res.SubexpNames()
					matches := res.FindAllStringSubmatch(comment.Text, -1)[0]
					keymap := map[string]string{}
					mac := &Macro{}
					for i, n := range matches {
						if i < 2 {
							continue
						}
						keymap[keys[i]] = n
					}
					mac.Key = name
					mac.Params = keymap
					macros = append(macros, mac)
				}
			}
		}
	}
	return macros
}

// used to create the composable macro regular expressions at compile time
func generateMacroRegexp(name, valRegExp string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("\\A(?P<key>%s):%s", name, valRegExp))
}
