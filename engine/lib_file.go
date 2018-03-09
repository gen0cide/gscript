package engine

import (
	"io/ioutil"
	"os"
)

// func (e *Engine) VMGetDirsInPath(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// }

// func (e *Engine) VMCanWriteFile(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// 	//Following code breaks building for windows :(
// 	/*
// 	   filePath := call.Argument(0)
// 	   filePathString, err := filePath.Export()
// 	   if err != nil {
// 	     e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
// 	     return otto.FalseValue()
// 	   }
// 	   result := LocalFileWritable(filePathString.(string))
// 	   if result == true {
// 	     return otto.TrueValue()
// 	   } else {
// 	     return otto.FalseValue()
// 	   }
// 	*/
// }

// func (e *Engine) VMCanExecFile(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// 	//Following code breaks building for windows :(
// 	/*
// 	   filePath := call.Argument(0)
// 	   filePathString, err := filePath.Export()
// 	   if err != nil {
// 	     e.Logger.WithField("trace", "true").Errorf("Function Error: error=ARY_ARG_NOT_String arg=%s", err.Error())
// 	     return otto.FalseValue()
// 	   }
// 	   result := LocalFileExecutable(filePathString.(string))
// 	   if result {
// 	     return otto.TrueValue()
// 	   } else {
// 	     return otto.FalseValue()
// 	   }
// 	*/
// }

// func (e *Engine) VMFileExists(call otto.FunctionCall) otto.Value {
// 	filePath := call.Argument(0)
// 	filePathString, err := filePath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	if LocalFileExists(filePathString.(string)) {
// 		return otto.TrueValue()
// 	} else {
// 		return otto.FalseValue()
// 	}
// }

// func (e *Engine) VMDirExists(call otto.FunctionCall) otto.Value {
// 	filePath := call.Argument(0)
// 	filePathString, err := filePath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	if LocalFileExists(filePathString.(string)) {
// 		return otto.TrueValue()
// 	} else {
// 		return otto.FalseValue()
// 	}
// }

// func (e *Engine) VMFileContains(call otto.FunctionCall) otto.Value {
// 	filePath := call.Argument(0)
// 	filePathString, err := filePath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	search := call.Argument(1)
// 	searchString, err := search.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	fileData, err := LocalFileRead(filePathString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	fileStrings := string(fileData)
// 	if strings.Contains(fileStrings, searchString.(string)) {
// 		return otto.TrueValue()
// 	} else {
// 		return otto.FalseValue()
// 	}
// }

// func (e *Engine) VMCanReadFile(call otto.FunctionCall) otto.Value {
// 	filePath := call.Argument(0)
// 	filePathString, err := filePath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	data, err := LocalFileRead(filePathString.(string))
// 	if data != nil && err == nil {
// 		return otto.TrueValue()
// 	} else if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
// 		return otto.FalseValue()
// 	} else {
// 		return otto.FalseValue()
// 	}
// }

// func (e *Engine) VMDeleteFile(call otto.FunctionCall) otto.Value {
// 	filePath := call.Argument(0)
// 	filePathAsString, err := filePath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	err = LocalFileDelete(filePathAsString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Error deleting the file: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return otto.TrueValue()
// }

func (e *Engine) WriteFile(path string, fileData []byte, perms int64) (int, error) {
	err := ioutil.WriteFile(path, fileData, os.FileMode(uint32(perms)))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Error writing the file: %s", err.Error())
		return 0, err
	}
	return len(fileData), nil
}

// func (e *Engine) VMReadFile(call otto.FunctionCall) otto.Value {
// 	filePath := call.Argument(0)
// 	filePathAsString, err := filePath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	bytes, err := LocalFileRead(filePathAsString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Error reading the file: %s", err.Error())
// 		return otto.FalseValue()
// 	}

// 	vmResponse, err := e.VM.ToValue(string(bytes))

// 	return vmResponse
// }

// func (e *Engine) VMCopyFile(call otto.FunctionCall) otto.Value {
// 	readPath, err := call.ArgumentList[0].ToString()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	writePath, err := call.ArgumentList[1].ToString()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	bytes, err := LocalFileRead(readPath)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Error reading the file: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	filePerms, err := os.Stat(readPath)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("OS Error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	err = LocalFileCreate(writePath, bytes, int(filePerms.Mode()))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Error writing the file: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return otto.TrueValue()
// }

