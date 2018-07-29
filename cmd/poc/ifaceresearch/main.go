package main

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/Jeffail/gabs"

	"github.com/gen0cide/gscript/compiler"

	"github.com/davecgh/go-spew/spew"
	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var (
	logger *logrus.Logger
	vm     *otto.Otto
)

const testval = `
function fappy(e) {
	return 5;
}

var a = {
	"Content-Type": "application/json",
	"Happy": 1337,
	"Robert": function(e) {
		return 5;
	},
	"Becky": fappy
}

var x = {}
`

func init() {
	logger = logrus.New()
	logger.Formatter = new(prefixed.TextFormatter)
	logger.SetLevel(logrus.DebugLevel)
	vm = otto.New()
	vm.Eval(testval)
}

func test(name string, v otto.Value) {
	isObj := v.IsObject()
	logger.Infof("%s is an object: %v", name, isObj)
	isPrim := v.IsPrimitive()
	logger.Infof("%s is a primative: %v", name, isPrim)
	val, err := v.Export()
	if err != nil {
		panic(err)
	}
	ttype(name, val)
}

func ttype(name string, val interface{}) {
	rval := reflect.ValueOf(val)
	rtype := rval.Type()
	logger.Infof("%s name is %s", name, rtype.Name())
	logger.Infof("%s string is %s", name, rtype.String())
	logger.Infof("%s kind is %s", name, rtype.Kind().String())
	if rtype.Kind() == reflect.Map {
		logger.Infof("%s is a map", name)
		real, ok := val.(map[string]interface{})
		if !ok {
			logger.Infof("%s is not a map[string]interface{}", name)
		} else {
			logger.Infof("%s is a map[string]interface{}", name)
			jsonObj, err := gabs.Consume(val)
			if err != nil {
				logger.Errorf("ERROR IN JSON SERIALIZATION OF %s: %v", name, err)
			}
			if jsonObj != nil {
				logger.Infof("JSON VERSION of %s:\n%s\n", name, jsonObj.StringIndent("", "  "))
			} else {
				logger.Errorf("Could not produce JSON for %s", name)
			}
			hedr, ok2 := val.(http.Header)
			if !ok2 {
				logger.Infof("%s is not an http.Header", name)
				newHeaders := map[string]string{}
				for k, item := range real {
					ttype(fmt.Sprintf("%s[%s]", name, k), item)
					newHeaders[k] = fmt.Sprintf("%v", item)
				}
				spew.Dump(newHeaders)
			} else {
				logger.Infof("%s is a http.Header", name)
				spew.Dump(hedr)
			}
		}
	}
}

func testvmret(name string) {
	v, err := vm.Get(name)
	if err != nil {
		panic(err)
	}
	test(name, v)
}

func main() {
	b := []string{
		"hello",
		"world",
	}
	vm.Set("b", b)

	c := compiler.NewWithDefault()
	vm.Set("c", c)

	d := map[string]interface{}{
		"Content-Type": []string{
			"English",
			"Mother",
		},
	}
	vm.Set("d", d)

	e := map[string]interface{}{}

	e["dog"] = "house"
	e["barf"] = []string{"bobz"}
	e["george"] = &http.Transport{}
	vm.Set("e", e)

	testvmret("a")
	testvmret("b")
	testvmret("c")
	testvmret("d")
	testvmret("e")
	testvmret("x")
}
