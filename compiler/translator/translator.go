package translator

// type BuiltInConverter func(t string) (string, error)

var (
	//BuiltInMap defines various mappings that need to exist between the virtual machine and golang
	BuiltInMap = map[string]string{
		"int":     "int64",
		"uintptr": "int64",
	}

	// TypeAliasMap maps a package to it's various type alias conversion types
	TypeAliasMap = map[string]TypeAliasToBuiltIn{
		"syscall": TypeAliasToBuiltIn{
			"Signal": "int",
		},
	}
)

// TypeAliasToBuiltIn maps package type aliases to built in types for conversion
type TypeAliasToBuiltIn map[string]string
