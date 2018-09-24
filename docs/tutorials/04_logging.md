# Logging
## About Logging

By default, gscript will disable and suppress all javascript logging, using 

```sh
>  console.log("something to log") . 
```

Logging is done by outputting the logged statements, and their associated javascript VM to the stdout where the compiled Genesis binary is running. 

This is done to protect the final binaries from reverse engineering and logging will only be included in the final binary when enabled.

## Enabling Logging

To enable logging simple provide the `gscript compile` subcommand the `--enable-logging` switch. This will by default disable obfuscation and won't work when obfuscation has been explicitly enabled. Obfuscation disabled is equivalent to the cli flag `--obfuscation-level 3`.

### Example

An example of this would look like:

```sh
> gscript compile --obfuscation-level 3 --enable-logging /path/to/gscript.gs
```

### Compiler Logging

If you would like to get more verbose logs from the compiler itself, while building your GSCRIPTs and bundling your assets then you should use the `--debug` flag on the GSCRIPT binary before the subcomand. This will produce verbose logs on all GSCRIPT subcommands. 
