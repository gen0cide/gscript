package translator

// type BuiltInConverter func(t string) (string, error)

var (
	//BuiltInMap defines various mappings that need to exist between the virtual machine and golang
	BuiltInMap = map[string]string{
		"int":     "int64",
		"uintptr": "int64",
	}
)

// type Translator struct {
// 	Logger logger.Logger
// }

// func convertInt()
