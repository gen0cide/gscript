package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteFileFromBytes(t *testing.T) {
	err := WriteFileFromBytes("test_file.txt", []byte("just another test"))
	if err != nil {
		assert.Nil(t, "Write file failed")
	}
	contents, err := ioutil.ReadFile("test_file.txt")
	if err == nil {
		assert.Equal(t, "just another test", string(contents), "they should be equal")
		err = os.Remove("test_file.txt")
		if err != nil {
			assert.Nil(t, "file not deleted")
		}
	} else {
		assert.Nil(t, "bad news")
	}
}

func TestWriteFileFromString(t *testing.T) {

	err := WriteFileFromString("test_file.txt", "this is a test")
	if err != nil {
		assert.Nil(t, "Write file failed")
	}
	contents, err := ioutil.ReadFile("test_file.txt")
	if err == nil {
		assert.Equal(t, "this is a test", string(contents), "they should be equal")
		err = os.Remove("test_file.txt")
		if err != nil {
			assert.Nil(t, "file not deleted")
		}
	} else {
		assert.Nil(t, "bad news")
	}
}
