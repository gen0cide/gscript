package engine

import (
	"os"
	"strconv"
	"strings"

	"github.com/djherbis/times"
	"github.com/robertkrimen/otto"
)

func (e *Engine) VMGetDirsInPath(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMCanWriteFile(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
	//Following code breaks building for windows :(
	/*
	   filePath := call.Argument(0)
	   filePathString, err := filePath.Export()
	   if err != nil {
	     e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
	     return otto.FalseValue()
	   }
	   result := LocalFileWritable(filePathString.(string))
	   if result == true {
	     return otto.TrueValue()
	   } else {
	     return otto.FalseValue()
	   }
	*/
}

func (e *Engine) VMCanExecFile(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
	//Following code breaks building for windows :(
	/*
	   filePath := call.Argument(0)
	   filePathString, err := filePath.Export()
	   if err != nil {
	     e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
	     return otto.FalseValue()
	   }
	   result := LocalFileExecutable(filePathString.(string))
	   if result {
	     return otto.TrueValue()
	   } else {
	     return otto.FalseValue()
	   }
	*/
}

func (e *Engine) VMFileExists(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	if LocalFileExists(filePathString.(string)) {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMDirExists(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	if LocalFileExists(filePathString.(string)) {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMFileContains(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	search := call.Argument(1)
	searchString, err := search.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	fileData, err := LocalFileRead(filePathString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return otto.FalseValue()
	}
	fileStrings := string(fileData)
	if strings.Contains(fileStrings, searchString.(string)) {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMCanReadFile(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	data, err := LocalFileRead(filePathString.(string))
	if data != nil && err == nil {
		return otto.TrueValue()
	} else if err != nil {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return otto.FalseValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMDeleteFile(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	err = LocalFileDelete(filePathAsString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Error deleting the file: %s", err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMWriteFile(call otto.FunctionCall) otto.Value {
	filePath, err := call.Argument(0).ToString()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	fileData := call.Argument(1)
	fileMode, err := call.Argument(2).ToString()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	fileBytes := e.ValueToByteSlice(fileData)
	err = LocalFileCreate(filePath, fileBytes, fileMode)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Error writing the file: %s", err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMReadFile(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	bytes, err := LocalFileRead(filePathAsString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Error reading the file: %s", err.Error())
		return otto.FalseValue()
	}

	vmResponse, err := e.VM.ToValue(string(bytes))

	return vmResponse
}

func (e *Engine) VMCopyFile(call otto.FunctionCall) otto.Value {
	readPath, err := call.ArgumentList[0].ToString()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	writePath, err := call.ArgumentList[1].ToString()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	bytes, err := LocalFileRead(readPath)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Error reading the file: %s", err.Error())
		return otto.FalseValue()
	}
	filePerms, err := os.Stat(readPath)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("OS Error: %s", err.Error())
		return otto.FalseValue()
	}
	err = LocalFileCreate(writePath, bytes, strconv.Itoa(int(filePerms.Mode())))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Error writing the file: %s", err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMExecuteFile(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.ToString()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	cmdArgs := call.Argument(1)
	argList := []string{}
	if !cmdArgs.IsNull() {
		argArray, err := cmdArgs.Export()
		if err != nil {
			e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
			return otto.FalseValue()
		}
		argList = argArray.([]string)
	}

	cmdOutput := ExecuteCommand(filePathAsString, argList...)
	vmResponse, err := e.VM.ToValue(cmdOutput)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMAppendFile(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	fileData := call.Argument(1)
	fileBytes := e.ValueToByteSlice(fileData)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	err = LocalFileAppendBytes(filePathAsString.(string), fileBytes)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Error appending the file: %s", err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMReplaceInFile(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	sFind := call.Argument(1)
	sFindAsString, err := sFind.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	sReplace := call.Argument(2)
	sReplaceAsString, err := sReplace.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	err = LocalFileReplace(filePathAsString.(string), sFindAsString.(string), sReplaceAsString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Error editing the file: %s", err.Error())
		return otto.FalseValue()
	}
	return otto.TrueValue()
}

func (e *Engine) VMFileChangeTime(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	t, err := times.Stat(filePathAsString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return otto.FalseValue()
	}
	if t.HasChangeTime() {
		vmResponse, err := e.VM.ToValue(t.ChangeTime().String())
		if err != nil {
			e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
			return otto.FalseValue()
		}
		return vmResponse
	} else {
		e.Logger.WithField("trace", "true").Errorf("File error: no ctime")
		return otto.FalseValue()
	}
}

func (e *Engine) VMFileBirthTime(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	t, err := times.Stat(filePathAsString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return otto.FalseValue()
	}
	if t.HasBirthTime() {
		vmResponse, err := e.VM.ToValue(t.BirthTime().String())
		if err != nil {
			e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
			return otto.FalseValue()
		}
		return vmResponse
	} else {
		e.Logger.WithField("trace", "true").Errorf("File error: no ctime")
		return otto.FalseValue()
	}
}

func (e *Engine) VMFileModifyTime(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	t, err := times.Stat(filePathAsString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return otto.FalseValue()
	}
	vmResponse, err := e.VM.ToValue(t.ModTime().String())
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMFileAccessTime(call otto.FunctionCall) otto.Value {
	filePath := call.Argument(0)
	filePathAsString, err := filePath.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	t, err := times.Stat(filePathAsString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
		return otto.FalseValue()
	}
	vmResponse, err := e.VM.ToValue(t.AccessTime().String())
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return otto.FalseValue()
	}
	return vmResponse
}