// func (e *Engine) VMExecuteFile(call otto.FunctionCall) otto.Value {
// 	filePath := call.Argument(0)
// 	filePathAsString, err := filePath.ToString()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	cmdArgs := call.Argument(1)
// 	argList := []string{}
// 	if !cmdArgs.IsNull() {
// 		argArray, err := cmdArgs.Export()
// 		if err != nil {
// 			e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 			return otto.FalseValue()
// 		}
// 		argList = argArray.([]string)
// 	}

// 	cmdOutput := ExecuteCommand(filePathAsString, argList...)
// 	vmResponse, err := e.VM.ToValue(cmdOutput)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return vmResponse
// }

// func (e *Engine) VMAppendFile(call otto.FunctionCall) otto.Value {
// 	filePath := call.Argument(0)
// 	fileData := call.Argument(1)
// 	fileBytes := e.ValueToByteSlice(fileData)
// 	filePathAsString, err := filePath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	err = LocalFileAppendBytes(filePathAsString.(string), fileBytes)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Error appending the file: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return otto.TrueValue()
// }

// func (e *Engine) VMReplaceInFile(call otto.FunctionCall) otto.Value {
// 	filePath := call.Argument(0)
// 	filePathAsString, err := filePath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	sFind := call.Argument(1)
// 	sFindAsString, err := sFind.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	sReplace := call.Argument(2)
// 	sReplaceAsString, err := sReplace.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	err = LocalFileReplace(filePathAsString.(string), sFindAsString.(string), sReplaceAsString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Error editing the file: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return otto.TrueValue()
// }

// func (e *Engine) VMFileChangeTime(call otto.FunctionCall) otto.Value {
// 	filePath := call.Argument(0)
// 	filePathAsString, err := filePath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	t, err := times.Stat(filePathAsString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	if t.HasChangeTime() {
// 		vmResponse, err := e.VM.ToValue(t.ChangeTime().String())
// 		if err != nil {
// 			e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
// 			return otto.FalseValue()
// 		}
// 		return vmResponse
// 	} else {
// 		e.Logger.WithField("trace", "true").Errorf("File error: no ctime")
// 		return otto.FalseValue()
// 	}
// }

// func (e *Engine) VMFileBirthTime(call otto.FunctionCall) otto.Value {
// 	filePath := call.Argument(0)
// 	filePathAsString, err := filePath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	t, err := times.Stat(filePathAsString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	if t.HasBirthTime() {
// 		vmResponse, err := e.VM.ToValue(t.BirthTime().String())
// 		if err != nil {
// 			e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
// 			return otto.FalseValue()
// 		}
// 		return vmResponse
// 	} else {
// 		e.Logger.WithField("trace", "true").Errorf("File error: no ctime")
// 		return otto.FalseValue()
// 	}
// }

// func (e *Engine) VMFileModifyTime(call otto.FunctionCall) otto.Value {
// 	filePath := call.Argument(0)
// 	filePathAsString, err := filePath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	t, err := times.Stat(filePathAsString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	vmResponse, err := e.VM.ToValue(t.ModTime().String())
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return vmResponse
// }

// func (e *Engine) VMFileAccessTime(call otto.FunctionCall) otto.Value {
// 	filePath := call.Argument(0)
// 	filePathAsString, err := filePath.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	t, err := times.Stat(filePathAsString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("File error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	vmResponse, err := e.VM.ToValue(t.AccessTime().String())
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return vmResponse
// }

// func LocalFileExists(path string) bool {
// 	_, err := os.Stat(path)
// 	if err == nil {
// 		return true
// 	}
// 	return false
// }

// //These two functions break compiling on windows
// /*
// func LocalFileWritable(path string) bool {
//   return unix.Access(path, unix.W_OK) == nil
// }

// func LocalFileExecutable(path string) bool {
//   return unix.Access(path, unix.X_OK) == nil
// }
// */

