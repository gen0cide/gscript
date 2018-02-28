package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var g_file_1 = fmt.Sprintf("/tmp/%s", RandString(6))
var g_file_2 = fmt.Sprintf("/tmp/%s", RandString(6))
var g_file_3 = fmt.Sprintf("/tmp/%s", RandString(6))

func TestVMCopyFile(t *testing.T) {
	file_2 := g_file_1
	testScript := fmt.Sprintf(`
        var file_1 = "/etc/passwd";
        var file_2 = "%s";
        var return_value = CopyFile(file_1, file_2);
  `, file_2)

	e := New("CopyFile")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString)
}

func TestVMAppendFile(t *testing.T) {
	bytes := "60,104,116,109,108,62,10,32,32,60,98,111,100,121,62,10,32,32,32,32"
	testScript := fmt.Sprintf(`
    var file_1 = "%s";
    var file_2 = "%s";
    var bytes = [%s];
    var return_value1 = AppendFile(file_1, bytes);
    var return_value2 = AppendFile(file_2, bytes);
  `, g_file_1, g_file_2, bytes)

	e := New("AppendFile")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value1")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString)
	retVal2, err := e.VM.Get("return_value2")
	assert.Nil(t, err)
	retValAsString2, err := retVal2.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString2)
}

func TestVMReplaceInFile(t *testing.T) {
	string01 := "root"
	string02 := "lol"
	testScript := fmt.Sprintf(`
    var file_1 = "%s";
    var string01 = "%s";
    var string02 = "%s";
    var return_value1 = ReplaceInFile(file_1, string01, string02);
  `, g_file_1, string01, string02)

	e := New("ReplaceInFile")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value1")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString)
}

func TestVMRetrieveFileFromURL(t *testing.T) {
	url := "https://alexlevinson.com/"
	file_3 := g_file_3
	testScript2 := fmt.Sprintf(`
    var url = "%s";
    var file_3 = "%s";
    var response2 = RetrieveFileFromURL(url);
    var return_value2 = response2;
    var response3 = WriteFile(file_3, response2);
  `, url, file_3)
	e := New("RetrieveFileFromURL")
	e.CreateVM()

	e.VM.Run(testScript2)
	retVal, err := e.VM.Get("return_value2")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "60,104,116,109,108,62,10,32,32,60,98,111,100,121,62,10,32,32,32,32,60,99,101,110,116,101,114,62,10,32,32,32,32,32,32,60,105,109,103,32,115,114,99,61,34,114,111,111,116,46,106,112,103,34,32,47,62,10,32,32,32,32,60,47,99,101,110,116,101,114,62,10,32,32,60,47,98,111,100,121,62,10,60,47,104,116,109,108,62,10", retValAsString)
}

func TestVMDeleteFile(t *testing.T) {
	testScript := fmt.Sprintf(`
    var file_1 = "%s";
    var return_value1 = DeleteFile(file_1);
        var file_2 = "%s";
        var return_value2 = DeleteFile(file_2);
    var file_3 = "%s";
        var return_value3 = DeleteFile(file_3);
    `, g_file_1, g_file_2, g_file_3)

	e := New("DeleteFile")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value1")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString)
	retVal2, err := e.VM.Get("return_value2")
	assert.Nil(t, err)
	retValAsString2, err := retVal2.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString2)
	retVal3, err := e.VM.Get("return_value3")
	assert.Nil(t, err)
	retValAsString3, err := retVal3.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString3)
}

func TestVMExecuteFile(t *testing.T) {
	testScript := `
    var file_path = "uname";
    var args = ["-o"];
          var results = ExecuteFile(file_path, args);
      `
	e := New("ExecuteFileTest")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("results")
	assert.Nil(t, err)
	assert.True(t, retVal.IsObject())
	retValAsInterface, err := retVal.Export()
	assert.Nil(t, err)
	realRetVal := retValAsInterface.(VMExecResponse)
	assert.Nil(t, err)
	assert.Equal(t, "GNU/Linux", realRetVal.Stdout[0])
}

func TestVMFileAccessTime(t *testing.T) {
	testScript := `
    var file = "/etc/passwd";
    var results1 = FileAccessTime(file);
      `
	e := New("GetAccessTimeTest")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal1, err := e.VM.Get("results1")
	assert.Nil(t, err)
	_, er := retVal1.ToString()
	assert.Nil(t, er)
	//assert.Equal(t, "asdf", retValAsString1)
}

func TestVMFileModifyTime(t *testing.T) {
	testScript := `
    var file = "/etc/passwd";
    var results1 = FileModifyTime(file);
      `
	e := New("GetModifyTimeTest")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal1, err := e.VM.Get("results1")
	assert.Nil(t, err)
	_, er := retVal1.ToString()
	assert.Nil(t, er)
	//assert.Equal(t, "asdf", retValAsString1)
}

func TestVMFileChangeTime(t *testing.T) {
	testScript := `
    var file = "/etc/passwd";
    var results1 = FileChangeTime(file);
      `
	e := New("GetChangeTimeTest")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal1, err := e.VM.Get("results1")
	assert.Nil(t, err)
	_, er := retVal1.ToString()
	assert.Nil(t, er)
	//assert.Equal(t, "asdf", retValAsString1)
}

func TestVMFileBirthTime(t *testing.T) {
	testScript := `
    var file = "/etc/passwd";
    var results1 = FileBirthTime(file);
      `
	e := New("GetChangeTimeTest")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal1, err := e.VM.Get("results1")
	assert.Nil(t, err)
	_, er := retVal1.ToString()
	assert.Nil(t, er)
	//assert.Equal(t, "asdf", retValAsString1)
}

func TestVMFileContains(t *testing.T) {
	testScript := `
    var wordz = "root";
    var file = "/etc/passwd";
    var return_value = FileContains(file, wordz);
  `
	e := New("FileContainsTest")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMDirExists(t *testing.T) {
	testScript := `
    var name = "/etc/";
    var return_value = DirExists(name);
  `
	e := New("TestDirExists")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMFileExists(t *testing.T) {
	testScript := `
    var name = "/etc/passwd";
    var return_value = FileExists(name);
  `
	e := New("TestFileExists")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMCanReadFile(t *testing.T) {
	testScript := `
    var name = "/etc/passwd";
    var return_value = CanReadFile(name);
  `
	e := New("TestCanReadFile")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMCanWriteFile(t *testing.T) {
	testScript := `
    var name = "/etc/passwd";
    var return_value = CanWriteFile(name);
  `
	e := New("TestCanWriteFile")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMCanExecFile(t *testing.T) {
	testScript := `
    var name = "/bin/bash";
    var return_value = CanExecFile(name);
  `
	e := New("TestCanExecFile")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}
