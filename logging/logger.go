package logging

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	TraceEnabled = false
)

type LogWriter struct {
	Name string
}

func (w LogWriter) Write(p []byte) (int, error) {
	output := fmt.Sprintf(
		"%s%s%s%s%s %s",
		color.HiWhiteString("["),
		color.HiRedString("GSCRIPT"),
		color.HiWhiteString(":"),
		color.HiYellowString(strings.ToLower(w.Name)),
		color.HiWhiteString("]"),
		string(p),
	)
	fmt.Fprintf(color.Output, "%s", output)
	return len(output), nil
}

type GSEFormatter struct{}

func (g *GSEFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var logLvl, logLine, tracer string
	var buffer bytes.Buffer
	var deleteTrace = false
	if val, ok := entry.Data["trace"]; ok {
		deleteTrace = true
		if val == "true" {
			buffer.WriteString(" ")
			if pc, file, line, ok := runtime.Caller(4); ok {
				fName := runtime.FuncForPC(pc).Name()
				baseFile := path.Base(file)
				buffer.WriteString(fmt.Sprintf("func=%s source=%s:%s", color.GreenString(fName), color.GreenString(baseFile), color.GreenString("%d", line)))
			} else {
				buffer.WriteString("")
			}
		}
	}

	if deleteTrace {
		delete(entry.Data, "trace")
	}

	names := make([]string, 0, len(entry.Data))
	for name := range entry.Data {
		names = append(names, name)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	for _, name := range names {
		val := entry.Data[name]
		buffer.WriteString(" ")
		buffer.WriteString(fmt.Sprintf("%s=", name))
		buffer.WriteString(color.GreenString(val.(string)))
	}

	tracer = buffer.String()

	switch entry.Level.String() {
	case "debug":
		logLvl = color.HiCyanString(strings.ToUpper(entry.Level.String()))
		logLine = color.CyanString(entry.Message)
	case "info":
		logLvl = color.HiWhiteString(strings.ToUpper(entry.Level.String()))
		logLine = color.WhiteString(entry.Message)
	case "warning":
		logLvl = color.HiBlueString(strings.ToUpper(entry.Level.String()))
		logLine = color.BlueString(entry.Message)
	case "error":
		logLvl = color.HiYellowString(strings.ToUpper(entry.Level.String()))
		logLine = color.YellowString(entry.Message)
	case "fatal":
		logLvl = color.HiRedString(strings.ToUpper(entry.Level.String()))
		logLine = color.RedString(entry.Message)
	default:
		logLvl = color.HiGreenString(strings.ToUpper(entry.Level.String()))
		logLine = color.GreenString(entry.Message)
	}
	line := fmt.Sprintf(
		"%s %s %s%s\n",
		entry.Time.Format(time.RFC822Z),
		logLvl,
		logLine,
		tracer,
	)
	return []byte(line), nil
}
