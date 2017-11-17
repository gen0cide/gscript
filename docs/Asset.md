# `Asset(filename)`

Retrieves the packed bytes of file `filename` from the GENESIS packed archive.

### Argument List

 * `file` (String) - Name of the file you imported via your gscript.

### Return Type

`array` - Array of bytes

## Example

```
//import:/tmp/foo.bin

var file_to_drop;

function BeforeDeploy() {
  file_to_drop = Asset("foo.bin")
}

// ....
```

This will automatically compile `foo.bin` into your compiled binary and retrieves it.