package file

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// WriteFileFromBytes writes data from a byte array to a dest filepath with the dest parent dirs permissions.
func WriteFileFromBytes(destPath string, fileData []byte) error {
	absDir, err := filepath.Abs(filepath.Dir(destPath))
	if err != nil {
		return err
	}
	dirInfo, err := os.Stat(absDir)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(destPath, fileData, os.FileMode(dirInfo.Mode()))
	if err != nil {
		return err
	}
	return nil
}

// WriteFileFromString writes data from a string to a dest filepath with the dest parent dirs permissions.
func WriteFileFromString(destPath string, fileData string) error {
	absDir, err := filepath.Abs(filepath.Dir(destPath))
	if err != nil {
		return err
	}
	dirInfo, err := os.Stat(absDir)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(destPath, []byte(fileData), os.FileMode(dirInfo.Mode()))
	if err != nil {
		return err
	}
	return nil
}

//ReadFileAsString takes a file path and reads that files contents and returns a string representation of the contents
func ReadFileAsString(readPath string) (string, error) {
	absPath, err := filepath.Abs(readPath)
	if err != nil {
		return "", err
	}
	contents, err := ioutil.ReadFile(absPath)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

//ReadFileAsBytes takes a file path and reads that files contents and returns a byte array of the contents
func ReadFileAsBytes(readPath string) ([]byte, error) {
	absPath, err := filepath.Abs(readPath)
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

// CopyFile copies a file from the src to the dest with the original files permissions
func CopyFile(srcPath, destPath string) (int, error) {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return 0, err
	}
	from, err := os.Open(srcPath)
	if err != nil {
		return 0, err
	}
	defer from.Close()

	to, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE, os.FileMode((srcInfo.Mode())))
	if err != nil {
		return 0, err
	}
	defer to.Close()

	dataCopied, err := io.Copy(to, from)
	if err != nil {
		return 0, err
	}
	return int(dataCopied), nil
}

// AppendFileBytes takes a file and adds a byte array to the end of it
func AppendFileBytes(targetPath string, addData []byte) error {
	fileInfo, err := os.Stat(targetPath)
	if err != nil {
		return err
	}
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return err
	}
	targetFile, err := os.OpenFile(absPath, os.O_APPEND|os.O_WRONLY, fileInfo.Mode())
	if err != nil {
		return err
	}
	if _, err = targetFile.Write(addData); err != nil {
		return err
	}
	targetFile.Close()
	return nil
}

// AppendFileString takes a file and adds a string to the end of it
func AppendFileString(targetPath, addString string) error {
	fileInfo, err := os.Stat(targetPath)
	if err != nil {
		return err
	}
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return err
	}
	targetFile, err := os.OpenFile(absPath, os.O_APPEND|os.O_WRONLY, fileInfo.Mode())
	if err != nil {
		return err
	}
	if _, err = targetFile.WriteString(addString); err != nil {
		return err
	}
	targetFile.Close()
	return nil
}

//ReplaceInFileWithString searches a file for a string and replaces each instance found of that string. Returns the amount of strings replaced
func ReplaceInFileWithString(file, match, replacement string) (int, error) {
	fileInfo, err := os.Stat(file)
	if os.IsNotExist(err) {
		return 0, err
	}
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}
	var count int = 0
	lines := strings.Split(string(contents), "\n")
	for index, line := range lines {
		if strings.Contains(line, match) {
			lines[index] = strings.Replace(line, match, replacement, 10)
			count++
		}
	}
	ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), fileInfo.Mode())
	return count, nil
}

//ReplaceInFileWithRegex searches a file for a string and replaces each instance found of that string. Returns the amount of strings replaced
func ReplaceInFileWithRegex(file string, regexString string, replaceWith string) (int, error) {
	re := regexp.MustCompile(regexString)
	fileInfo, err := os.Stat(file)
	if os.IsNotExist(err) {
		return 0, err
	}
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}
	var count int = 0
	lines := strings.Split(string(contents), "\n")
	for index, line := range lines {
		if re.MatchString(line) {
			lines[index] = re.ReplaceAllString(line, replaceWith)
			count++
		}
	}
	ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), fileInfo.Mode())
	return count, nil
}

//SetPerms changes the file permissions of a givin file
func SetPerms(targetPath string, perms int64) error {
	err := os.Chmod(targetPath, os.FileMode(uint32(perms)))
	if err != nil {
		return err
	}
	return nil
}
