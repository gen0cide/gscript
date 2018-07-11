package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteFileFromBytes(t *testing.T) {
	err := WriteFileFromBytes("test_file.txt", []byte("just another test"))
	assert.Nil(t, err)
	contents, err := ioutil.ReadFile("test_file.txt")
	assert.Nil(t, err)
	assert.Equal(t, "just another test", string(contents), "they should be equal")
	err = os.Remove("test_file.txt")
	assert.Nil(t, err)
}

func TestWriteFileFromString(t *testing.T) {

	err := WriteFileFromString("test_file.txt", "this is a test")
	assert.Nil(t, err)
	contents, err := ioutil.ReadFile("test_file.txt")
	if err == nil {
		assert.Equal(t, "this is a test", string(contents), "they should be equal")
		//err = os.Remove("test_file.txt")
		//if err != nil {
		//		assert.Nil(t, "file not deleted")
		//}
	} else {
		assert.Nil(t, "bad news")
	}
}

func TestReadFileAsString(t *testing.T) {
	fileString, err := ReadFileAsString("test_file.txt")
	assert.Nil(t, err)
	assert.NotNil(t, fileString)
	assert.Equal(t, "this is a test", fileString, "they should be equal")
}

func TestReadFileAsBytes(t *testing.T) {
	bytes, err := ReadFileAsBytes("test_file.txt")
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	assert.Equal(t, "this is a test", string(bytes), "they should be equal")
}

func TestCopyFile(t *testing.T) {
	dataCopied, err := CopyFile("test_file.txt", "new_test_file.txt")
	assert.Nil(t, err)
	assert.Equal(t, 14, dataCopied, "they should be equal")
}

func TestAppendFileBytes(t *testing.T) {
	firstData, _ := ioutil.ReadFile("test_file.txt")
	fmt.Println(string(firstData))
	err := AppendFileBytes("test_file.txt", []byte("some new data"))
	assert.Nil(t, err)
	secondData, _ := ioutil.ReadFile("test_file.txt")
	fmt.Println(string(secondData))
	assert.NotEqual(t, firstData, secondData, "Should not be equal")
}

func TestAppendFileString(t *testing.T) {
	firstData, _ := ioutil.ReadFile("test_file.txt")
	fmt.Println(string(firstData))
	err := AppendFileString("test_file.txt", "\nsome new data\n")
	assert.Nil(t, err)
	secondData, _ := ioutil.ReadFile("test_file.txt")
	fmt.Println(string(secondData))
	assert.NotEqual(t, firstData, secondData, "Should not be equal")
}

func TestReplaceInFileWithString(t *testing.T) {
	replaced, err := ReplaceInFileWithString("test_file.txt", "data", "replaced")
	assert.Nil(t, err)
	assert.NotNil(t, replaced)
}

func TestReplaceInFileWithRegex(t *testing.T) {
	replaced, err := ReplaceInFileWithRegex("test_file.txt", "(place)", "regexed")
	assert.Nil(t, err)
	assert.NotNil(t, replaced)
}

func TestSetPerms(t *testing.T) {
	err := SetPerms("test_file.txt", 0777)
	assert.Nil(t, err)
}
