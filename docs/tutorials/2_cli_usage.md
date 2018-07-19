# CLI Usage
## GSCRIPT Basic Flags
`--help : prints the help menu`
`--version : prints the version of gscript`
`--debug : prints verbose output when the cli tool runs`

## GSCRIPT Basic Command
`docs : will generate the docs localy for gscript`
`templates : will generate a starter gscript`
`shell : will drop the user into the gscript repl`
`compile : the essential command for compiling gscripts into binaries`

### Docs Options
`docs macros : will generate macro docs`
`docs scripts : will generate script docs`
`docs stdlib : will generate stdlib docs`
### Templates Options
`templates list : will list all available templates`
`templates show : will show the specified template`
### Shell Options
`shell : will drop the user into the gscript repl`
### Compile Options
`compile --os : for specifying the platform to compile the target binary to`
`compile --arch : for specifying the archatecture to compile the target binary to`
`compile --output-file : for specifying an output file to save the final binary as opposed to a temp path`
`compile --keep-build-dir : for keeping the build directory that contains the intermediate golang to debug any gscript compiler issues`
`compile --enable-human-readable-names : this options is useful for debugging any potential linker issues.`
`compile --enable-debugging : This option only works with obfuscation disabled, it will drop the running vm into an interactive debugger with the fist instance of DebugConsole()`
`compile --enable-import-all-native-funcs : this options is useful for importing and linking all of the golang native libraries for when the user specifies the --enable-debugging option`
`compile --enable-logging : This flag is helpful for logging all of the console.log() functions output to stdout. Will only work when obfuscation is disabled.`
`compile --obfuscation-level : This flag takes a vale of 0-3 with 0 being the highest level obuscation and 3 being no obfuscation applied to the final binary.`
`compile --disable-native-compilation : This flag will specifically not compile the intermediate representation to a native binary (default is false)`
`compile --enable-test-build : This will enable the test harness in the build - for testing only! (default: false)`
`compile --enable-upx-compression : This will ompress the final binary using UPX (default: false)`


## GSCRIPT Basic Usage
With lots of debugging and logging enabled:
`gscript --debug compile --enable-logging --obfuscation-level 3 /path/to/gscript.gs`
With max obfuscation and no debugging:
`gscript c /path/to/gscript.gs`
With multiple gscripts:
`gscript c /path/to/gscripts/*.gs`


