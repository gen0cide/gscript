package compiler

type Options struct {
	Scripts            []string
	Outfile            string
	OS                 string
	Arch               string
	ShowSource         bool
	CompressionEnabled bool
	LoggingEnabled     bool
	PreserveBuildDir   bool
}