// func LocalDirCreate(path string) error {
// 	err := os.MkdirAll(path, 0700)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func LocalDirRemoveAll(dir string) error {
// 	d, err := os.Open(dir)
// 	if err != nil {
// 		return err
// 	}
// 	defer d.Close()
// 	names, err := d.Readdirnames(-1)
// 	if err != nil {
// 		return err
// 	}
// 	for _, name := range names {
// 		err = os.RemoveAll(filepath.Join(dir, name))
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	err = os.RemoveAll(dir)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func LocalFileDelete(path string) error {
// 	if LocalFileExists(path) {
// 		err := os.Remove(path)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	return errors.New("The file dosn't exist to delete")
// }

// func LocalFileCreate(path string, bytes []byte, perms int) error {
// 	if LocalFileExists(path) {
// 		return errors.New("The file to create already exists so we won't overwite it")
// 	}
// 	var p os.FileMode
// 	// pInt, err := strconv.Atoi(perms)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	p = os.FileMode(perms)
// 	spew.Dump(perms)
// 	spew.Dump(p)
// 	err := ioutil.WriteFile(path, bytes, p)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func LocalFileAppendBytes(filename string, bytes []byte) error {
// 	if LocalFileExists(filename) {
// 		fileInfo, err := os.Stat(filename)
// 		if err != nil {
// 			return err
// 		}
// 		file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, fileInfo.Mode())
// 		if err != nil {
// 			return err
// 		}
// 		if _, err = file.Write(bytes); err != nil {
// 			return err
// 		}
// 		file.Close()
// 		return nil
// 	}
// 	err := LocalFileCreate(filename, bytes, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func LocalFileAppendString(input, filename string) error {
// 	fileInfo, err := os.Stat(filename)
// 	if err != nil {
// 		return err
// 	}
// 	file, err := os.OpenFile(filename, os.O_APPEND, fileInfo.Mode())
// 	if err != nil {
// 		return err
// 	}
// 	if _, err = file.WriteString(input); err != nil {
// 		return err
// 	}
// 	file.Close()
// 	return nil
// }

// func LocalFileReplace(file, match, replacement string) error {
// 	if LocalFileExists(file) {
// 		fileInfo, err := os.Stat(file)
// 		if err != nil {
// 			return err
// 		}
// 		contents, err := ioutil.ReadFile(file)
// 		if err != nil {
// 			return err
// 		}
// 		lines := strings.Split(string(contents), "\n")
// 		for index, line := range lines {
// 			if strings.Contains(line, match) {
// 				lines[index] = strings.Replace(line, match, replacement, 10)
// 			}
// 		}

// 		ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), fileInfo.Mode())
// 		return nil
// 	} else {
// 		return errors.New("The file to read does not exist")
// 	}
// }

// func LocalFileReplaceMulti(file string, matches []string, replacement string) error {
// 	if LocalFileExists(file) {
// 		fileInfo, err := os.Stat(file)
// 		if err != nil {
// 			return err
// 		}
// 		contents, err := ioutil.ReadFile(file)
// 		if err != nil {
// 			return err
// 		}
// 		lines := strings.Split(string(contents), "\n")
// 		for index, line := range lines {
// 			for _, match := range matches {
// 				if strings.Contains(line, match) {
// 					lines[index] = replacement
// 				}
// 			}
// 		}
// 		ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), fileInfo.Mode())
// 		return nil
// 	} else {
// 		return errors.New("The file to read does not exist")
// 	}
// }

// func LocalFileRead(path string) ([]byte, error) {
// 	if LocalFileExists(path) {
// 		dat, err := ioutil.ReadFile(path)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return dat, nil
// 	}
// 	return nil, errors.New("The file to read does not exist")
// }

// // func XorFiles(file1 string, file2 string, outPut string) error {
// // 	dat1, err := ioutil.ReadFile(file1)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	dat2, err := ioutil.ReadFile(file2)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	dat3 := XorBytes(dat1[:], dat2[:])
// // 	err = LocalFileCreate(outPut, dat3[:], 0644)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	return nil
// // }

// func LocalCopyFile(src, dst string) error {
// 	from, err := os.Open(src)
// 	if err != nil {
// 		return err
// 	}
// 	defer from.Close()

// 	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
// 	if err != nil {
// 		return err
// 	}
// 	defer to.Close()

// 	_, err = io.Copy(to, from)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
