package main

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"unsafe"

	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"
)

var (
	to    = testObj{Name: "JOSE"}
	tests = map[string]interface{}{
		"bool":       true,
		"byte":       byte(42),
		"complex128": complex(float64(8), float64(27)),
		"complex64":  complex(float32(6), float32(7)),
		"error":      errors.New("test error"),
		"float32":    math.Float32frombits(uint32(42949600)),
		"float64":    math.Float64frombits(uint64(18446744070009551615)),
		"int":        int(-9023372036854770000),
		"int8":       int8(-112),
		"int16":      int16(-30103),
		"int32":      int32(-2147402248),
		"int64":      int64(-9223372036854220481),
		"rune":       '本',
		"string":     "日本語",
		"uint":       uint(13446744073709551600),
		"uint8":      uint8(98),
		"uint16":     uint16(59228),
		"uint32":     uint32(4214967495),
		"uint64":     uint64(10446744073709551615),
		"uintptr":    ptr(to),
		"nil":        nil,
	}
	logger   = logrus.New()
	goodtest = map[string]*mapper{}
	badtest  = map[string]*mapper{}
	alltests = map[string]*mapper{}
)

type testObj struct {
	Name string
}

type mapper struct {
	Label             string
	TestVal           interface{}
	RetVal            interface{}
	GoBeginType       reflect.Value
	GoFuncArgType     reflect.Value
	JSGetRetType      otto.Value
	JSFuncExpType     otto.Value
	JSFuncRetType     otto.Value
	JSArgType         otto.Value
	GoReturnType      reflect.Value
	SuccessfulFuncArg bool
	SuccessfulFuncRet bool
	VM                *otto.Otto
	Log               *logrus.Entry
}

func newMapper(label string, tv interface{}, l *logrus.Logger) *mapper {
	m := &mapper{
		Label:       label,
		TestVal:     tv,
		VM:          otto.New(),
		GoBeginType: reflect.ValueOf(tv),
		Log:         l.WithField("test", label),
	}
	m.VM.Set("TypeAssert", m.vmTestFunc)
	return m
}

func ptr(i interface{}) unsafe.Pointer {
	return unsafe.Pointer(&i)
}

func main() {
	newTests := map[string]interface{}{}
	for name, item := range tests {
		newTests[name] = item
		sliceName := fmt.Sprintf("[]%s", name)
		newItem := slicemaker(item)
		newTests[sliceName] = newItem
	}
	for name, item := range tests {
		mapName := fmt.Sprintf("map[string]%s", name)
		newItem := mapmaker(item)
		newTests[mapName] = newItem
	}
	for name, item := range tests {
		mapSliceName := fmt.Sprintf("map[string][]%s", name)
		newItem := slicemaker(item)
		newMapItem := mapmaker(newItem)
		newTests[mapSliceName] = newMapItem
	}
	for name, item := range newTests {
		m := newMapper(name, item, logger)
		alltests[name] = m
		m.Log.Infof("##### TEST SUITE [%s] #####", name)
		m.typeinfo("(init)", item)
		m.Log.Info(">>> SETTING JAVASCRIPT VARIABLE")
		err := m.VM.Set("testval", item)
		if err != nil {
			m.Log.WithError(err).Errorf("coult not SET testval")
			continue
		}
		retval, err := m.VM.Get("testval")
		if err != nil {
			m.Log.WithError(err).Errorf("coult not GET testval")
			continue
		}
		m.jstypeinfo("(js decl)", retval)
		m.JSGetRetType = retval
		m.Log.Info(">>> EXECUTING TYPE ASSERT FUNCTION")
		m.JSFuncRetType, err = m.VM.Eval("TypeAssert(testval)")
		if err != nil {
			m.Log.WithError(err).Errorf("coult not EVAL testval")
			continue
		}
		m.Log.Info(">>> CHECKING FUNCTION RETURN")
		m.jstypeinfo("(function ret - js)", m.JSFuncRetType)
		goret, err := m.JSFuncRetType.Export()
		if err != nil {
			m.Log.WithError(err).Errorf("coult not EXPORT function return val")
		}
		m.typeinfo("(function ret - go)", goret)
		m.SuccessfulFuncRet = reflect.DeepEqual(m.TestVal, goret)
		m.Log.Infof("ARG TYPE ASSERTION IS THE SAME: %v", m.SuccessfulFuncArg)
		if m.SuccessfulFuncArg && m.SuccessfulFuncRet {
			goodtest[name] = m
			continue
		}
		badtest[name] = m
	}
	logger.Info("############## RESULTS ##############")
	for n, m := range alltests {
		res := "** FAILED **"
		if m.SuccessfulFuncArg && m.SuccessfulFuncRet {
			res = "PASSED"
		}
		logger.Infof("  %40s: [%s] ARG=%v RET=%v", n, res, m.SuccessfulFuncArg, m.SuccessfulFuncRet)
	}
	return
}

