# Using the Standard Library

## Crypto

- GetMD5FromString(data string) string

- GetMD5FromBytes(data []byte) string

- GetSHA1FromString(data string) string

- GetSHA1FromBytes(data []byte) string

- GetSHA256FromString(data string) string

- GetSHA256FromBytes(data []byte) string

- GenerateRSASSHKeyPair(size int) (pubkey string, privkey string, error)

## Encoding

- DecodeBase64(data string) (string, error)

- EncodeBase64(data string) string

- EncodeStringAsBytes(data string) []byte

- EncodeBytesAsString(data []byte) string

## Exec 

- ExecuteCommand(progname string, args []string) (pid int, stdout string, stderr string, exitCode int, err error)

- ExecuteCommandAsync(progname string, args []string) (proc *exec.Cmd, err error)

## File 

- WriteFileFromBytes(filepath string, data []byte) error

- WriteFileFromString(filepath string, data string) error

- ReadFileAsBytes(filepath string) ([]byte, error)

- ReadFileAsString(filepath string) (string, error)

- AppendBytesToFile(filepath string, data []byte) error

- AppendStringToFile(filepath string, data string) error

- CopyFile(srcpath string, dstpath string, perms string) (bytesWritten int, err error)

- ReplaceInFileWithString(match string, new string) error

- ReplaceInFileWithRegex(regexString string, replaceWith string) error

- SetPerms(filepath string, perms string) error

- CheckExists(targetPath string) bool

## Net 

- CheckForInUseTCP(port int) (bool, error)

- CheckForInUseUDP(port int) (bool, error)

## Os 

- TerminateSelf() error

- TerminateVM()

## rand 

- GetInt(min int, max int) int

- GetAlphaNumericString(len int) string

- GetAlphaString(len int) string

- GetAlphaNumericSpecialString(len int) string

- GetBool() bool

## requests 

- PostURL(url string, data string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)

- PostJSON(url string, jsondata string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)

- PostFile(url string, filepath string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)

- GetURLAsString(url string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)

- GetURLAsBytes(url string, headers map[string]string, ignoresslerrors bool) (resp *http.Response, body string, err error)

## time

- GetUnix() int
