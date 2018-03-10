package engine

import (
	"io/ioutil"

	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
)

type Engine struct {
	VM              *otto.Otto
	Logger          *logrus.Logger
	Imports         map[string]func() []byte
	Name            string
	DebuggerEnabled bool
	Timeout         int
	Halted          bool
}

func New(name string) *Engine {
	logger := logrus.New()
	logger.Formatter = new(logrus.TextFormatter)
	logger.Out = ioutil.Discard
	return &Engine{
		Name:            name,
		Imports:         map[string]func() []byte{},
		Logger:          logger,
		DebuggerEnabled: false,
		Halted:          false,
		Timeout:         30,
	}
}

func (e *Engine) SetLogger(logger *logrus.Logger) {
	e.Logger = logger
}

func (e *Engine) SetTimeout(timeout int) {
	e.Timeout = timeout
}

func (e *Engine) AddImport(name string, data func() []byte) {
	e.Imports[name] = data
}

func (e *Engine) SetName(name string) {
	e.Name = name
}
