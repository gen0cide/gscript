package standard

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/vigneshuvi/GoDateFormat"
)

var (
	traceEnabled = false
	reg          regexp.Regexp
	debugLevel   = color.New(color.FgCyan, color.Bold)
	infoLevel    = color.New(color.FgHiWhite, color.Bold)
	warningLevel = color.New(color.FgHiBlue, color.Bold)
	errorLevel   = color.New(color.FgHiYellow, color.Bold)
	fatalLevel   = color.New(color.FgRed, color.Bold)
	defaultLevel = color.New(color.FgGreen, color.Bold)

	debugMsg   = color.New(color.FgCyan)
	infoMsg    = color.New(color.FgHiWhite)
	warningMsg = color.New(color.FgHiBlue)
	errorMsg   = color.New(color.FgHiYellow)
	fatalMsg   = color.New(color.FgRed)
	defaultMsg = color.New(color.FgGreen)
)

func init() {
	reg, err := regexp.Compile("[[:^alpha:]]+")
	if err != nil {
		panic(err)
	}
	_ = reg
}

func strip(src string) string {
	return reg.ReplaceAllString(src, "")
}

type standardLogWriter struct {
	name string
	prog string
}

type standardStrippedFormatter struct{}
type standardDefaultFormatter struct{}
type standardDebugFormatter struct{}

