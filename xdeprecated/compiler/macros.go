package compiler

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_PRIORITY = 100
	DEFAULT_TIMEOUT  = 30
)

type MacroList struct {
	Priority    int
	Timeout     int
	LocalFiles  []string
	RemoteFiles []string
}

func ParseMacros(script string, logger *logrus.Entry) *MacroList {
	ml := &MacroList{
		Priority:    DEFAULT_PRIORITY,
		Timeout:     DEFAULT_TIMEOUT,
		LocalFiles:  []string{},
		RemoteFiles: []string{},
	}

	lines := strings.Split(script, "\n")
	for _, l := range lines {
		r := regexp.MustCompile(`\A//import:(.*?)\z`)
		matches := r.FindAllString(l, -1)
		for _, rawF := range matches {
			f := strings.TrimSpace(strings.Replace(rawF, "//import:", "", -1))
			if _, err := os.Stat(f); os.IsNotExist(err) {
				logger.Fatalf("Asset file does not exist: %s", f)
				continue
			}
			ml.LocalFiles = append(ml.LocalFiles, f)
		}
		r = regexp.MustCompile(`\A//priority:(.*?)\z`)
		matches = r.FindAllString(l, -1)
		for _, rawF := range matches {
			f := strings.TrimSpace(strings.Replace(rawF, "//priority:", "", -1))
			p, err := strconv.Atoi(f)
			if err != nil {
				logger.Fatalf("Priority macro is not a number")
			} else {
				ml.Priority = p
			}
		}
		r = regexp.MustCompile(`\A//timeout:(.*?)\z`)
		matches = r.FindAllString(l, -1)
		for _, rawF := range matches {
			f := strings.TrimSpace(strings.Replace(rawF, "//timeout:", "", -1))
			p, err := strconv.Atoi(f)
			if err != nil {
				logger.Fatalf("Priority macro is not a number")
			} else {
				ml.Timeout = p
			}
		}
		r = regexp.MustCompile(`\A//url_import:(.*?)\z`)
		matches = r.FindAllString(l, -1)
		for _, rawF := range matches {
			f := strings.TrimSpace(strings.Replace(rawF, "//url_import:", "", -1))
			u, err := url.Parse(f)
			if err != nil {
				logger.Fatalf("Could not parse URL: %s", err.Error())
			}
			filename := path.Base(u.Path)

			randVector := RandUpperAlphaString(12)

			dir, err := ioutil.TempDir("", randVector)
			if err != nil {
				logger.Fatalf("Could create temp directory: %s", err.Error())
			}
			filePath := filepath.Join(dir, filename)
			out, err := os.Create(filePath)
			if err != nil {
				logger.Fatalf("Could create temp file: %s", err.Error())
			}
			defer out.Close()

			resp, err := http.Get(u.String())
			if err != nil {
				logger.Fatalf("Could not retreive URL: %s", err.Error())
			}
			defer resp.Body.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				logger.Fatalf("Could not save temp file: %s", err.Error())
			}

			ml.RemoteFiles = append(ml.RemoteFiles, filePath)
		}
	}
	return ml
}