func (m *mapper) typeinfo(phase string, i interface{}) {
	v := reflect.ValueOf(i)
	m.Log.Infof(">>> GOLANG TYPE DETAILS FOR PHASE - %s", phase)
	m.Log.Infof("         type: %T", i)
	if i == nil {
		return
	}
	m.Log.Infof("    type.name: %s", v.Type().Name())
	m.Log.Infof("  type.string: %s", v.Type().String())
	m.Log.Infof("    type.kind: %T", v.Type().Kind().String())
}

func (m *mapper) jstypeinfo(phase string, v otto.Value) {
	m.Log.Infof(">>> JS TYPE DETAILS FOR PHASE - %s", phase)
	m.Log.Infof("        v.class: %s", v.Class())
	m.Log.Infof("       v.string: %s", v.String())
	//m.Log.Infof("           v.is: boolean=%v defined=%v function=%v nan=%v null=%v number=%v object=%v primative=%v string=%v undefined=%v", v.IsBoolean(), v.IsDefined(), v.IsFunction(), v.IsNaN(), v.IsNull(), v.IsNumber(), v.IsObject(), v.IsPrimitive(), v.IsString(), v.IsUndefined())
	bret, err := v.ToBoolean()
	m.Log.Infof("    v.toboolean: %v (%v)", bret, (err == nil))
	fret, err := v.ToFloat()
	m.Log.Infof("      v.tofloat: %v (%v)", fret, (err == nil))
	iret, err := v.ToInteger()
	m.Log.Infof("    v.toboolean: %v (%v)", iret, (err == nil))
	sret, err := v.ToString()
	m.Log.Infof("    v.toboolean: %v (%v)", sret, (err == nil))
}

func (m *mapper) vmTestFunc(call otto.FunctionCall) otto.Value {
	arg := call.Argument(0)
	m.jstypeinfo("(function arg - js)", arg)
	m.JSArgType = arg
	goarg, err := arg.Export()
	if err != nil {
		m.Log.WithError(err).Errorf("coult not export arg to golang")
		return call.Otto.MakeTypeError(err.Error())
	}
	m.typeinfo("(function arg - go)", goarg)
	m.GoReturnType = reflect.ValueOf(goarg)
	m.SuccessfulFuncArg = reflect.DeepEqual(m.TestVal, goarg)
	m.Log.Infof("ARG TYPE ASSERTION IS THE SAME: %v", m.SuccessfulFuncArg)
	finalret, err := call.Otto.ToValue(m.TestVal)
	if err != nil {
		m.Log.WithError(err).Errorf("coult not export golang to return val")
		return call.Otto.MakeTypeError(err.Error())
	}
	m.jstypeinfo("(function exp - js)", finalret)
	m.JSFuncExpType = finalret
	return finalret
}

func slicemaker(i interface{}) []interface{} {
	slc := []interface{}{}
	for x := 0; x < 5; x++ {
		slc = append(slc, i)
	}
	return slc
}

func mapmaker(i interface{}) map[string]interface{} {
	slc := map[string]interface{}{}
	for x := 0; x < 5; x++ {
		slc[randstr.Hex(2)] = i
	}
	return slc
}
