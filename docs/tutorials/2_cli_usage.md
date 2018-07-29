# GSCRIPT CLI Usage

## Basic Usage

With lots of debugging and logging enabled:

```sh
gscript --debug compile --enable-logging --obfuscation-level 3 /path/to/gscript.gs
```

With max obfuscation and no debugging:

```sh
gscript c /path/to/gscript.gs
```

With multiple gscripts:

```sh
gscript c /path/to/gscripts/*.gs
```

## Basic Flags

```sh
--help     # prints the help menu
--version  # prints the version of gscript
--debug    # prints verbose output when the cli tool runs
```

## Basic Commands

```sh
gscript compile    # The essential command for compiling gscripts into binaries
gscript docs       # Generate the docs localy for gscript
gscript shell      # Drop the user into the gscript repl
gscript templates  # Generate a starter gscript
```

### Docs Options

```sh
gscript docs macros   # Generate macro docs
gscript docs scripts  # Generate script docs
gscript docs stdlib   # Generate stdlib docs
```

### Templates Options

```sh
gscript templates list  # List all available templates
gscript templates show  # Show the specified template
```

### Shell Options

```sh
gscript shell  # will drop the user into the gscript repl
```

### Compile Options

```sh
gscript compile --os                              # For specifying the platform to compile the target binary to
gscript compile --arch                            # For specifying the archatecture to compile the target binary to
gscript compile --output-file                     # For specifying an output file to save the final binary as opposed to a temp path
gscript compile --keep-build-dir                  # For keeping the build directory that contains the intermediate golang to debug any gscript compiler issues
gscript compile --enable-human-readable-names     # This options is useful for debugging any potential linker issues.
gscript compile --enable-debugging                # This option only works with obfuscation disabled, it will drop the running vm into an interactive debugger with the fist instance of DebugConsole()
gscript compile --enable-import-all-native-funcs  # This options is useful for importing and linking all of the golang native libraries for when the user specifies the --enable-debugging option
gscript compile --enable-logging                  # This flag is helpful for logging all of the console.log() functions output to stdout. Will only work when obfuscation is disabled.
gscript compile --obfuscation-level               # This flag takes a vale of 0-3 with 0 being the highest level obuscation and 3 being no obfuscation applied to the final binary.
gscript compile --disable-native-compilation      # This flag will specifically not compile the intermediate representation to a native binary (default: false)
gscript compile --enable-test-build               # This will enable the test harness in the build - for testing only! (default: false)
gscript compile --enable-upx-compression          # This will ompress the final binary using UPX (default: false)
```

