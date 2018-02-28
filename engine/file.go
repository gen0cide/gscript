package engine

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func LocalFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

//These two functions break compiling on windows
/*
func LocalFileWritable(path string) bool {
  return unix.Access(path, unix.W_OK) == nil
}

func LocalFileExecutable(path string) bool {
  return unix.Access(path, unix.X_OK) == nil
}
*/

func LocalDirCreate(path string) error {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		return err
	}
	return nil
}

func LocalDirRemoveAll(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	err = os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}

func LocalFileDelete(path string) error {
	if LocalFileExists(path) {
		err := os.Remove(path)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("The file dosn't exist to delete")
}

func LocalFileCreate(path string, bytes []byte) error {
	if LocalFileExists(path) {
		return errors.New("The file to create already exists so we won't overwite it")
	}
	err := ioutil.WriteFile(path, bytes, 0700)
	if err != nil {
		return err
	}
	return nil
}

func LocalFileAppendBytes(filename string, bytes []byte) error {
	if LocalFileExists(filename) {
		fileInfo, err := os.Stat(filename)
		if err != nil {
			return err
		}
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, fileInfo.Mode())
		if err != nil {
			return err
		}
		if _, err = file.Write(bytes); err != nil {
			return err
		}
		file.Close()
		return nil
	}
	err := LocalFileCreate(filename, bytes)
	if err != nil {
		return err
	}
	return nil
}

func LocalFileAppendString(input, filename string) error {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_APPEND, fileInfo.Mode())
	if err != nil {
		return err
	}
	if _, err = file.WriteString(input); err != nil {
		return err
	}
	file.Close()
	return nil
}

func LocalFileReplace(file, match, replacement string) error {
	if LocalFileExists(file) {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return err
		}
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		lines := strings.Split(string(contents), "\n")
		for index, line := range lines {
			if strings.Contains(line, match) {
				lines[index] = strings.Replace(line, match, replacement, 10)
			}
		}

		ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), fileInfo.Mode())
		return nil
	} else {
		return errors.New("The file to read does not exist")
	}
}

func LocalFileReplaceMulti(file string, matches []string, replacement string) error {
	if LocalFileExists(file) {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return err
		}
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		lines := strings.Split(string(contents), "\n")
		for index, line := range lines {
			for _, match := range matches {
				if strings.Contains(line, match) {
					lines[index] = replacement
				}
			}
		}
		ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), fileInfo.Mode())
		return nil
	} else {
		return errors.New("The file to read does not exist")
	}
}

func LocalFileRead(path string) ([]byte, error) {
	if LocalFileExists(path) {
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return dat, nil
	}
	return nil, errors.New("The file to read does not exist")
}

func XorFiles(file1 string, file2 string, outPut string) error {
	dat1, err := ioutil.ReadFile(file1)
	if err != nil {
		return err
	}
	dat2, err := ioutil.ReadFile(file2)
	if err != nil {
		return err
	}
	dat3 := XorBytes(dat1[:], dat2[:])
	err = LocalFileCreate(outPut, dat3[:])
	if err != nil {
		return err
	}
	return nil
}

func LocalCopyFile(src, dst string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}