func (g *standardStrippedFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(color.HiWhiteString("%s", entry.Message))

	names := make([]string, 0, len(entry.Data))
	for name := range entry.Data {
		names = append(names, name)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	for _, name := range names {
		val := entry.Data[name]
		buf.WriteString(" ")
		buf.WriteString(fmt.Sprintf("%s=", name))
		buf.WriteString(color.GreenString(val.(string)))
	}

	buf.WriteString("\n")

	return buf.Bytes(), nil
}

func (g *standardDebugFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var logLvl, logLine, tracer string
	var buffer bytes.Buffer
	var deleteTrace = false

	if val, ok := entry.Data["stripped"]; ok {
		if val == "true" {
			buffer.WriteString(fmt.Sprintf("%s", entry.Message))
			return buffer.Bytes(), nil
		}
	}

	if val, ok := entry.Data["trace"]; ok {
		deleteTrace = true
		if val == "true" {
			buffer.WriteString(" ")
			if pc, file, line, ok := runtime.Caller(5); ok {
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
		logLvl = color.HiCyanString("DEBUG")
		logLine = color.CyanString(entry.Message)
	case "info":
		logLvl = color.HiWhiteString("INFO ")
		logLine = color.WhiteString(entry.Message)
	case "warning":
		logLvl = color.HiBlueString("WARN ")
		logLine = color.HiBlueString(entry.Message)
	case "error":
		logLvl = color.HiYellowString("ERROR")
		logLine = color.YellowString(entry.Message)
	case "fatal":
		logLvl = color.HiRedString("FATAL")
		logLine = color.RedString(entry.Message)
	default:
		logLvl = color.HiGreenString(strings.ToUpper(entry.Level.String()))
		logLine = color.GreenString(entry.Message)
	}
	line := fmt.Sprintf(
		"%s %s %s%s\n",
		entry.Time.Format(GoDateFormat.ConvertFormat("yyyy-MM-dd HH:MM:SS tt (Z)")),
		logLvl,
		logLine,
		tracer,
	)
	return []byte(line), nil
}

func (g *standardDefaultFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var logLvl, logLine string

	switch entry.Level.String() {
	case "debug":
		logLvl = debugLevel.Sprint("DEBUG")
		logLine = debugMsg.Sprint(entry.Message)
	case "info":
		logLvl = infoLevel.Sprint(" INFO")
		logLine = infoMsg.Sprint(entry.Message)
	case "warning":
		logLvl = warningLevel.Sprint(" WARN")
		logLine = warningMsg.Sprint(entry.Message)
	case "error":
		logLvl = errorLevel.Sprint("ERROR")
		logLine = errorMsg.Sprint(entry.Message)
	case "fatal":
		logLvl = fatalLevel.Sprint("FATAL")
		logLine = fatalMsg.Sprint(entry.Message)
	default:
		logLvl = defaultLevel.Sprintf(" %s ", strings.ToUpper(entry.Level.String()))
		logLine = defaultMsg.Sprint(entry.Message)
	}
	line := fmt.Sprintf(
		"%s %s\n",
		logLvl,
		logLine,
	)
	return []byte(line), nil
}

func (w standardLogWriter) Write(p []byte) (int, error) {
	output := fmt.Sprintf(
		"%s%s%s%s%s %s",
		color.HiWhiteString("["),
		color.HiRedString(w.prog),
		color.HiWhiteString(":"),
		color.HiYellowString(strings.ToLower(w.name)),
		color.HiWhiteString("]"),
		string(p),
	)
	written, err := io.Copy(color.Output, strings.NewReader(output))
	return int(written), err
	//fmt.Fprintf(color.Output, "%s", output)
	//return len(output), nil
}

// Logger implements the Logger interface and is the standard logger for all genesis output
type Logger struct {
	Name   string
	Logger *logrus.Logger
	Prog   string
}

// NewStandardLogger returns a new Logger
func NewStandardLogger(l *logrus.Logger, prog, name string, customOutput, stripped bool) *Logger {
	if l == nil {
		l = logrus.New()
	}
	if customOutput == false {
		l.Out = standardLogWriter{name: name, prog: prog}
		if stripped == true {
			l.Formatter = &standardStrippedFormatter{}
		} else {
			l.Formatter = &standardDefaultFormatter{}
		}
	}
	return &Logger{
		Logger: l,
		Name:   name,
		Prog:   prog,
	}
}

// Print wraps the corrasponding logrus logging function for console output
func (l *Logger) Print(args ...interface{}) {
	l.Logger.Print(args...)
	return
}

// Printf wraps the corrasponding logrus logging function for console output
func (l *Logger) Printf(format string, args ...interface{}) {
	l.Logger.Printf(format, args...)
	return
}

// Println wraps the corrasponding logrus logging function for console output
func (l *Logger) Println(args ...interface{}) {
	l.Logger.Println(args...)
	return
}

// Debug wraps the corrasponding logrus logging function for console output
func (l *Logger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
	return
}

// Debugf wraps the corrasponding logrus logging function for console output
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
	return
}

// Debugln wraps the corrasponding logrus logging function for console output
func (l *Logger) Debugln(args ...interface{}) {
	l.Logger.Debugln(args...)
	return
}

// Info wraps the corrasponding logrus logging function for console output
func (l *Logger) Info(args ...interface{}) {
	l.Logger.Info(args...)
	return
}

// Infof wraps the corrasponding logrus logging function for console output
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
	return
}

// Infoln wraps the corrasponding logrus logging function for console output
func (l *Logger) Infoln(args ...interface{}) {
	l.Logger.Infoln(args...)
	return
}

// Warn wraps the corrasponding logrus logging function for console output
func (l *Logger) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
	return
}

// Warnf wraps the corrasponding logrus logging function for console output
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args...)
	return
}

// Warnln wraps the corrasponding logrus logging function for console output
func (l *Logger) Warnln(args ...interface{}) {
	l.Logger.Warnln(args...)
	return
}

// Error wraps the corrasponding logrus logging function for console output
func (l *Logger) Error(args ...interface{}) {
	l.Logger.Error(args...)
	return
}

// Errorf wraps the corrasponding logrus logging function for console output
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
	return
}

// Errorln wraps the corrasponding logrus logging function for console output
func (l *Logger) Errorln(args ...interface{}) {
	l.Logger.Errorln(args...)
	return
}

// Fatal wraps the corrasponding logrus logging function for console output
func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
	return
}

// Fatalf wraps the corrasponding logrus logging function for console output
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
	return
}

// Fatalln wraps the corrasponding logrus logging function for console output
func (l *Logger) Fatalln(args ...interface{}) {
	l.Logger.Fatalln(args...)
	return
}
