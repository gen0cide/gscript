package engine

import (
	"io/ioutil"
	"os"
	"os/user"
	"runtime"
	"strings"

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
	Paused          bool
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
		Paused:          false,
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

func (e *Engine) LoadScript(source []byte) error {
	_, err := e.VM.Run(string(source))
	return err
}

func (e *Engine) currentUser() map[string]string {
	userInfo := map[string]string{}
	u, err := user.Current()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("User Loading Error: %s", err.Error())
		return userInfo
	}
	userInfo["uid"] = u.Uid
	userInfo["gid"] = u.Gid
	userInfo["username"] = u.Username
	userInfo["name"] = u.Name
	userInfo["home_dir"] = u.HomeDir
	groups, err := u.GroupIds()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Group Loading Error: %s", err.Error())
		return userInfo
	}
	userInfo["groups"] = strings.Join(groups, ":")
	return userInfo
}

func (e *Engine) injectVars() {
	userInfo, err := e.VM.ToValue(e.currentUser())
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Could not inject user info into VM: %s", err.Error())
	} else {
		e.VM.Set("USER_INFO", userInfo)
	}
	osVal, err := e.VM.ToValue(runtime.GOOS)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Could not inject os info into VM: %s", err.Error())
	} else {
		e.VM.Set("OS", osVal)
	}
	hn, err := os.Hostname()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Could not obtain hostname info: %s", err.Error())
	} else {
		hostnameVal, err := e.VM.ToValue(hn)
		if err != nil {
			e.Logger.Errorf("Could not inject hostname info into VM: %s", err.Error())
		} else {
			e.VM.Set("HOSTNAME", hostnameVal)
		}
	}
	archVal, err := e.VM.ToValue(runtime.GOARCH)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Could not inject arch info into VM: %s", err.Error())
	} else {
		e.VM.Set("ARCH", archVal)
	}
	ipVals, err := e.VM.ToValue(getLocalIPs())
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Could not inject ip info into VM: %s", err.Error())
	} else {
		e.VM.Set("IP_ADDRS", ipVals)
	}
}
