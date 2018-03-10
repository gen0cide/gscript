package logging

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/gen0cide/gscript"
	"github.com/sirupsen/logrus"
	"github.com/vigneshuvi/GoDateFormat"
)

var (
	TraceEnabled = false

	reg regexp.Regexp
)

func init() {
	reg, err := regexp.Compile("[[:^alpha:]]+")
	if err != nil {
		panic(err)
	}
	_ = reg
}

type LogWriter struct {
	Name string
}

func strip(src string) string {
	return reg.ReplaceAllString(src, "")
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

type GSEStrippedFormatter struct{}

type GSEFormatter struct{}

func (g *GSEStrippedFormatter) Format(entry *logrus.Entry) ([]byte, error) {
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

func (g *GSEFormatter) Format(entry *logrus.Entry) ([]byte, error) {
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
		logLvl = color.HiCyanString(strings.ToUpper(entry.Level.String()))
		logLine = color.CyanString(entry.Message)
	case "info":
		logLvl = color.HiWhiteString(strings.ToUpper(entry.Level.String()))
		logLine = color.WhiteString(entry.Message)
	case "warning":
		logLvl = color.HiBlueString(strings.ToUpper(entry.Level.String()))
		logLine = color.HiBlueString(entry.Message)
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
		entry.Time.Format(GoDateFormat.ConvertFormat("yyyy-MM-dd HH:MM:SS tt (Z)")),
		logLvl,
		logLine,
		tracer,
	)
	return []byte(line), nil
}

func PrintLogo() {
	fmt.Fprintf(color.Output, "%s\n", AsciiLogo())
}

func AsciiLogo() string {
	lines := []string{
		color.YellowString("***********************************************************"),
		color.HiWhiteString("                             ____                         "),
		color.HiWhiteString("                     __,-~~/~    `---.                    "),
		color.HiWhiteString("                   _/_,---(      ,    )                   "),
		color.HiWhiteString("               __ /        <    /   )  \\___               "),
		color.HiWhiteString("- ------===;;;'====------------------===;;;===----- -  -  "),
		color.HiWhiteString("                  \\/  ~\"~\"~\"~\"~\"~\\~\"~)~\"/                 "),
		color.HiWhiteString("                  (_ (   \\  (     >    \\)                 "),
		color.HiWhiteString("                   \\_( _ <         >_>'                   "),
		color.HiWhiteString("                      ~ `-i' ::>|--\"                      "),
		color.HiWhiteString("                          I;|.|.|                         "),
		color.HiWhiteString("                         <|i::|i|`.                       "),
		fmt.Sprintf("            %s          %s          %s  ", color.HiGreenString("uL"), color.HiWhiteString("(` ^'\"`-' \")"), color.HiYellowString(")")),
		fmt.Sprintf("        %s          %s  ", color.HiGreenString(".ue888Nc.."), color.HiYellowString("(   (          ( /(")),
		fmt.Sprintf("       %s  %s  ", color.HiGreenString("d88E`\"888E`"), color.HiYellowString("(    (  )(  )\\  `  )   )\\())")),
		fmt.Sprintf("       %s   %s  ", color.HiGreenString("888E  888E"), color.YellowString(")\\   )\\(()\\((_) /(/(  (_))/")),
		fmt.Sprintf("       %s  %s   ", color.HiGreenString("888E  888E"), color.HiRedString("((_) ((_)((_)(_)((_)_\\ | |_")),
		fmt.Sprintf("       %s  %s  ", color.HiGreenString("888E  888E"), color.RedString("(_-</ _|| '_|| || '_ \\)|  _|")),
		fmt.Sprintf("       %s  %s %s ", color.HiGreenString("888& .888E"), color.RedString("/__/\\__||_|  |_|| .__/  \\__|"), color.WhiteString(gscript.Version)),
		fmt.Sprintf("       %s                  %s           ", color.HiGreenString("*888\" 888&"), color.RedString("|_|")),
		fmt.Sprintf("        %s  %s        -- By --", color.HiGreenString("`\"   \"888E"), color.HiWhiteString("G E N I S I S")),
		fmt.Sprintf("       %s   %s       %s", color.HiGreenString(".dWi   `88E"), color.HiWhiteString("S C R I P T I N G"), color.CyanString("gen0cide")),
		fmt.Sprintf("       %s    %s            %s", color.HiGreenString("4888~  J8%%"), color.HiWhiteString("E N G I N E"), color.CyanString("ahhh")),
		fmt.Sprintf("        %s             ", color.HiGreenString("^\"===*\"`")),
		"                github.com/gen0cide/gscript",
		color.YellowString("***********************************************************"),
	}

	return strings.Join(lines, "\n")
}
