package compiler

// StringDef represents an instance of an obfuscated string within
// the gscript compilers intermediate representation
type StringDef struct {
	// unique ID of the string def in relation to the target source tree
	ID string `json:"id"`

	// original value of the string as defined in source
	Value string `json:"value"`

	// key used to encrypt string with
	Key rune `json:"key"`

	// the encrypted data to represent this string
	Data []rune `json:"data"`
}
