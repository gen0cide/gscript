
# Debugging

## Using The Interactive Debugger

GSCRIPT includes the ability to drop into an interactive debugger using `DebugConsole()` function within a GSCRIPT. 

To enable this feature the `--enable-debugging` must be specified on the command line. Otherwise this function will be skipped in the final binary. Further, obfuscation can not be used when using the `--enable-debugging` command line flag. 

### Useful Debugging Compiler Flags
`--enable-debugging` - This is compiler flag that is used for dropping into the interactive GSCRIPT interpreter, and testing functions in GSCRIPT.
`enable-import-all-native-funcs` - This can be extremely useful for testing native GoLang funcs in the interactive debugger, before writing a script that calls them. Normally the vm will only be compiled with the functions that are used in the script or namespace. 
`keep-build-dir` - This can be helpful for reviewing the intermediate representation before the source code is compiled down to a native binary. This is what comes from the templating and is useful for debugging compiler errors

### Useful Debugging functions
`DebugConsole()` - This is the GSCRIPT function that will drop a running GSCRIPT into the interactive debugging console. It is often helpful to put these where you want to start debugging, after importing and using some functions. 
`TypeOf(object)` - This is a function only exposed in the interactive GSCRIPT debugger, but it will dump the GoLang types of the target object, letting one debug type confusion issues between javascript and GoLang native functions. 

### Debugger Examples
