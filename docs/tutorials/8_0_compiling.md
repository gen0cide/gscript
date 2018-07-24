
# Compiling
Compiling is as simple as running the command line utility, with any number of syntactically valid gscripts at the end. 

### Compiler options
`--os : for specifying the platform to compile the target binary to`

`--arch : for specifying the archatecture to compile the target binary to`

`--output-file : for specifying an output file to save the final binary as opposed to a temp path`

`--keep-build-dir : for keeping the build directory that contains the intermediate golang to debug any gscript compiler issues`

`--enable-human-readable-names : this options is useful for debugging any potential linker issues.`

`--enable-debugging : This option only works with obfuscation disabled, it will drop the running vm into an interactive debugger with the fist instance of DebugConsole()`

`--enable-import-all-native-funcs : this options is useful for importing and linking all of the golang native libraries for when the user specifies the --enable-debugging option`

`--enable-logging : This flag is helpful for logging all of the console.log() functions output to stdout. Will only work when obfuscation is disabled.`

`--obfuscation-level : This flag takes a vale of 0-3 with 0 being the highest level obuscation and 3 being no obfuscation applied to the final binary.`

`--disable-native-compilation : This flag will specifically not compile the intermediate representation to a native binary (default is false)`

`--enable-test-build : This will enable the test harness in the build - for testing only! (default: false)`

`--enable-upx-compression : This will ompress the final binary using UPX (default: false)`

## Compiler examples
`gscript compile /path/to/gscript.gs`
`gscript c /path/to/gscript.gs`
`gscirpt c --enable-logging --obfuscation-level 3 /path/to/gscript.gs`
`gscirpt c --enable-debugging --obfuscation-level 3 /path/to/gscript.gs`
`gscript compile /gscripts/technique1.gs /gscripts/technique2.gs`
`gscript c /gscripts/technique1.gs /gscripts/technique2.gs`
`gscript compile /gscripts/*.gs`
`gscript c /gscripts/*.gs`


## Compiler details
The compiler executes a complex series of tasks, as defined below:
- CheckForConfigErrors
Validated compiler configuration

- CreateBuildDir
Creates the build dir

- ProcessMacros
Parses and processes compiler macros

- InitializeImports
 checks the import references in golang

- DetectVersions
 entry points located within scripts

- GatherAssets
asset tree built

 - WalkGenesisASTs
genesis scripts analyzed

 - LocateGoDependencies
native dependencies resolved

- BuildGolangASTs
native code bundles mapped to the virtual machine

 - SanityCheckScriptToNativeMapping
script callers for native code validated

- SwizzleNativeCalls
native code dynamically linked to the genesis virtual machine

 - SanityCheckSwizzles
dynamic link correctness validated

- WritePreloads
built in genesis helper library injected

- WriteScripts
scripts staged for compilation

 - EncodeAssets
assets encrypted and embedded into the genesis VMs

- WriteVMBundles
virtual machines compiled into intermediate representation

 - CreateEntryPoint
genesis vm callers embedded into final binary entry point

- PerformPreCompileObfuscation
pre-obfuscation completed (stylist tangled all hairs)

- BuildNativeBinary
statically linked native binary built

- PerformPostCompileObfuscation
post-obfuscation completed (mordor has assaulted the binary)