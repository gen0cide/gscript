package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
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
		color.HiYellowString(w.Name),
		color.HiWhiteString("]"),
		string(p),
	)
	fmt.Fprintf(color.Output, "%s", output)
	return len(output), nil
}

type GSEFormatter struct{}

func (g *GSEFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	line := fmt.Sprintf("%s %s> %s\n", color.RedString(strings.ToUpper(entry.Level.String())), entry.Time.Format(time.RFC822Z), color.GreenString(entry.Message))
	return []byte(line), nil
}
