package main

import (
	"os"

	"github.com/happierall/l"
)

type VType struct {
	VarName string `yaml:"var"`
	Type    string `yaml:"type"`
	IsArray bool   `yaml:"is_array"`
	Desc    string `yaml:"desc"`
}

type Function struct {
	Method  string  `yaml:"method"`
	Desc    string  `yaml:"desc"`
	Args    []VType `yaml:"args"`
	Rets    []VType `yaml:"returns"`
	Example string  `yaml:"example"`
}

func main() {
	logger := l.New()
	logger.DisabledInfo = false
	if len(os.Args) < 2 {
		logger.Crit("You did not specify a YAML configuration file.")
	}

}
