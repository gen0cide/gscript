package null

// Logger is a built in type that implements the Logger interface to prevent scripts
// from writing output to the screen during execution (default logging behavior of binary)
type Logger struct{}

// Print implements the Logger interface type to prevent debug output
func (n *Logger) Print(args ...interface{}) {
	_ = args
	return
}

// Printf implements the Logger interface type to prevent debug output
func (n *Logger) Printf(format string, args ...interface{}) {
	_ = format
	_ = args
	return
}

// Println implements the Logger interface type to prevent debug output
func (n *Logger) Println(args ...interface{}) {
	_ = args
	return
}

// Debug implements the Logger interface type to prevent debug output
func (n *Logger) Debug(args ...interface{}) {
	_ = args
	return
}

// Debugf implements the Logger interface type to prevent debug output
func (n *Logger) Debugf(format string, args ...interface{}) {
	_ = format
	_ = args
	return
}

// Debugln implements the Logger interface type to prevent debug output
func (n *Logger) Debugln(args ...interface{}) {
	_ = args
	return
}

// Info implements the Logger interface type to prevent debug output
func (n *Logger) Info(args ...interface{}) {
	_ = args
	return
}

// Infof implements the Logger interface type to prevent debug output
func (n *Logger) Infof(format string, args ...interface{}) {
	_ = format
	_ = args
	return
}

// Infoln implements the Logger interface type to prevent debug output
func (n *Logger) Infoln(args ...interface{}) {
	_ = args
	return
}

// Warn implements the Logger interface type to prevent debug output
func (n *Logger) Warn(args ...interface{}) {
	_ = args
	return
}

// Warnf implements the Logger interface type to prevent debug output
func (n *Logger) Warnf(format string, args ...interface{}) {
	_ = format
	_ = args
	return
}

// Warnln implements the Logger interface type to prevent debug output
func (n *Logger) Warnln(args ...interface{}) {
	_ = args
	return
}

// Error implements the Logger interface type to prevent debug output
func (n *Logger) Error(args ...interface{}) {
	_ = args
	return
}

// Errorf implements the Logger interface type to prevent debug output
func (n *Logger) Errorf(format string, args ...interface{}) {
	_ = format
	_ = args
	return
}

// Errorln implements the Logger interface type to prevent debug output
func (n *Logger) Errorln(args ...interface{}) {
	_ = args
	return
}

// Fatal implements the Logger interface type to prevent debug output
func (n *Logger) Fatal(args ...interface{}) {
	_ = args
	return
}

// Fatalf implements the Logger interface type to prevent debug output
func (n *Logger) Fatalf(format string, args ...interface{}) {
	_ = format
	_ = args
	return
}

// Fatalln implements the Logger interface type to prevent debug output
func (n *Logger) Fatalln(args ...interface{}) {
	_ = args
	return
}
